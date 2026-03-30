<script setup>
import { ref, nextTick, watch, computed, onMounted, onUnmounted } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { useImageCompress } from '../composables/useImageCompress'
import { useNotificationSound } from '../composables/useNotificationSound'

const chatStore = useChatStore()
const authStore = useAuthStore()
const router = useRouter()
const { compress } = useImageCompress()
const { play: playNotification } = useNotificationSound()

const messageInput = ref('')
const messagesContainer = ref(null)
const textareaRef = ref(null)
const fileInputRef = ref(null)
const previewImage = ref(null)
const previewBlob = ref(null)
const fullscreenImage = ref(null)
const draftRestored = ref(false)
const messagesLoaded = ref(false)
const prevMsgCount = ref(0)

const canSend = computed(() => messageInput.value.trim().length > 0 || previewImage.value)

const displayMessages = computed(() => {
  const result = []
  const msgs = chatStore.messages
  for (let i = 0; i < msgs.length; i++) {
    const msg = msgs[i]
    const prev = i > 0 ? msgs[i - 1] : null
    if (!prev || (new Date(msg.createdAt) - new Date(prev.createdAt)) > 5 * 60 * 1000) {
      result.push({ type: 'divider', time: msg.createdAt, id: `d-${msg.id}` })
    }
    result.push({ type: 'message', msg, id: msg.id })
  }
  return result
})

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  const isYesterday = date.toDateString() === yesterday.toDateString()
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  if (isToday) return `${hh}:${mm}`
  if (isYesterday) return `昨天 ${hh}:${mm}`
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${month}-${day} ${hh}:${mm}`
}

const formatDividerTime = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  const isYesterday = date.toDateString() === yesterday.toDateString()
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  if (isToday) return `${hh}:${mm}`
  if (isYesterday) return `昨天 ${hh}:${mm}`
  const month = String(date.getMonth() + 1)
  const day = date.getDate()
  return `${month}月${day}日 ${hh}:${mm}`
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const handleOpen = async () => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }
  if (chatStore.isOpen) {
    chatStore.closePanel()
  } else {
    const draft = await chatStore.openPanel()
    if (draft && !draftRestored.value) {
      messageInput.value = draft
      draftRestored.value = true
    }
    scrollToBottom()
  }
}

const handleFileSelect = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    alert('图片大小不能超过 5MB')
    return
  }
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
    alert('只支持 JPG、PNG、WebP 格式')
    return
  }
  try {
    const blob = await compress(file, { maxWidth: 800, maxHeight: 800, quality: 0.7 })
    previewBlob.value = blob
    previewImage.value = URL.createObjectURL(blob)
    nextTick(() => scrollToBottom())
  } catch (err) {
    console.error('Image compression failed:', err)
  }
  e.target.value = ''
}

const removePreview = () => {
  if (previewImage.value) {
    URL.revokeObjectURL(previewImage.value)
  }
  previewImage.value = null
  previewBlob.value = null
}

const handleSend = async () => {
  if (!canSend.value) return

  if (previewBlob.value) {
    await chatStore.sendImageMessage(previewBlob.value, messageInput.value.trim())
    removePreview()
    messageInput.value = ''
    scrollToBottom()
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
    }
    return
  }

  const content = messageInput.value
  messageInput.value = ''
  localStorage.removeItem('chat_draft')
  await chatStore.sendMessage(content)
  scrollToBottom()
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
  }
}

const handleKeydown = (e) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

const autoResize = () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
    textareaRef.value.style.height = Math.min(textareaRef.value.scrollHeight, 72) + 'px'
  }
  chatStore.notifyUserActivity()
}

const onInput = () => {
  autoResize()
  localStorage.setItem('chat_draft', messageInput.value)
  chatStore.sendTyping()
}

const openFullscreen = (url) => {
  fullscreenImage.value = url
}

const closeFullscreen = () => {
  fullscreenImage.value = null
}

const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const keepConversation = () => {
  chatStore.resetAutoCloseTimer()
}

const handleEndChat = async () => {
  if (!chatStore.currentConversation) return
  await chatStore.closeConversation()
}

const handleNewChat = async () => {
  await chatStore.fetchOrCreateConversation()
  scrollToBottom()
}

watch(() => chatStore.messages.length, (newLen) => {
  if (messagesLoaded.value && newLen > prevMsgCount.value) {
    const lastMsg = chatStore.messages[chatStore.messages.length - 1]
    if (lastMsg && lastMsg.senderType === 'service' && chatStore.soundEnabled && (!chatStore.isOpen || !document.hasFocus())) {
      playNotification()
    }
  }
  prevMsgCount.value = newLen
  scrollToBottom()
  nextTick(() => {
    messagesLoaded.value = true
  })
})

onMounted(() => {
  if (authStore.isLoggedIn) {
    chatStore.init()
  }
})

onUnmounted(() => {
  chatStore.cleanup()
})

watch(() => authStore.isLoggedIn, (val) => {
  if (val) {
    chatStore.init()
  } else {
    chatStore.cleanup()
  }
})
</script>

<template>
  <div class="chat-widget" v-if="authStore.isLoggedIn">
    <Transition name="chat-slide">
      <div v-if="chatStore.isOpen" class="chat-panel">
        <div class="chat-header">
          <div class="header-left">
            <span class="status-dot"></span>
            <span class="header-title">在线客服</span>
          </div>
          <div class="header-right">
            <button class="end-btn" @click="handleEndChat" v-if="chatStore.currentConversation && chatStore.currentConversation.status === 'open'" title="结束对话">结束对话</button>
            <button class="sound-btn" @click="chatStore.toggleSound()" :title="chatStore.soundEnabled ? '关闭提示音' : '开启提示音'">
              <svg v-if="chatStore.soundEnabled" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
                <path d="M19.07 4.93a10 10 0 0 1 0 14.14"></path>
                <path d="M15.54 8.46a5 5 0 0 1 0 7.07"></path>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
                <line x1="23" y1="9" x2="17" y2="15"></line>
                <line x1="17" y1="9" x2="23" y2="15"></line>
              </svg>
            </button>
            <button class="close-btn" @click="chatStore.closePanel()">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
        </div>

        <div v-if="chatStore.autoCloseWarning" class="auto-close-warning">
          <span>由于长时间未操作，对话将在 {{ chatStore.autoCloseCountdown }} 秒后暂停</span>
          <button class="keep-btn" @click="keepConversation">继续对话</button>
        </div>

        <div class="messages-area" ref="messagesContainer" :class="{ 'animate-msgs': messagesLoaded }">
          <div v-if="chatStore.isLoading" class="loading-state">
            <div class="loading-dots">
              <span></span><span></span><span></span>
            </div>
          </div>
          <template v-else>
            <TransitionGroup name="msg">
              <template v-for="item in displayMessages" :key="item.id">
                <div v-if="item.type === 'divider'" class="time-divider">
                  <span class="divider-text">{{ formatDividerTime(item.time) }}</span>
                </div>
                <div
                  v-else
                  class="message-wrapper"
                  :class="{
                    'is-customer': item.msg.senderType === 'customer',
                    'is-service': item.msg.senderType === 'service',
                    'is-system': item.msg.senderType === 'system',
                  }"
                >
                  <div class="message-bubble" v-if="item.msg.senderType === 'system'">
                    <div class="message-text system-text">{{ item.msg.content }}</div>
                  </div>
                  <div class="message-bubble" v-else>
                    <template v-if="item.msg.messageType === 'image'">
                      <img
                        :src="item.msg.thumbnailUrl || item.msg.content"
                        class="chat-image"
                        @click="openFullscreen(item.msg.content)"
                        loading="lazy"
                      />
                    </template>
                    <template v-else>
                      <div class="message-text">{{ item.msg.content }}</div>
                    </template>
                    <div class="message-time">{{ formatTime(item.msg.createdAt) }}</div>
                  </div>
                </div>
              </template>
            </TransitionGroup>
            <div v-if="chatStore.isServiceTyping" class="typing-indicator">
              <div class="typing-dots">
                <span></span><span></span><span></span>
              </div>
              <span class="typing-text">客服正在输入...</span>
            </div>
            <div v-if="chatStore.messages.length === 0 && !chatStore.isServiceTyping" class="empty-state">
              <p>您好，有什么可以帮助您的吗？</p>
            </div>
          </template>
        </div>

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

        <div class="input-area">
          <template v-if="chatStore.currentConversation && chatStore.currentConversation.status === 'open'">
            <button class="attach-btn" @click="triggerFileInput" title="发送图片">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21.44 11.05l-9.19 9.19a6 6 0 01-8.49-8.49l9.19-9.19a4 4 0 015.66 5.66l-9.2 9.19a2 2 0 01-2.83-2.83l8.49-8.48"></path>
              </svg>
            </button>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/jpeg,image/png,image/webp"
              style="display:none"
              @change="handleFileSelect"
            />
            <textarea
              ref="textareaRef"
              v-model="messageInput"
              placeholder="输入消息..."
              rows="1"
              maxlength="500"
              @keydown="handleKeydown"
              @input="onInput"
            ></textarea>
            <button
              class="send-btn"
              :class="{ active: canSend }"
              :disabled="!canSend"
              @click="handleSend"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="22" y1="2" x2="11" y2="13"></line>
                <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
              </svg>
            </button>
          </template>
          <div v-else class="closed-hint">
            <span>对话已结束</span>
            <button class="new-chat-btn" @click="handleNewChat">开启新对话</button>
          </div>
        </div>
      </div>
    </Transition>

    <Transition name="fullscreen-fade">
      <div v-if="fullscreenImage" class="fullscreen-overlay" @click="closeFullscreen">
        <img :src="fullscreenImage" class="fullscreen-img" @click.stop />
        <button class="fullscreen-close" @click="closeFullscreen">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </Transition>

    <button
      class="chat-fab"
      :class="{ 'has-unread': chatStore.unreadCount > 0 && !chatStore.isOpen }"
      @click="handleOpen"
    >
      <svg v-if="!chatStore.isOpen" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
      <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
      <span v-if="chatStore.unreadCount > 0 && !chatStore.isOpen" class="unread-dot"></span>
    </button>
  </div>
</template>

<style scoped>
.chat-widget {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1200;
}

.chat-widget > * {
  pointer-events: auto;
}

.chat-fab {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: none;
  background: #d4a574;
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(212, 165, 116, 0.4);
  transition: transform 0.2s, box-shadow 0.2s;
  z-index: 1201;
  user-select: none;
}

.chat-fab:hover {
  transform: scale(1.08);
  box-shadow: 0 6px 20px rgba(212, 165, 116, 0.5);
}

.chat-fab.has-unread {
  animation: pulse-badge 2s ease-in-out 3;
}

@keyframes pulse-badge {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.08); }
}

.unread-dot {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #e74c3c;
}

.chat-panel {
  position: fixed;
  right: 24px;
  bottom: 84px;
  width: 400px;
  height: 560px;
  background: #1a1a1a;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  z-index: 1200;
}

.chat-slide-enter-active {
  transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.chat-slide-leave-active {
  transition: all 0.2s ease-in;
}

.chat-slide-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.chat-slide-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.98);
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: #222;
  border-bottom: 1px solid #333;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #4caf50;
}

.header-title {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
}

.close-btn {
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  padding: 4px;
  display: flex;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #fff;
}

.sound-btn {
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  padding: 4px;
  display: flex;
  transition: color 0.2s;
}

.sound-btn:hover {
  color: #d4a574;
}

.end-btn {
  background: none;
  border: 1px solid #555;
  color: #999;
  cursor: pointer;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  display: flex;
  align-items: center;
  transition: all 0.2s;
}

.end-btn:hover {
  border-color: #e74c3c;
  color: #e74c3c;
}

.auto-close-warning {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 14px;
  background: rgba(255, 152, 0, 0.15);
  border-bottom: 1px solid rgba(255, 152, 0, 0.3);
  flex-shrink: 0;
}

.auto-close-warning span {
  color: #ffb74d;
  font-size: 12px;
}

.keep-btn {
  background: #d4a574;
  color: #0f0f0f;
  border: none;
  border-radius: 4px;
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.2s;
}

.keep-btn:hover {
  background: #e0b88a;
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.messages-area::-webkit-scrollbar {
  width: 4px;
}

.messages-area::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 2px;
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
  background: #888;
  animation: typing-bounce 1.2s infinite ease-in-out;
}

.typing-dots span:nth-child(2) { animation-delay: 0.2s; }
.typing-dots span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typing-bounce {
  0%, 80%, 100% { transform: translateY(0); opacity: 0.4; }
  40% { transform: translateY(-4px); opacity: 1; }
}

.typing-text {
  color: #888;
  font-size: 12px;
}

.messages-area.animate-msgs .msg-enter-active {
  transition: all 0.3s ease-out;
}

.messages-area.animate-msgs .msg-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.time-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
}

.divider-text {
  color: #666;
  font-size: 11px;
  background: #1a1a1a;
  padding: 0 12px;
  position: relative;
}

.time-divider::before,
.time-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #333;
}

.message-wrapper {
  display: flex;
}

.message-wrapper.is-customer {
  justify-content: flex-end;
}

.message-wrapper.is-service {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 75%;
  padding: 10px 14px;
  border-radius: 12px;
  position: relative;
}

.is-customer .message-bubble {
  background: #d4a574;
  color: #0f0f0f;
  border-bottom-right-radius: 4px;
}

.is-service .message-bubble {
  background: #2a2a2a;
  color: #fff;
  border-bottom-left-radius: 4px;
}

.message-wrapper.is-system {
  justify-content: center;
}

.message-wrapper.is-system .message-bubble {
  background: transparent;
  color: #888;
  padding: 4px 12px;
  max-width: 90%;
  text-align: center;
}

.system-text {
  font-size: 12px !important;
  color: #888 !important;
}

.message-text {
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
  white-space: pre-wrap;
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

.message-time {
  font-size: 11px;
  opacity: 0.6;
  margin-top: 4px;
  text-align: right;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  font-size: 14px;
}

.loading-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-dots {
  display: flex;
  gap: 6px;
}

.loading-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #666;
  animation: bounce 1.2s infinite ease-in-out;
}

.loading-dots span:nth-child(2) { animation-delay: 0.2s; }
.loading-dots span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}

.preview-area {
  padding: 8px 16px;
  background: #1a1a1a;
  border-top: 1px solid #333;
  flex-shrink: 0;
}

.preview-thumb {
  position: relative;
  display: inline-block;
}

.preview-thumb img {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid #444;
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

.input-area {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  padding: 12px 16px;
  background: #1a1a1a;
  border-top: 1px solid #333;
  flex-shrink: 0;
}

.attach-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: #2a2a2a;
  color: #999;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: color 0.2s, background 0.2s;
}

.attach-btn:hover {
  color: #d4a574;
  background: #333;
}

.input-area textarea {
  flex: 1;
  background: #2a2a2a;
  border: 1px solid #444;
  border-radius: 8px;
  color: #fff;
  padding: 10px 12px;
  font-size: 14px;
  font-family: inherit;
  resize: none;
  outline: none;
  max-height: 72px;
  line-height: 1.4;
  transition: border-color 0.2s;
}

.input-area textarea:focus {
  border-color: #d4a574;
}

.input-area textarea::placeholder {
  color: #666;
}

.send-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: #333;
  color: #666;
  cursor: not-allowed;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
}

.send-btn.active {
  background: #d4a574;
  color: #fff;
  cursor: pointer;
}

.send-btn.active:hover {
  background: #e0b88a;
}

.closed-hint {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #666;
  font-size: 13px;
}

.new-chat-btn {
  background: #d4a574;
  color: #0f0f0f;
  border: none;
  border-radius: 6px;
  padding: 6px 16px;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.2s;
}

.new-chat-btn:hover {
  background: #e0b88a;
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

.fullscreen-fade-enter-active,
.fullscreen-fade-leave-active {
  transition: opacity 0.2s ease;
}

.fullscreen-fade-enter-from,
.fullscreen-fade-leave-to {
  opacity: 0;
}

@media (max-width: 480px) {
  .chat-panel {
    right: 0;
    bottom: 0;
    width: 100vw;
    height: 100vh;
    border-radius: 0;
  }

  .chat-fab {
    right: 16px;
    bottom: 16px;
  }
}
</style>
