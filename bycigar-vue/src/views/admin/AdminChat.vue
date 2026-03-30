<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useToastStore } from '../../stores/toast'
import { useImageCompress } from '../../composables/useImageCompress'
import { useNotificationSound } from '../../composables/useNotificationSound'

const API_BASE = '/api'
const toast = useToastStore()
const { compress } = useImageCompress()
const { play: playNotification } = useNotificationSound()

const conversations = ref([])
const selectedConversation = ref(null)
const messages = ref([])
const filterStatus = ref('')
const loading = ref(false)
const messagesLoading = ref(false)
const replyContent = ref('')
const messagesContainer = ref(null)
const textareaRef = ref(null)
const fileInputRef = ref(null)
const previewImage = ref(null)
const previewBlob = ref(null)
const fullscreenImage = ref(null)
const ws = ref(null)
const wsReconnectTimer = ref(null)
const wsReconnectDelay = ref(3000)
const wsConnected = ref(false)
const isCustomerTyping = ref(false)
const typingTimeout = ref(null)
const lastTypingSent = ref(0)
const messagesLoaded = ref(false)
const prevMsgCount = ref(0)
const serviceOnline = ref(false)
const quickReplies = ref([])
const showQuickReplies = ref(false)
const quickReplySearch = ref('')

const filteredQuickReplies = computed(() => {
  if (!quickReplySearch.value) return quickReplies.value
  const q = quickReplySearch.value.toLowerCase()
  return quickReplies.value.filter(r => r.title.toLowerCase().includes(q) || r.content.toLowerCase().includes(q))
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const connectWebSocket = () => {
  disconnectWebSocket()
  const token = localStorage.getItem('token')
  if (!token) return
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const url = `${protocol}//${window.location.host}/api/admin/chat/ws?token=${encodeURIComponent(token)}`

  try {
    ws.value = new WebSocket(url)

    ws.value.onopen = () => {
      wsConnected.value = true
      wsReconnectDelay.value = 3000
    }

    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        handleWSMessage(data)
      } catch (e) {
        console.error('Parse WS message failed:', e)
      }
    }

    ws.value.onclose = () => {
      wsConnected.value = false
      scheduleReconnect()
    }

    ws.value.onerror = () => {
      wsConnected.value = false
    }
  } catch (e) {
    console.error('WebSocket connect failed:', e)
    scheduleReconnect()
  }
}

const disconnectWebSocket = () => {
  if (wsReconnectTimer.value) {
    clearTimeout(wsReconnectTimer.value)
    wsReconnectTimer.value = null
  }
  if (ws.value) {
    ws.value.onclose = null
    ws.value.close()
    ws.value = null
  }
  wsConnected.value = false
}

const scheduleReconnect = () => {
  if (wsReconnectTimer.value) return
  wsReconnectTimer.value = setTimeout(() => {
    wsReconnectTimer.value = null
    connectWebSocket()
  }, wsReconnectDelay.value)
  wsReconnectDelay.value = Math.min(wsReconnectDelay.value * 2, 30000)
}

const wsSend = (data) => {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify(data))
  }
}

const handleWSMessage = (data) => {
  switch (data.type) {
    case 'new_message':
      if (data.message && selectedConversation.value && data.conversationId === selectedConversation.value.id) {
        const exists = messages.value.some(m => m.id === data.message.id)
        if (!exists) {
          messages.value.push(data.message)
          scrollToBottom()
          if (data.message.senderType === 'customer') {
            playNotification()
          }
        }
        wsSend({ type: 'mark_read', conversationId: selectedConversation.value.id })
      }
      fetchConversations()
      break
    case 'conversation_updated':
      if (data.conversation) {
        const idx = conversations.value.findIndex(c => c.id === data.conversation.id)
        if (idx >= 0) {
          conversations.value[idx] = { ...conversations.value[idx], ...data.conversation }
        }
        if (selectedConversation.value && data.conversation.id === selectedConversation.value.id) {
          selectedConversation.value = { ...selectedConversation.value, ...data.conversation }
        }
      }
      break
    case 'unread_stats':
      break
    case 'typing':
      if (selectedConversation.value && data.conversationId === selectedConversation.value.id) {
        isCustomerTyping.value = true
        if (typingTimeout.value) clearTimeout(typingTimeout.value)
        typingTimeout.value = setTimeout(() => {
          isCustomerTyping.value = false
        }, 3000)
      }
      break
    case 'ack':
      if (data.messageId) {
        const m = messages.value.find(m => m.id === data.messageId)
        if (m) m.status = data.status || 'sent'
      }
      break
    case 'message_status':
      if (data.messageId && data.status) {
        const m = messages.value.find(m => m.id === data.messageId)
        if (m) m.status = data.status
      }
      break
    case 'message_recalled':
      if (data.message) {
        const idx = messages.value.findIndex(m => m.id === data.message.id)
        if (idx >= 0) {
          messages.value[idx] = { ...messages.value[idx], ...data.message }
        }
      }
      break
  }
}

const fetchConversations = async () => {
  try {
    const params = new URLSearchParams()
    if (filterStatus.value) params.append('status', filterStatus.value)
    const res = await fetch(`${API_BASE}/admin/chat/conversations?${params}`, { headers: authHeaders() })
    const data = await res.json()
    conversations.value = data.conversations || []
  } catch (e) {
    console.error('Fetch conversations failed:', e)
  }
}

const selectConversation = async (conv) => {
  selectedConversation.value = conv
  messages.value = []
  await fetchMessages()
  scrollToBottom()
  wsSend({ type: 'mark_read', conversationId: conv.id })
}

const fetchMessages = async (afterId) => {
  if (!selectedConversation.value) return
  try {
    const convId = selectedConversation.value.id
    let url = `${API_BASE}/admin/chat/conversations/${convId}/messages`
    if (afterId) url += `?after=${afterId}`
    const res = await fetch(url, { headers: authHeaders() })
    const data = await res.json()
    if (afterId) {
      if (data.messages && data.messages.length > 0) {
        messages.value = [...messages.value, ...data.messages]
      }
    } else {
      messages.value = data.messages || []
    }
    if (selectedConversation.value) {
      const updated = conversations.value.find(c => c.id === selectedConversation.value.id)
      if (updated) {
        updated.unreadCount = 0
      }
    }
  } catch (e) {
    console.error('Fetch messages failed:', e)
  }
}

const sendReply = async () => {
  if (!selectedConversation.value) return
  if (!replyContent.value.trim() && !previewBlob.value) return

  const convId = selectedConversation.value.id

  if (previewBlob.value) {
    try {
      const formData = new FormData()
      formData.append('file', previewBlob.value, 'chat_image.jpg')
      const token = localStorage.getItem('token')
      const uploadRes = await fetch(`${API_BASE}/admin/upload`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: formData
      })
      const uploadData = await uploadRes.json()
      if (!uploadData.success) { toast.error('图片上传失败'); return }

      if (wsConnected.value) {
        wsSend({
          type: 'send_message',
          conversationId: convId,
          content: uploadData.url,
          messageType: 'image',
          thumbnailUrl: uploadData.thumbnailUrl || ''
        })
      } else {
        const res = await fetch(`${API_BASE}/admin/chat/conversations/${convId}/messages`, {
          method: 'POST',
          headers: authHeaders(),
          body: JSON.stringify({
            content: uploadData.url,
            messageType: 'image',
            thumbnailUrl: uploadData.thumbnailUrl || ''
          })
        })
        const data = await res.json()
        if (data.message) { messages.value.push(data.message); scrollToBottom() }
      }
      removePreview()
    } catch (e) {
      toast.error('图片发送失败')
    }
  }

  if (replyContent.value.trim()) {
    const content = replyContent.value.trim()
    replyContent.value = ''
    if (textareaRef.value) textareaRef.value.style.height = 'auto'

    if (wsConnected.value) {
      wsSend({ type: 'send_message', conversationId: convId, content })
    } else {
      try {
        const res = await fetch(`${API_BASE}/admin/chat/conversations/${convId}/messages`, {
          method: 'POST',
          headers: authHeaders(),
          body: JSON.stringify({ content })
        })
        const data = await res.json()
        if (data.message) { messages.value.push(data.message); scrollToBottom() }
      } catch (e) {
        toast.error('发送失败')
      }
    }
  }
}

const closeConversation = async () => {
  if (!selectedConversation.value) return
  try {
    const convId = selectedConversation.value.id
    const res = await fetch(`${API_BASE}/admin/chat/conversations/${convId}/close`, {
      method: 'PUT',
      headers: authHeaders()
    })
    const data = await res.json()
    if (data.success) {
      selectedConversation.value.status = 'closed'
      const conv = conversations.value.find(c => c.id === convId)
      if (conv) conv.status = 'closed'
      toast.success('对话已关闭')
    }
  } catch (e) {
    toast.error('操作失败')
  }
}

const handleKeydown = (e) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendReply()
  }
}

const autoResize = () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
    textareaRef.value.style.height = Math.min(textareaRef.value.scrollHeight, 72) + 'px'
  }
  if (selectedConversation.value && replyContent.value.trim()) {
    const now = Date.now()
    if (now - lastTypingSent.value >= 2000) {
      lastTypingSent.value = now
      wsSend({ type: 'typing', conversationId: selectedConversation.value.id })
    }
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  if (isToday) return `${hh}:${mm}`
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${month}-${day} ${hh}:${mm}`
}

const truncate = (str, len) => {
  if (!str) return ''
  return str.length > len ? str.substring(0, len) + '...' : str
}

const getPreviewText = (msg) => {
  if (!msg) return ''
  if (msg.messageType === 'image') return '[图片]'
  return truncate(msg.content, 20)
}

const handleFileSelect = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 5 * 1024 * 1024) { toast.error('图片大小不能超过 5MB'); return }
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) { toast.error('只支持 JPG、PNG、WebP 格式'); return }
  try {
    const blob = await compress(file, { maxWidth: 800, maxHeight: 800, quality: 0.7 })
    previewBlob.value = blob
    previewImage.value = URL.createObjectURL(blob)
  } catch (err) {
    console.error('Image compression failed:', err)
  }
  e.target.value = ''
}

const removePreview = () => {
  if (previewImage.value) URL.revokeObjectURL(previewImage.value)
  previewImage.value = null
  previewBlob.value = null
}

const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const openFullscreen = (url) => {
  fullscreenImage.value = url
}

const closeFullscreen = () => {
  fullscreenImage.value = null
}

const canSend = computed(() => replyContent.value.trim().length > 0 || previewBlob.value)

watch(filterStatus, () => {
  fetchConversations()
})

watch(() => messages.value.length, (newLen) => {
  scrollToBottom()
  if (messagesLoaded.value && newLen > prevMsgCount.value) {
    nextTick(() => scrollToBottom())
  }
  prevMsgCount.value = newLen
  nextTick(() => {
    messagesLoaded.value = true
  })
})

const fetchServiceStatus = async () => {
  try {
    const res = await fetch(`${API_BASE}/chat/service-status`)
    const data = await res.json()
    serviceOnline.value = data.online || false
  } catch (e) {
    console.error('Fetch service status failed:', e)
  }
}

const toggleServiceStatus = () => {
  const newStatus = !serviceOnline.value
  wsSend({ type: newStatus ? 'service_online' : 'service_offline' })
  serviceOnline.value = newStatus
}

const fetchQuickReplies = async () => {
  try {
    const res = await fetch(`${API_BASE}/admin/quick-replies`, { headers: authHeaders() })
    const data = await res.json()
    quickReplies.value = data.quickReplies || []
  } catch (e) {
    console.error('Fetch quick replies failed:', e)
  }
}

const insertQuickReply = (reply) => {
  replyContent.value = reply.content
  showQuickReplies.value = false
  quickReplySearch.value = ''
  if (textareaRef.value) {
    textareaRef.value.focus()
    autoResize()
  }
}

const recallMessage = (msg) => {
  wsSend({
    type: 'recall_message',
    conversationId: selectedConversation.value.id,
    messageId: msg.id
  })
}

const assignConversation = async (conv, userId) => {
  try {
    await fetch(`${API_BASE}/admin/chat/conversations/${conv.id}/assign`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ assignedTo: userId })
    })
    toast.success('已分配客服')
    fetchConversations()
  } catch (e) {
    toast.error('分配失败')
  }
}

const getAssignedName = (conv) => {
  if (!conv.assignedUser) return '未分配'
  return conv.assignedUser.name || conv.assignedUser.email
}

onMounted(() => {
  fetchConversations()
  connectWebSocket()
  fetchServiceStatus()
  fetchQuickReplies()
})

onUnmounted(() => {
  disconnectWebSocket()
  if (typingTimeout.value) clearTimeout(typingTimeout.value)
})
</script>

<template>
  <div class="admin-chat">
    <div class="chat-sidebar">
      <div class="sidebar-header">
        <div class="service-status-bar">
          <div class="status-indicator" :class="{ online: serviceOnline }">
            <span class="status-dot-admin"></span>
            <span>{{ serviceOnline ? '在线接客中' : '当前离线' }}</span>
          </div>
          <button class="toggle-status-btn" :class="{ online: serviceOnline }" @click="toggleServiceStatus">
            {{ serviceOnline ? '下线' : '上线' }}
          </button>
        </div>
        <div class="filter-tabs">
          <button
            :class="{ active: filterStatus === '' }"
            @click="filterStatus = ''"
          >全部</button>
          <button
            :class="{ active: filterStatus === 'open' }"
            @click="filterStatus = 'open'"
          >进行中</button>
          <button
            :class="{ active: filterStatus === 'closed' }"
            @click="filterStatus = 'closed'"
          >已关闭</button>
        </div>
      </div>
      <div class="conversation-list">
        <div
          v-for="conv in conversations"
          :key="conv.id"
          class="conversation-item"
          :class="{
            active: selectedConversation?.id === conv.id,
            unread: conv.unreadCount > 0,
            closed: conv.status === 'closed'
          }"
          @click="selectConversation(conv)"
        >
          <div class="conv-avatar">{{ conv.user?.name?.charAt(0) || '?' }}</div>
          <div class="conv-info">
            <div class="conv-top">
              <span class="conv-name">{{ conv.user?.name || conv.user?.email }}</span>
              <span class="conv-time">{{ formatTime(conv.lastMessageAt) }}</span>
            </div>
            <div class="conv-bottom">
              <span class="conv-preview">{{ getPreviewText(conv.lastMessage) }}</span>
              <span v-if="conv.unreadCount > 0" class="conv-badge">{{ conv.unreadCount }}</span>
            </div>
          </div>
        </div>
        <div v-if="conversations.length === 0" class="empty-list">暂无对话</div>
      </div>
    </div>

    <div class="chat-main">
      <template v-if="selectedConversation">
        <div class="main-header">
          <div class="header-user">
            <span class="user-name">{{ selectedConversation.user?.name || selectedConversation.user?.email }}</span>
            <span class="user-status" :class="selectedConversation.status">
              {{ selectedConversation.status === 'open' ? '进行中' : '已关闭' }}
            </span>
          </div>
          <button
            v-if="selectedConversation.status === 'open'"
            class="close-conv-btn"
            @click="closeConversation"
          >关闭对话</button>
        </div>

        <div class="main-messages" ref="messagesContainer" :class="{ 'animate-msgs': messagesLoaded }">
          <TransitionGroup name="msg">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="message-wrapper"
              :class="{
                'is-customer': msg.senderType === 'customer',
                'is-service': msg.senderType === 'service',
                'is-system': msg.senderType === 'system',
              }"
            >
              <div class="message-bubble" v-if="msg.senderType === 'system'">
                <div class="message-text system-text">{{ msg.content }}</div>
              </div>
              <div class="message-bubble recalled" v-else-if="msg.recalledAt">
                <div class="message-text recalled-text">该消息已撤回</div>
              </div>
              <div class="message-bubble" v-else>
                <template v-if="msg.messageType === 'image'">
                  <img
                    :src="msg.thumbnailUrl || msg.content"
                    class="chat-image"
                    @click="openFullscreen(msg.content)"
                    loading="lazy"
                  />
                </template>
                <template v-else>
                  <div class="message-text">{{ msg.content }}</div>
                </template>
                <div class="message-meta">
                  <span class="message-time">{{ formatTime(msg.createdAt) }}</span>
                  <template v-if="msg.senderType === 'service'">
                    <button v-if="!msg.recalledAt" class="recall-btn" @click="recallMessage(msg)">撤回</button>
                    <span class="msg-status">
                      <svg v-if="msg.status === 'read'" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#4caf50" stroke-width="2"><polyline points="18 6 7 17 2 12"></polyline><polyline points="22 6 11 17"></polyline></svg>
                      <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#999" stroke-width="2"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    </span>
                  </template>
                </div>
              </div>
            </div>
          </TransitionGroup>
          <div v-if="isCustomerTyping" class="typing-indicator">
            <div class="typing-dots">
              <span></span><span></span><span></span>
            </div>
            <span class="typing-text">客户正在输入...</span>
          </div>
          <div v-if="messages.length === 0 && !isCustomerTyping" class="empty-messages">暂无消息</div>
        </div>

        <div v-if="selectedConversation.status === 'open'" class="main-input">
          <div v-if="previewImage" class="preview-area">
            <div class="preview-thumb">
              <img :src="previewImage" alt="preview" />
              <button class="preview-remove" @click="removePreview">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>
          <button class="attach-btn" @click="triggerFileInput" title="发送图片">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21.44 11.05l-9.19 9.19a6 6 0 01-8.49-8.49l9.19-9.19a4 4 0 015.66 5.66l-9.2 9.19a2 2 0 01-2.83-2.83l8.49-8.48"></path>
            </svg>
          </button>
          <div class="quick-reply-wrapper">
            <button class="attach-btn" @click="showQuickReplies = !showQuickReplies" title="快捷回复">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-5.78 2.32l-1.52 1.5a.5.5 0 0 1-.86-.28L9.3 16.5a.5.5 0 0 1 0-1l1.52-1.5A8.5 8.5 0 0 1 11.5 3z"></path>
              </svg>
            </button>
            <div v-if="showQuickReplies" class="quick-reply-panel">
              <input v-model="quickReplySearch" class="quick-reply-search" placeholder="搜索快捷回复..." />
              <div class="quick-reply-list">
                <div
                  v-for="r in filteredQuickReplies"
                  :key="r.id"
                  class="quick-reply-item"
                  @click="insertQuickReply(r)"
                >
                  <div class="qr-title">{{ r.title }}</div>
                  <div class="qr-content">{{ r.content }}</div>
                </div>
                <div v-if="filteredQuickReplies.length === 0" class="qr-empty">暂无快捷回复</div>
              </div>
            </div>
          </div>
          <input
            ref="fileInputRef"
            type="file"
            accept="image/jpeg,image/png,image/webp"
            style="display:none"
            @change="handleFileSelect"
          />
          <textarea
            ref="textareaRef"
            v-model="replyContent"
            placeholder="输入回复..."
            rows="1"
            maxlength="500"
            @keydown="handleKeydown"
            @input="autoResize"
          ></textarea>
          <button
            class="send-btn"
            :disabled="!canSend"
            @click="sendReply"
          >发送</button>
        </div>
      </template>
      <div v-else class="empty-main">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#ccc" stroke-width="1.5" stroke-linecap="round">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
        <p>选择一个对话开始回复</p>
      </div>
    </div>

    <div v-if="fullscreenImage" class="fullscreen-overlay" @click="closeFullscreen">
      <img :src="fullscreenImage" class="fullscreen-img" @click.stop />
      <button class="fullscreen-close" @click="closeFullscreen">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.admin-chat {
  display: flex;
  height: calc(100vh - 120px);
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.chat-sidebar {
  width: 320px;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.service-status-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 12px;
}
.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #999;
}
.status-dot-admin {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ccc;
  transition: background 0.3s;
}
.status-indicator.online .status-dot-admin {
  background: #4caf50;
}
.status-indicator.online {
  color: #4caf50;
  font-weight: 500;
}
.toggle-status-btn {
  padding: 4px 14px;
  border-radius: 14px;
  font-size: 12px;
  cursor: pointer;
  border: 1px solid #ddd;
  background: #fff;
  color: #666;
  transition: all 0.2s;
}
.toggle-status-btn.online {
  background: #4caf50;
  color: #fff;
  border-color: #4caf50;
}
.toggle-status-btn:hover {
  opacity: 0.85;
}

.filter-tabs {
  display: flex;
  gap: 4px;
}

.filter-tabs button {
  flex: 1;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  color: #666;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-tabs button.active {
  background: #d4a574;
  color: #fff;
  border-color: #d4a574;
}

.filter-tabs button:hover:not(.active) {
  background: #f5f5f5;
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid #f5f5f5;
}

.conversation-item:hover {
  background: #f9f9f9;
}

.conversation-item.active {
  background: #f0ebe4;
}

.conversation-item.unread {
  background: #fef9f3;
}

.conversation-item.closed {
  opacity: 0.6;
}

.conversation-item.closed .conv-name,
.conversation-item.closed .conv-preview {
  text-decoration: line-through;
}

.conv-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #d4a574;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.conv-info {
  flex: 1;
  min-width: 0;
}

.conv-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.conv-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.conv-time {
  font-size: 12px;
  color: #999;
}

.conv-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conv-preview {
  font-size: 13px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.conv-badge {
  background: #e74c3c;
  color: #fff;
  font-size: 11px;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 5px;
  margin-left: 8px;
  flex-shrink: 0;
}

.empty-list {
  padding: 40px;
  text-align: center;
  color: #999;
  font-size: 14px;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.main-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
}

.header-user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.user-status {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

.user-status.open {
  background: #e8f5e9;
  color: #4caf50;
}

.user-status.closed {
  background: #f5f5f5;
  color: #999;
}

.close-conv-btn {
  padding: 6px 14px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  color: #666;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.close-conv-btn:hover {
  background: #f5f5f5;
  border-color: #ccc;
}

.main-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: #fafafa;
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 0;
}

.typing-dots {
  display: flex;
  gap: 3px;
}

.typing-dots span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #999;
  animation: typing-bounce 1.2s infinite ease-in-out;
}

.typing-dots span:nth-child(2) { animation-delay: 0.2s; }
.typing-dots span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typing-bounce {
  0%, 80%, 100% { transform: translateY(0); opacity: 0.4; }
  40% { transform: translateY(-4px); opacity: 1; }
}

.typing-text {
  color: #999;
  font-size: 12px;
}

.main-messages.animate-msgs .msg-enter-active {
  transition: all 0.3s ease-out;
}

.main-messages.animate-msgs .msg-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.message-wrapper {
  display: flex;
}

.message-wrapper.is-customer {
  justify-content: flex-start;
}

.message-wrapper.is-service {
  justify-content: flex-end;
}

.message-bubble {
  max-width: 60%;
  padding: 10px 14px;
  border-radius: 12px;
}

.is-customer .message-bubble {
  background: #fff;
  color: #333;
  border-bottom-left-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.is-service .message-bubble {
  background: #d4a574;
  color: #fff;
  border-bottom-right-radius: 4px;
}

.message-wrapper.is-system {
  justify-content: center;
}

.message-wrapper.is-system .message-bubble {
  background: transparent;
  color: #888;
  padding: 4px 12px;
  max-width: 90%;
}

.system-text {
  font-size: 12px !important;
  color: #888 !important;
  text-align: center;
}

.message-text {
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
  white-space: pre-wrap;
}

.message-time {
  font-size: 11px;
  opacity: 0.6;
  margin-top: 4px;
  text-align: right;
}

.empty-messages {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 14px;
}

.main-input {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 12px 20px;
  border-top: 1px solid #eee;
  background: #fff;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.preview-area {
  width: 100%;
  margin-bottom: 4px;
}

.preview-thumb {
  position: relative;
  display: inline-block;
}

.preview-thumb img {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid #ddd;
}

.preview-remove {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #e74c3c;
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.attach-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 1px solid #ddd;
  background: #fff;
  color: #999;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: color 0.2s;
}

.attach-btn:hover {
  color: #d4a574;
  border-color: #d4a574;
}

.main-input textarea {
  flex: 1;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 14px;
  font-family: inherit;
  resize: none;
  outline: none;
  max-height: 72px;
  line-height: 1.4;
  transition: border-color 0.2s;
}

.main-input textarea:focus {
  border-color: #d4a574;
}

.main-input textarea::placeholder {
  color: #bbb;
}

.send-btn {
  padding: 10px 20px;
  background: #d4a574;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: #e0b88a;
}

.send-btn:disabled {
  background: #ddd;
  color: #999;
  cursor: not-allowed;
}

.empty-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #ccc;
}

.empty-main p {
  font-size: 14px;
}

@media (max-width: 768px) {
  .admin-chat {
    flex-direction: column;
    height: auto;
    min-height: calc(100vh - 120px);
  }

  .chat-sidebar {
    width: 100%;
    max-height: 300px;
  }
}

.chat-image {
  max-width: 200px;
  max-height: 200px;
  border-radius: 8px;
  cursor: pointer;
  display: block;
  object-fit: contain;
  transition: opacity 0.2s;
}

.chat-image:hover {
  opacity: 0.85;
}

.fullscreen-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1300;
  cursor: pointer;
}

.fullscreen-img {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
  border-radius: 4px;
}

.fullscreen-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: #fff;
  cursor: pointer;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.fullscreen-close:hover {
  background: rgba(255, 255, 255, 0.2);
}

.message-bubble.recalled {
  background: transparent !important;
}

.recalled-text {
  color: #999 !important;
  font-style: italic;
  font-size: 12px !important;
}

.message-meta {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 4px;
}

.recall-btn {
  background: none;
  border: none;
  color: #999;
  font-size: 11px;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
}

.recall-btn:hover {
  color: #e74c3c;
}

.msg-status {
  display: flex;
  align-items: center;
}

.quick-reply-wrapper {
  position: relative;
}

.quick-reply-panel {
  position: absolute;
  bottom: 100%;
  left: 0;
  width: 300px;
  max-height: 320px;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin-bottom: 8px;
  z-index: 10;
}

.quick-reply-search {
  padding: 10px 12px;
  border: none;
  border-bottom: 1px solid #eee;
  font-size: 13px;
  outline: none;
}

.quick-reply-list {
  overflow-y: auto;
  max-height: 260px;
}

.quick-reply-item {
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid #f5f5f5;
}

.quick-reply-item:hover {
  background: #f9f5f0;
}

.quick-reply-item .qr-title {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.quick-reply-item .qr-content {
  font-size: 12px;
  color: #888;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.qr-empty {
  padding: 20px;
  text-align: center;
  color: #999;
  font-size: 13px;
}
</style>
