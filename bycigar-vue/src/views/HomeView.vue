<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useCarousel } from '../composables/useCarousel'
import { marked } from 'marked'
import ProductCard from '../components/ProductCard.vue'

const API_BASE = '/api'

const loading = ref(true)
const error = ref(null)
const config = ref({})
const featuredProducts = ref([])
const newProducts = ref([])
const topSellingProducts = ref([])
const categoryProducts = ref([])
const brandStory = ref(null)
const brandPhilosophy = ref(null)

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
    const [configRes, featuredRes, bannersRes, categoriesRes, newRes, topSellingRes, storyRes, philosophyRes] = await Promise.all([
      fetch(`${API_BASE}/config`),
      fetch(`${API_BASE}/products?featured=true&limit=12`),
      fetch(`${API_BASE}/banners`),
      fetch(`${API_BASE}/categories`),
      fetch(`${API_BASE}/products?sortBy=createdAt&sortOrder=desc&limit=8`),
      fetch(`${API_BASE}/products/top-selling?limit=8`),
      fetch(`${API_BASE}/pages/brand-story`),
      fetch(`${API_BASE}/pages/brand-philosophy`)
    ])
    
    config.value = await configRes.json()
    
    const featuredData = await featuredRes.json()
    featuredProducts.value = featuredData.products || []

    const newData = await newRes.json()
    newProducts.value = newData.products || []

    const topSellingData = await topSellingRes.json()
    topSellingProducts.value = topSellingData.products || []
    
    const bannersData = await bannersRes.json()
    banners.value = bannersData.length > 0 ? bannersData : [
      { imageUrl: '/media/bycigar/微信图片_20260303152810_1_341(2).jpg', link: '/brand-gl-pease' },
      { imageUrl: '/media/bycigar/banner-4.png', link: '#' },
      { imageUrl: '/media/bycigar/banner-5.jpg', link: '#' }
    ]

    const categoriesData = await categoriesRes.json()
    const categoriesWithProducts = (categoriesData || []).filter(c => c._count > 0)
    
    if (categoriesWithProducts.length > 0) {
      const productResults = await Promise.all(
        categoriesWithProducts.map(async (cat) => {
          const res = await fetch(`${API_BASE}/products?categoryId=${cat.id}&limit=8`)
          const data = await res.json()
          return { category: cat, products: data.products || [] }
        })
      )
      categoryProducts.value = productResults.filter(item => item.products.length > 0)
    }

    if (storyRes.ok) {
      const storyData = await storyRes.json()
      if (storyData.content) {
        storyData.htmlContent = marked(storyData.content)
        brandStory.value = storyData
      }
    }

    if (philosophyRes.ok) {
      const philosophyData = await philosophyRes.json()
      if (philosophyData.content) {
        philosophyData.htmlContent = marked(philosophyData.content)
        brandPhilosophy.value = philosophyData
      }
    }
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

    <section class="promo-section" v-if="config.home_promo_left_image || config.home_promo_right_image">
      <div class="container">
        <div class="promo-grid">
          <a v-if="config.home_promo_left_image" :href="config.home_promo_left_link || '#'" class="promo-card">
            <img :src="config.home_promo_left_image" alt="Promo Left" loading="lazy">
          </a>
          <a v-if="config.home_promo_right_image" :href="config.home_promo_right_link || '#'" class="promo-card">
            <img :src="config.home_promo_right_image" alt="Promo Right" loading="lazy">
          </a>
        </div>
      </div>
    </section>

    <section class="new-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ config.home_new_title || '新品上架' }}</h2>
          <router-link to="/products?sortBy=createdAt&sortOrder=desc" class="view-more">
            查看更多 <span class="arrow">&rarr;</span>
          </router-link>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="products-grid grid-6">
          <ProductCard v-for="product in newProducts" :key="'new-' + product.id" :product="product" />
        </div>
      </div>
    </section>

    <section class="banner-section" v-if="config.home_banner_1">
      <div class="container">
        <img :src="config.home_banner_1" alt="Banner" class="full-width-banner">
      </div>
    </section>

    <section
      v-for="item in categoryProducts"
      :key="item.category.id"
      class="category-section"
    >
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ item.category.name }}</h2>
          <router-link :to="'/category/' + item.category.slug" class="view-more">
            查看更多 <span class="arrow">&rarr;</span>
          </router-link>
        </div>
        <div class="products-grid grid-6">
          <ProductCard v-for="product in item.products" :key="product.id" :product="product" />
        </div>
      </div>
    </section>

    <section class="top-selling-section" v-if="topSellingProducts.length > 0">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ config.home_topselling_title || '热销排行' }}</h2>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="products-grid grid-6">
          <ProductCard v-for="product in topSellingProducts" :key="'top-' + product.id" :product="product" />
        </div>
      </div>
    </section>

    <section class="services-section">
      <div class="container">
        <div class="services-grid">
          <div class="service-item">
            <div class="service-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#d4a574" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
            </div>
            <h3>正品保证</h3>
            <p>100%正品 假一赔十</p>
          </div>
          <div class="service-item">
            <div class="service-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#d4a574" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>
            </div>
            <h3>极速配送</h3>
            <p>下单即发 极速送达</p>
          </div>
          <div class="service-item">
            <div class="service-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#d4a574" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 12 20 22 4 22 4 12"/><rect x="2" y="7" width="20" height="5"/><line x1="12" y1="22" x2="12" y2="7"/><path d="M12 7H7.5a2.5 2.5 0 0 1 0-5C11 2 12 7 12 7z"/><path d="M12 7h4.5a2.5 2.5 0 0 0 0-5C13 2 12 7 12 7z"/></svg>
            </div>
            <h3>售后无忧</h3>
            <p>7天无理由 退换便捷</p>
          </div>
          <div class="service-item">
            <div class="service-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#d4a574" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
            </div>
            <h3>专业客服</h3>
            <p>在线答疑 贴心服务</p>
          </div>
        </div>
      </div>
    </section>

    <section class="brand-section" v-if="brandStory">
      <div class="container">
        <h2 class="section-title centered">{{ brandStory.title }}</h2>
        <div class="brand-content" v-html="brandStory.htmlContent"></div>
      </div>
    </section>

    <section class="philosophy-section" v-if="brandPhilosophy">
      <div class="container">
        <h2 class="section-title centered">{{ brandPhilosophy.title }}</h2>
        <div class="brand-content" v-html="brandPhilosophy.htmlContent"></div>
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

.products-section .section-title {
  text-align: center;
  color: #fff;
  font-size: 24px;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #d4a574;
  display: inline-block;
  width: 100%;
}

.promo-section {
  padding: 20px 0;
}

.promo-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.promo-card {
  display: block;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.3s;
}

.promo-card:hover {
  transform: translateY(-3px);
}

.promo-card img {
  width: 100%;
  height: auto;
  aspect-ratio: 2/1;
  object-fit: cover;
  display: block;
}

.new-section {
  padding: 30px 0;
  border-top: 1px solid rgba(212, 165, 116, 0.15);
}

.new-section .section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 25px;
}

.new-section .section-title {
  color: #fff;
  font-size: 20px;
  margin: 0;
  position: relative;
  padding-left: 14px;
}

.new-section .section-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 20px;
  background: #d4a574;
  border-radius: 2px;
}

.category-section {
  padding: 30px 0;
  border-top: 1px solid rgba(212, 165, 116, 0.15);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 25px;
}

.section-header .section-title {
  color: #fff;
  font-size: 20px;
  margin: 0;
  padding: 0;
  border-bottom: none;
  position: relative;
  padding-left: 14px;
}

.section-header .section-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 20px;
  background: #d4a574;
  border-radius: 2px;
}

.view-more {
  color: #999;
  font-size: 14px;
  text-decoration: none;
  transition: color 0.3s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.view-more:hover {
  color: #d4a574;
}

.view-more .arrow {
  transition: transform 0.3s;
}

.view-more:hover .arrow {
  transform: translateX(3px);
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

.top-selling-section {
  padding: 30px 0;
  border-top: 1px solid rgba(212, 165, 116, 0.15);
}

.top-selling-section .section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 25px;
}

.top-selling-section .section-title {
  color: #fff;
  font-size: 20px;
  margin: 0;
  position: relative;
  padding-left: 14px;
}

.top-selling-section .section-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 20px;
  background: #d4a574;
  border-radius: 2px;
}

.services-section {
  padding: 50px 0;
  border-top: 1px solid rgba(212, 165, 116, 0.15);
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  text-align: center;
}

.service-item {
  padding: 30px 15px;
  background: #1a1a1a;
  border-radius: 8px;
  transition: transform 0.3s;
}

.service-item:hover {
  transform: translateY(-3px);
}

.service-icon {
  margin-bottom: 15px;
}

.service-item h3 {
  color: #d4a574;
  font-size: 16px;
  margin: 0 0 8px;
}

.service-item p {
  color: #999;
  font-size: 13px;
  margin: 0;
}

.brand-section,
.philosophy-section {
  padding: 50px 0;
  border-top: 1px solid rgba(212, 165, 116, 0.15);
}

.section-title.centered {
  text-align: center;
  color: #fff;
  font-size: 24px;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #d4a574;
  display: inline-block;
  width: 100%;
}

.brand-content {
  color: #ccc;
  line-height: 1.8;
  font-size: 15px;
  max-width: 900px;
  margin: 0 auto;
}

.brand-content :deep(h1),
.brand-content :deep(h2),
.brand-content :deep(h3) {
  color: #d4a574;
}

.brand-content :deep(img) {
  max-width: 100%;
  border-radius: 8px;
  margin: 15px 0;
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
    bottom: 10px;
    gap: 6px;
  }

  .dot {
    width: 8px;
    height: 8px;
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
    width: 8px;
    height: 8px;
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

  .section-header .section-title {
    font-size: 18px;
  }

  .view-more {
    font-size: 13px;
  }

  .promo-grid {
    grid-template-columns: 1fr;
  }

  .services-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .brand-content {
    font-size: 14px;
  }
}
</style>
