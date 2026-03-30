<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotificationsStore } from '../stores/notifications'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const router = useRouter()
const store = useNotificationsStore()
const toast = useToastStore()

const notification = computed(() => store.currentNotification)
const loading = computed(() => store.detailLoading)

const typeConfig = computed(() => {
  if (!notification.value) return {}
  const map = {
    order_status: {
      label: '订单状态',
      color: '#6b9fff',
      bg: '#1a2440',
      actionText: '查看相关订单',
      actionLink: notification.value.link || '/orders'
    },
    back_in_stock: {
      label: '到货通知',
      color: '#6bdf8f',
      bg: '#1a3a24',
      actionText: '查看商品详情',
      actionLink: notification.value.link || (notification.value.productId ? `/products/${notification.value.productId}` : null)
    },
    price_drop: {
      label: '降价提醒',
      color: '#f59e42',
      bg: '#3a2a1a',
      actionText: '查看商品详情',
      actionLink: notification.value.link || (notification.value.productId ? `/products/${notification.value.productId}` : null)
    }
  }
  return map[notification.value.type] || { label: '系统通知', color: '#888', bg: '#252525', actionText: null, actionLink: null }
})

function formatFullTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  return `${y}年${m}月${d}日 ${h}:${min}`
}

function goBack() {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}

function handleAction() {
  if (typeConfig.value.actionLink) {
    router.push(typeConfig.value.actionLink)
  }
}

async function handleDelete() {
  if (!notification.value) return
  const success = await store.deleteNotification(notification.value.id)
  if (success) {
    toast.success('通知已删除')
    router.replace('/orders')
  } else {
    toast.error('删除失败')
  }
}

onMounted(async () => {
  try {
    await store.fetchNotification(route.params.id)
  } catch {
    toast.error('通知不存在或已被删除')
    router.replace('/orders')
  }
})
</script>

<template>
  <div class="notification-detail-page">
    <div class="container">
      <div class="detail-header">
        <button class="back-btn" @click="goBack">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
          <span>返回</span>
        </button>
        <h1 class="page-title">通知详情</h1>
      </div>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else-if="notification" class="detail-card">
        <div class="detail-type-badge" :style="{ color: typeConfig.color, background: typeConfig.bg }">
          {{ typeConfig.label }}
        </div>

        <h2 class="detail-title">{{ notification.title }}</h2>

        <div class="detail-divider"></div>

        <div class="detail-content">{{ notification.content }}</div>

        <div class="detail-divider"></div>

        <div class="detail-time">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          <span>{{ formatFullTime(notification.createdAt) }}</span>
        </div>

        <div v-if="typeConfig.actionLink" class="detail-divider"></div>

        <button v-if="typeConfig.actionLink" class="action-btn" @click="handleAction">
          <span>{{ typeConfig.actionText }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>

        <div class="detail-divider"></div>

        <button class="delete-btn" @click="handleDelete">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
          <span>删除通知</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notification-detail-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.container {
  max-width: 640px;
  margin: 0 auto;
  padding: 0 15px;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 30px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  border: none;
  color: #999;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: color 0.2s, background 0.2s;
  font-size: 14px;
}

.back-btn:hover {
  color: #d4a574;
  background: rgba(255, 255, 255, 0.05);
}

.page-title {
  color: #d4a574;
  font-size: 22px;
  font-weight: 500;
  margin: 0;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #888;
  font-size: 14px;
}

.detail-card {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 32px;
}

.detail-type-badge {
  display: inline-block;
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 20px;
}

.detail-title {
  color: #eee;
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 0;
  line-height: 1.4;
}

.detail-divider {
  height: 1px;
  background: #2a2a2a;
  margin: 20px 0;
}

.detail-content {
  color: #bbb;
  font-size: 15px;
  line-height: 1.8;
  white-space: pre-wrap;
}

.detail-time {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 13px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 14px 24px;
  background: #d4a574;
  color: #0f0f0f;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s, transform 0.1s;
}

.action-btn:hover {
  background: #e0b88a;
}

.action-btn:active {
  transform: scale(0.98);
}

.delete-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 12px 24px;
  background: transparent;
  color: #e74c3c;
  border: 1px solid #3a2020;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.delete-btn:hover {
  background: rgba(231, 76, 60, 0.1);
  border-color: #e74c3c;
}

@media (max-width: 480px) {
  .notification-detail-page {
    padding: 20px 0 40px;
  }

  .detail-card {
    padding: 24px 20px;
  }

  .detail-title {
    font-size: 18px;
  }
}
</style>
