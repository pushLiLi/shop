<script setup>
import { ref, onMounted, computed } from 'vue'
import { useToastStore } from '../../stores/toast'

const API_BASE = '/api'
const toast = useToastStore()

const orders = ref([])
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 20

const filterStatus = ref('')
const search = ref('')

const showDetailModal = ref(false)
const detailOrder = ref(null)
const detailUser = ref(null)

const showStatusModal = ref(false)
const statusOrder = ref(null)
const selectedStatus = ref('')

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

const nextStatuses = computed(() => {
  if (!statusOrder.value) return []
  const transitions = {
    pending: ['processing', 'cancelled'],
    processing: ['shipped', 'cancelled'],
    shipped: ['completed'],
    completed: [],
    cancelled: []
  }
  return transitions[statusOrder.value.status] || []
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const fetchOrders = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({ page: currentPage.value, limit })
    if (search.value) params.append('search', search.value)
    if (filterStatus.value) params.append('status', filterStatus.value)

    const res = await fetch(`${API_BASE}/admin/orders?${params}`, { headers: authHeaders() })
    const data = await res.json()
    orders.value = data.orders || []
    totalPages.value = data.totalPages || 1
  } catch (e) {
    console.error('Error fetching orders:', e)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchOrders()
}

const resetFilters = () => {
  search.value = ''
  filterStatus.value = ''
  currentPage.value = 1
  fetchOrders()
}

const openDetail = async (order) => {
  try {
    const res = await fetch(`${API_BASE}/admin/orders/${order.id}`, { headers: authHeaders() })
    const data = await res.json()
    detailOrder.value = data.order
    detailUser.value = data.user
    showDetailModal.value = true
  } catch (e) {
    toast.error('获取订单详情失败')
  }
}

const openStatusModal = (order) => {
  statusOrder.value = order
  selectedStatus.value = ''
  showStatusModal.value = true
}

const updateStatus = async () => {
  if (!selectedStatus.value) {
    toast.error('请选择新状态')
    return
  }

  try {
    const res = await fetch(`${API_BASE}/admin/orders/${statusOrder.value.id}/status`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ status: selectedStatus.value })
    })
    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '更新失败')
    }
    toast.success('订单状态已更新')
    showStatusModal.value = false
    fetchOrders()
  } catch (e) {
    toast.error(e.message)
  }
}

const formatPrice = (price) => `¥${parseFloat(price).toFixed(2)}`
const formatDate = (date) => new Date(date).toLocaleString('zh-CN')

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchOrders()
  }
}
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchOrders()
  }
}

onMounted(() => fetchOrders())
</script>

<template>
  <div class="admin-orders">
    <div class="toolbar">
      <div class="toolbar-left">
        <div class="search-box">
          <input v-model="search" type="text" placeholder="搜索订单号..." @keyup.enter="handleSearch">
          <button class="btn-search" @click="handleSearch">搜索</button>
        </div>
        <select v-model="filterStatus" @change="handleSearch">
          <option value="">全部状态</option>
          <option v-for="(label, key) in statusLabels" :key="key" :value="key">{{ label }}</option>
        </select>
        <button class="btn-reset" @click="resetFilters">重置</button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th style="width: 60px">ID</th>
            <th style="width: 180px">订单号</th>
            <th>用户</th>
            <th style="width: 100px">商品数</th>
            <th style="width: 120px">总金额</th>
            <th style="width: 100px">状态</th>
            <th style="width: 170px">下单时间</th>
            <th style="width: 160px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td>{{ order.id }}</td>
            <td class="order-no">{{ order.orderNo }}</td>
            <td>{{ order.user?.name || order.user?.email || '-' }}</td>
            <td>{{ order.items?.length || 0 }} 件</td>
            <td>{{ formatPrice(order.total) }}</td>
            <td>
              <span class="badge" :class="statusBadgeClass(order.status)">
                {{ statusLabels[order.status] || order.status }}
              </span>
            </td>
            <td class="time-cell">{{ formatDate(order.createdAt) }}</td>
            <td>
              <button class="btn-view" @click="openDetail(order)">详情</button>
              <button
                v-if="order.status !== 'completed' && order.status !== 'cancelled'"
                class="btn-status"
                @click="openStatusModal(order)"
              >状态</button>
            </td>
          </tr>
          <tr v-if="orders.length === 0">
            <td colspan="8" class="empty-text">暂无订单数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="totalPages > 1" class="pagination">
      <button :disabled="currentPage === 1" @click="prevPage">上一页</button>
      <span>{{ currentPage }} / {{ totalPages }}</span>
      <button :disabled="currentPage === totalPages" @click="nextPage">下一页</button>
    </div>

    <div v-if="showDetailModal" class="modal-overlay" @click.self="showDetailModal = false">
      <div class="modal detail-modal">
        <div class="modal-header">
          <h3>订单详情 - {{ detailOrder?.orderNo }}</h3>
          <button class="modal-close" @click="showDetailModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="detail-section">
            <div class="section-title">订单信息</div>
            <div class="detail-grid">
              <div class="detail-item"><span class="label">订单号</span><span>{{ detailOrder?.orderNo }}</span></div>
              <div class="detail-item"><span class="label">状态</span><span class="badge" :class="statusBadgeClass(detailOrder?.status)">{{ statusLabels[detailOrder?.status] }}</span></div>
              <div class="detail-item"><span class="label">总金额</span><span class="price">{{ formatPrice(detailOrder?.total) }}</span></div>
              <div class="detail-item"><span class="label">下单时间</span><span>{{ formatDate(detailOrder?.createdAt) }}</span></div>
              <div class="detail-item" v-if="detailOrder?.remark"><span class="label">备注</span><span>{{ detailOrder.remark }}</span></div>
            </div>
          </div>

          <div class="detail-section">
            <div class="section-title">用户信息</div>
            <div class="detail-grid">
              <div class="detail-item"><span class="label">姓名</span><span>{{ detailUser?.name || '-' }}</span></div>
              <div class="detail-item"><span class="label">邮箱</span><span>{{ detailUser?.email || '-' }}</span></div>
            </div>
          </div>

          <div class="detail-section" v-if="detailOrder?.address">
            <div class="section-title">收货地址</div>
            <div class="detail-grid">
              <div class="detail-item"><span class="label">收件人</span><span>{{ detailOrder.address.fullName }}</span></div>
              <div class="detail-item"><span class="label">电话</span><span>{{ detailOrder.address.phone }}</span></div>
              <div class="detail-item full"><span class="label">地址</span><span>{{ detailOrder.address.addressLine1 }} {{ detailOrder.address.addressLine2 }} {{ detailOrder.address.city }} {{ detailOrder.address.state }} {{ detailOrder.address.zipCode }}</span></div>
            </div>
          </div>

          <div class="detail-section">
            <div class="section-title">商品列表</div>
            <div class="items-list">
              <div v-for="item in detailOrder?.items" :key="item.id" class="order-item">
                <img v-if="item.product?.imageUrl" :src="item.product.imageUrl" class="item-thumb">
                <div class="item-info">
                  <div class="item-name">{{ item.product?.name || '商品已删除' }}</div>
                  <div class="item-meta">¥{{ parseFloat(item.price).toFixed(2) }} x {{ item.quantity }}</div>
                </div>
                <div class="item-total">¥{{ (item.price * item.quantity).toFixed(2) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showStatusModal" class="modal-overlay" @click.self="showStatusModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>更新订单状态</h3>
          <button class="modal-close" @click="showStatusModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="current-status">
            当前状态：<span class="badge" :class="statusBadgeClass(statusOrder?.status)">{{ statusLabels[statusOrder?.status] }}</span>
          </div>
          <div class="form-group">
            <label>变更至</label>
            <select v-model="selectedStatus">
              <option value="">请选择</option>
              <option v-for="s in nextStatuses" :key="s" :value="s">{{ statusLabels[s] }}</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showStatusModal = false">取消</button>
          <button class="btn-save" @click="updateStatus">确认变更</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-orders {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  flex-wrap: wrap;
  gap: 10px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.search-box {
  display: flex;
  gap: 5px;
}

.search-box input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 200px;
}

select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
}

.btn-search,
.btn-reset {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-search {
  background: #1a1a1a;
  color: #fff;
}

.btn-search:hover {
  background: #333;
}

.btn-reset {
  background: #f5f5f5;
  color: #666;
}

.btn-reset:hover {
  background: #e0e0e0;
}

.loading {
  padding: 40px;
  text-align: center;
  color: #666;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f9f9f9;
  font-weight: 600;
  color: #333;
  white-space: nowrap;
}

.data-table td {
  color: #666;
}

.order-no {
  font-family: monospace;
  font-size: 13px;
}

.time-cell {
  font-size: 13px;
  white-space: nowrap;
}

.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.badge-success { background: #e8f5e9; color: #2e7d32; }
.badge-danger { background: #ffebee; color: #c62828; }
.badge-warning { background: #fff8e1; color: #f57c00; }
.badge-info { background: #e3f2fd; color: #1565c0; }
.badge-primary { background: #e8eaf6; color: #283593; }
.badge-default { background: #f5f5f5; color: #999; }

.btn-view,
.btn-status {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  margin-right: 5px;
  transition: all 0.2s;
}

.btn-view {
  background: #e3f2fd;
  color: #1565c0;
}

.btn-view:hover {
  background: #bbdefb;
}

.btn-status {
  background: #fff8e1;
  color: #f57c00;
}

.btn-status:hover {
  background: #ffecb3;
}

.empty-text {
  text-align: center;
  color: #999;
  padding: 40px !important;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  padding: 20px;
  border-top: 1px solid #eee;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: #fff;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  border-color: #d4a574;
  color: #d4a574;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.detail-modal {
  max-width: 700px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  font-size: 16px;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 15px 20px;
  border-top: 1px solid #eee;
}

.btn-cancel,
.btn-save {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-cancel {
  background: #f5f5f5;
  color: #666;
}

.btn-cancel:hover {
  background: #e0e0e0;
}

.btn-save {
  background: #d4a574;
  color: #fff;
}

.btn-save:hover {
  background: #c49464;
}

.current-status {
  margin-bottom: 15px;
  font-size: 14px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #333;
}

.form-group select {
  width: 100%;
  padding: 10px 12px;
}

.detail-section {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.detail-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
  padding-left: 10px;
  border-left: 3px solid #d4a574;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.detail-item {
  font-size: 14px;
  color: #666;
}

.detail-item.full {
  grid-column: 1 / -1;
}

.detail-item .label {
  color: #999;
  margin-right: 8px;
  min-width: 50px;
  display: inline-block;
}

.price {
  color: #d4a574;
  font-weight: 600;
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
  background: #fafafa;
  border-radius: 6px;
}

.item-thumb {
  width: 48px;
  height: 48px;
  object-fit: cover;
  border-radius: 4px;
}

.item-info {
  flex: 1;
}

.item-name {
  font-size: 14px;
  color: #333;
}

.item-meta {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
}

.item-total {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}
</style>
