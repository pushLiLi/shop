<script setup>
import { computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationsStore } from '../stores/notifications'

const store = useNotificationsStore()
const router = useRouter()

const items = computed(() => store.items)
const unreadCount = computed(() => store.unreadCount)
const loading = computed(() => store.loading)

function formatTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return minutes + '分钟前'
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return hours + '小时前'
  const days = Math.floor(hours / 24)
  if (days < 30) return days + '天前'
  return date.toLocaleDateString('zh-CN')
}

function getTypeIcon(type) {
  switch (type) {
    case 'order_status':
      return `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>`
    case 'back_in_stock':
      return `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`
    case 'price_drop':
      return `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"></polyline><polyline points="17 6 23 6 23 12"></polyline></svg>`
    default:
      return `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`
  }
}

async function handleClick(item) {
  store.closePanel()
  router.push(`/notifications/${item.id}`)
}

async function handleMarkAllRead() {
  await store.markAllRead()
}

function closePanel() {
  store.closePanel()
}

function handleOverlayClick(e) {
  if (e.target === e.currentTarget) {
    closePanel()
  }
}

watch(() => store.isOpen, (isOpen) => {
  if (isOpen) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="store.isOpen" class="notification-overlay" @click="handleOverlayClick">
        <div class="notification-panel">
          <div class="panel-header">
            <h2>通知</h2>
            <div class="panel-header-actions">
              <button v-if="unreadCount > 0" class="mark-all-btn" @click="handleMarkAllRead">全部已读</button>
              <button class="close-btn" @click="closePanel">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>

          <div class="panel-content">
            <div v-if="loading" class="loading">加载中...</div>

            <div v-else-if="items.length === 0" class="empty-state">
              <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
              </svg>
              <p>暂无通知</p>
            </div>

            <div v-else class="notification-list">
              <div
                v-for="item in items"
                :key="item.id"
                class="notification-item"
                :class="{ unread: !item.isRead }"
                @click="handleClick(item)"
              >
                <div class="item-dot" v-if="!item.isRead"></div>
                <div class="item-icon" :class="item.type" v-html="getTypeIcon(item.type)"></div>
                <div class="item-body">
                  <div class="item-title">{{ item.title }}</div>
                  <div class="item-content">{{ item.content }}</div>
                  <div class="item-time">{{ formatTime(item.createdAt) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.notification-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  justify-content: flex-end;
}

.notification-panel {
  width: 400px;
  max-width: 100%;
  height: 100%;
  background: #1a1a1a;
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.3);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #333;
}

.panel-header h2 {
  color: #d4a574;
  font-size: 20px;
  font-weight: 500;
  margin: 0;
}

.panel-header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mark-all-btn {
  background: transparent;
  border: none;
  color: #d4a574;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px;
  transition: opacity 0.2s;
}

.mark-all-btn:hover {
  opacity: 0.8;
}

.close-btn {
  background: transparent;
  border: none;
  color: #888;
  cursor: pointer;
  padding: 5px;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #fff;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #888;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #888;
}

.empty-state svg {
  opacity: 0.3;
  margin-bottom: 20px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.notification-list {
  display: flex;
  flex-direction: column;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  cursor: pointer;
  transition: background 0.2s;
  position: relative;
  border-bottom: 1px solid #252525;
}

.notification-item:hover {
  background: #252525;
}

.notification-item.unread {
  background: #1f1f1f;
}

.notification-item.unread:hover {
  background: #282828;
}

.item-dot {
  position: absolute;
  top: 20px;
  left: 10px;
  width: 6px;
  height: 6px;
  background: #d4a574;
  border-radius: 50%;
  flex-shrink: 0;
}

.item-icon {
  color: #888;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: #252525;
  border-radius: 50%;
}

.item-icon.order_status {
  color: #6b9fff;
  background: #1a2440;
}

.item-icon.back_in_stock {
  color: #6bdf8f;
  background: #1a3a24;
}

.item-icon.price_drop {
  color: #f59e42;
  background: #3a2a1a;
}

.item-body {
  flex: 1;
  min-width: 0;
}

.item-title {
  color: #eee;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.item-content {
  color: #999;
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 6px;
}

.item-time {
  color: #666;
  font-size: 12px;
}

.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-enter-active .notification-panel,
.drawer-leave-active .notification-panel {
  transition: transform 0.3s ease;
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from .notification-panel,
.drawer-leave-to .notification-panel {
  transform: translateX(100%);
}

@media (max-width: 480px) {
  .notification-panel {
    width: 100%;
  }
}
</style>
