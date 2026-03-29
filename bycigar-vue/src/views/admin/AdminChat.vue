<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useToastStore } from '../../stores/toast'

const API_BASE = '/api'
const toast = useToastStore()

const conversations = ref([])
const selectedConversation = ref(null)
const messages = ref([])
const filterStatus = ref('')
const loading = ref(false)
const messagesLoading = ref(false)
const replyContent = ref('')
const messagesContainer = ref(null)
const textareaRef = ref(null)
const pollTimer = ref(null)

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

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
  if (!replyContent.value.trim() || !selectedConversation.value) return
  const content = replyContent.value
  replyContent.value = ''
  try {
    const convId = selectedConversation.value.id
    const res = await fetch(`${API_BASE}/admin/chat/conversations/${convId}/messages`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ content: content.trim() })
    })
    const data = await res.json()
    if (data.message) {
      messages.value.push(data.message)
      scrollToBottom()
    }
  } catch (e) {
    toast.error('发送失败')
  }
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
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

const startPolling = () => {
  stopPolling()
  const poll = () => {
    fetchConversations()
    if (selectedConversation.value) {
      const lastId = messages.value.length > 0 ? messages.value[messages.value.length - 1].id : 0
      fetchMessages(lastId)
    }
    pollTimer.value = setTimeout(poll, 5000)
  }
  pollTimer.value = setTimeout(poll, 5000)
}

const stopPolling = () => {
  if (pollTimer.value) {
    clearTimeout(pollTimer.value)
    pollTimer.value = null
  }
}

watch(filterStatus, () => {
  fetchConversations()
})

watch(() => messages.value.length, () => {
  scrollToBottom()
})

onMounted(() => {
  fetchConversations()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<template>
  <div class="admin-chat">
    <div class="chat-sidebar">
      <div class="sidebar-header">
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
            unread: conv.unreadCount > 0
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
              <span class="conv-preview">{{ truncate(conv.lastMessage?.content, 20) }}</span>
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

        <div class="main-messages" ref="messagesContainer">
          <div
            v-for="msg in messages"
            :key="msg.id"
            class="message-wrapper"
            :class="{
              'is-customer': msg.senderType === 'customer',
              'is-service': msg.senderType === 'service'
            }"
          >
            <div class="message-bubble">
              <div class="message-text">{{ msg.content }}</div>
              <div class="message-time">{{ formatTime(msg.createdAt) }}</div>
            </div>
          </div>
          <div v-if="messages.length === 0" class="empty-messages">暂无消息</div>
        </div>

        <div v-if="selectedConversation.status === 'open'" class="main-input">
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
            :disabled="!replyContent.trim()"
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
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid #eee;
  background: #fff;
  flex-shrink: 0;
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
</style>
