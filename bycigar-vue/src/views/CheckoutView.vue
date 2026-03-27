<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { getStateName } from '../utils/states'

const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const toast = useToastStore()

const items = computed(() => cartStore.items)
const total = computed(() => cartStore.total)
const totalQuantity = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
const loading = ref(false)
const error = ref(null)
const addresses = ref([])
const selectedAddressId = ref(null)
const remark = ref('')

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
})

async function fetchAddresses() {
  try {
    const res = await fetch('http://localhost:3000/api/addresses', {
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

async function createOrder() {
  if (items.value.length === 0) {
    error.value = '购物车是空的'
    return
  }
  if (!selectedAddressId.value) {
    error.value = '请选择收货地址'
    return
  }
  try {
    loading.value = true
    error.value = null
    const orderData = {
      addressId: selectedAddressId.value,
      remark: remark.value
    }
    const res = await fetch('http://localhost:3000/api/orders', {
      method: 'POST',
      headers: authStore.getAuthHeaders(),
      body: JSON.stringify(orderData)
    })
    await res.json()
    cartStore.clear()
    toast.success('订单创建成功！')
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
</style>
