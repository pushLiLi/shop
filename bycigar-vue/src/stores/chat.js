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
    ws: null,
    wsReconnectTimer: null,
    wsReconnectDelay: 3000,
    wsConnected: false,
    onMessage: null
  }),

  actions: {
    getWsUrl() {
      const token = localStorage.getItem('token')
      if (!token) return null
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      return `${protocol}//${window.location.host}/api/chat/ws?token=${encodeURIComponent(token)}`
    },

    connectWebSocket() {
      this.disconnectWebSocket()
      const url = this.getWsUrl()
      if (!url) return

      try {
        this.ws = new WebSocket(url)

        this.ws.onopen = () => {
          this.wsConnected = true
          this.wsReconnectDelay = 3000
          this.fetchUnreadCount()
        }

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            this.handleWSMessage(data)
          } catch (e) {
            console.error('Parse WS message failed:', e)
          }
        }

        this.ws.onclose = () => {
          this.wsConnected = false
          this.scheduleReconnect()
        }

        this.ws.onerror = () => {
          this.wsConnected = false
        }
      } catch (e) {
        console.error('WebSocket connect failed:', e)
        this.scheduleReconnect()
      }
    },

    disconnectWebSocket() {
      if (this.wsReconnectTimer) {
        clearTimeout(this.wsReconnectTimer)
        this.wsReconnectTimer = null
      }
      if (this.ws) {
        this.ws.onclose = null
        this.ws.close()
        this.ws = null
      }
      this.wsConnected = false
    },

    scheduleReconnect() {
      if (this.wsReconnectTimer) return
      this.wsReconnectTimer = setTimeout(() => {
        this.wsReconnectTimer = null
        this.connectWebSocket()
      }, this.wsReconnectDelay)
      this.wsReconnectDelay = Math.min(this.wsReconnectDelay * 2, 30000)
    },

    handleWSMessage(data) {
      switch (data.type) {
        case 'new_message':
          if (data.message && this.currentConversation && data.message.conversationId === this.currentConversation.id) {
            const exists = this.messages.some(m => m.id === data.message.id)
            if (!exists) {
              this.messages.push(data.message)
            }
            this.wsSend({ type: 'mark_read', conversationId: this.currentConversation.id })
          }
          break
        case 'unread_count':
          this.unreadCount = data.count || 0
          break
        case 'conversation_updated':
          if (data.conversation) {
            const idx = this.conversations.findIndex(c => c.id === data.conversation.id)
            if (idx >= 0) {
              this.conversations[idx] = { ...this.conversations[idx], ...data.conversation }
            }
            if (this.currentConversation && data.conversation.id === this.currentConversation.id) {
              this.currentConversation = { ...this.currentConversation, ...data.conversation }
            }
          }
          break
      }
      if (this.onMessage) {
        this.onMessage(data)
      }
    },

    wsSend(data) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(data))
      }
    },

    async openPanel() {
      this.isOpen = true
      await this.fetchOrCreateConversation()
    },

    closePanel() {
      this.isOpen = false
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
          this.wsSend({ type: 'mark_read', conversationId: openConv.id })
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
      const convId = this.currentConversation.id
      const payload = { content: content.trim() }

      if (this.wsConnected) {
        this.wsSend({ type: 'send_message', conversationId: convId, content: content.trim() })
      } else {
        try {
          const res = await fetch(`${API_BASE}/chat/conversations/${convId}/messages`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(payload)
          })
          const data = await res.json()
          if (data.message) {
            this.messages.push(data.message)
          }
        } catch (e) {
          console.error('Send message failed:', e)
        }
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

    init() {
      this.connectWebSocket()
    },

    cleanup() {
      this.disconnectWebSocket()
      this.isOpen = false
      this.currentConversation = null
      this.messages = []
      this.conversations = []
      this.unreadCount = 0
      this.onMessage = null
    }
  }
})
