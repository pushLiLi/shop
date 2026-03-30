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
const detailProof = ref(null)

const showQuickReview = ref(false)
const quickReviewOrder = ref(null)
const rejectReason = ref('')
const proofImageError = ref(false)

const selectedIds = ref([])

const showShipModal = ref(false)
const shipOrder = ref(null)
const shipTrackingCompany = ref('')
const shipTrackingNumber = ref('')

const statusLabels = {
  pending: '待处理',
  paid: '已支付',
  processing: '处理中',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消'
}

const statusTransitions = {
  pending: ['cancelled'],
  paid: [],
  processing: ['shipped', 'cancelled'],
  shipped: ['completed'],
  completed: [],
  cancelled: []
}

const proofStatusLabels = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已驳回'
}

const statusBadgeClass = (status) => {
  const map = {
    pending: 'badge-warning',
    paid: 'badge-paid',
    processing: 'badge-info',
    shipped: 'badge-primary',
    completed: 'badge-success',
    cancelled: 'badge-danger'
  }
  return map[status] || 'badge-default'
}

const nextStatuses = computed(() => {
  if (!detailOrder.value) return []
  return statusTransitions[detailOrder.value.status] || []
})

const pendingProofCount = computed(() =>
  orders.value.filter(o => o.paymentProof?.status === 'pending').length
)

const allSelected = computed(() =>
  orders.value.length > 0 && orders.value.every(o =>
    o.paymentProof?.status === 'pending' ? selectedIds.value.includes(o.id) : true
  )
)

const hasSelected = computed(() => selectedIds.value.length > 0)

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
    selectedIds.value = []
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

const filterPendingProofs = () => {
  filterStatus.value = 'pending'
  currentPage.value = 1
  loading.value = true
  fetch(`${API_BASE}/admin/orders?page=${currentPage.value}&limit=${limit}&proof_status=pending`, { headers: authHeaders() })
    .then(r => r.json())
    .then(data => {
      orders.value = data.orders || []
      totalPages.value = data.totalPages || 1
      selectedIds.value = []
    })
    .catch(() => toast.error('获取订单失败'))
    .finally(() => { loading.value = false })
}

const openDetail = async (order) => {
  try {
    const res = await fetch(`${API_BASE}/admin/orders/${order.id}`, { headers: authHeaders() })
    const data = await res.json()
    detailOrder.value = data.order
    detailUser.value = data.user
    detailProof.value = data.paymentProof || order.paymentProof || null

    showDetailModal.value = true
  } catch (e) {
    toast.error('获取订单详情失败')
  }
}

const quickStatusChange = async (order, newStatus) => {
  if (!newStatus) return

  if (newStatus === 'shipped') {
    shipOrder.value = order
    shipTrackingCompany.value = ''
    shipTrackingNumber.value = ''
    showShipModal.value = true
    return
  }

  if (newStatus === 'cancelled' && !confirm('确定要取消该订单吗？')) return

  try {
    const res = await fetch(`${API_BASE}/admin/orders/${order.id}/status`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ status: newStatus })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '更新失败')
    toast.success('订单状态已更新')
    fetchOrders()
  } catch (e) {
    toast.error(e.message)
  }
}

const openQuickReview = (order) => {
  quickReviewOrder.value = order
  rejectReason.value = ''
  showQuickReview.value = true
}

const submitReview = async (action) => {
  const order = quickReviewOrder.value
  if (!order?.paymentProof) return

  try {
    const body = { action }
    if (action === 'reject' && rejectReason.value) {
      body.rejectReason = rejectReason.value
    }
    const res = await fetch(`${API_BASE}/admin/payment-proofs/${order.paymentProof.id}/review`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(body)
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '审核失败')
    toast.success(action === 'approve' ? '已通过审核' : '已驳回凭证')
    showQuickReview.value = false
    fetchOrders()
  } catch (e) {
    toast.error(e.message)
  }
}

const reviewProofFromDetail = async (action) => {
  if (!detailProof.value) return

  try {
    const body = { action }
    if (action === 'reject' && rejectReason.value) {
      body.rejectReason = rejectReason.value
    }
    const res = await fetch(`${API_BASE}/admin/payment-proofs/${detailProof.value.id}/review`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(body)
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '审核失败')
    toast.success(action === 'approve' ? '已通过审核' : '已驳回凭证')
    showDetailModal.value = false
    fetchOrders()
  } catch (e) {
    toast.error(e.message)
  }
}

const toggleSelect = (id) => {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) selectedIds.value.splice(idx, 1)
  else selectedIds.value.push(id)
}

const toggleSelectAll = () => {
  const pendingOrders = orders.value.filter(o => o.paymentProof?.status === 'pending')
  if (allSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = pendingOrders.map(o => o.id)
  }
}

const batchReview = async (action) => {
  if (selectedIds.value.length === 0) return
  if (action === 'reject' && !confirm(`确定要驳回 ${selectedIds.value.length} 个凭证吗？`)) return

  try {
    const body = { ids: selectedIds.value, action }
    if (action === 'reject' && rejectReason.value) {
      body.rejectReason = rejectReason.value
    }
    const res = await fetch(`${API_BASE}/admin/payment-proofs/batch-review`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(body)
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '批量审核失败')
    toast.success(`已${action === 'approve' ? '通过' : '驳回'} ${data.reviewed} 个凭证`)
    selectedIds.value = []
    fetchOrders()
  } catch (e) {
    toast.error(e.message)
  }
}

const formatPrice = (price) => `¥${parseFloat(price).toFixed(2)}`
const formatDate = (date) => new Date(date).toLocaleString('zh-CN')

const submitShip = async () => {
  const order = shipOrder.value
  if (!order) return
  if (!shipTrackingCompany.value.trim() || !shipTrackingNumber.value.trim()) {
    toast.error('请填写物流平台和快递单号')
    return
  }
  try {
    const res = await fetch(`${API_BASE}/admin/orders/${order.id}/status`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        status: 'shipped',
        trackingCompany: shipTrackingCompany.value.trim(),
        trackingNumber: shipTrackingNumber.value.trim()
      })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '发货失败')
    toast.success('已发货')
    showShipModal.value = false
    fetchOrders()
  } catch (e) {
    toast.error(e.message)
  }
}

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
        <button class="btn-refresh" :class="{ spinning: loading }" @click="fetchOrders" title="刷新订单列表">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
        </button>
        <div class="search-box">
          <input v-model="search" type="text" placeholder="搜索订单号..." @keyup.enter="handleSearch">
          <button class="btn-search" @click="handleSearch">搜索</button>
        </div>
        <select v-model="filterStatus" @change="handleSearch">
          <option value="">全部状态</option>
          <option v-for="(label, key) in statusLabels" :key="key" :value="key">{{ label }}</option>
        </select>
        <button class="btn-proof-filter" @click="filterPendingProofs">
          待审核凭证<span v-if="pendingProofCount"> ({{ pendingProofCount }})</span>
        </button>
        <button class="btn-reset" @click="resetFilters">重置</button>
      </div>
    </div>

    <div v-if="hasSelected" class="batch-bar">
      <span class="batch-info">已选 {{ selectedIds.length }} 个待审核凭证</span>
      <div class="batch-actions">
        <button class="btn-batch-approve" @click="batchReview('approve')">批量通过</button>
        <button class="btn-batch-reject" @click="batchReview('reject')">批量驳回</button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th style="width: 40px">
              <input type="checkbox" :checked="allSelected" @change="toggleSelectAll" title="全选待审核" />
            </th>
            <th style="width: 60px">ID</th>
            <th style="width: 170px">订单号</th>
            <th>用户</th>
            <th style="width: 100px">总金额</th>
            <th style="width: 100px">状态</th>
            <th style="width: 100px">付款凭证</th>
            <th style="width: 160px">下单时间</th>
            <th style="width: 180px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id"
              :class="{ 'has-pending-proof': order.paymentProof?.status === 'pending' }">
            <td>
              <input v-if="order.paymentProof?.status === 'pending'"
                     type="checkbox" :checked="selectedIds.includes(order.id)"
                     @change="toggleSelect(order.id)" />
            </td>
            <td>{{ order.id }}</td>
            <td class="order-no">{{ order.orderNo }}</td>
            <td>{{ order.user?.name || order.user?.email || '-' }}</td>
            <td>{{ formatPrice(order.total) }}</td>
            <td>
              <span class="badge" :class="statusBadgeClass(order.status)">
                {{ statusLabels[order.status] || order.status }}
              </span>
            </td>
            <td>
              <span v-if="!order.paymentProof" class="proof-badge proof-none">未提交</span>
              <span v-else class="proof-badge"
                    :class="{
                      'proof-pending': order.paymentProof.status === 'pending',
                      'proof-approved': order.paymentProof.status === 'approved',
                      'proof-rejected': order.paymentProof.status === 'rejected'
                    }"
                    @click="order.paymentProof.status !== 'approved' && openQuickReview(order)">
                {{ proofStatusLabels[order.paymentProof.status] }}
              </span>
            </td>
            <td class="time-cell">{{ formatDate(order.createdAt) }}</td>
            <td>
              <button class="btn-view" @click="openDetail(order)">详情</button>
              <button v-if="order.paymentProof?.status === 'pending'"
                      class="btn-review-proof" @click="openQuickReview(order)">审核</button>
              <select v-if="statusTransitions[order.status]?.length"
                      class="inline-status-select"
                      @change="quickStatusChange(order, $event.target.value)"
                      :value="''">
                <option value="" disabled>状态 ▾</option>
                <option v-for="s in statusTransitions[order.status]" :key="s" :value="s">
                  {{ statusLabels[s] }}
                </option>
              </select>
            </td>
          </tr>
          <tr v-if="orders.length === 0">
            <td colspan="9" class="empty-text">暂无订单数据</td>
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

          <div class="detail-section" v-if="detailOrder?.trackingCompany || detailOrder?.trackingNumber">
            <div class="section-title">物流信息</div>
            <div class="detail-grid">
              <div class="detail-item"><span class="label">物流平台</span><span>{{ detailOrder.trackingCompany }}</span></div>
              <div class="detail-item"><span class="label">快递单号</span><span class="mono">{{ detailOrder.trackingNumber }}</span></div>
            </div>
          </div>

          <div class="detail-section">
            <div class="section-title">商品列表</div>
            <div class="items-list">
              <div v-for="item in detailOrder?.items" :key="item.id" class="order-item">
                <img v-if="item.product?.thumbnailUrl || item.product?.imageUrl" :src="item.product.imageUrl" class="item-thumb" loading="lazy" />
                <div class="item-total">¥{{ (item.price * item.quantity).toFixed(2) }}</div>
              </div>
            </div>
          </div>

          <div v-if="detailProof" class="detail-section">
            <div class="section-title">付款凭证</div>
            <div class="proof-detail-grid">
              <div class="detail-item"><span class="label">付款方式</span><span>{{ typeof detailProof.paymentMethod === 'object' ? detailProof.paymentMethod?.name : detailProof.paymentMethod || '-' }}</span></div>
              <div class="detail-item">
                <span class="label">凭证状态</span>
                <span :class="['badge', detailProof.status === 'pending' ? 'badge-warning' : detailProof.status === 'approved' ? 'badge-success' : 'badge-danger']">
                  {{ detailProof.status === 'pending' ? '待审核' : detailProof.status === 'approved' ? '已通过' : '已驳回' }}
                </span>
              </div>
              <div v-if="detailProof.imageUrl" class="detail-item full">
                <span class="label">付款截图</span>
                <div class="proof-image-wrapper">
                  <a :href="detailProof.imageUrl" target="_blank">
                    <img :src="detailProof.imageUrl" class="proof-image" @error="$event.target.style.display='none';$event.target.parentElement.nextElementSibling && ($event.target.parentElement.nextElementSibling.style.display='block')" />
                  </a>
                  <div class="image-load-error" style="display:none">图片加载失败，<a :href="detailProof.imageUrl" target="_blank">点击在新窗口打开</a></div>
                </div>
              </div>
              <div v-if="detailProof.rejectReason" class="detail-item full">
                <span class="label">驳回原因</span><span style="color:#c62828">{{ detailProof.rejectReason }}</span>
              </div>
            </div>
            <div v-if="detailProof.status === 'pending'" class="proof-actions">
              <textarea v-model="rejectReason" placeholder="驳回原因（可选）" class="reject-textarea"></textarea>
              <div class="proof-actions-buttons">
                <button class="btn-reject" @click="reviewProofFromDetail('reject')">驳回</button>
                <button class="btn-approve" @click="reviewProofFromDetail('approve')">通过</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showQuickReview" class="modal-overlay" @click.self="showQuickReview = false">
      <div class="modal review-modal">
        <div class="modal-header">
          <h3>审核付款凭证</h3>
          <button class="modal-close" @click="showQuickReview = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="review-order-info">
            <div class="review-info-item">
              <span class="label">订单号</span>
              <span class="order-no">{{ quickReviewOrder?.orderNo }}</span>
            </div>
            <div class="review-info-item">
              <span class="label">用户</span>
              <span>{{ quickReviewOrder?.user?.name || quickReviewOrder?.user?.email || '-' }}</span>
            </div>
            <div class="review-info-item">
              <span class="label">金额</span>
              <span class="price">{{ formatPrice(quickReviewOrder?.total) }}</span>
            </div>
            <div class="review-info-item">
              <span class="label">付款方式</span>
              <span>{{ quickReviewOrder?.paymentProof?.paymentMethod || '-' }}</span>
            </div>
          </div>

          <div v-if="quickReviewOrder?.paymentProof?.imageUrl" class="review-proof-image">
            <a :href="quickReviewOrder.paymentProof.imageUrl" target="_blank">
              <img :src="quickReviewOrder.paymentProof.imageUrl" @error="$event.target.style.display='none';$event.target.nextElementSibling && ($event.target.nextElementSibling.style.display='block')" />
              <div class="image-load-error" style="display:none">图片加载失败，<a :href="quickReviewOrder.paymentProof.imageUrl" target="_blank">点击在新窗口打开</a></div>
            </a>
          </div>

          <div v-if="quickReviewOrder?.paymentProof?.rejectReason" class="review-reject-reason">
            <span class="label">上次驳回原因：</span>{{ quickReviewOrder.paymentProof.rejectReason }}
          </div>

          <div class="form-group">
            <label>驳回原因（驳回时填写）</label>
            <textarea v-model="rejectReason" placeholder="可选" class="reject-textarea"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showQuickReview = false">取消</button>
          <button class="btn-reject" @click="submitReview('reject')">驳回</button>
          <button class="btn-save" @click="submitReview('approve')">通过</button>
        </div>
      </div>
    </div>

    <div v-if="showShipModal" class="modal-overlay" @click.self="showShipModal = false">
      <div class="modal ship-modal">
        <div class="modal-header">
          <h3>发货 - {{ shipOrder?.orderNo }}</h3>
          <button class="modal-close" @click="showShipModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>物流平台 <span style="color:#c62828">*</span></label>
            <input v-model="shipTrackingCompany" type="text" placeholder="如：顺丰、中通、圆通" class="form-input" />
          </div>
          <div class="form-group">
            <label>快递单号 <span style="color:#c62828">*</span></label>
            <input v-model="shipTrackingNumber" type="text" placeholder="请输入快递单号" class="form-input" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showShipModal = false">取消</button>
          <button class="btn-save" @click="submitShip">确认发货</button>
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

.btn-proof-filter {
  padding: 8px 16px;
  border: 2px solid #ff9800;
  border-radius: 4px;
  background: #fff8e1;
  color: #e65100;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-proof-filter:hover {
  background: #ffecb3;
}

.batch-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background: #fff3e0;
  border-bottom: 1px solid #ffe0b2;
}

.batch-info {
  font-size: 14px;
  color: #e65100;
  font-weight: 500;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.btn-batch-approve,
.btn-batch-reject {
  padding: 6px 14px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.btn-batch-approve {
  background: #d4a574;
  color: #fff;
}

.btn-batch-approve:hover {
  background: #c49464;
}

.btn-batch-reject {
  background: #ffebee;
  color: #c62828;
}

.btn-batch-reject:hover {
  background: #ffcdd2;
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

.data-table tr.has-pending-proof {
  border-left: 3px solid #ff9800;
  background: #fffdf7;
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
.badge-paid { background: #e8f5e9; color: #388e3c; }
.badge-default { background: #f5f5f5; color: #999; }

.proof-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.proof-none {
  background: #f5f5f5;
  color: #bbb;
}

.proof-pending {
  background: #fff3e0;
  color: #e65100;
  cursor: pointer;
  font-weight: 600;
  animation: pulse-glow 2s ease-in-out infinite;
}

.proof-pending:hover {
  background: #ffe0b2;
}

.proof-approved {
  background: #e8f5e9;
  color: #2e7d32;
}

.proof-rejected {
  background: #ffebee;
  color: #c62828;
  cursor: pointer;
}

.proof-rejected:hover {
  background: #ffcdd2;
}

@keyframes pulse-glow {
  0%, 100% { box-shadow: 0 0 0 0 rgba(255, 152, 0, 0); }
  50% { box-shadow: 0 0 0 3px rgba(255, 152, 0, 0.3); }
}

.btn-view,
.btn-review-proof {
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

.btn-review-proof {
  background: #fff3e0;
  color: #e65100;
  font-weight: 600;
  animation: pulse-glow 2s ease-in-out infinite;
}

.btn-review-proof:hover {
  background: #ffe0b2;
}

.inline-status-select {
  padding: 5px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  font-size: 12px;
  cursor: pointer;
  color: #666;
  transition: border-color 0.2s;
}

.inline-status-select:hover {
  border-color: #d4a574;
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

.review-modal {
  max-width: 600px;
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

.btn-reject {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  background: #ffebee;
  color: #c62828;
  transition: all 0.2s;
}

.btn-reject:hover {
  background: #ffcdd2;
}

.btn-approve {
  padding: 8px 16px;
  background: #d4a574;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.btn-approve:hover {
  background: #c49464;
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

.reject-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  min-height: 70px;
  resize: vertical;
  box-sizing: border-box;
  font-size: 14px;
}

.reject-textarea:focus {
  outline: none;
  border-color: #d4a574;
}

.review-order-info {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 16px;
  padding: 12px;
  background: #fafafa;
  border-radius: 6px;
}

.review-info-item {
  font-size: 14px;
  color: #666;
}

.review-info-item .label {
  color: #999;
  margin-right: 8px;
}

.review-proof-image {
  margin-bottom: 16px;
  text-align: center;
}

.review-proof-image img {
  max-width: 100%;
  max-height: 400px;
  border-radius: 6px;
  border: 1px solid #eee;
  cursor: pointer;
}

.review-reject-reason {
  font-size: 13px;
  color: #c62828;
  padding: 8px 12px;
  background: #fff5f5;
  border-radius: 4px;
  margin-bottom: 12px;
}

.review-reject-reason .label {
  color: #999;
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

.item-total {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.proof-detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.proof-image-wrapper {
  margin-top: 8px;
}

.image-load-error {
  padding: 12px 16px;
  background: #fff3e0;
  border: 1px solid #ffe0b2;
  border-radius: 6px;
  color: #e65100;
  font-size: 13px;
}

.proof-image {
  max-width: 300px;
  max-height: 300px;
  border-radius: 6px;
  border: 1px solid #eee;
  cursor: pointer;
}

.proof-actions {
  margin-top: 15px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.proof-actions-buttons {
  display: flex;
  gap: 10px;
}

.ship-modal {
  max-width: 480px;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  transition: border-color 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #d4a574;
}

.mono {
  font-family: monospace;
}
</style>
