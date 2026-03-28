<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useCarousel } from '../composables/useCarousel'
import ProductCard from '../components/ProductCard.vue'

const API_BASE = 'http://localhost:3000/api'

const loading = ref(true)
const error = ref(null)
const config = ref({})
const featuredProducts = ref([])

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

const {
  currentIndex,
  slides: banners,
  next: nextSlide,
  prev: prevSlide,
  goTo: goToSlide,
  onMouseEnter,
  onMouseLeave,
  onTouchStart,
  onTouchEnd
} = useCarousel({ autoplay: true, interval: 4000, pauseOnHover: true })

async function fetchData() {
  try {
    loading.value = true
    const [configRes, featuredRes, bannersRes] = await Promise.all([
      fetch(`${API_BASE}/config`),
      fetch(`${API_BASE}/products?featured=true&limit=12`),
      fetch(`${API_BASE}/banners`)
    ])
    
    config.value = await configRes.json()
    
    const featuredData = await featuredRes.json()
    featuredProducts.value = featuredData.products || []
    
    const bannersData = await bannersRes.json()
    banners.value = bannersData.length > 0 ? bannersData : [
      { imageUrl: '/media/bycigar/微信图片_20260303152810_1_341(2).jpg', link: '/brand-gl-pease' },
      { imageUrl: '/media/bycigar/banner-4.png', link: '#' },
      { imageUrl: '/media/bycigar/banner-5.jpg', link: '#' }
    ]
  } catch (e) {
    error.value = e.message
    console.error('Error:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => { fetchData() })
</script>

<template>
  <main class="home-page">
    <section class="hero-slider">
      <div
        class="slider-container"
        @mouseenter="onMouseEnter"
        @mouseleave="onMouseLeave"
        @touchstart="onTouchStart"
        @touchend="onTouchEnd"
      >
        <div
          class="slider-track"
          :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
        >
          <div
            v-for="(banner, index) in banners"
            :key="index"
            class="slide"
          >
            <router-link :to="banner.link || '#'">
              <img
                :src="banner.imageUrl"
                :alt="banner.title || 'Banner ' + (index + 1)"
                :loading="index > 1 ? 'lazy' : 'eager'"
                :fetchpriority="index === 0 ? 'high' : 'auto'"
              >
            </router-link>
          </div>
        </div>
        <button class="slider-btn prev" @click="prevSlide">&#10094;</button>
        <button class="slider-btn next" @click="nextSlide">&#10095;</button>
        <div class="slider-dots">
          <button
            v-for="(_, index) in banners"
            :key="index"
            class="dot"
            :class="{ active: index === currentIndex }"
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
          <ProductCard v-for="product in featuredProducts" :key="product.id" :product="product" />
        </div>
      </div>
    </section>

    <section class="banner-section" v-if="config.home_banner_1">
      <div class="container">
        <img :src="config.home_banner_1" alt="Banner" class="full-width-banner">
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
  touch-action: pan-y;
}

.slider-track {
  display: flex;
  transition: transform 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  will-change: transform;
}

.slide {
  flex: 0 0 100%;
  width: 100%;
}

.slide img {
  width: 100%;
  height: auto;
  aspect-ratio: 3/1;
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
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
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

@media (max-width: 768px) {
  .slide img {
    aspect-ratio: 2/1;
  }

  .slider-btn {
    padding: 12px 14px;
    font-size: 18px;
    min-width: 44px;
    min-height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .slider-dots {
    bottom: 12px;
  }

  .dot {
    width: 10px;
    height: 10px;
    min-width: 44px;
    min-height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    position: relative;
  }

  .dot::after {
    content: '';
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: rgba(255,255,255,0.5);
    position: absolute;
  }

  .dot.active::after {
    background: #d4a574;
  }

  .products-grid.grid-6 {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
