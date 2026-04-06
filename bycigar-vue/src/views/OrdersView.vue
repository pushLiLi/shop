<script setup>
import { ref, onMounted } from 'vue'
import { useToastStore } from '../stores/toast'
import { useImageCompress } from '../composables/useImageCompress'
import { formatPriceByCurrency } from '../composables/useFormatPrice'

const toast = useToastStore()
const { compress } = useImageCompress()
const orders = ref([])
const loading = ref(false)
const error = ref(null)
const paymentProofs = ref({})
const reuploadFile = ref(null)
const reuploadingOrderId = ref(null)

const statusMap = {
  'pending': '待处理',
  'paid': '已支付',
  'processing': '处理中',
  'shipped': '已发货',
  'completed': '已完成',
  'cancelled': '已取消'
}

const statusClass = {
  'pending': 'status-pending',
  'paid': 'status-paid',
  'processing': 'status-processing',
  'shipped': 'status-shipped',
  'completed': 'status-completed',
  'cancelled': 'status-cancelled'
}

const proofStatusMap = {
  'pending': '待审核',
  'approved': '已通过',
  'rejected': '已驳回'
}

async function fetchOrders() {
  try {
    loading.value = true
    const token = localStorage.getItem('token')
    const res = await fetch('/api/orders', {
      headers: {
        'Authorization': token ? `Bearer ${token}` : ''
      }
    })
    const data = await res.json()
    orders.value = data.orders || []
    fetchPendingProofs()
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function fetchPendingProofs() {
  const token = localStorage.getItem('token')
  const headers = { 'Authorization': token ? `Bearer ${token}` : '' }
  const pendingOrders = orders.value.filter(o => o.status === 'pending')
  await Promise.all(pendingOrders.map(async (order) => {
    try {
      const res = await fetch(`/api/orders/${order.id}/payment-proof`, { headers })
      const data = await res.json()
      if (data.paymentProof) {
        paymentProofs.value[order.id] = data.paymentProof
      }
    } catch (e) {
      // ignore
    }
  }))
}

async function handleReupload(orderId) {
  if (!reuploadFile.value) return
  try {
    reuploadingOrderId.value = orderId
    const token = localStorage.getItem('token')
    const proof = paymentProofs.value[orderId]
    const formData = new FormData()
    const compressed = await compress(reuploadFile.value, { maxWidth: 1920, maxHeight: 1920, quality: 0.9 })
    formData.append('file', compressed, 'proof.jpg')
    formData.append('paymentMethodId', proof?.paymentMethodId || proof?.payment_method_id || '')

    const res = await fetch(`/api/orders/${orderId}/payment-proof`, {
      method: 'POST',
      headers: { 'Authorization': token ? `Bearer ${token}` : '' },
      body: formData
    })
    const data = await res.json()
    if (res.ok) {
      toast.success('付款截图已重新上传')
      paymentProofs.value[orderId] = data.paymentProof
      reuploadFile.value = null
    } else {
      toast.error(data.error || '上传失败')
    }
  } catch (e) {
    toast.error('上传失败')
  } finally {
    reuploadingOrderId.value = null
  }
}

function onReuploadFileChange(e) {
  reuploadFile.value = e.target.files[0] || null
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}


onMounted(() => {
  fetchOrders()
})
</script>

<template>
  <div class="orders-page">
    <div class="container">
      <h1 class="page-title">我的订单</h1>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else-if="error" class="error">{{ error }}</div>

      <div v-else-if="orders.length === 0" class="empty">
        <p>暂无订单</p>
        <router-link to="/" class="link">去购物</router-link>
      </div>

      <div v-else class="orders-list">
        <div v-for="order in orders" :key="order.id" class="order-card">
          <div class="order-header">
            <div class="order-info">
              <span class="order-id">订单号: {{ order.orderNo }}</span>
              <span class="order-date">{{ formatDate(order.createdAt) }}</span>
            </div>
            <span class="order-status" :class="statusClass[order.status]">
              {{ statusMap[order.status] || order.status }}
            </span>
          </div>

          <div class="order-items-list">
            <div v-for="item in order.items" :key="item.id" class="order-item-row">
              <span class="item-name">{{ item.product?.name || '商品' }}</span>
              <span class="item-qty">x{{ item.quantity }}</span>
              <span class="item-price">{{ formatPriceByCurrency(item.price, item.currency) }}</span>
            </div>
          </div>

          <div v-if="order.trackingCompany || order.trackingNumber" class="tracking-bar">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>
            <span class="tracking-text">{{ order.trackingCompany }}：{{ order.trackingNumber }}</span>
          </div>

          <div v-if="paymentProofs[order.id]" class="proof-status-bar">
            <div class="proof-info">
              <span class="proof-status-label">付款凭证：</span>
              <span :class="['proof-badge', `proof-${paymentProofs[order.id].status}`]">
                {{ proofStatusMap[paymentProofs[order.id].status] }}
              </span>
              <span v-if="paymentProofs[order.id].rejectReason" class="reject-reason">
                （{{ paymentProofs[order.id].rejectReason }}）
              </span>
            </div>
            <div v-if="paymentProofs[order.id].status === 'rejected'" class="reupload-area">
              <input type="file" accept="image/*" @change="onReuploadFileChange" class="reupload-input" />
              <button
                class="reupload-btn"
                :disabled="!reuploadFile || reuploadingOrderId === order.id"
                @click="handleReupload(order.id)"
              >
                {{ reuploadingOrderId === order.id ? '上传中...' : '重新上传' }}
              </button>
            </div>
          </div>

          <div class="order-footer">
            <router-link :to="`/orders/${order.id}`" class="btn-detail">查看详情</router-link>
            <span class="order-total">总计: {{ formatPriceByCurrency(order.total, 'CNY') }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.orders-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}
.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 15px;
}
.page-title {
  color: #d4a574;
  font-size: 28px;
  margin-bottom: 30px;
  border-bottom: 2px solid #d4a574;
  padding-bottom: 10px;
}
.loading, .empty {
  text-align: center;
  padding: 60px 20px;
  color: #888;
}
.error {
  color: #e74;
  text-align: center;
  padding: 20px;
}
.link {
  color: #d4a574;
  text-decoration: none;
}
.link:hover {
  text-decoration: underline;
}
.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.order-card {
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
}
.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #2a2a2a;
  border-bottom: 1px solid #333;
}
.order-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.order-id {
  color: #fff;
  font-weight: bold;
}
.order-date {
  color: #888;
  font-size: 13px;
}
.order-status {
  padding: 5px 12px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: bold;
}
.status-pending {
  background: #f0ad4e;
  color: #1a1a1a;
}
.status-paid {
  background: #5cb85c;
  color: #1a1a1a;
}
.status-processing {
  background: #6c757d;
  color: #fff;
}
.status-shipped{
  background: #5bc0de;
  color: #1a1a1a;
}
.status-completed{
  background: #d4a574;
  color: #1a1a1a;
}
.status-cancelled{
  background: #d9534f;
  color: #fff;
}
.order-items-list {
  padding: 15px 20px;
}
.order-item-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #2a2a2a;
  color: #ccc;
  font-size: 14px;
}
.order-item-row:last-child {
  border-bottom: none;
}
.item-qty {
  color: #888;
}
.item-price {
  color: #d4a574;
}

.tracking-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-top: 1px solid #2a2a2a;
  color: #5bc0de;
  font-size: 13px;
}

.tracking-text {
  color: #ccc;
  font-family: monospace;
  font-size: 13px;
}

.proof-status-bar {
  padding: 12px 20px;
  border-top: 1px solid #2a2a2a;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
.proof-info {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}
.proof-status-label {
  color: #888;
  font-size: 13px;
}
.proof-badge {
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}
.proof-pending {
  background: #f0ad4e;
  color: #1a1a1a;
}
.proof-approved {
  background: #5cb85c;
  color: #1a1a1a;
}
.proof-rejected {
  background: #d9534f;
  color: #fff;
}
.reject-reason {
  color: #d9534f;
  font-size: 12px;
}
.reupload-area {
  display: flex;
  align-items: center;
  gap: 8px;
}
.reupload-input {
  font-size: 13px;
  color: #ccc;
}
.reupload-input::file-selector-button {
  background: #2a2a2a;
  color: #ccc;
  border: 1px solid #444;
  padding: 4px 10px;
  border-radius: 4px;
  cursor: pointer;
}
.reupload-btn {
  background: #d4a574;
  color: #1a1a1a;
  border: none;
  padding: 5px 14px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  font-weight: 600;
}
.reupload-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.reupload-btn:hover:not(:disabled) {
  background: #e5b584;
}

.order-footer {
  padding: 15px 20px;
  background: #2a2a2a;
  border-top: 1px solid #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.btn-detail {
  color: #d4a574;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s;
}

.btn-detail:hover {
  color: #e5b584;
}
.order-total {
  color: #d4a574;
  font-size: 18px;
  font-weight: bold;
}

@media (max-width: 768px) {
  .orders-page {
    padding: 20px 0 40px;
  }

  .page-title {
    font-size: 22px;
    margin-bottom: 20px;
  }

  .order-header {
    flex-wrap: wrap;
    gap: 8px;
    padding: 12px 15px;
  }

  .order-items-list {
    padding: 12px 15px;
  }

  .order-item-row {
    flex-wrap: wrap;
    gap: 4px;
  }

  .tracking-text {
    word-break: break-all;
    font-size: 12px;
  }

  .order-footer {
    flex-wrap: wrap;
    gap: 10px;
    padding: 12px 15px;
  }

  .order-total {
    font-size: 16px;
  }
}

@media (max-width: 576px) {
  .orders-page {
    padding: 15px 0 30px;
  }

  .page-title {
    font-size: 20px;
  }

  .order-header {
    padding: 10px 12px;
  }

  .order-items-list {
    padding: 10px 12px;
  }

  .order-footer {
    padding: 10px 12px;
  }

  .proof-status-bar {
    padding: 10px 12px;
    flex-direction: column;
    align-items: flex-start;
  }

  .reupload-area {
    flex-wrap: wrap;
  }
}
</style>
