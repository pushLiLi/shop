<script setup>
import { ref, onMounted, computed, reactive } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useCarousel } from '../composables/useCarousel'
import ProductCard from '../components/ProductCard.vue'

const API_BASE = '/api'

const loading = ref(true)
const error = ref(null)
const config = ref({})
const featuredProducts = ref([])
const newProducts = ref([])
const topSellingProducts = ref([])
const categoryProducts = ref([])

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

const featuredScroll = ref(null)
const newScroll = ref(null)
const topSellingScroll = ref(null)
const categoryScrolls = ref([])

const scrollStates = reactive({
  featured: { canLeft: false, canRight: true },
  new: { canLeft: false, canRight: true },
  topSelling: { canLeft: false, canRight: true },
  categories: {}
})

function scrollSection(container, direction) {
  if (!container) return
  const item = container.querySelector('.product-scroll-item')
  const scrollAmount = item ? item.offsetWidth + 15 : 280
  container.scrollBy({
    left: direction === 'left' ? -scrollAmount : scrollAmount,
    behavior: 'smooth'
  })
}

function updateScrollState(event, key) {
  const el = event.target
  const state = key.startsWith('cat-') ? scrollStates.categories[key] : scrollStates[key]
  if (!state) return
  state.canLeft = el.scrollLeft > 5
  state.canRight = el.scrollLeft < el.scrollWidth - el.clientWidth - 5
}

async function fetchData() {
  try {
    loading.value = true
    const [configRes, featuredRes, bannersRes, categoriesRes, newRes, topSellingRes] = await Promise.all([
      fetch(`${API_BASE}/config`),
      fetch(`${API_BASE}/products?featured=true&limit=12`),
      fetch(`${API_BASE}/banners`),
      fetch(`${API_BASE}/categories`),
      fetch(`${API_BASE}/products?sortBy=createdAt&sortOrder=desc&limit=12`),
      fetch(`${API_BASE}/products/top-selling?limit=12`)
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
      { imageUrl: 'https://picsum.photos/seed/cigar1/1400/500', link: '#' },
      { imageUrl: 'https://picsum.photos/seed/cigar2/1400/500', link: '#' },
      { imageUrl: 'https://picsum.photos/seed/cigar3/1400/500', link: '#' },
      { imageUrl: 'https://picsum.photos/seed/cigar4/1400/500', link: '#' },
      { imageUrl: 'https://picsum.photos/seed/cigar5/1400/500', link: '#' }
    ]

    const categoriesData = await categoriesRes.json()
    const categoriesWithProducts = (categoriesData || []).filter(c => c._count > 0)

    if (categoriesWithProducts.length > 0) {
      const productResults = await Promise.all(
        categoriesWithProducts.map(async (cat) => {
          const res = await fetch(`${API_BASE}/products?categoryId=${cat.id}&limit=12`)
          const data = await res.json()
          return { category: cat, products: data.products || [] }
        })
      )
      categoryProducts.value = productResults.filter(item => item.products.length > 0)
      categoryProducts.value.forEach((_, i) => {
        scrollStates.categories[`cat-${i}`] = { canLeft: false, canRight: true }
      })
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
        <div class="section-header">
          <h2 class="section-title">{{ config.home_featured_title || '特别推荐' }}</h2>
          <router-link to="/products?featured=true" class="view-more">
            查看更多 <span class="arrow">&rarr;</span>
          </router-link>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="products-scroll-wrapper">
          <div class="products-scroll" ref="featuredScroll" @scroll="updateScrollState($event, 'featured')">
            <div class="product-scroll-item" v-for="product in featuredProducts" :key="product.id">
              <ProductCard :product="product" />
            </div>
          </div>
          <button class="scroll-btn left" v-show="scrollStates.featured.canLeft" @click="scrollSection(featuredScroll, 'left')">&#10094;</button>
          <button class="scroll-btn right" v-show="scrollStates.featured.canRight" @click="scrollSection(featuredScroll, 'right')">&#10095;</button>
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
        <div v-else class="products-scroll-wrapper">
          <div class="products-scroll" ref="newScroll" @scroll="updateScrollState($event, 'new')">
            <div class="product-scroll-item" v-for="product in newProducts" :key="'new-' + product.id">
              <ProductCard :product="product" />
            </div>
          </div>
          <button class="scroll-btn left" v-show="scrollStates.new.canLeft" @click="scrollSection(newScroll, 'left')">&#10094;</button>
          <button class="scroll-btn right" v-show="scrollStates.new.canRight" @click="scrollSection(newScroll, 'right')">&#10095;</button>
        </div>
      </div>
    </section>

    <section class="banner-section" v-if="config.home_banner_1">
      <div class="container">
        <img :src="config.home_banner_1" alt="Banner" class="full-width-banner">
      </div>
    </section>

    <section
      v-for="(item, index) in categoryProducts"
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
        <div class="products-scroll-wrapper">
          <div class="products-scroll" :ref="el => categoryScrolls[index] = el" @scroll="updateScrollState($event, 'cat-' + index)">
            <div class="product-scroll-item" v-for="product in item.products" :key="product.id">
              <ProductCard :product="product" />
            </div>
          </div>
          <button class="scroll-btn left" v-show="scrollStates.categories['cat-' + index]?.canLeft" @click="scrollSection(categoryScrolls[index], 'left')">&#10094;</button>
          <button class="scroll-btn right" v-show="scrollStates.categories['cat-' + index]?.canRight" @click="scrollSection(categoryScrolls[index], 'right')">&#10095;</button>
        </div>
      </div>
    </section>

    <section class="top-selling-section" v-if="topSellingProducts.length > 0">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ config.home_topselling_title || '热销排行' }}</h2>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="products-scroll-wrapper">
          <div class="products-scroll" ref="topSellingScroll" @scroll="updateScrollState($event, 'topSelling')">
            <div class="product-scroll-item" v-for="product in topSellingProducts" :key="'top-' + product.id">
              <ProductCard :product="product" />
            </div>
          </div>
          <button class="scroll-btn left" v-show="scrollStates.topSelling.canLeft" @click="scrollSection(topSellingScroll, 'left')">&#10094;</button>
          <button class="scroll-btn right" v-show="scrollStates.topSelling.canRight" @click="scrollSection(topSellingScroll, 'right')">&#10095;</button>
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
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#d4a574" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" y1="9" x2="9.01" y2="9"/><line x1="15" y1="9" x2="15.01" y2="9"/></svg>
            </div>
            <h3>品质甄选</h3>
            <p>严格筛选 匠心之选</p>
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

.products-scroll-wrapper {
  position: relative;
}

.products-scroll {
  display: flex;
  gap: 15px;
  overflow-x: auto;
  scroll-snap-type: x mandatory;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
  padding-bottom: 5px;
}

.products-scroll::-webkit-scrollbar {
  display: none;
}

.products-scroll {
  scrollbar-width: none;
}

.product-scroll-item {
  flex: 0 0 calc((100% - 60px) / 5);
  scroll-snap-align: start;
  min-width: calc((100% - 60px) / 5);
}

.scroll-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.7);
  color: #d4a574;
  border: 1px solid rgba(212, 165, 116, 0.3);
  cursor: pointer;
  z-index: 5;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.scroll-btn:hover {
  background: rgba(0, 0, 0, 0.9);
  border-color: #d4a574;
}

.scroll-btn.left {
  left: -5px;
}

.scroll-btn.right {
  right: -5px;
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

.loading {
  text-align: center;
  color: #d4a574;
  padding: 40px;
}

@media (max-width: 1200px) {
  .product-scroll-item {
    flex: 0 0 calc((100% - 30px) / 3);
    min-width: calc((100% - 30px) / 3);
  }
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

  .product-scroll-item {
    flex: 0 0 calc((100% - 15px) / 2);
    min-width: calc((100% - 15px) / 2);
  }

  .scroll-btn {
    display: none;
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
}
</style>
