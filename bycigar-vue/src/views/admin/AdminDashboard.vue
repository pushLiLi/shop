<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'

const API_BASE = '/api'
const authStore = useAuthStore()
const showRevenue = computed(() => authStore.isSuperAdmin)

const stats = ref({})
const recentOrders = ref([])
const lowStockProducts = ref([])
const topProducts = ref([])
const loading = ref(true)

const collapsedOrders = ref(false)
const collapsedLowStock = ref(false)
const collapsedTopProducts = ref(false)

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const statusLabels = {
  pending: '待处理',
  processing: '处理中',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消'
}

const statusBadgeClass = (status) => {
  const map = {
    pending: 'badge-warning',
    processing: 'badge-info',
    shipped: 'badge-primary',
    completed: 'badge-success',
    cancelled: 'badge-danger'
  }
  return map[status] || 'badge-default'
}

const fetchStats = async () => {
  const res = await fetch(`${API_BASE}/admin/dashboard/stats`, { headers: authHeaders() })
  stats.value = await res.json()
}

const fetchRecentOrders = async () => {
  const res = await fetch(`${API_BASE}/admin/dashboard/recent-orders`, { headers: authHeaders() })
  const data = await res.json()
  recentOrders.value = data.orders || []
}

const fetchLowStock = async () => {
  const res = await fetch(`${API_BASE}/admin/dashboard/low-stock`, { headers: authHeaders() })
  const data = await res.json()
  lowStockProducts.value = data.products || []
}

const fetchTopProducts = async () => {
  const res = await fetch(`${API_BASE}/admin/dashboard/top-products`, { headers: authHeaders() })
  const data = await res.json()
  topProducts.value = data.products || []
}

const formatPrice = (price) => `¥${parseFloat(price || 0).toFixed(2)}`
const formatDate = (date) => new Date(date).toLocaleString('zh-CN')

onMounted(async () => {
  try {
    await Promise.all([fetchStats(), fetchRecentOrders(), fetchLowStock(), fetchTopProducts()])
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="admin-dashboard">
    <div v-if="loading" class="loading">加载中...</div>

    <template v-else>
      <div class="stats-grid" :class="{ 'stats-grid-3': !showRevenue }">
        <div class="stat-card" v-if="showRevenue">
          <div class="stat-icon revenue">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="1" x2="12" y2="23"></line><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path></svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatPrice(stats.totalRevenue) }}</div>
            <div class="stat-label">总营收</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon orders">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path><line x1="3" y1="6" x2="21" y2="6"></line><path d="M16 10a4 4 0 0 1-8 0"></path></svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalOrders }}</div>
            <div class="stat-label">总订单</div>
          </div>
          <div class="stat-extra" v-if="stats.pendingOrders > 0">
            <span class="pending-badge">{{ stats.pendingOrders }} 待处理</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon users">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalUsers }}</div>
            <div class="stat-label">注册用户</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon today">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.todayOrders }}</div>
            <div class="stat-label">今日订单</div>
          </div>
          <div class="stat-extra" v-if="showRevenue">
            <span class="today-revenue">{{ formatPrice(stats.todayRevenue) }}</span>
          </div>
        </div>
      </div>

      <div class="dashboard-grid">
        <div class="dashboard-section">
          <div class="section-header" @click="collapsedOrders = !collapsedOrders">
            <h3>近期订单</h3>
            <svg class="collapse-icon" :class="{ collapsed: collapsedOrders }" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#999" stroke-width="2"><polyline points="6 9 12 15 18 9"></polyline></svg>
          </div>
          <div class="section-body" v-show="!collapsedOrders">
            <table class="data-table">
              <thead>
                <tr>
                  <th>订单号</th>
                  <th>用户</th>
                  <th>金额</th>
                  <th>状态</th>
                  <th>时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="order in recentOrders" :key="order.id">
                  <td class="order-no">{{ order.orderNo }}</td>
                  <td>{{ order.userName || '-' }}</td>
                  <td>{{ formatPrice(order.total) }}</td>
                  <td><span class="badge" :class="statusBadgeClass(order.status)">{{ statusLabels[order.status] }}</span></td>
                  <td class="time-cell">{{ formatDate(order.createdAt) }}</td>
                </tr>
                <tr v-if="recentOrders.length === 0">
                  <td colspan="5" class="empty-text">暂无订单</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="dashboard-section">
          <div class="section-header" @click="collapsedLowStock = !collapsedLowStock">
            <h3>低库存预警</h3>
            <div class="section-header-right">
              <span class="section-sub">库存 ≤ 10</span>
              <svg class="collapse-icon" :class="{ collapsed: collapsedLowStock }" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#999" stroke-width="2"><polyline points="6 9 12 15 18 9"></polyline></svg>
            </div>
          </div>
          <div class="section-body" v-show="!collapsedLowStock">
            <div class="low-stock-list">
              <div v-for="product in lowStockProducts" :key="product.id" class="low-stock-item">
                <div class="product-info">
                  <span class="product-name">{{ product.name }}</span>
                </div>
                <div class="stock-value" :class="{ critical: product.stock <= 5 }">
                  {{ product.stock }} 件
                </div>
              </div>
              <div v-if="lowStockProducts.length === 0" class="empty-text">库存充足</div>
            </div>
          </div>
        </div>

        <div class="dashboard-section full-width">
          <div class="section-header" @click="collapsedTopProducts = !collapsedTopProducts">
            <h3>热销商品 TOP 10</h3>
            <svg class="collapse-icon" :class="{ collapsed: collapsedTopProducts }" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#999" stroke-width="2"><polyline points="6 9 12 15 18 9"></polyline></svg>
          </div>
          <div class="section-body" v-show="!collapsedTopProducts">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 60px">排名</th>
                  <th>商品名称</th>
                  <th style="width: 120px">销量</th>
                  <th v-if="showRevenue" style="width: 150px">营收</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(product, index) in topProducts" :key="product.productId">
                  <td><span class="rank" :class="'rank-' + (index + 1)">{{ index + 1 }}</span></td>
                  <td>{{ product.productName }}</td>
                  <td>{{ product.totalSold }} 件</td>
                  <td v-if="showRevenue" class="price">{{ formatPrice(product.revenue) }}</td>
                </tr>
                <tr v-if="topProducts.length === 0">
                  <td :colspan="showRevenue ? 4 : 3" class="empty-text">暂无销售数据</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.admin-dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.loading {
  padding: 40px;
  text-align: center;
  color: #666;
  background: #fff;
  border-radius: 8px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stats-grid-3 {
  grid-template-columns: repeat(3, 1fr);
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon.revenue { background: #fff8e1; color: #f57c00; }
.stat-icon.orders { background: #e3f2fd; color: #1565c0; }
.stat-icon.users { background: #e8f5e9; color: #2e7d32; }
.stat-icon.today { background: #fce4ec; color: #c62828; }

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
}

.stat-label {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
}

.stat-extra {
  position: absolute;
  top: 10px;
  right: 12px;
}

.pending-badge {
  background: #fff8e1;
  color: #f57c00;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
}

.today-revenue {
  font-size: 12px;
  color: #999;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.dashboard-section {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  overflow: hidden;
}

.dashboard-section.full-width {
  grid-column: 1 / -1;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  user-select: none;
}

.section-header:hover {
  background: #fafafa;
}

.section-header h3 {
  margin: 0;
  font-size: 15px;
  color: #1a1a1a;
}

.section-sub {
  font-size: 13px;
  color: #999;
}

.section-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.collapse-icon {
  transition: transform 0.2s ease;
  flex-shrink: 0;
}

.collapse-icon.collapsed {
  transform: rotate(-90deg);
}

.section-body {
  padding: 0;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 10px 15px;
  text-align: left;
  border-bottom: 1px solid #f5f5f5;
}

.data-table th {
  background: #fafafa;
  font-weight: 600;
  color: #333;
  font-size: 13px;
}

.data-table td {
  color: #666;
  font-size: 13px;
}

.order-no {
  font-family: monospace;
  font-size: 12px;
}

.time-cell {
  font-size: 12px;
  white-space: nowrap;
}

.badge {
  display: inline-block;
  padding: 3px 7px;
  border-radius: 4px;
  font-size: 11px;
}

.badge-success { background: #e8f5e9; color: #2e7d32; }
.badge-danger { background: #ffebee; color: #c62828; }
.badge-warning { background: #fff8e1; color: #f57c00; }
.badge-info { background: #e3f2fd; color: #1565c0; }
.badge-primary { background: #e8eaf6; color: #283593; }
.badge-default { background: #f5f5f5; color: #999; }

.empty-text {
  text-align: center;
  color: #999;
  padding: 30px 15px !important;
}

.low-stock-list {
  display: flex;
  flex-direction: column;
}

.low-stock-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  border-bottom: 1px solid #f5f5f5;
}

.low-stock-item:last-child {
  border-bottom: none;
}

.product-name {
  font-size: 13px;
  color: #333;
}

.stock-value {
  font-size: 13px;
  font-weight: 600;
  color: #f57c00;
}

.stock-value.critical {
  color: #c62828;
}

.rank {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
  background: #f5f5f5;
  color: #666;
}

.rank-1 { background: #fff8e1; color: #f57c00; }
.rank-2 { background: #f5f5f5; color: #666; }
.rank-3 { background: #fff3e0; color: #e65100; }

.price {
  color: #d4a574;
  font-weight: 600;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style>
