<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useContactMethodsStore } from '../stores/contactMethods'
import { useToastStore } from '../stores/toast'

const store = useContactMethodsStore()
const toast = useToastStore()

const expanded = ref(false)
const qrModal = ref(null)
const widgetRef = ref(null)

const activeMethods = computed(() => store.methods.filter(m => m.isActive))

function toggle() {
  expanded.value = !expanded.value
}

function handleClick(method) {
  const type = method.type
  const value = method.value

  switch (type) {
    case 'whatsapp':
      window.open(`https://wa.me/${value}`, '_blank')
      break
    case 'telegram':
      window.open(`https://t.me/${value}`, '_blank')
      break
    case 'phone':
      window.location.href = `tel:${value}`
      break
    case 'email':
      window.location.href = `mailto:${value}`
      break
    case 'wechat':
    case 'qq':
      qrModal.value = method
      break
    case 'custom':
      if (value.startsWith('http')) {
        window.open(value, '_blank')
      } else {
        copyToClipboard(value)
      }
      break
  }
  expanded.value = false
}

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
    toast.success('已复制到剪贴板')
  } catch {
    const ta = document.createElement('textarea')
    ta.value = text
    ta.style.position = 'fixed'
    ta.style.left = '-9999px'
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    toast.success('已复制到剪贴板')
  }
}

function closeQRModal() {
  qrModal.value = null
}

function handleClickOutside(e) {
  if (widgetRef.value && !widgetRef.value.contains(e.target)) {
    expanded.value = false
  }
}

onMounted(() => {
  store.fetchMethods()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div v-if="activeMethods.length > 0" class="contact-widget" ref="widgetRef">
    <Transition name="expand">
      <div v-if="expanded" class="contact-menu">
        <button
          v-for="method in activeMethods"
          :key="method.id"
          class="contact-item"
          @click="handleClick(method)"
          :title="method.label"
        >
          <span class="contact-icon" :class="method.type">
            <svg v-if="method.type === 'phone'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
            <svg v-else-if="method.type === 'email'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            <svg v-else-if="method.type === 'whatsapp'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
            <svg v-else-if="method.type === 'wechat'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8.5 2C5.46 2 3 4.46 3 7.5c0 1.68.75 3.18 1.94 4.2L4 14l2.6-1.3c.6.2 1.24.3 1.9.3.34 0 .67-.03 1-.08"></path><path d="M15.5 8c-3.04 0-5.5 2.46-5.5 5.5 0 3.04 2.46 5.5 5.5 5.5.66 0 1.3-.1 1.9-.3L20 20l-.94-2.3A5.47 5.47 0 0 0 21 13.5C21 10.46 18.54 8 15.5 8z"></path></svg>
            <svg v-else-if="method.type === 'qq'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><path d="M8 14s1.5 2 4 2 4-2 4-2"></path><line x1="9" y1="9" x2="9.01" y2="9"></line><line x1="15" y1="9" x2="15.01" y2="9"></line></svg>
            <svg v-else-if="method.type === 'telegram'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
            <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path></svg>
          </span>
          <span class="contact-label">{{ method.label }}</span>
        </button>
      </div>
    </Transition>

    <button class="contact-fab" :class="{ expanded }" @click.stop="toggle" aria-label="联系方式">
      <svg v-if="!expanded" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"></path>
      </svg>
      <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>

    <Transition name="qr-fade">
      <div v-if="qrModal" class="qr-overlay" @click.self="closeQRModal">
        <div class="qr-modal">
          <button class="qr-close" @click="closeQRModal">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
          <h3 class="qr-title">{{ qrModal.label }}</h3>
          <img v-if="qrModal.qrCodeUrl" :src="qrModal.qrCodeUrl" alt="二维码" class="qr-image" />
          <div class="qr-value">
            <span>{{ qrModal.type === 'wechat' ? '微信号' : 'QQ号' }}: {{ qrModal.value }}</span>
            <button class="qr-copy-btn" @click="copyToClipboard(qrModal.value)">复制</button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.contact-widget {
  position: fixed;
  left: 24px;
  bottom: 24px;
  z-index: 1100;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.contact-fab {
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
  user-select: none;
}

.contact-fab:hover {
  transform: scale(1.08);
  box-shadow: 0 6px 20px rgba(212, 165, 116, 0.5);
}

.contact-fab.expanded {
  background: #333;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.contact-menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
  align-items: flex-start;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 24px;
  cursor: pointer;
  color: #fff;
  font-size: 13px;
  white-space: nowrap;
  transition: all 0.2s;
  min-width: 120px;
}

.contact-item:hover {
  background: #2a2a2a;
  border-color: #d4a574;
}

.contact-icon {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #fff;
}

.contact-icon.phone { background: #4caf50; }
.contact-icon.email { background: #2196f3; }
.contact-icon.whatsapp { background: #25d366; }
.contact-icon.wechat { background: #07c160; }
.contact-icon.qq { background: #12b7f5; }
.contact-icon.telegram { background: #0088cc; }
.contact-icon.custom { background: #ff9800; }

.contact-label {
  color: #eee;
  font-weight: 500;
}

.expand-enter-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.expand-leave-active {
  transition: all 0.15s ease-in;
}

.expand-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.expand-leave-to {
  opacity: 0;
  transform: translateY(5px);
}

.qr-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1150;
  cursor: pointer;
}

.qr-modal {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  position: relative;
  min-width: 280px;
  cursor: default;
}

.qr-close {
  position: absolute;
  top: 12px;
  right: 12px;
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  padding: 4px;
  display: flex;
  transition: color 0.2s;
}

.qr-close:hover {
  color: #fff;
}

.qr-title {
  color: #d4a574;
  font-size: 16px;
  margin: 0 0 16px;
}

.qr-image {
  width: 200px;
  height: 200px;
  object-fit: contain;
  border-radius: 8px;
  border: 1px solid #333;
  margin-bottom: 16px;
}

.qr-value {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #ccc;
  font-size: 14px;
}

.qr-copy-btn {
  padding: 4px 12px;
  background: #d4a574;
  color: #0f0f0f;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: background 0.2s;
}

.qr-copy-btn:hover {
  background: #e0b88a;
}

.qr-fade-enter-active,
.qr-fade-leave-active {
  transition: opacity 0.2s ease;
}

.qr-fade-enter-from,
.qr-fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .contact-widget {
    left: 16px;
    bottom: 16px;
  }

  .contact-fab {
    width: 44px;
    height: 44px;
  }

  .contact-item {
    padding: 8px 14px;
    font-size: 12px;
    min-width: 100px;
  }

  .contact-icon {
    width: 24px;
    height: 24px;
  }
}
</style>
