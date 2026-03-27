<script setup>
import { computed, onMounted } from 'vue'
import { useCartStore } from '../stores/cart'

const cartStore = useCartStore()

const items = computed(() => cartStore.items)
const total = computed(() => cartStore.total)
const loading = computed(() => cartStore.loading)

function formatPrice(price) {
  return `$${Number(price).toFixed(2)}`
}

onMounted(() => {
  cartStore.fetchCart()
})
</script>

<template>
  <div class="cart-page">
    <div class="container">
      <h1 class="page-title">购物车</h1>
      
      <div v-if="loading" class="loading">加载中...</div>
      
      <div v-else-if="items.length === 0" class="empty-cart">
        <p>购物车是空的</p>
        <router-link to="/" class="continue-shopping">继续购物</router-link>
      </div>
      
      <div v-else>
        <div class="cart-items">
          <div v-for="item in items" :key="item.productId" class="cart-item">
            <router-link :to="'/products/' + item.productId" class="item-image">
              <img :src="item.product?.imageUrl" :alt="item.product?.name">
            </router-link>
            <div class="item-info">
              <h3 class="item-name">
                <router-link :to="'/products/' + item.productId">{{ item.product?.name }}</router-link>
              </h3>
              <div class="item-price">单价: {{ formatPrice(item.product?.price) }}</div>
              <div class="item-actions">
                <div class="quantity-control">
                  <button @click="cartStore.updateQuantity(item.productId, item.quantity - 1)">-</button>
                  <span>{{ item.quantity }}</span>
                  <button @click="cartStore.updateQuantity(item.productId, item.quantity + 1)">+</button>
                </div>
                <button @click="cartStore.removeItem(item.productId)" class="remove-btn">删除</button>
              </div>
              <div class="item-total">
                小计: {{ formatPrice((item.product?.price || 0) * item.quantity) }}
              </div>
            </div>
          </div>
        </div>
        
        <div class="cart-summary">
          <div class="summary-row">
            <span class="summary-label">总计:</span>
            <span class="summary-value">{{ formatPrice(total) }}</span>
          </div>
          <router-link to="/checkout" class="checkout-btn">去结算</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cart-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-title {
  color: #d4a574;
  font-size: 28px;
  margin-bottom: 30px;
  border-bottom: 2px solid #d4a574;
}

.loading, .empty-cart {
  text-align: center;
  padding: 80px 20px;
  color: #888;
}

.continue-shopping {
  color: #d4a574;
  text-decoration: none;
}

.continue-shopping:hover {
  text-decoration: underline;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 30px;
}

.cart-item {
  display: flex;
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
}

.item-image {
  width: 120px;
  height: 120px;
  flex-shrink: 0;
  background: #fff;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.item-info {
  flex: 1;
  padding: 15px;
}

.item-name {
  margin: 0 0 10px;
}

.item-name a {
  color: #fff;
  text-decoration: none;
  font-size: 16px;
}

.item-name a:hover {
  color: #d4a574;
}

.item-price {
  color: #888;
  font-size: 14px;
  margin-bottom: 10px;
}

.item-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 10px;
}

.quantity-control button {
  background: #2a2a2a;
  border: 1px solid #444;
  color: #fff;
  width: 30px;
  height: 30px;
  cursor: pointer;
  border-radius: 4px;
}

.quantity-control button:hover {
  background: #3a3a3a;
}

.quantity-control span {
  color: #fff;
  min-width: 40px;
  text-align: center;
}

.remove-btn {
  background: transparent;
  border: 1px solid #e74;
  color: #e74;
  padding: 6px 12px;
  cursor: pointer;
  border-radius: 4px;
}

.remove-btn:hover {
  background: #e74;
  color: #1a1a1a;
}

.item-total {
  color: #d4a574;
  font-weight: bold;
  font-size: 16px;
  margin-top: 10px;
}

.cart-summary {
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.summary-label {
  color: #888;
  font-size: 14px;
}

.summary-value {
  color: #d4a574;
  font-size: 24px;
  font-weight: bold;
}

.checkout-btn {
  background: #d4a574;
  color: #1a1a1a;
  padding: 12px 30px;
  font-size: 16px;
  border-radius: 4px;
  text-decoration: none;
  font-weight: bold;
}

.checkout-btn:hover {
  background: #e5b584;
}

@media (max-width: 768px) {
  .cart-item {
    flex-direction: column;
  }
  
  .item-image {
    width: 100%;
    height: 200px;
  }
  
  .cart-summary {
    flex-direction: column;
    gap: 20px;
  }
}
</style>
