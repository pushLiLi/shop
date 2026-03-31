<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useToastStore } from '../../stores/toast'

const API_BASE = '/api'
const toast = useToastStore()
const authStore = useAuthStore()

const users = ref([])
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 20
const sortBy = ref('createdAt')
const sortOrder = ref('desc')
const search = ref('')
const activeTab = ref('')
let fetchController = null

const showDetailModal = ref(false)
const detailUser = ref(null)
const detailOrders = ref([])
const detailAddresses = ref([])

const showRoleModal = ref(false)
const roleUser = ref(null)
const selectedRole = ref('')

const showResetModal = ref(false)
const resetUser = ref(null)

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const roleLabels = {
  admin: '超级管理员',
  service: '管理员',
  customer: '客户'
}

const fetchUsers = async () => {
  if (fetchController) fetchController.abort()
  fetchController = new AbortController()
  loading.value = true
  try {
    const params = new URLSearchParams({ page: currentPage.value, limit })
    if (search.value) params.append('search', search.value)
    if (activeTab.value) params.append('role', activeTab.value)
    if (sortBy.value) params.append('sortBy', sortBy.value)
    if (sortOrder.value) params.append('sortOrder', sortOrder.value)

    const res = await fetch(`${API_BASE}/admin/users?${params}`, { headers: authHeaders(), signal: fetchController.signal })
    const data = await res.json()
    users.value = data.users || []
    totalPages.value = data.totalPages || 1
  } catch (e) {
    if (e.name !== 'AbortError') {
      console.error('Error fetching users:', e)
    }
  } finally {
    loading.value = false
  }
}

const handleSort = (field) => {
  if (sortBy.value === field) {
    if (sortOrder.value === 'desc') {
      sortOrder.value = 'asc'
    } else {
      sortBy.value = ''
      sortOrder.value = 'desc'
    }
  } else {
    sortBy.value = field
    sortOrder.value = 'desc'
  }
  currentPage.value = 1
  fetchUsers()
}

const sortIcon = (field) => {
  if (sortBy.value !== field) return ''
  return sortOrder.value === 'desc' ? ' ↓' : ' ↑'
}

const handleSearch = () => {
  currentPage.value = 1
  fetchUsers()
}

const resetSearch = () => {
  search.value = ''
  currentPage.value = 1
  fetchUsers()
}

const switchTab = (role) => {
  activeTab.value = role
  currentPage.value = 1
  fetchUsers()
}

const openDetail = async (user) => {
  try {
    const res = await fetch(`${API_BASE}/admin/users/${user.id}`, { headers: authHeaders() })
    const data = await res.json()
    detailUser.value = data.user
    detailOrders.value = data.orders || []
    detailAddresses.value = data.addresses || []
    showDetailModal.value = true
  } catch (e) {
    toast.error('获取用户详情失败')
  }
}

const openRoleModal = (user) => {
  roleUser.value = user
  selectedRole.value = user.role
  showRoleModal.value = true
}

const openResetModal = (user) => {
  resetUser.value = user
  showResetModal.value = true
}

const updateRole = async () => {
  try {
    const res = await fetch(`${API_BASE}/admin/users/${roleUser.value.id}/role`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ role: selectedRole.value })
    })
    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '更新失败')
    }
    toast.success('用户角色已更新')
    showRoleModal.value = false
    fetchUsers()
  } catch (e) {
    toast.error(e.message)
  }
}

const resetPassword = async () => {
  try {
    const res = await fetch(`${API_BASE}/admin/users/${resetUser.value.id}/reset-password`, {
      method: 'POST',
      headers: authHeaders()
    })
    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '重置失败')
    }
    toast.success('密码已重置为 123456，请通知用户尽快修改')
    showResetModal.value = false
  } catch (e) {
    toast.error(e.message)
  }
}

const formatPrice = (price) => `¥${parseFloat(price || 0).toFixed(2)}`
const formatDate = (date) => new Date(date).toLocaleString('zh-CN')
const formatAddress = (addr) => `${addr.fullName} ${addr.phone} ${addr.state}${addr.city}${addr.addressLine1}${addr.addressLine2 ? ' ' + addr.addressLine2 : ''} ${addr.zipCode}`

const roleBadgeClass = (role) => {
  if (role === 'admin') return 'badge-warning'
  if (role === 'service') return 'badge-info'
  return 'badge-default'
}

const statusLabels = {
  pending: '待处理',
  processing: '处理中',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消'
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchUsers()
  }
}
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchUsers()
  }
}

onMounted(() => fetchUsers())
</script>

<template>
  <div class="admin-users">
    <div class="toolbar">
      <div class="toolbar-left">
        <button class="btn-refresh" :class="{ spinning: loading }" @click="fetchUsers" title="刷新用户列表">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
        </button>
        <div class="search-box">
          <input v-model="search" type="text" placeholder="搜索邮箱或姓名..." @keyup.enter="handleSearch">
          <button class="btn-search" @click="handleSearch">搜索</button>
        </div>
        <button class="btn-reset" @click="resetSearch">重置</button>
      </div>
    </div>

    <div class="tab-bar">
      <button class="tab-item" :class="{ active: activeTab === '' }" @click="switchTab('')">全部</button>
      <button class="tab-item" :class="{ active: activeTab === 'customer' }" @click="switchTab('customer')">客户</button>
      <button class="tab-item" :class="{ active: activeTab === 'admin,service' }" @click="switchTab('admin,service')">管理员</button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th style="width: 60px" class="sortable-th" @click="handleSort('id')">ID{{ sortIcon('id') }}</th>
            <th class="sortable-th" @click="handleSort('email')">邮箱{{ sortIcon('email') }}</th>
            <th class="sortable-th" @click="handleSort('name')">姓名{{ sortIcon('name') }}</th>
            <th style="width: 100px" class="sortable-th" @click="handleSort('role')">角色{{ sortIcon('role') }}</th>
            <th style="width: 100px">订单数</th>
            <th style="width: 120px">消费金额</th>
            <th style="width: 170px" class="sortable-th" @click="handleSort('createdAt')">注册时间{{ sortIcon('createdAt') }}</th>
            <th style="width: 140px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.name || '-' }}</td>
            <td>
              <span class="badge" :class="roleBadgeClass(user.role)">
                {{ roleLabels[user.role] || user.role }}
              </span>
            </td>
            <td>{{ user.orderCount }}</td>
            <td>{{ formatPrice(user.totalSpent) }}</td>
            <td class="time-cell">{{ formatDate(user.createdAt) }}</td>
            <td>
              <button class="btn-view" @click="openDetail(user)">详情</button>
              <button v-if="authStore.isSuperAdmin || user.role !== 'admin'" class="btn-reset-pwd" @click="openResetModal(user)">重置密码</button>
              <button v-if="authStore.isSuperAdmin" class="btn-role" @click="openRoleModal(user)">角色</button>
            </td>
          </tr>
          <tr v-if="users.length === 0">
            <td colspan="8" class="empty-text">暂无用户数据</td>
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
          <h3>用户详情</h3>
          <button class="modal-close" @click="showDetailModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="detail-section">
            <div class="section-title">基本信息</div>
            <div class="detail-grid">
              <div class="detail-item"><span class="label">ID</span><span>{{ detailUser?.id }}</span></div>
              <div class="detail-item"><span class="label">邮箱</span><span>{{ detailUser?.email }}</span></div>
              <div class="detail-item"><span class="label">姓名</span><span>{{ detailUser?.name || '-' }}</span></div>
              <div class="detail-item"><span class="label">角色</span><span class="badge" :class="roleBadgeClass(detailUser?.role)">{{ roleLabels[detailUser?.role] }}</span></div>
              <div class="detail-item"><span class="label">注册时间</span><span>{{ formatDate(detailUser?.createdAt) }}</span></div>
            </div>
          </div>

          <div class="detail-section" v-if="detailOrders.length > 0">
            <div class="section-title">近期订单</div>
            <table class="data-table compact">
              <thead>
                <tr>
                  <th>订单号</th>
                  <th>金额</th>
                  <th>状态</th>
                  <th>时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="order in detailOrders" :key="order.id">
                  <td class="order-no">{{ order.orderNo }}</td>
                  <td>{{ formatPrice(order.total) }}</td>
                  <td><span class="badge" :class="'badge-' + order.status">{{ statusLabels[order.status] }}</span></td>
                  <td class="time-cell">{{ formatDate(order.createdAt) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="no-orders">该用户暂无订单</div>

          <div class="detail-section" v-if="detailAddresses.length > 0">
            <div class="section-title">收货地址</div>
            <div class="address-list">
              <div v-for="addr in detailAddresses" :key="addr.id" class="address-item">
                <span v-if="addr.isDefault" class="badge badge-primary">默认</span>
                <span class="address-text">{{ formatAddress(addr) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showRoleModal" class="modal-overlay" @click.self="showRoleModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>修改用户角色</h3>
          <button class="modal-close" @click="showRoleModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="current-role">
            当前角色：<span class="badge" :class="roleBadgeClass(roleUser?.role)">{{ roleLabels[roleUser?.role] }}</span>
          </div>
          <div class="form-group">
            <label>新角色</label>
            <select v-model="selectedRole">
              <option value="customer">客户</option>
              <option value="service">管理员</option>
              <option value="admin">超级管理员</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showRoleModal = false">取消</button>
          <button class="btn-save" @click="updateRole">确认修改</button>
        </div>
      </div>
    </div>

    <div v-if="showResetModal" class="modal-overlay" @click.self="showResetModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>重置密码</h3>
          <button class="modal-close" @click="showResetModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <p>确认将用户 <strong>{{ resetUser?.email }}</strong> 的密码重置为 <strong>123456</strong>？</p>
          <p class="reset-hint">请通知用户登录后尽快修改密码</p>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showResetModal = false">取消</button>
          <button class="btn-save btn-danger" @click="resetPassword">确认重置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-users {
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

.btn-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  color: #999;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
  flex-shrink: 0;
}

.btn-refresh:hover {
  background: #f0f0f0;
  color: #333;
}

.btn-refresh.spinning svg {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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

.data-table.compact th,
.data-table.compact td {
  padding: 8px 12px;
}

.data-table th {
  background: #f9f9f9;
  font-weight: 600;
  color: #333;
  white-space: nowrap;
}

.sortable-th {
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.sortable-th:hover {
  background: #f0f0f0;
}

.data-table td {
  color: #666;
}

.time-cell {
  font-size: 13px;
  white-space: nowrap;
}

.order-no {
  font-family: monospace;
  font-size: 12px;
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
.btn-role,
.btn-reset-pwd {
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

.btn-role {
  background: #fff8e1;
  color: #f57c00;
}

.btn-role:hover {
  background: #ffecb3;
}

.btn-reset-pwd {
  background: #fce4ec;
  color: #c62828;
}

.btn-reset-pwd:hover {
  background: #f8bbd0;
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

.current-role {
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
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
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

.detail-item .label {
  color: #999;
  margin-right: 8px;
  min-width: 50px;
  display: inline-block;
}

.no-orders {
  text-align: center;
  color: #999;
  padding: 20px;
}

.address-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.address-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f9f9f9;
  border-radius: 4px;
  font-size: 13px;
  color: #666;
}

.address-text {
  line-height: 1.5;
}

.reset-hint {
  color: #f57c00;
  font-size: 13px;
  margin-top: 10px;
}

.btn-danger {
  background: #c62828 !important;
}

.btn-danger:hover {
  background: #b71c1c !important;
}

.tab-bar {
  display: flex;
  border-bottom: 1px solid #eee;
  padding: 0 20px;
  gap: 0;
}

.tab-item {
  padding: 12px 20px;
  border: none;
  background: none;
  color: #999;
  font-size: 14px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: color 0.2s, border-color 0.2s;
}

.tab-item:hover {
  color: #333;
}

.tab-item.active {
  color: #d4a574;
  border-bottom-color: #d4a574;
  font-weight: 600;
}
</style>
