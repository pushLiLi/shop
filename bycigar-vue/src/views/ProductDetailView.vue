<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useFavoritesStore } from '../stores/favorites'
import { useToastStore } from '../stores/toast'
import ProductCard from '../components/ProductCard.vue'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const favoritesStore = useFavoritesStore()
const toast = useToastStore()
const product = ref(null)
const relatedProducts = ref([])
const loading = ref(true)
const error = ref(null)
const quantity = ref(1)

const productId = computed(() => route.params.id)

const isFavorite = computed(() => {
  if (!product.value) return false
  return favoritesStore.items.some(item => item.productId === product.value.id)
})

const formatPrice = (price) => {
  return `$${Number(price).toFixed(2)}`
}

async function fetchProduct() {
  try {
    loading.value = true
    error.value = null
    
    const res = await fetch(`/api/products/${productId.value}`)
    if (!res.ok) throw new Error('产品不存在')
    
    const data = await res.json()
    product.value = data
    
    if (data.category?.id) {
      const relatedRes = await fetch(`/api/products?categoryId=${data.category.id}&limit=4`)
      const relatedData = await relatedRes.json()
      relatedProducts.value = relatedData.products?.filter(p => p.id !== data.id) || []
    }
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function increaseQty() {
  quantity.value++
}

function decreaseQty() {
  if (quantity.value > 1) quantity.value--
}

function handleQtyInput() {
  if (!quantity.value || quantity.value < 1) {
    quantity.value = 1
  }
}

async function addToCart() {
  if (!product.value) return
  await cartStore.addItem(product.value, quantity.value)
  toast.success(`已添加 ${quantity.value} 件到购物车`)
}

async function toggleFavorite() {
  if (!product.value) return
  if (isFavorite.value) {
    await favoritesStore.removeItem(product.value.id)
  } else {
    await favoritesStore.addItem(product.value)
  }
}

onMounted(() => {
  fetchProduct()
})
</script>

<template>
  <div class="product-detail-page">
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <router-link to="/" class="back-link">返回首页</router-link>
    </div>
    <template v-else-if="product">
      <div class="container">
        <nav class="breadcrumb">
          <router-link to="/">首页</router-link>
          <span class="separator">/</span>
          <router-link v-if="product.category" :to="'/category/' + product.category.slug">
            {{ product.category.name }}
          </router-link>
          <span class="separator">/</span>
          <span class="current">{{ product.name }}</span>
        </nav>

        <div class="product-main">
          <div class="product-gallery">
            <div class="main-image">
              <img :src="product.imageUrl" :alt="product.name">
            </div>
          </div>

          <div class="product-info">
            <h1 class="product-title">{{ product.name }}</h1>
            <div class="product-brand" v-if="product.brand">{{ product.brand }}</div>
            <div class="product-price-main">{{ formatPrice(product.price) }}</div>
            
            <div class="product-description" v-if="product.description">
              <p>{{ product.description }}</p>
            </div>

            <div class="purchase-section">
              <div class="quantity-selector">
                <button class="qty-btn" @click="decreaseQty">-</button>
                <input type="number" v-model.number="quantity" min="1" class="qty-input" @input="handleQtyInput">
                <button class="qty-btn" @click="increaseQty">+</button>
              </div>
              <button class="buy-btn" @click="addToCart">加入购物车</button>
              <button class="favorite-btn" @click="toggleFavorite" :class="{ active: isFavorite }">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" :fill="isFavorite ? '#d4a574' : 'none'" stroke="currentColor" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <section class="related-section" v-if="relatedProducts.length > 0">
          <h2 class="section-title">相关产品</h2>
          <div class="products-grid">
            <ProductCard v-for="p in relatedProducts" :key="p.id" :product="p" />
          </div>
        </section>
      </div>
    </template>
  </div>
</template>

<style scoped>
.product-detail-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 20px 0 60px;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 15px;
}

.loading, .error {
  text-align: center;
  padding: 100px 20px;
  color: #d4a574;
  font-size: 18px;
}

.back-link {
  color: #d4a574;
  text-decoration: underline;
}

.breadcrumb {
  padding: 20px 0;
  font-size: 14px;
}

.breadcrumb a {
  color: #888;
  text-decoration: none;
}

.breadcrumb a:hover {
  color: #d4a574;
}

.breadcrumb .separator {
  margin: 0 10px;
  color: #555;
}

.breadcrumb .current {
  color: #d4a574;
}

.product-main {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  margin-bottom: 60px;
}

.product-gallery {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
}

.main-image {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.main-image img {
  width: 100%;
  height: auto;
  aspect-ratio: 1;
  object-fit: cover;
}

.product-info {
  padding: 20px 0;
}

.product-title {
  color: #fff;
  font-size: 28px;
  margin: 0 0 10px;
  line-height: 1.3;
}

.product-brand {
  color: #888;
  font-size: 16px;
  margin-bottom: 20px;
}

.product-price-main {
  color: #d4a574;
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 30px;
}

.product-description {
  color: #aaa;
  line-height: 1.8;
  margin-bottom: 30px;
  padding: 20px;
  background: #1a1a1a;
  border-radius: 8px;
}

.purchase-section {
  display: flex;
  gap: 20px;
  align-items: center;
}

.quantity-selector {
  display: flex;
  align-items: center;
  border: 1px solid #333;
  border-radius: 4px;
  overflow: hidden;
}

.qty-btn {
  background: #1a1a1a;
  border: none;
  color: #d4a574;
  width: 40px;
  height: 40px;
  cursor: pointer;
  font-size: 18px;
}

.qty-btn:hover {
  background: #2a2a2a;
}

.qty-input {
  background: #2a2a2a;
  border: 1px solid #444;
  border-radius: 4px;
  color: #fff;
  width: 60px;
  height: 40px;
  text-align: center;
  font-size: 16px;
  -moz-appearance: textfield;
}

.qty-input::-webkit-inner-spin-button,
.qty-input::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.qty-input:focus {
  outline: none;
  box-shadow: 0 0 8px rgba(212, 165, 116, 0.5);
}

.qty-input::-webkit-inner-spin-button,
.qty-input::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.buy-btn {
  flex: 1;
  background: #d4a574;
  color: #1a1a1a;
  border: none;
  padding: 15px 40px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.buy-btn:hover {
  background: #e5b584;
}

.favorite-btn {
  background: transparent;
  border: 1px solid #333;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  color: #888;
}

.favorite-btn:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.favorite-btn.active {
  color: #d4a574;
}

.related-section {
  padding-top: 40px;
  border-top: 1px solid #2a2a2a;
}

.section-title {
  color: #fff;
  font-size: 24px;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #d4a574;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
}

@media (max-width: 992px) {
  .product-main {
    grid-template-columns: 1fr;
  }
  
  .products-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .products-grid {
    grid-template-columns: 1fr;
  }
  
  .purchase-section {
    flex-direction: column;
  }
  
  .buy-btn {
    width: 100%;
  }
}
</style>
