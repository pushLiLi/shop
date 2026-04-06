<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToastStore } from '../stores/toast'
import { useImageCompress } from '../composables/useImageCompress'
import { formatPriceByCurrency } from '../composables/useFormatPrice'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const { compress } = useImageCompress()

const order = ref(null)
const paymentProof = ref(null)
const loading = ref(true)
const reuploadFile = ref(null)
const reuploadPreview = ref(null)
const reuploading = ref(false)

const statusMap = {
  pending: '待处理',
  paid: '已支付',
  processing: '处理中',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消'
}

const statusClass = {
  pending: 'status-pending',
  paid: 'status-paid',
  processing: 'status-processing',
  shipped: 'status-shipped',
  completed: 'status-completed',
  cancelled: 'status-cancelled'
}

const proofStatusMap = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已驳回'
}

const authHeaders = () => ({
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

async function fetchOrder() {
  try {
    loading.value = true
    const id = route.params.id
    const res = await fetch(`/api/orders/${id}`, { headers: authHeaders() })
    if (!res.ok) {
      toast.error('订单不存在')
      router.push('/orders')
      return
    }
    const data = await res.json()
    order.value = data.order
    paymentProof.value = data.paymentProof || null
  } catch (e) {
    toast.error('加载失败')
    router.push('/orders')
  } finally {
    loading.value = false
  }
}

async function handleReupload() {
  if (!reuploadFile.value || !order.value) return
  try {
    reuploading.value = true
    const compressed = await compress(reuploadFile.value, { maxWidth: 1920, maxHeight: 1920, quality: 0.9 })
    const formData = new FormData()
    formData.append('file', compressed, 'proof.jpg')
    formData.append('paymentMethodId', paymentProof.value?.paymentMethodId || '')

    const res = await fetch(`/api/orders/${order.value.id}/payment-proof`, {
      method: 'POST',
      headers: authHeaders(),
      body: formData
    })
    const data = await res.json()
    if (res.ok) {
      toast.success('付款截图已重新上传')
      paymentProof.value = data.paymentProof
      reuploadFile.value = null
      reuploadPreview.value = null
    } else {
      toast.error(data.error || '上传失败')
    }
  } catch (e) {
    toast.error('上传失败')
  } finally {
    reuploading.value = false
  }
}

function onReuploadFileChange(e) {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 10 * 1024 * 1024) {
    toast.error('图片大小不能超过 10MB')
    return
  }
  reuploadFile.value = file
  const reader = new FileReader()
  reader.onload = (ev) => {
    reuploadPreview.value = ev.target.result
  }
  reader.readAsDataURL(file)
}

function removeReuploadFile() {
  reuploadFile.value = null
  reuploadPreview.value = null
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleString('zh-CN')
}


onMounted(() => fetchOrder())
</script>

<template>
  <div class="order-detail-page">
    <div class="container">
      <div class="page-header">
        <button class="btn-back" @click="router.push('/orders')">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"></polyline></svg>
          返回订单列表
        </button>
        <h1 class="page-title">订单详情</h1>
      </div>

      <div v-if="loading" class="loading">加载中...</div>

      <template v-else-if="order">
        <div class="section order-info-section">
          <div class="section-header">
            <h2>订单信息</h2>
            <span class="order-status" :class="statusClass[order.status]">
              {{ statusMap[order.status] || order.status }}
            </span>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">订单号</span>
              <span class="value mono">{{ order.orderNo }}</span>
            </div>
            <div class="info-item">
              <span class="label">下单时间</span>
              <span class="value">{{ formatDate(order.createdAt) }}</span>
            </div>
            <div class="info-item">
              <span class="label">备注</span>
              <span class="value">{{ order.remark || '无' }}</span>
            </div>
          </div>
        </div>

        <div v-if="order.address" class="section">
          <h2 class="section-title">收货地址</h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">收件人</span>
              <span class="value">{{ order.address.fullName }}</span>
            </div>
            <div class="info-item">
              <span class="label">电话</span>
              <span class="value">{{ order.address.phone }}</span>
            </div>
            <div class="info-item full">
              <span class="label">地址</span>
              <span class="value">{{ order.address.addressLine1 }} {{ order.address.addressLine2 }} {{ order.address.city }} {{ order.address.state }} {{ order.address.zipCode }}</span>
            </div>
          </div>
        </div>

        <div v-if="order.trackingCompany || order.trackingNumber" class="section">
          <h2 class="section-title">物流信息</h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">物流平台</span>
              <span class="value">{{ order.trackingCompany }}</span>
            </div>
            <div class="info-item">
              <span class="label">快递单号</span>
              <span class="value mono">{{ order.trackingNumber }}</span>
            </div>
          </div>
        </div>

        <div class="section">
          <h2 class="section-title">商品列表</h2>
          <div class="items-list">
            <div v-for="item in order.items" :key="item.id" class="item-row">
              <img v-if="item.product?.imageUrl || item.product?.thumbnailUrl"
                   :src="item.product.imageUrl || item.product.thumbnailUrl"
                   class="item-image" loading="lazy" />
              <div v-else class="item-image-placeholder"></div>
              <div class="item-info">
                <span class="item-name">{{ item.product?.name || '商品' }}</span>
                <span class="item-qty">x{{ item.quantity }}</span>
              </div>
              <span class="item-price">{{ formatPriceByCurrency(item.price * item.quantity, item.currency) }}</span>
            </div>
          </div>
          <div class="order-total">
            <span>合计</span>
            <span class="total-price">{{ formatPriceByCurrency(order.total, 'CNY') }}</span>
          </div>
        </div>

        <div v-if="paymentProof" class="section">
          <h2 class="section-title">付款凭证</h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">凭证状态</span>
              <span :class="['proof-badge', `proof-${paymentProof.status}`]">
                {{ proofStatusMap[paymentProof.status] }}
              </span>
            </div>
            <div class="info-item">
              <span class="label">付款方式</span>
              <span class="value">{{ typeof paymentProof.paymentMethod === 'object' ? paymentProof.paymentMethod?.name : paymentProof.paymentMethod || '-' }}</span>
            </div>
            <div v-if="paymentProof.imageUrl" class="info-item full proof-image-item">
              <span class="label">付款截图</span>
              <div class="proof-image-wrapper">
                <a :href="paymentProof.imageUrl" target="_blank">
                  <img :src="paymentProof.imageUrl" class="proof-image"
                       @error="$event.target.style.display='none';$event.target.nextElementSibling && ($event.target.nextElementSibling.style.display='block')" />
                  <div class="image-load-error" style="display:none">图片加载失败，<a :href="paymentProof.imageUrl" target="_blank">点击在新窗口打开</a></div>
                </a>
              </div>
            </div>
            <div v-if="paymentProof.rejectReason" class="info-item full">
              <span class="label">驳回原因</span>
              <span class="value reject-reason">{{ paymentProof.rejectReason }}</span>
            </div>
          </div>

          <div v-if="paymentProof.status === 'rejected'" class="reupload-area">
            <p class="reupload-hint">凭证被驳回，请重新上传付款截图</p>
            <div v-if="!reuploadPreview" class="reupload-drop" @click="$refs.reuploadInput.click()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="upload-icon">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
              <span>点击选择付款截图</span>
              <span class="upload-hint">支持 JPG、PNG、GIF、WebP，最大 10MB</span>
            </div>
            <div v-else class="reupload-preview-wrapper">
              <img :src="reuploadPreview" class="reupload-preview-img" />
              <button class="remove-reupload-btn" @click="removeReuploadFile">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
            <input ref="reuploadInput" type="file" accept="image/*" @change="onReuploadFileChange" hidden />
            <button
              v-if="reuploadPreview"
              class="btn-reupload"
              :disabled="reuploading"
              @click="handleReupload"
            >
              {{ reuploading ? '上传中...' : '提交重新上传' }}
            </button>
          </div>
        </div>

        <div v-else class="section">
          <h2 class="section-title">付款凭证</h2>
          <p class="no-proof">暂未上传付款凭证</p>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.order-detail-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 30px;
}

.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #1a1a1a;
  border: 1px solid #333;
  color: #ccc;
  padding: 8px 14px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-back:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.page-title {
  color: #d4a574;
  font-size: 24px;
  margin: 0;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #888;
}

.section {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h2 {
  color: #fff;
  font-size: 16px;
  margin: 0;
}

.section-title {
  color: #fff;
  font-size: 16px;
  margin: 0 0 16px;
  padding-bottom: 10px;
  border-bottom: 1px solid #2a2a2a;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.info-item {
  font-size: 14px;
  color: #ccc;
}

.info-item.full {
  grid-column: 1 / -1;
}

.info-item .label {
  color: #888;
  margin-right: 8px;
}

.info-item .value {
  color: #ccc;
}

.mono {
  font-family: monospace;
  font-size: 13px;
}

.order-status {
  padding: 5px 14px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
}

.status-pending { background: #f0ad4e; color: #1a1a1a; }
.status-paid { background: #5cb85c; color: #1a1a1a; }
.status-processing { background: #6c757d; color: #fff; }
.status-shipped { background: #5bc0de; color: #1a1a1a; }
.status-completed { background: #d4a574; color: #1a1a1a; }
.status-cancelled { background: #d9534f; color: #fff; }

.items-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px;
  background: #222;
  border-radius: 6px;
}

.item-image {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.item-image-placeholder {
  width: 56px;
  height: 56px;
  background: #2a2a2a;
  border-radius: 4px;
  flex-shrink: 0;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-name {
  color: #ccc;
  font-size: 14px;
}

.item-qty {
  color: #888;
  font-size: 13px;
}

.item-price {
  color: #d4a574;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.order-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #2a2a2a;
  font-size: 16px;
  color: #ccc;
}

.total-price {
  color: #d4a574;
  font-size: 20px;
  font-weight: 700;
}

.proof-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.proof-pending { background: #f0ad4e; color: #1a1a1a; }
.proof-approved { background: #5cb85c; color: #1a1a1a; }
.proof-rejected { background: #d9534f; color: #fff; }

.proof-image-item {
  flex-direction: column;
  align-items: flex-start;
}

.proof-image-wrapper {
  margin-top: 8px;
}

.proof-image {
  max-width: 100%;
  max-height: 400px;
  border-radius: 6px;
  border: 1px solid #333;
  cursor: pointer;
}

.image-load-error {
  padding: 12px 16px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  color: #f0ad4e;
  font-size: 13px;
}

.image-load-error a {
  color: #d4a574;
}

.reject-reason {
  color: #d9534f !important;
}

.no-proof {
  color: #888;
  font-size: 14px;
}

.reupload-area {
  margin-top: 16px;
  padding: 16px;
  background: #222;
  border-radius: 6px;
  border: 1px solid #d9534f33;
}

.reupload-hint {
  color: #f0ad4e;
  font-size: 14px;
  margin: 0 0 12px;
}

.reupload-drop {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 30px;
  border: 2px dashed #444;
  border-radius: 8px;
  cursor: pointer;
  color: #888;
  transition: all 0.3s;
}

.reupload-drop:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.upload-icon {
  width: 36px;
  height: 36px;
}

.upload-hint {
  font-size: 12px;
  color: #666;
}

.reupload-preview-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 12px;
}

.reupload-preview-img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
  border: 1px solid #444;
}

.remove-reupload-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.7);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-reupload-btn svg {
  width: 16px;
  height: 16px;
  stroke: #fff;
}

.btn-reupload {
  width: 100%;
  background: #d4a574;
  color: #1a1a1a;
  border: none;
  padding: 12px;
  border-radius: 4px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
}

.btn-reupload:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-reupload:hover:not(:disabled) {
  background: #e5b584;
}

@media (max-width: 768px) {
  .order-detail-page {
    padding: 20px 0 40px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 20px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .section {
    padding: 15px;
  }

  .item-row {
    flex-wrap: wrap;
  }

  .mono {
    word-break: break-all;
  }
}

@media (max-width: 576px) {
  .order-detail-page {
    padding: 15px 0 30px;
  }

  .page-title {
    font-size: 20px;
  }

  .btn-back {
    padding: 6px 10px;
    font-size: 13px;
  }

  .proof-image {
    max-height: 250px;
  }

  .total-price {
    font-size: 18px;
  }
}
</style>
