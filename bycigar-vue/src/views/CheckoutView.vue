<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { getStateName } from '../utils/states'
import { useImageCompress } from '../composables/useImageCompress'

const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const toast = useToastStore()
const { compress } = useImageCompress()

const items = computed(() => cartStore.items)
const total = computed(() => cartStore.total)
const totalQuantity = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
const loading = ref(false)
const error = ref(null)
const addresses = ref([])
const selectedAddressId = ref(null)
const remark = ref('')
const paymentMethods = ref([])
const selectedPaymentMethodId = ref(null)
const proofFile = ref(null)
const proofPreview = ref(null)

const selectedPaymentMethod = computed(() =>
  paymentMethods.value.find(m => m.id === selectedPaymentMethodId.value)
)

onMounted(async () => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }

  const isValid = await authStore.validateToken()
  if (!isValid) {
    router.push('/login')
    return
  }

  fetchAddresses()
  fetchPaymentMethods()
})

async function fetchAddresses() {
  try {
    const res = await fetch('/api/addresses', {
      headers: authStore.getAuthHeaders()
    })
    if (res.ok) {
      const data = await res.json()
      addresses.value = data.addresses || []
      const defaultAddr = addresses.value.find(a => a.isDefault)
      if (defaultAddr) {
        selectedAddressId.value = defaultAddr.id
      } else if (addresses.value.length > 0) {
        selectedAddressId.value = addresses.value[0].id
      }
    }
  } catch (e) {
    console.error('获取地址失败:', e)
  }
}

async function fetchPaymentMethods() {
  try {
    const res = await fetch('/api/payment-methods')
    if (res.ok) {
      const data = await res.json()
      paymentMethods.value = data.paymentMethods || []
      if (paymentMethods.value.length > 0) {
        selectedPaymentMethodId.value = paymentMethods.value[0].id
      }
    }
  } catch (e) {
    console.error('获取付款方式失败:', e)
  }
}

function handleProofUpload(e) {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 10 * 1024 * 1024) {
    toast.error('图片大小不能超过 10MB')
    return
  }
  proofFile.value = file
  const reader = new FileReader()
  reader.onload = (ev) => {
    proofPreview.value = ev.target.result
  }
  reader.readAsDataURL(file)
}

function removeProof() {
  proofFile.value = null
  proofPreview.value = null
}

async function createOrder() {
  if (items.value.length === 0) {
    error.value = '购物车是空的'
    return
  }
  if (!selectedAddressId.value) {
    error.value = '请选择收货地址'
    return
  }
  if (!selectedPaymentMethodId.value) {
    error.value = '请选择付款方式'
    return
  }
  if (!proofFile.value) {
    error.value = '请上传付款截图'
    return
  }
  try {
    loading.value = true
    error.value = null

    const orderData = {
      addressId: selectedAddressId.value,
      remark: remark.value
    }
    const res = await fetch('/api/orders', {
      method: 'POST',
      headers: {
        ...authStore.getAuthHeaders(),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(orderData)
    })
    const orderResult = await res.json()
    if (!res.ok) {
      error.value = orderResult.error || '创建订单失败'
      return
    }

    const orderId = orderResult.orderId
    if (!orderId) {
      error.value = '订单创建异常，请查看订单列表'
      cartStore.clear()
      router.push('/orders')
      return
    }

    const formData = new FormData()
    const compressedProof = await compress(proofFile.value, { maxWidth: 1920, maxHeight: 1920, quality: 0.8 })
    formData.append('file', compressedProof, 'proof.jpg')
    formData.append('paymentMethodId', selectedPaymentMethodId.value)

    const proofRes = await fetch(`/api/orders/${orderId}/payment-proof`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: formData
    })
    if (!proofRes.ok) {
      const proofData = await proofRes.json()
      toast.warning('订单已创建，但付款凭证上传失败：' + (proofData.error || '未知错误'))
    }

    cartStore.clear()
    toast.success('订单提交成功！')
    router.push('/orders')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function formatPrice(price) {
  return '$' + Number(price).toFixed(2)
}
</script>

<template>
  <div class="checkout-page">
    <div class="container">
      <h1 class="page-title">结算</h1>
      <div v-if="error" class="error-message">{{ error }}</div>
      <div class="checkout-content">
        <div class="order-items">
          <h2 class="section-title">订单商品</h2>
          <div class="items-list">
            <div v-for="item in items" :key="item.productId" class="order-item">
              <span class="item-name">{{ item.product?.name }}</span>
              <span class="item-quantity">x{{ item.quantity }}</span>
              <span class="item-price">{{ formatPrice((item.product?.price || 0) * item.quantity) }}</span>
            </div>
          </div>
        </div>
        <div class="shipping-info">
          <h2 class="section-title">收货地址</h2>

          <div v-if="addresses.length === 0" class="no-address">
            <p>您还没有保存收货地址</p>
            <p class="hint">请先添加收货地址后再进行结算</p>
            <router-link to="/profile" class="btn-add-address">去添加地址</router-link>
          </div>

          <div v-else class="address-list">
            <div
              v-for="addr in addresses"
              :key="addr.id"
              :class="['address-option', { selected: selectedAddressId === addr.id, default: addr.isDefault }]"
              @click="selectedAddressId = addr.id"
            >
              <div class="address-radio">
                <span :class="['radio', { checked: selectedAddressId === addr.id }]"></span>
              </div>
              <div class="address-details">
                <div class="address-name">
                  {{ addr.fullName }}
                  <span v-if="addr.isDefault" class="default-tag">默认</span>
                </div>
                <div class="address-text">
                  <svg class="addr-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                  <span class="addr-label">详细地址：</span>{{ addr.addressLine1 }}{{ addr.addressLine2 ? ', ' + addr.addressLine2 : '' }}
                </div>
                <div class="address-text">
                  <svg class="addr-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 21h18M5 21V7l8-4v18M19 21V11l-6-4"/></svg>
                  <span class="addr-label">城市邮编：</span>{{ addr.city }}, {{ getStateName(addr.state) }} {{ addr.zipCode }}
                </div>
                <div class="address-phone">
                  <svg class="addr-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                  <span class="addr-label">联系电话：</span>{{ addr.phone }}
                </div>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>备注</label>
            <textarea v-model="remark" placeholder="可选"></textarea>
          </div>
        </div>

        <div class="payment-section">
          <h2 class="section-title">付款方式</h2>

          <div v-if="paymentMethods.length === 0" class="no-payment">
            <p>暂无可用的付款方式，请联系客服</p>
          </div>

          <div v-else class="payment-method-list">
            <div
              v-for="method in paymentMethods"
              :key="method.id"
              :class="['payment-method-option', { selected: selectedPaymentMethodId === method.id }]"
              @click="selectedPaymentMethodId = method.id"
            >
              <div class="payment-radio">
                <span :class="['radio', { checked: selectedPaymentMethodId === method.id }]"></span>
              </div>
              <span class="payment-name">{{ method.name }}</span>
            </div>
          </div>

          <div v-if="selectedPaymentMethod" class="payment-detail">
            <div v-if="selectedPaymentMethod.qrCodeUrl" class="qrcode-wrapper">
              <img :src="selectedPaymentMethod.qrCodeUrl" :alt="selectedPaymentMethod.name" class="qrcode-img" />
            </div>
            <div v-if="selectedPaymentMethod.paymentUrl" class="payment-link-wrapper">
              <a :href="selectedPaymentMethod.paymentUrl" target="_blank" class="payment-link-btn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="paypal-icon"><path d="M21 8V12.5C21 16.5 18.5 20 14 20H10C8.9 20 8 19.1 8 18V14C8 13.4 8.4 13 9 13H11.5C12.1 13 12.5 12.6 12.5 12C12.5 11.4 12.1 11 11.5 11H10C9.4 11 9 10.6 9 10V8H11.5C12.1 8 12.5 7.6 12.5 7C12.5 6.4 12.1 6 11.5 6H9C7.3 6 6 7.3 6 9V12C6 12.6 6.4 13 7 13H9V16C9 16.6 9.4 17 10 17H14C14.6 17 15 16.6 15 16V13H16C16.6 13 17 12.6 17 12V9C17 7.3 15.7 6 14 6H10V7.3C10.6 6.7 11.3 6.4 12 6.4H14C15.7 6.4 17 7.7 17 9.4V10H14C13.4 10 13 10.4 13 11C13 11.6 13.4 12 14 12H17V13H14C13.4 13 13 13.4 13 14C13 14.6 13.4 15 14 15H17V16H14C12.3 16 11 17.3 11 19V20C11 20.6 11.4 21 12 21H14C17.3 21 20 18.3 20 15V8H21Z"/></svg>
                前往 PayPal 付款
              </a>
              <p class="payment-link-hint">点击按钮跳转到 PayPal 完成付款</p>
            </div>
            <div v-if="selectedPaymentMethod.instructions" class="instructions">
              <p class="instructions-label">付款说明：</p>
              <p class="instructions-text">{{ selectedPaymentMethod.instructions }}</p>
            </div>
          </div>

          <div class="proof-upload">
            <p class="proof-label">上传付款截图</p>
            <div v-if="!proofPreview" class="upload-area" @click="$refs.proofInput.click()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="upload-icon">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
              <span>点击选择付款截图</span>
              <span class="upload-hint">支持 JPG、PNG、GIF、WebP，最大 10MB</span>
            </div>
            <div v-else class="proof-preview-wrapper">
              <img :src="proofPreview" class="proof-preview-img" />
              <button class="remove-proof-btn" @click="removeProof">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
            <input ref="proofInput" type="file" accept="image/*" @change="handleProofUpload" hidden />
          </div>
        </div>

        <div class="order-summary">
          <h2 class="section-title">订单汇总</h2>
          <div class="summary-row">
            <span>商品数量:</span>
            <span>{{ totalQuantity }}</span>
          </div>
          <div class="summary-row total-row">
            <span>总计:</span>
            <span class="total-value">{{ formatPrice(total) }}</span>
          </div>
          <button
            class="submit-btn"
            @click="createOrder"
            :disabled="loading || addresses.length === 0"
          >
            {{ loading ? '提交中...' : '提交订单' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.checkout-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.container {
  max-width: 800px;
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
.error-message {
  background: #3a1a1a;
  color: #e74;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}
.checkout-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}
.section-title {
  color: #fff;
  font-size: 18px;
  margin-bottom: 15px;
}
.order-items {
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
}
.items-list{
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.order-item{
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #2a2a2a;
}
.order-item:last-child {
  border-bottom: none;
}
.item-name{
  color: #fff;
  flex: 1;
}
.item-quantity{
  color: #888;
  margin: 0 20px;
}
.item-price{
  color: #d4a574;
  font-weight: bold;
}
.shipping-info{
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
}
.no-address{
  text-align: center;
  padding: 40px 20px;
  color: #888;
}
.no-address p {
  margin-bottom: 10px;
}
.hint{
  font-size: 13px;
  color: #666;
}
.btn-add-address{
  display: inline-block;
  background: #d4a574;
  color: #1a1a1a;
  padding: 12px 24px;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 600;
  margin-top: 15px;
}
.btn-add-address:hover{
  background: #c49564;
}
.address-list{
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.address-option{
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 15px;
  background: #2a2a2a;
  border-radius: 8px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}
.address-option:hover{
  border-color: #555;
}
.address-option.selected{
  border-color: #d4a574;
  background: rgba(212, 165, 116, 0.1);
}
.address-radio{
  flex-shrink: 0;
}
.radio{
  width: 20px;
  height: 20px;
  border: 2px solid #666;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.radio.checked{
  border-color: #d4a574;
}
.radio.checked::after{
  content: '';
  width: 10px;
  height: 10px;
  background: #d4a574;
  border-radius: 50%;
}
.address-details{
  flex: 1;
}
.address-name{
  color: #fff;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.default-tag{
  background: #d4a574;
  color: #1a1a1a;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}
.address-text{
  color: #ccc;
  font-size: 14px;
  line-height: 1.6;
  display: flex;
  align-items: flex-start;
}
.addr-icon {
  width: 16px;
  height: 16px;
  margin-right: 6px;
  flex-shrink: 0;
  margin-top: 2px;
}
.addr-label {
  color: #888;
  margin-right: 4px;
  flex-shrink: 0;
}
.address-phone{
  color: #d4a574;
  font-size: 14px;
  margin-top: 8px;
  display: flex;
  align-items: center;
}
.form-group{
  margin-bottom: 15px;
}
.form-group label{
  display: block;
  color: #888;
  margin-bottom: 8px;
  font-size: 14px;
}
.form-group textarea{
  width: 100%;
  padding: 12px;
  background: #2a2a2a;
  border: 1px solid #444;
  border-radius: 4px;
  color: #fff;
  font-size: 14px;
  box-sizing: border-box;
  min-height: 80px;
  resize: vertical;
}
.form-group textarea:focus{
  border-color: #d4a574;
  outline: none;
}

.payment-section{
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
}
.no-payment{
  text-align: center;
  padding: 30px 20px;
  color: #888;
}
.payment-method-list{
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}
.payment-method-option{
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 15px;
  background: #2a2a2a;
  border-radius: 8px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}
.payment-method-option:hover{
  border-color: #555;
}
.payment-method-option.selected{
  border-color: #d4a574;
  background: rgba(212, 165, 116, 0.1);
}
.payment-radio{
  flex-shrink: 0;
}
.payment-name{
  color: #fff;
  font-weight: 500;
}
.payment-detail{
  background: #2a2a2a;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}
.qrcode-wrapper{
  max-width: 220px;
}
.qrcode-img{
  width: 100%;
  border-radius: 8px;
  border: 1px solid #444;
}
.payment-link-wrapper{
  text-align: center;
}
.payment-link-btn{
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: #0070ba;
  color: #fff;
  text-decoration: none;
  border-radius: 6px;
  font-weight: 600;
  font-size: 15px;
  transition: background 0.3s;
}
.payment-link-btn:hover{
  background: #005ea6;
}
.paypal-icon{
  width: 20px;
  height: 20px;
}
.payment-link-hint{
  color: #888;
  font-size: 12px;
  margin-top: 8px;
}
.instructions{
  width: 100%;
  text-align: left;
}
.instructions-label{
  color: #d4a574;
  font-weight: 600;
  margin-bottom: 6px;
  font-size: 14px;
}
.instructions-text{
  color: #ccc;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.proof-upload{
  margin-top: 5px;
}
.proof-label{
  color: #888;
  margin-bottom: 10px;
  font-size: 14px;
}
.upload-area{
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
.upload-area:hover{
  border-color: #d4a574;
  color: #d4a574;
}
.upload-icon{
  width: 36px;
  height: 36px;
}
.upload-hint{
  font-size: 12px;
  color: #666;
}
.proof-preview-wrapper{
  position: relative;
  display: inline-block;
}
.proof-preview-img{
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
  border: 1px solid #444;
}
.remove-proof-btn{
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
.remove-proof-btn svg{
  width: 16px;
  height: 16px;
  stroke: #fff;
}

.order-summary{
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
}
.summary-row{
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  color: #888;
}
.total-row{
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #2a2a2a;
}
.total-value{
  color: #d4a574;
  font-size: 24px;
  font-weight: bold;
}
.submit-btn{
  width: 100%;
  background: #d4a574;
  color: #1a1a1a;
  border: none;
  padding: 15px;
  font-size: 16px;
  cursor: pointer;
  border-radius: 4px;
  font-weight: bold;
  margin-top: 20px;
}
.submit-btn:hover:not(:disabled){
  background: #e5b584;
}
.submit-btn:disabled{
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .checkout-page {
    padding: 20px 0 40px;
  }

  .page-title {
    font-size: 22px;
    margin-bottom: 20px;
  }

  .order-items,
  .shipping-info,
  .payment-section,
  .order-summary {
    padding: 15px;
  }

  .order-item {
    flex-wrap: wrap;
    gap: 4px;
  }

  .item-quantity {
    margin: 0 10px;
  }

  .radio {
    width: 24px;
    height: 24px;
  }

  .radio.checked::after {
    width: 12px;
    height: 12px;
  }

  .address-option {
    padding: 12px;
  }

  .address-text {
    flex-wrap: wrap;
  }

  .total-value {
    font-size: 20px;
  }

  .payment-method-option {
    padding: 10px 12px;
  }

  .qrcode-wrapper {
    max-width: 180px;
  }
}

@media (max-width: 576px) {
  .checkout-page {
    padding: 15px 0 30px;
  }

  .page-title {
    font-size: 20px;
  }

  .upload-area {
    padding: 20px 15px;
  }

  .proof-preview-img {
    max-height: 200px;
  }
}
</style>
