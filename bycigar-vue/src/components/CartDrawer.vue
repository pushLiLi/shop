<script setup>
import { computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'

const cartStore = useCartStore()
const router = useRouter()

const items = computed(() => cartStore.items)
const total = computed(() => cartStore.total)
const loading = computed(() => cartStore.loading)

function formatPrice(price) {
  return `$${Number(price).toFixed(2)}`
}

function handleInput(item, e) {
  const val = parseInt(e.target.value) || 0
  const idx = cartStore.items.findIndex(i => i.id === item.id)
  if (idx !== -1) {
    cartStore.items[idx] = { ...cartStore.items[idx], quantity: val }
  }
}

function handleBlur(item) {
  let qty = parseInt(item.quantity) || 0
  if (qty < 1) qty = 1
  const idx = cartStore.items.findIndex(i => i.id === item.id)
  if (idx !== -1) {
    cartStore.items[idx] = { ...cartStore.items[idx], quantity: qty }
  }
  cartStore.updateQuantity(item.id, qty)
}

function closeDrawer() {
  cartStore.closeCart()
}

function goToCheckout() {
  cartStore.closeCart()
  router.push('/checkout')
}

function handleOverlayClick(e) {
  if (e.target === e.currentTarget) {
    closeDrawer()
  }
}

watch(() => cartStore.isOpen, (isOpen) => {
  if (isOpen) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="cartStore.isOpen" class="cart-drawer-overlay" @click="handleOverlayClick">
        <div class="cart-drawer">
          <div class="drawer-header">
            <h2>购物车</h2>
            <button class="close-btn" @click="closeDrawer">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="drawer-content">
            <div v-if="loading" class="loading">加载中...</div>

            <div v-else-if="items.length === 0" class="empty-cart">
              <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                <circle cx="9" cy="21" r="1"></circle>
                <circle cx="20" cy="21" r="1"></circle>
                <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path>
              </svg>
              <p>购物车是空的</p>
              <button class="continue-btn" @click="closeDrawer">继续购物</button>
            </div>

            <div v-else class="cart-items">
              <div v-for="item in items" :key="item.productId" class="cart-item">
                <div class="item-image">
                  <img :src="item.product?.imageUrl" :alt="item.product?.name">
                </div>
                <div class="item-info">
                  <h3 class="item-name">{{ item.product?.name }}</h3>
                  <div class="item-price">{{ formatPrice(item.product?.price) }}</div>
                  <div class="item-actions">
                    <div class="quantity-control">
                      <button @click="cartStore.updateQuantity(item.id, item.quantity - 1)">-</button>
                      <input 
                        type="number"
                        :value="item.quantity"
                        @input="handleInput(item, $event)"
                        @blur="handleBlur(item)"
                        @keyup.enter="handleBlur(item)"
                        min="1"
                      />
                      <button @click="cartStore.updateQuantity(item.id, item.quantity + 1)">+</button>
                    </div>
                    <button 
                      @click="cartStore.removeItem(item.id)" 
                      class="remove-btn"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      </svg>
                    </button>
                  </div>
                  <div class="item-total">
                    小计: {{ formatPrice((item.product?.price || 0) * item.quantity) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="items.length > 0" class="drawer-footer">
            <div class="summary-row">
              <span class="summary-label">总计</span>
              <span class="summary-value">{{ formatPrice(total) }}</span>
            </div>
            <button class="checkout-btn" @click="goToCheckout">去结算</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.cart-drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  justify-content: flex-end;
}

.cart-drawer {
  width: 400px;
  max-width: 100%;
  height: 100%;
  background: #1a1a1a;
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.3);
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #333;
}

.drawer-header h2 {
  color: #d4a574;
  font-size: 20px;
  font-weight: 500;
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: #888;
  cursor: pointer;
  padding: 5px;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #fff;
}

.drawer-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #888;
}

.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #888;
}

.empty-cart svg {
  opacity: 0.3;
  margin-bottom: 20px;
}

.empty-cart p {
  margin-bottom: 20px;
}

.continue-btn {
  background: transparent;
  border: 1px solid #d4a574;
  color: #d4a574;
  padding: 10px 24px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
}

.continue-btn:hover {
  background: #d4a574;
  color: #1a1a1a;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.cart-item {
  display: flex;
  gap: 15px;
  background: #252525;
  border-radius: 8px;
  padding: 12px;
}

.item-image {
  width: 80px;
  height: 80px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 4px;
  overflow: hidden;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  color: #fff;
  font-size: 14px;
  font-weight: 400;
  margin: 0 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-price {
  color: #888;
  font-size: 13px;
  margin-bottom: 8px;
}

.item-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 4px;
}

.quantity-control button {
  background: #2a2a2a;
  border: 1px solid #444;
  color: #fff;
  width: 24px;
  height: 24px;
  cursor: pointer;
  border-radius: 4px;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quantity-control button:hover {
  background: #3a3a3a;
}

.quantity-control input {
  background: #2a2a2a;
  border: 1px solid #444;
  color: #fff;
  width: 40px;
  height: 24px;
  text-align: center;
  font-size: 13px;
  border-radius: 4px;
}

.quantity-control input::-webkit-outer-spin-button,
.quantity-control input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.quantity-control input[type=number] {
  -moz-appearance: textfield;
}

.remove-btn {
  background: transparent;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s;
}

.remove-btn:hover {
  color: #e74;
}

.item-total {
  color: #d4a574;
  font-size: 14px;
  font-weight: 500;
}

.drawer-footer {
  padding: 20px;
  border-top: 1px solid #333;
  background: #1a1a1a;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.summary-label {
  color: #888;
  font-size: 14px;
}

.summary-value {
  color: #d4a574;
  font-size: 22px;
  font-weight: bold;
}

.checkout-btn {
  width: 100%;
  background: #d4a574;
  color: #1a1a1a;
  border: none;
  padding: 14px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
}

.checkout-btn:hover {
  background: #e5b584;
}

.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-enter-active .cart-drawer,
.drawer-leave-active .cart-drawer {
  transition: transform 0.3s ease;
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from .cart-drawer,
.drawer-leave-to .cart-drawer {
  transform: translateX(100%);
}

@media (max-width: 480px) {
  .cart-drawer {
    width: 100%;
  }
}
</style>
