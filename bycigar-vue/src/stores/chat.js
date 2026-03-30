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
    onMessage: null,
    autoCloseTimer: null,
    autoCloseWarning: false,
    autoCloseCountdown: 0,
    isServiceTyping: false,
    typingTimeout: null,
    lastTypingSent: 0,
    soundEnabled: localStorage.getItem('chat_sound_enabled') !== 'false'
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
          this.resetAutoCloseTimer()
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
        case 'typing':
          if (this.currentConversation && data.conversationId === this.currentConversation.id) {
            this.isServiceTyping = true
            if (this.typingTimeout) clearTimeout(this.typingTimeout)
            this.typingTimeout = setTimeout(() => {
              this.isServiceTyping = false
            }, 3000)
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

    sendTyping() {
      if (!this.currentConversation) return
      const now = Date.now()
      if (now - this.lastTypingSent < 2000) return
      this.lastTypingSent = now
      this.wsSend({ type: 'typing', conversationId: this.currentConversation.id })
    },

    toggleSound() {
      this.soundEnabled = !this.soundEnabled
      localStorage.setItem('chat_sound_enabled', String(this.soundEnabled))
    },

    async openPanel() {
      this.isOpen = true
      const draft = localStorage.getItem('chat_draft')
      this.resetAutoCloseTimer()
      await this.fetchOrCreateConversation()
      return draft
    },

    closePanel() {
      this.clearAutoCloseTimer()
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
      this.resetAutoCloseTimer()
    },

    async sendImageMessage(imageBlob, caption) {
      if (!this.currentConversation) return null
      try {
        const formData = new FormData()
        formData.append('file', imageBlob, 'chat_image.jpg')

        const token = localStorage.getItem('token')
        const uploadRes = await fetch(`${API_BASE}/chat/upload-image`, {
          method: 'POST',
          headers: { 'Authorization': token ? `Bearer ${token}` : '' },
          body: formData
        })
        const uploadData = await uploadRes.json()
        if (!uploadData.success) {
          console.error('Upload failed:', uploadData.error)
          return null
        }

        const convId = this.currentConversation.id
        const payload = {
          content: uploadData.url,
          messageType: 'image',
          thumbnailUrl: uploadData.thumbnailUrl || ''
        }

        if (this.wsConnected) {
          this.wsSend({
            type: 'send_message',
            conversationId: convId,
            content: uploadData.url,
            messageType: 'image',
            thumbnailUrl: uploadData.thumbnailUrl || ''
          })
        } else {
          const res = await fetch(`${API_BASE}/chat/conversations/${convId}/messages`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(payload)
          })
          const data = await res.json()
          if (data.message) {
            this.messages.push(data.message)
          }
        }
        this.resetAutoCloseTimer()
        return true
      } catch (e) {
        console.error('Send image failed:', e)
        return null
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

    resetAutoCloseTimer() {
      if (this.autoCloseTimer) {
        clearTimeout(this.autoCloseTimer)
        this.autoCloseTimer = null
      }
      this.autoCloseWarning = false
      this.autoCloseCountdown = 0
      if (!this.isOpen) return

      const tenMinutes = 10 * 60 * 1000
      this.autoCloseTimer = setTimeout(() => {
        this.autoCloseWarning = true
        this.autoCloseCountdown = 30
        const countdownInterval = setInterval(() => {
          this.autoCloseCountdown--
          if (this.autoCloseCountdown <= 0) {
            clearInterval(countdownInterval)
            const draft = ''
            localStorage.setItem('chat_draft', draft)
            this.autoCloseWarning = false
            this.isOpen = false
            this.autoCloseTimer = null
          }
        }, 1000)
      }, tenMinutes - 30000)
    },

    clearAutoCloseTimer() {
      if (this.autoCloseTimer) {
        clearTimeout(this.autoCloseTimer)
        this.autoCloseTimer = null
      }
      this.autoCloseWarning = false
      this.autoCloseCountdown = 0
    },

    notifyUserActivity() {
      if (this.isOpen) {
        this.resetAutoCloseTimer()
      }
    },

    init() {
      this.connectWebSocket()
    },

    cleanup() {
      this.disconnectWebSocket()
      this.clearAutoCloseTimer()
      if (this.typingTimeout) {
        clearTimeout(this.typingTimeout)
        this.typingTimeout = null
      }
      this.isOpen = false
      this.currentConversation = null
      this.messages = []
      this.conversations = []
      this.unreadCount = 0
      this.onMessage = null
      this.isServiceTyping = false
    }
  }
})
