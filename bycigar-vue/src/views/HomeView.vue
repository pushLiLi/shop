<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const API_BASE = 'http://localhost:3000/api'

const currentSlide = ref(0)
const loading = ref(true)
const error = ref(null)
const config = ref({})
const featuredProducts = ref([])
const cubanProducts = ref([])
const banners = ref([])

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

let slideInterval = null

const nextSlide = () => {
  if (banners.value.length > 0) {
    currentSlide.value = (currentSlide.value + 1) % banners.value.length
  }
}

const prevSlide = () => {
  if (banners.value.length > 0) {
    currentSlide.value = (currentSlide.value - 1 + banners.value.length) % banners.value.length
  }
}

const goToSlide = (index) => {
  currentSlide.value = index
}

async function fetchData() {
  try {
    loading.value = true
    const [configRes, featuredRes, cubanRes, bannersRes] = await Promise.all([
      fetch(`${API_BASE}/config`),
      fetch(`${API_BASE}/products?featured=true&limit=6`),
      fetch(`${API_BASE}/products?category=cuban&limit=3`),
      fetch(`${API_BASE}/banners`)
    ])
    
    config.value = await configRes.json()
    
    const featuredData = await featuredRes.json()
    featuredProducts.value = featuredData.products || []
    
    const cubanData = await cubanRes.json()
    cubanProducts.value = cubanData.products || []
    
    const bannersData = await bannersRes.json()
    banners.value = bannersData.length > 0 ? bannersData : [
      { imageUrl: '/static/media/微信图片_20260303152810_1_341(2).jpg', link: '/brand-gl-pease' },
      { imageUrl: '/static/media/banner-4.png', link: '#' },
      { imageUrl: '/static/media/banner-5.jpg', link: '#' }
    ]
  } catch (e) {
    error.value = e.message
    console.error('Error:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
  slideInterval = setInterval(nextSlide, 3000)
})

onUnmounted(() => {
  if (slideInterval) {
    clearInterval(slideInterval)
  }
})

const formatPrice = (price) => {
  return `$${price.toFixed(2)}`
}

</script>

<template>
  <main class="home-page">
    <section class="hero-slider">
      <div class="slider-container">
        <div 
          v-for="(banner, index) in banners" 
          :key="index"
          class="slide"
          :class="{ active: index === currentSlide }"
        >
          <router-link :to="banner.link || '#'">
            <img :src="banner.imageUrl" :alt="banner.title || 'Banner ' + (index + 1)">
          </router-link>
        </div>
        <button class="slider-btn prev" @click="prevSlide">&#10094;</button>
        <button class="slider-btn next" @click="nextSlide">&#10095;</button>
        <div class="slider-dots">
          <button 
            v-for="(_, index) in banners" 
            :key="index"
            class="dot"
            :class="{ active: index === currentSlide }"
            @click="goToSlide(index)"
          ></button>
        </div>
      </div>
    </section>

    <section class="products-section">
      <div class="container">
        <h2 class="section-title">{{ config.home_featured_title || '特别推荐' }}</h2>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="products-grid grid-6">
          <div v-for="product in featuredProducts" :key="product.id" class="product-card">
            <router-link :to="'/products/' + product.id" class="product-image">
              <img :src="product.imageUrl" :alt="product.name">
            </router-link>
            <div class="product-info">
              <h3 class="product-name">
                <router-link :to="'/products/' + product.id">{{ product.name }}</router-link>
              </h3>
              <div class="product-bottom">
                <button class="add-cart-btn">加入购物车</button>
                <div class="product-price">{{ formatPrice(product.price) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="banner-section">
      <div class="container">
        <img src="/static/media/banner-1.png" alt="Banner" class="full-width-banner">
      </div>
    </section>

    <section class="products-section">
      <div class="container">
        <h2 class="section-title">{{ config.home_cuban_title || '古巴雪茄推荐' }}</h2>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="products-grid grid-3">
          <div v-for="product in cubanProducts" :key="product.id" class="product-card">
            <router-link :to="'/products/' + product.id" class="product-image">
              <img :src="product.imageUrl" :alt="product.name">
            </router-link>
            <div class="product-info">
              <h3 class="product-name">
                <router-link :to="'/products/' + product.id">{{ product.name }}</router-link>
              </h3>
              <div class="product-bottom">
                <button class="add-cart-btn">加入购物车</button>
                <div class="product-price">{{ formatPrice(product.price) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="banner-section">
      <div class="container">
        <img src="/static/media/banner-3.png" alt="Banner" class="full-width-banner">
      </div>
    </section>

    <section class="banner-section">
      <div class="container">
        <img src="/static/media/banner-2.png" alt="Banner" class="full-width-banner">
      </div>
    </section>
  </main>
</template>

<style scoped>
.home-page {
  background: #0f0f0f;
  min-height: 100vh;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 15px;
}

.hero-slider {
  margin-bottom: 40px;
}

.slider-container {
  position: relative;
  overflow: hidden;
}

.slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  opacity: 0;
  transition: opacity 0.5s ease;
}

.slide.active {
  position: relative;
  opacity: 1;
}

.slide img {
  width: 100%;
  height: auto;
  max-height: 600px;
  object-fit: cover;
  display: block;
}

.slider-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0,0,0,0.5);
  color: #fff;
  border: none;
  padding: 15px 20px;
  cursor: pointer;
  font-size: 24px;
  transition: background 0.3s;
  z-index: 10;
}

.slider-btn:hover {
  background: rgba(0,0,0,0.8);
}

.slider-btn.prev {
  left: 0;
}

.slider-btn.next {
  right: 0;
}

.slider-dots {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
  z-index: 10;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: rgba(255,255,255,0.5);
  border: none;
  cursor: pointer;
  transition: background 0.3s;
}

.dot.active {
  background: #d4a574;
}

.products-section {
  padding: 40px 0;
}

.section-title {
  text-align: center;
  color: #fff;
  font-size: 24px;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #d4a574;
  display: inline-block;
  width: 100%;
}

.products-grid {
  display: grid;
  gap: 15px;
}

.products-grid.grid-6 {
  grid-template-columns: repeat(6, 1fr);
}

.products-grid.grid-3 {
  grid-template-columns: repeat(3, 1fr);
}

.product-card {
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(212, 165, 116, 0.1);
}

.product-image {
  display: block;
  background: #fff;
  padding: 10px;
}

.product-image img {
  width: 100%;
  height: auto;
  aspect-ratio: 1;
  object-fit: cover;
}

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 13px;
  margin: 0 0 10px;
  line-height: 1.4;
}

.product-name a {
  color: #ccc;
  text-decoration: none;
  transition: color 0.3s;
}

.product-name a:hover {
  color: #d4a574;
}

.product-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.add-cart-btn {
  background: transparent;
  border: 1px solid #d4a574;
  color: #d4a574;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 4px;
}

.add-cart-btn:hover {
  background: #d4a574;
  color: #1a1a1a;
}

.product-price {
  color: #d4a574;
  font-weight: bold;
  font-size: 14px;
}

.banner-section {
  padding: 20px 0;
}

.full-width-banner {
  width: 100%;
  height: auto;
  display: block;
  border-radius: 8px;
}

.loading {
  text-align: center;
  color: #d4a574;
  padding: 40px;
}

@media (max-width: 1200px) {
  .products-grid.grid-6 {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 992px) {
  .products-grid.grid-6 {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .products-grid.grid-3 {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
    .products-grid.grid-6 {
    grid-template-columns: repeat(2, 1fr);
  }
  
    .products-grid.grid-3 {
    grid-template-columns: 1fr;
  }
  
    .slide img {
    max-height: 300px;
  }
  
    .slider-btn {
    padding: 10px 15px;
    font-size: 18px;
  }
}
</style>
