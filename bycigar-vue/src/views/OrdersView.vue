<script setup>
import { ref, onMounted } from 'vue'

const orders = ref([])
const loading = ref(false)
const error = ref(null)

const statusMap = {
  'pending': '待处理',
  'paid': '已支付',
  'shipped': '已发货',
  'completed': '已完成',
  'cancelled': '已取消'
}

const statusClass = {
  'pending': 'status-pending',
  'paid': 'status-paid',
  'shipped': 'status-shipped',
  'completed': 'status-completed',
  'cancelled': 'status-cancelled'
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
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

function formatPrice(price) {
  return '$' + Number(price).toFixed(2)
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
              <span class="item-price">{{ formatPrice(item.price) }}</span>
            </div>
          </div>
          
          <div class="order-footer">
            <span class="order-total">总计: {{ formatPrice(order.total) }}</span>
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
.order-footer {
  padding: 15px 20px;
  background: #2a2a2a;
  border-top: 1px solid #333;
  text-align: right;
}
.order-total {
  color: #d4a574;
  font-size: 18px;
  font-weight: bold;
}
</style>
