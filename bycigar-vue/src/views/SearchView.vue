<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import ProductCard from '../components/ProductCard.vue'

const route = useRoute()

const isCompact = ref(window.innerWidth <= 992)

function onResize() {
  isCompact.value = window.innerWidth <= 992
}

onMounted(() => {
  window.addEventListener('resize', onResize)
  fetchCategories()
})
onUnmounted(() => window.removeEventListener('resize', onResize))
const products = ref([])
const loading = ref(false)
const error = ref(null)
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const sortBy = ref('createdAt')
const sortOrder = ref('desc')
const selectedCategorySlug = ref('')
const minPrice = ref('')
const maxPrice = ref('')
const categories = ref([])
const showMobileFilters = ref(false)

const searchQuery = ref('')

async function fetchCategories() {
  try {
    const res = await fetch('/api/categories')
    const data = await res.json()
    categories.value = Array.isArray(data) ? data : (data.data || [])
  } catch (e) {
    console.error('获取分类失败:', e)
  }
}

function getCategoryName(slug) {
  for (const cat of categories.value) {
    if (cat.slug === slug) return cat.name
    if (cat.children) {
      for (const child of cat.children) {
        if (child.slug === slug) return child.name
      }
    }
  }
  return slug
}

async function searchProducts() {
  if (!searchQuery.value.trim()) {
    products.value = []
    totalCount.value = 0
    return
  }
  
  try {
    loading.value = true
    error.value = null
    
    const params = new URLSearchParams({
      search: searchQuery.value,
      page: currentPage.value,
      limit: pageSize.value,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value
    })
    
    if (selectedCategorySlug.value) {
      params.append('category', selectedCategorySlug.value)
    }
    if (minPrice.value) {
      params.append('minPrice', minPrice.value)
    }
    if (maxPrice.value) {
      params.append('maxPrice', maxPrice.value)
    }
    
    const res = await fetch(`/api/products?${params}`)
    const data = await res.json()
    
    products.value = data.products || []
    totalCount.value = data.total || 0
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

watch(() => route.query.q, (newQuery) => {
  searchQuery.value = newQuery || ''
  selectedCategorySlug.value = ''
  minPrice.value = ''
  maxPrice.value = ''
  currentPage.value = 1
  searchProducts()
}, { immediate: true })

function applyPriceFilter() {
  currentPage.value = 1
  searchProducts()
}

function clearPriceFilter() {
  minPrice.value = ''
  maxPrice.value = ''
  currentPage.value = 1
  searchProducts()
}

function handleCategorySelect(slug) {
  if (selectedCategorySlug.value === slug) {
    selectedCategorySlug.value = ''
  } else {
    selectedCategorySlug.value = slug
  }
  currentPage.value = 1
  searchProducts()
}

function changePage(page) {
  currentPage.value = page
  searchProducts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function changeSort(newSortBy) {
  if (sortBy.value === newSortBy) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = newSortBy
    sortOrder.value = 'desc'
  }
  currentPage.value = 1
  searchProducts()
}

const totalPages = ref(1)
watch([totalCount, pageSize], () => {
  totalPages.value = Math.ceil(totalCount.value / pageSize.value)
})
</script>

<template>
  <div class="search-page">
    <div class="container">
      <div class="search-header">
        <h1 class="page-title">搜索结果</h1>
        <p class="search-info" v-if="searchQuery">
          关键词: "<span class="keyword">{{ searchQuery }}</span>" 
          <span v-if="!loading">共 {{ totalCount }} 个结果</span>
        </p>
      </div>

      <div class="search-layout" v-if="searchQuery">
        <aside class="search-sidebar">
          <div class="sidebar-section" v-if="!isCompact">
            <div class="sidebar-title">分类筛选</div>
            <ul class="filter-category-list">
              <li>
                <button
                  class="filter-category-link"
                  :class="{ active: !selectedCategorySlug }"
                  @click="selectedCategorySlug = ''; currentPage = 1; searchProducts()"
                >全部分类</button>
              </li>
              <li v-for="cat in categories" :key="cat.id">
                <button
                  class="filter-category-link"
                  :class="{ active: selectedCategorySlug === cat.slug }"
                  @click="handleCategorySelect(cat.slug)"
                >{{ cat.name }}</button>
                <ul v-if="cat.children && cat.children.length" class="filter-subcategory-list">
                  <li v-for="child in cat.children" :key="child.id">
                    <button
                      class="filter-category-link sub"
                      :class="{ active: selectedCategorySlug === child.slug }"
                      @click="handleCategorySelect(child.slug)"
                    >{{ child.name }}</button>
                  </li>
                </ul>
              </li>
            </ul>
          </div>
          <div class="sidebar-section" v-if="!isCompact">
            <div class="sidebar-section-title">价格区间</div>
            <div class="price-filter">
              <input
                v-model="minPrice"
                type="number"
                class="price-input"
                placeholder="最低价"
                min="0"
                @keydown.enter="applyPriceFilter"
              >
              <span class="price-sep">—</span>
              <input
                v-model="maxPrice"
                type="number"
                class="price-input"
                placeholder="最高价"
                min="0"
                @keydown.enter="applyPriceFilter"
              >
            </div>
            <div class="price-actions">
              <button class="price-apply-btn" @click="applyPriceFilter">确定</button>
              <button class="price-clear-btn" @click="clearPriceFilter" v-if="minPrice || maxPrice">清除</button>
            </div>
          </div>
          <div class="sidebar-mobile-filters" v-if="isCompact">
            <button class="mobile-filter-toggle" @click="showMobileFilters = !showMobileFilters">
              <span>筛选</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"></polyline></svg>
            </button>
            <Transition name="filter-slide">
              <div v-if="showMobileFilters" class="mobile-filter-content">
                <div class="mobile-filter-group">
                  <span class="mobile-filter-label">分类筛选</span>
                  <div class="mobile-category-chips">
                    <button
                      class="category-chip"
                      :class="{ active: !selectedCategorySlug }"
                      @click="selectedCategorySlug = ''; currentPage = 1; searchProducts()"
                    >全部</button>
                    <button
                      v-for="cat in categories"
                      :key="cat.id"
                      class="category-chip"
                      :class="{ active: selectedCategorySlug === cat.slug }"
                      @click="handleCategorySelect(cat.slug)"
                    >{{ cat.name }}</button>
                  </div>
                </div>
                <div class="mobile-filter-group">
                  <span class="mobile-filter-label">价格区间</span>
                  <div class="price-filter">
                    <input
                      v-model="minPrice"
                      type="number"
                      class="price-input"
                      placeholder="最低价"
                      min="0"
                      @keydown.enter="applyPriceFilter"
                    >
                    <span class="price-sep">—</span>
                    <input
                      v-model="maxPrice"
                      type="number"
                      class="price-input"
                      placeholder="最高价"
                      min="0"
                      @keydown.enter="applyPriceFilter"
                    >
                  </div>
                  <div class="price-actions">
                    <button class="price-apply-btn" @click="applyPriceFilter">确定</button>
                    <button class="price-clear-btn" @click="clearPriceFilter" v-if="minPrice || maxPrice">清除</button>
                  </div>
                </div>
              </div>
            </Transition>
          </div>
        </aside>

        <div class="search-main">
          <div class="search-controls" v-if="products.length > 0 || selectedCategorySlug || minPrice || maxPrice">
            <div class="sort-options">
              <span class="sort-label">排序:</span>
              <button 
                class="sort-btn" 
                :class="{ active: sortBy === 'createdAt' }"
                @click="changeSort('createdAt')"
              >最新<span v-if="sortBy === 'createdAt'" class="sort-arrow">{{ sortOrder === 'asc' ? ' ↑' : ' ↓' }}</span></button>
              <button 
                class="sort-btn" 
                :class="{ active: sortBy === 'price' }"
                @click="changeSort('price')"
              >价格<span v-if="sortBy === 'price'" class="sort-arrow">{{ sortOrder === 'asc' ? ' ↑' : ' ↓' }}</span></button>
              <button 
                class="sort-btn" 
                :class="{ active: sortBy === 'name' }"
                @click="changeSort('name')"
              >名称<span v-if="sortBy === 'name'" class="sort-arrow">{{ sortOrder === 'asc' ? ' ↑' : ' ↓' }}</span></button>
            </div>
            <div class="active-filters" v-if="selectedCategorySlug || minPrice || maxPrice">
              <span class="filter-tag" v-if="selectedCategorySlug">
                {{ getCategoryName(selectedCategorySlug) }}
                <button class="filter-tag-close" @click="selectedCategorySlug = ''; currentPage = 1; searchProducts()">&times;</button>
              </span>
              <span class="filter-tag" v-if="minPrice || maxPrice">
                {{ minPrice ? '¥' + minPrice : '¥0' }} - {{ maxPrice ? '¥' + maxPrice : '¥∞' }}
                <button class="filter-tag-close" @click="clearPriceFilter">&times;</button>
              </span>
            </div>
          </div>

          <div v-if="loading" class="loading">搜索中...</div>
          
          <div v-else-if="products.length === 0" class="no-results">
            <p>未找到相关产品</p>
            <p class="hint">请尝试其他关键词或调整筛选条件</p>
          </div>
          
          <div v-else>
            <div class="products-grid">
              <ProductCard v-for="product in products" :key="product.id" :product="product" :horizontal="isCompact" />
            </div>

            <div class="pagination" v-if="totalPages > 1">
              <button 
                class="page-btn" 
                :disabled="currentPage === 1"
                @click="changePage(currentPage - 1)"
              >上一页</button>
              
              <button 
                v-for="page in totalPages" 
                :key="page"
                class="page-btn"
                :class="{ active: page === currentPage }"
                @click="changePage(page)"
              >{{ page }}</button>
              
              <button 
                class="page-btn" 
                :disabled="currentPage === totalPages"
                @click="changePage(currentPage + 1)"
              >下一页</button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <p>请输入搜索关键词</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 15px;
}

.search-header {
  margin-bottom: 30px;
}

.page-title {
  color: #fff;
  font-size: 32px;
  margin: 0 0 10px;
}

.search-info {
  color: #888;
  font-size: 16px;
}

.keyword {
  color: #d4a574;
}

.search-layout {
  display: flex;
  gap: 30px;
  align-items: flex-start;
}

.search-sidebar {
  flex-shrink: 0;
}

.sidebar-section {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
  min-width: 220px;
  max-width: 240px;
  margin-bottom: 15px;
}

.sidebar-title {
  color: #d4a574;
  font-size: 16px;
  font-weight: 600;
  border-bottom: 1px solid #333;
  padding-bottom: 12px;
  margin-bottom: 12px;
}

.sidebar-section-title {
  color: #d4a574;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

.filter-category-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.filter-category-list li {
  margin-bottom: 2px;
}

.filter-category-link {
  display: block;
  width: 100%;
  text-align: left;
  padding: 8px 12px;
  color: #ccc;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.filter-category-link:hover {
  color: #d4a574;
  background: rgba(212, 165, 116, 0.08);
}

.filter-category-link.active {
  color: #d4a574;
  background: rgba(212, 165, 116, 0.15);
  font-weight: 500;
}

.filter-subcategory-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.filter-subcategory-list li {
  margin-bottom: 2px;
}

.filter-category-link.sub {
  padding-left: 28px;
  font-size: 13px;
  color: #999;
}

.search-main {
  flex: 1;
  min-width: 0;
}

.price-filter {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.price-input {
  width: 100%;
  background: #2d2d2d;
  border: 1px solid #333;
  border-radius: 4px;
  padding: 8px 10px;
  color: #fff;
  font-size: 13px;
  outline: none;
}

.price-input:focus {
  border-color: #d4a574;
}

.price-input::-webkit-inner-spin-button,
.price-input::-webkit-outer-spin-button {
  -webkit-appearance: none;
}

.price-sep {
  color: #555;
  flex-shrink: 0;
}

.price-actions {
  display: flex;
  gap: 8px;
}

.price-apply-btn {
  flex: 1;
  background: #d4a574;
  color: #1a1a1a;
  border: none;
  padding: 7px 12px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  font-weight: 500;
  transition: opacity 0.2s;
}

.price-apply-btn:hover {
  opacity: 0.9;
}

.price-clear-btn {
  background: transparent;
  border: 1px solid #444;
  color: #888;
  padding: 7px 12px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.price-clear-btn:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.mobile-filter-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  color: #d4a574;
  padding: 12px 16px;
  width: 100%;
  cursor: pointer;
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 15px;
}

.mobile-filter-toggle svg {
  transition: transform 0.2s;
}

.mobile-filter-content {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 15px;
}

.mobile-filter-group {
  margin-bottom: 16px;
}

.mobile-filter-group:last-child {
  margin-bottom: 0;
}

.mobile-filter-label {
  display: block;
  color: #d4a574;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
}

.mobile-category-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-chip {
  background: #2d2d2d;
  border: 1px solid #333;
  color: #ccc;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-chip:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.category-chip.active {
  background: #d4a574;
  border-color: #d4a574;
  color: #1a1a1a;
}

.active-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.filter-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(212, 165, 116, 0.15);
  border: 1px solid rgba(212, 165, 116, 0.3);
  color: #d4a574;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 13px;
}

.filter-tag-close {
  background: transparent;
  border: none;
  color: #d4a574;
  cursor: pointer;
  padding: 0 2px;
  font-size: 16px;
  line-height: 1;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.filter-tag-close:hover {
  opacity: 1;
}

.search-controls {
  margin-bottom: 30px;
  padding: 15px 0;
  border-bottom: 1px solid #2a2a2a;
}

.sort-options {
  display: flex;
  align-items: center;
  gap: 10px;
}

.sort-label {
  color: #888;
  font-size: 14px;
}

.sort-btn {
  background: transparent;
  border: 1px solid #333;
  color: #888;
  padding: 8px 16px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.sort-btn:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.sort-btn.active {
  background: #d4a574;
  border-color: #d4a574;
  color: #1a1a1a;
}

.sort-arrow {
  font-size: 12px;
  margin-left: 2px;
}

.loading, .no-results, .empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #888;
  font-size: 18px;
}

.no-results .hint {
  color: #555;
  font-size: 14px;
  margin-top: 10px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
  margin-bottom: 40px;
}

.pagination {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.page-btn {
  background: #1a1a1a;
  border: 1px solid #333;
  color: #888;
  padding: 10px 16px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  border-color: #d4a574;
  color: #d4a574;
}

.page-btn.active {
  background: #d4a574;
  border-color: #d4a574;
  color: #1a1a1a;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.filter-slide-enter-active,
.filter-slide-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.filter-slide-enter-from,
.filter-slide-leave-to {
  opacity: 0;
  max-height: 0;
  margin-bottom: 0;
}

.filter-slide-enter-to,
.filter-slide-leave-from {
  opacity: 1;
  max-height: 400px;
}

@media (max-width: 992px) {
  .search-layout {
    flex-direction: column;
  }

  .sidebar-section {
    display: none;
  }

  .products-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .search-page {
    padding: 20px 0 40px;
  }

  .page-title {
    font-size: 24px;
  }
}

@media (max-width: 768px) {
  .products-grid {
    grid-template-columns: 1fr;
    gap: 10px;
  }

  .search-page {
    padding: 15px 0 30px;
  }

  .page-title {
    font-size: 22px;
  }

  .sort-options {
    flex-wrap: wrap;
    gap: 6px;
  }

  .sort-btn {
    padding: 6px 10px;
    font-size: 13px;
  }
}
</style>
