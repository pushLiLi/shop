import { defineStore } from 'pinia'

const API_BASE = '/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    conversations: [],
    currentConversation: null,
    messages: [],
    isOpen: false,
    unreadCount: 0,
    isLoading: false,
    pollTimer: null,
    slowPollTimer: null
  }),

  actions: {
    async openPanel() {
      this.isOpen = true
      this.stopPolling()
      await this.fetchOrCreateConversation()
      this.startPolling()
    },

    closePanel() {
      this.isOpen = false
      this.stopPolling()
      this.startSlowPolling()
    },

    async fetchOrCreateConversation() {
      try {
        this.isLoading = true
        const res = await fetch(`${API_BASE}/chat/conversations`, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        this.conversations = data.conversations || []

        const openConv = this.conversations.find(c => c.status === 'open')
        if (openConv) {
          this.currentConversation = openConv
          await this.fetchMessages()
        } else {
          const createRes = await fetch(`${API_BASE}/chat/conversations`, {
            method: 'POST',
            headers: getAuthHeaders()
          })
          const createData = await createRes.json()
          this.currentConversation = createData.conversation
          this.conversations.unshift(createData.conversation)
          this.messages = []
        }
      } catch (e) {
        console.error('Fetch conversation failed:', e)
      } finally {
        this.isLoading = false
      }
    },

    async fetchMessages(afterId) {
      if (!this.currentConversation) return
      try {
        const convId = this.currentConversation.id
        let url = `${API_BASE}/chat/conversations/${convId}/messages`
        if (afterId) {
          url += `?after=${afterId}`
        }
        const res = await fetch(url, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (afterId) {
          if (data.messages && data.messages.length > 0) {
            this.messages = [...this.messages, ...data.messages]
          }
        } else {
          this.messages = data.messages || []
        }
      } catch (e) {
        console.error('Fetch messages failed:', e)
      }
    },

    async sendMessage(content) {
      if (!this.currentConversation || !content.trim()) return
      try {
        const convId = this.currentConversation.id
        const res = await fetch(`${API_BASE}/chat/conversations/${convId}/messages`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({ content: content.trim() })
        })
        const data = await res.json()
        if (data.message) {
          this.messages.push(data.message)
        }
      } catch (e) {
        console.error('Send message failed:', e)
      }
    },

    async fetchUnreadCount() {
      try {
        const res = await fetch(`${API_BASE}/chat/unread-count`, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        this.unreadCount = data.count || 0
      } catch (e) {
        console.error('Fetch unread count failed:', e)
      }
    },

    startPolling() {
      this.stopPolling()
      const poll = () => {
        if (!this.isOpen || !this.currentConversation) return
        const lastId = this.messages.length > 0 ? this.messages[this.messages.length - 1].id : 0
        this.fetchMessages(lastId)
        this.pollTimer = setTimeout(poll, 3000)
      }
      this.pollTimer = setTimeout(poll, 3000)
    },

    stopPolling() {
      if (this.pollTimer) {
        clearTimeout(this.pollTimer)
        this.pollTimer = null
      }
    },

    startSlowPolling() {
      this.stopSlowPolling()
      if (!this.currentConversation) return
      const poll = () => {
        if (this.isOpen) return
        this.fetchUnreadCount()
        this.slowPollTimer = setTimeout(poll, 10000)
      }
      this.slowPollTimer = setTimeout(poll, 10000)
    },

    stopSlowPolling() {
      if (this.slowPollTimer) {
        clearTimeout(this.slowPollTimer)
        this.slowPollTimer = null
      }
    },

    initPolling() {
      this.fetchUnreadCount()
      this.startSlowPolling()
    },

    cleanup() {
      this.stopPolling()
      this.stopSlowPolling()
      this.isOpen = false
      this.currentConversation = null
      this.messages = []
      this.conversations = []
      this.unreadCount = 0
    }
  }
})
