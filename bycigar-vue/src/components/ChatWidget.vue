<script setup>
import { ref, nextTick, watch, computed, onMounted, onUnmounted } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const chatStore = useChatStore()
const authStore = useAuthStore()
const router = useRouter()

const messageInput = ref('')
const messagesContainer = ref(null)
const textareaRef = ref(null)
const fabRef = ref(null)

const fabX = ref(window.innerWidth - 86)
const fabY = ref(window.innerHeight - 86)
const isDragging = ref(false)
const dragStartX = ref(0)
const dragStartY = ref(0)
const dragOffsetX = ref(0)
const dragOffsetY = ref(0)
const hasMoved = ref(false)

const canSend = computed(() => messageInput.value.trim().length > 0)

const panelStyle = computed(() => {
  const panelW = window.innerWidth <= 480 ? window.innerWidth : 380
  const panelH = window.innerWidth <= 480 ? window.innerHeight : 520
  let left = fabX.value + 56 / 2 - panelW / 2
  left = Math.max(10, Math.min(window.innerWidth - panelW - 10, left))
  let top = fabY.value - panelH - 14
  if (top < 10) top = 10
  return { left: left + 'px', top: top + 'px' }
})

const clampPosition = (x, y) => {
  const btnSize = 56
  const margin = 10
  const maxX = window.innerWidth - btnSize - margin
  const maxY = window.innerHeight - btnSize - margin
  return {
    x: Math.max(margin, Math.min(maxX, x)),
    y: Math.max(margin, Math.min(maxY, y))
  }
}

const onDragStart = (clientX, clientY) => {
  isDragging.value = true
  hasMoved.value = false
  dragStartX.value = clientX
  dragStartY.value = clientY
  dragOffsetX.value = clientX - fabX.value
  dragOffsetY.value = clientY - fabY.value
}

const onDragMove = (clientX, clientY) => {
  if (!isDragging.value) return
  const dx = Math.abs(clientX - dragStartX.value)
  const dy = Math.abs(clientY - dragStartY.value)
  if (dx > 4 || dy > 4) {
    hasMoved.value = true
  }
  const pos = clampPosition(clientX - dragOffsetX.value, clientY - dragOffsetY.value)
  fabX.value = pos.x
  fabY.value = pos.y
}

const onDragEnd = () => {
  if (!isDragging.value) return
  isDragging.value = false
  if (!hasMoved.value) {
    handleOpen()
  }
}

const onMouseDown = (e) => {
  e.preventDefault()
  onDragStart(e.clientX, e.clientY)
}

const onTouchStart = (e) => {
  const t = e.touches[0]
  onDragStart(t.clientX, t.clientY)
}

const onMouseMove = (e) => {
  onDragMove(e.clientX, e.clientY)
}

const onTouchMove = (e) => {
  if (!isDragging.value) return
  e.preventDefault()
  const t = e.touches[0]
  onDragMove(t.clientX, t.clientY)
}

const onMouseUp = () => {
  onDragEnd()
}

const onTouchEnd = () => {
  onDragEnd()
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
    await chatStore.openPanel()
    scrollToBottom()
  }
}

const handleSend = async () => {
  if (!canSend.value) return
  const content = messageInput.value
  messageInput.value = ''
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
}

watch(() => chatStore.messages.length, () => {
  scrollToBottom()
})

onMounted(() => {
  if (authStore.isLoggedIn) {
    chatStore.initPolling()
  }
  const pos = clampPosition(window.innerWidth - 86, window.innerHeight - 86)
  fabX.value = pos.x
  fabY.value = pos.y
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.addEventListener('touchmove', onTouchMove, { passive: false })
  document.addEventListener('touchend', onTouchEnd)
})

onUnmounted(() => {
  chatStore.cleanup()
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
  document.removeEventListener('touchmove', onTouchMove)
  document.removeEventListener('touchend', onTouchEnd)
})

watch(() => authStore.isLoggedIn, (val) => {
  if (val) {
    chatStore.initPolling()
  } else {
    chatStore.cleanup()
  }
})
</script>

<template>
  <div class="chat-widget" v-if="authStore.isLoggedIn">
    <Transition name="chat-slide">
      <div v-if="chatStore.isOpen" class="chat-panel" :style="panelStyle">
        <div class="chat-header">
          <div class="header-left">
            <span class="status-dot"></span>
            <span class="header-title">在线客服</span>
          </div>
          <button class="close-btn" @click="chatStore.closePanel()">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="messages-area" ref="messagesContainer">
          <div v-if="chatStore.isLoading" class="loading-state">
            <div class="loading-dots">
              <span></span><span></span><span></span>
            </div>
          </div>
          <template v-else>
            <div
              v-for="msg in chatStore.messages"
              :key="msg.id"
              class="message-wrapper"
              :class="{
                'is-customer': msg.senderType === 'customer',
                'is-service': msg.senderType === 'service',
                'is-system': msg.senderType === 'system'
              }"
            >
              <div class="message-bubble">
                <div class="message-text">{{ msg.content }}</div>
                <div class="message-time">{{ formatTime(msg.createdAt) }}</div>
              </div>
            </div>
            <div v-if="chatStore.messages.length === 0" class="empty-state">
              <p>您好，有什么可以帮助您的吗？</p>
            </div>
          </template>
        </div>

        <div class="input-area">
          <textarea
            ref="textareaRef"
            v-model="messageInput"
            placeholder="输入消息..."
            rows="1"
            maxlength="500"
            @keydown="handleKeydown"
            @input="autoResize"
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
        </div>
      </div>
    </Transition>

    <button
      ref="fabRef"
      class="chat-fab"
      :class="{ 'has-unread': chatStore.unreadCount > 0, 'is-dragging': isDragging }"
      :style="{ left: fabX + 'px', top: fabY + 'px' }"
      @mousedown="onMouseDown"
      @touchstart="onTouchStart"
    >
      <svg v-if="!chatStore.isOpen" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
      <svg v-else width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
      <span v-if="chatStore.unreadCount > 0 && !chatStore.isOpen" class="unread-badge">
        {{ chatStore.unreadCount > 99 ? '99+' : chatStore.unreadCount }}
      </span>
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
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: none;
  background: #d4a574;
  color: #fff;
  cursor: grab;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(212, 165, 116, 0.4);
  transition: box-shadow 0.2s;
  z-index: 1201;
  user-select: none;
  touch-action: none;
}

.chat-fab:hover {
  box-shadow: 0 6px 20px rgba(212, 165, 116, 0.5);
}

.chat-fab.is-dragging {
  cursor: grabbing;
  transition: none;
}

.chat-fab.has-unread {
  animation: pulse-badge 2s ease-in-out 3;
}

@keyframes pulse-badge {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.unread-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #e74c3c;
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  min-width: 20px;
  height: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 5px;
}

.chat-panel {
  position: fixed;
  width: 380px;
  height: 520px;
  background: #1a1a1a;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  z-index: 1200;
}

.chat-slide-enter-active,
.chat-slide-leave-active {
  transition: all 0.3s ease;
}

.chat-slide-enter-from,
.chat-slide-leave-to {
  opacity: 0;
  transform: translateY(20px);
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

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.messages-area::-webkit-scrollbar {
  width: 4px;
}

.messages-area::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 2px;
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

.message-wrapper.is-system {
  justify-content: center;
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

.input-area {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 12px 16px;
  background: #1a1a1a;
  border-top: 1px solid #333;
  flex-shrink: 0;
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
  width: 38px;
  height: 38px;
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

@media (max-width: 480px) {
  .chat-panel {
    width: 100vw;
    height: 100vh;
    border-radius: 0;
  }
}
</style>
