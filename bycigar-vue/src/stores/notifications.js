import { defineStore } from 'pinia'

const API_BASE = '/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useNotificationsStore = defineStore('notifications', {
  state: () => ({
    items: [],
    unreadCount: 0,
    loading: false,
    isOpen: false,
    page: 1,
    totalPages: 1,
    currentNotification: null,
    detailLoading: false
  }),

  actions: {
    openPanel() {
      this.isOpen = true
      if (this.items.length === 0) {
        this.fetchNotifications()
      }
    },

    closePanel() {
      this.isOpen = false
    },

    async fetchNotifications() {
      try {
        this.loading = true
        const res = await fetch(`${API_BASE}/notifications?page=1&limit=50`, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        this.items = data.notifications || []
        this.page = data.page || 1
        this.totalPages = data.totalPages || 1
      } catch (e) {
        console.error('Fetch notifications failed:', e)
      } finally {
        this.loading = false
      }
    },

    async fetchUnreadCount() {
      try {
        const res = await fetch(`${API_BASE}/notifications/unread-count`, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        this.unreadCount = data.count || 0
      } catch (e) {
        console.error('Fetch unread count failed:', e)
      }
    },

    async markAsRead(id) {
      try {
        const res = await fetch(`${API_BASE}/notifications/${id}/read`, {
          method: 'PUT',
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (data.success) {
          const item = this.items.find(n => n.id === id)
          if (item) {
            item.isRead = true
          }
          this.unreadCount = Math.max(0, this.unreadCount - 1)
        }
      } catch (e) {
        console.error('Mark as read failed:', e)
      }
    },

    async markAllRead() {
      try {
        const res = await fetch(`${API_BASE}/notifications/read-all`, {
          method: 'PUT',
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (data.success) {
          this.items.forEach(n => { n.isRead = true })
          this.unreadCount = 0
        }
      } catch (e) {
        console.error('Mark all read failed:', e)
      }
    },

    clear() {
      this.items = []
      this.unreadCount = 0
    },

    async fetchNotification(id) {
      try {
        this.detailLoading = true
        const res = await fetch(`${API_BASE}/notifications/${id}`, {
          headers: getAuthHeaders()
        })
        if (!res.ok) throw new Error('通知不存在')
        const data = await res.json()
        this.currentNotification = data
        const item = this.items.find(n => n.id === Number(id))
        if (item && !item.isRead) {
          item.isRead = true
          this.unreadCount = Math.max(0, this.unreadCount - 1)
        }
      } catch (e) {
        this.currentNotification = null
        throw e
      } finally {
        this.detailLoading = false
      }
    }
  }
})
