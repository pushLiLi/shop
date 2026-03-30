import { defineStore } from 'pinia'
import { useChatStore } from './chat'
import { useNotificationSound } from '../composables/useNotificationSound'

const API_BASE = '/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

let notificationSoundPlay = null
function getNotificationSound() {
  if (!notificationSoundPlay) {
    const { play } = useNotificationSound()
    notificationSoundPlay = play
  }
  return notificationSoundPlay
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
        const newCount = data.count || 0
        if (newCount > this.unreadCount && this.items.length > 0) {
          this.fetchNotifications()
        }
        this.unreadCount = newCount
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

    async deleteNotification(id) {
      try {
        const res = await fetch(`${API_BASE}/notifications/${id}`, {
          method: 'DELETE',
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (data.success) {
          const item = this.items.find(n => n.id === id)
          if (item && !item.isRead) {
            this.unreadCount = Math.max(0, this.unreadCount - 1)
          }
          this.items = this.items.filter(n => n.id !== id)
        }
        return data.success
      } catch (e) {
        console.error('Delete notification failed:', e)
        return false
      }
    },

    async deleteReadNotifications() {
      try {
        const res = await fetch(`${API_BASE}/notifications/read`, {
          method: 'DELETE',
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (data.success) {
          this.items = this.items.filter(n => !n.isRead)
        }
        return data.success
      } catch (e) {
        console.error('Delete read notifications failed:', e)
        return false
      }
    },

    initWSListener() {
      const chatStore = useChatStore()
      chatStore.onMessage = (data) => {
        if (data.type === 'notification') {
          this.handleWSNotification(data.notification)
          try { getNotificationSound()() } catch {}
        }
      }
      this.fetchUnreadCount()
    },

    cleanupWSListener() {
      const chatStore = useChatStore()
      if (chatStore.onMessage) {
        chatStore.onMessage = null
      }
    },

    handleWSNotification(notification) {
      const exists = this.items.some(n => n.id === notification.id)
      if (!exists) {
        this.items.unshift(notification)
        this.unreadCount += 1
      }
    },

    clear() {
      this.cleanupWSListener()
      this.items = []
      this.unreadCount = 0
    },

    async fetchNotification(id) {
      try {
        this.currentNotification = null
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
        throw e
      } finally {
        this.detailLoading = false
      }
    }
  }
})
