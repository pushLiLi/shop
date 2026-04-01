<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useToastStore } from '../../stores/toast'

const authStore = useAuthStore()
const toast = useToastStore()

const options = ref({
  orders: false,
  users: false,
  conversations: false,
  products: false
})

const loading = ref(false)
const confirmText = ref('')
const showConfirm = ref(false)
const cleanupResult = ref(null)

const confirmRequired = '确认清理'

const hasSelected = computed(() => {
  return options.value.orders || options.value.users || options.value.conversations || options.value.products
})

const canConfirm = computed(() => {
  return confirmText.value === confirmRequired
})

const optionCards = [
  {
    key: 'orders',
    title: '订单数据',
    items: ['订单', '订单项', '支付凭证', '订单汇总']
  },
  {
    key: 'users',
    title: '用户数据',
    items: ['客户用户（保留管理员）', '购物车', '收藏', '收货地址', '通知']
  },
  {
    key: 'conversations',
    title: '会话数据',
    items: ['聊天会话', '聊天消息', '评分', '快捷回复']
  },
  {
    key: 'products',
    title: '商品数据',
    items: ['商品（含软删除）', '分类（含软删除）', '关联购物车/收藏']
  }
]

function requestCleanup() {
  if (!hasSelected.value) {
    toast.error('请至少选择一项清理内容')
    return
  }
  showConfirm.value = true
  confirmText.value = ''
  cleanupResult.value = null
}

function cancelConfirm() {
  showConfirm.value = false
  confirmText.value = ''
}

async function executeCleanup() {
  if (!canConfirm.value) return
  loading.value = true
  try {
    const response = await fetch('/api/admin/cleanup', {
      method: 'POST',
      headers: authStore.getAuthHeaders(),
      body: JSON.stringify(options.value)
    })
    const data = await response.json()
    if (!response.ok) {
      toast.error(data.error || '清理失败')
      return
    }
    cleanupResult.value = data.result
    toast.success(data.message || '清理完成')
    showConfirm.value = false
    confirmText.value = ''
    options.value = { orders: false, users: false, conversations: false, products: false }
  } catch (e) {
    toast.error('网络错误，请重试')
  } finally {
    loading.value = false
  }
}

function formatCount(count) {
  return count === 0 ? '0' : count.toLocaleString()
}
</script>

<template>
  <div class="cleanup-page">
    <div class="warning-banner">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
        <line x1="12" y1="9" x2="12" y2="13"></line>
        <line x1="12" y1="17" x2="12.01" y2="17"></line>
      </svg>
      <span>此功能将永久删除数据，操作不可逆，请谨慎使用</span>
    </div>

    <div class="option-grid">
      <div
        v-for="card in optionCards"
        :key="card.key"
        class="option-card"
        :class="{ active: options[card.key] }"
        @click="options[card.key] = !options[card.key]"
      >
        <div class="card-header">
          <label class="card-checkbox">
            <input type="checkbox" v-model="options[card.key]" />
            <span class="checkmark"></span>
          </label>
          <h3>{{ card.title }}</h3>
        </div>
        <ul class="card-items">
          <li v-for="item in card.items" :key="item">{{ item }}</li>
        </ul>
      </div>
    </div>

    <div class="action-bar">
      <button class="btn-cleanup" :disabled="!hasSelected" @click="requestCleanup">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="3 6 5 6 21 6"></polyline>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
        </svg>
        执行清理
      </button>
    </div>

    <div v-if="cleanupResult" class="result-panel">
      <h3>清理结果</h3>
      <div class="result-grid">
        <template v-if="cleanupResult.ordersDeleted > 0 || cleanupResult.orderItemsDeleted > 0 || cleanupResult.paymentProofsDeleted > 0 || cleanupResult.orderSummariesDeleted > 0">
          <div class="result-group">
            <div class="result-group-title">订单相关</div>
            <div v-if="cleanupResult.ordersDeleted > 0" class="result-item">
              <span>订单</span><span>{{ formatCount(cleanupResult.ordersDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.orderItemsDeleted > 0" class="result-item">
              <span>订单项</span><span>{{ formatCount(cleanupResult.orderItemsDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.paymentProofsDeleted > 0" class="result-item">
              <span>支付凭证</span><span>{{ formatCount(cleanupResult.paymentProofsDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.orderSummariesDeleted > 0" class="result-item">
              <span>订单汇总</span><span>{{ formatCount(cleanupResult.orderSummariesDeleted) }} 条</span>
            </div>
          </div>
        </template>
        <template v-if="cleanupResult.usersDeleted > 0 || cleanupResult.cartItemsDeleted > 0 || cleanupResult.favoritesDeleted > 0 || cleanupResult.addressesDeleted > 0 || cleanupResult.notificationsDeleted > 0">
          <div class="result-group">
            <div class="result-group-title">用户相关</div>
            <div v-if="cleanupResult.usersDeleted > 0" class="result-item">
              <span>客户用户</span><span>{{ formatCount(cleanupResult.usersDeleted) }} 个</span>
            </div>
            <div v-if="cleanupResult.cartItemsDeleted > 0" class="result-item">
              <span>购物车</span><span>{{ formatCount(cleanupResult.cartItemsDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.favoritesDeleted > 0" class="result-item">
              <span>收藏</span><span>{{ formatCount(cleanupResult.favoritesDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.addressesDeleted > 0" class="result-item">
              <span>收货地址</span><span>{{ formatCount(cleanupResult.addressesDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.notificationsDeleted > 0" class="result-item">
              <span>通知</span><span>{{ formatCount(cleanupResult.notificationsDeleted) }} 条</span>
            </div>
          </div>
        </template>
        <template v-if="cleanupResult.conversationsDeleted > 0 || cleanupResult.messagesDeleted > 0 || cleanupResult.ratingsDeleted > 0 || cleanupResult.quickRepliesDeleted > 0">
          <div class="result-group">
            <div class="result-group-title">会话相关</div>
            <div v-if="cleanupResult.conversationsDeleted > 0" class="result-item">
              <span>会话</span><span>{{ formatCount(cleanupResult.conversationsDeleted) }} 个</span>
            </div>
            <div v-if="cleanupResult.messagesDeleted > 0" class="result-item">
              <span>消息</span><span>{{ formatCount(cleanupResult.messagesDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.ratingsDeleted > 0" class="result-item">
              <span>评分</span><span>{{ formatCount(cleanupResult.ratingsDeleted) }} 条</span>
            </div>
            <div v-if="cleanupResult.quickRepliesDeleted > 0" class="result-item">
              <span>快捷回复</span><span>{{ formatCount(cleanupResult.quickRepliesDeleted) }} 条</span>
            </div>
          </div>
        </template>
        <template v-if="cleanupResult.productsDeleted > 0 || cleanupResult.categoriesDeleted > 0">
          <div class="result-group">
            <div class="result-group-title">商品相关</div>
            <div v-if="cleanupResult.productsDeleted > 0" class="result-item">
              <span>商品</span><span>{{ formatCount(cleanupResult.productsDeleted) }} 个</span>
            </div>
            <div v-if="cleanupResult.categoriesDeleted > 0" class="result-item">
              <span>分类</span><span>{{ formatCount(cleanupResult.categoriesDeleted) }} 个</span>
            </div>
          </div>
        </template>
      </div>
    </div>

    <div v-if="showConfirm" class="modal-overlay" @click.self="cancelConfirm">
      <div class="modal-content">
        <h3>确认清理数据</h3>
        <p class="modal-desc">此操作将永久删除选中的数据，且不可恢复。请输入 <strong>{{ confirmRequired }}</strong> 以确认操作。</p>
        <div class="selected-summary">
          <span v-if="options.orders" class="tag">订单数据</span>
          <span v-if="options.users" class="tag">用户数据</span>
          <span v-if="options.conversations" class="tag">会话数据</span>
          <span v-if="options.products" class="tag">商品数据</span>
        </div>
        <input
          v-model="confirmText"
          type="text"
          class="confirm-input"
          :placeholder="`请输入「${confirmRequired}」`"
          @keyup.enter="executeCleanup"
        />
        <div class="modal-actions">
          <button class="btn-cancel" @click="cancelConfirm">取消</button>
          <button class="btn-danger" :disabled="!canConfirm || loading" @click="executeCleanup">
            {{ loading ? '清理中...' : '确认清理' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cleanup-page {
  max-width: 900px;
}

.warning-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 18px;
  background: #fff3e0;
  border: 1px solid #ffb74d;
  border-radius: 8px;
  color: #e65100;
  font-size: 14px;
  margin-bottom: 24px;
}

.warning-banner svg {
  flex-shrink: 0;
}

.option-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.option-card {
  background: #fff;
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.option-card:hover {
  border-color: #bdbdbd;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.option-card.active {
  border-color: #d4a574;
  background: #fdf8f3;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.card-checkbox {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.card-checkbox input {
  width: 18px;
  height: 18px;
  accent-color: #d4a574;
  cursor: pointer;
}

.card-items {
  margin: 0;
  padding-left: 20px;
  list-style: disc;
  color: #666;
  font-size: 13px;
  line-height: 1.8;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
}

.btn-cleanup {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 28px;
  background: #d32f2f;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cleanup:hover:not(:disabled) {
  background: #b71c1c;
}

.btn-cleanup:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.result-panel {
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  padding: 20px;
  margin-top: 24px;
}

.result-panel h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 16px;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.result-group {
  background: #fafafa;
  border-radius: 8px;
  padding: 14px;
}

.result-group-title {
  font-size: 13px;
  font-weight: 600;
  color: #d4a574;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.result-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 14px;
  color: #555;
}

.result-item span:last-child {
  font-weight: 500;
  color: #1a1a1a;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  width: 480px;
  max-width: 90vw;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.modal-content h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 12px;
}

.modal-desc {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 16px;
}

.modal-desc strong {
  color: #d32f2f;
}

.selected-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.tag {
  padding: 4px 12px;
  background: #fce4ec;
  color: #c62828;
  border-radius: 4px;
  font-size: 13px;
}

.confirm-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
  margin-bottom: 20px;
}

.confirm-input:focus {
  border-color: #d4a574;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-cancel {
  padding: 10px 20px;
  background: #f5f5f5;
  border: none;
  border-radius: 6px;
  color: #666;
  cursor: pointer;
  font-size: 14px;
}

.btn-cancel:hover {
  background: #e0e0e0;
}

.btn-danger {
  padding: 10px 20px;
  background: #d32f2f;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-danger:hover:not(:disabled) {
  background: #b71c1c;
}

.btn-danger:disabled {
  background: #ccc;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .option-grid {
    grid-template-columns: 1fr;
  }

  .result-grid {
    grid-template-columns: 1fr;
  }
}
</style>
