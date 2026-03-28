<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ProductCard from '../components/ProductCard.vue'
import CategorySidebar from '../components/CategorySidebar.vue'

const route = useRoute()
const products = ref([])
const category = ref(null)
const loading = ref(true)
const error = ref(null)
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const sortBy = ref('createdAt')
const sortOrder = ref('desc')

const categorySlug = ref('')
const categoryName = ref('')

async function fetchProducts() {
  try {
    loading.value = true
    error.value = null
    
    const params = new URLSearchParams({
      page: currentPage.value,
      limit: pageSize.value,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value
    })
    
    if (categorySlug.value) {
      params.append('category', categorySlug.value)
    }
    
    const res = await fetch(`http://localhost:3000/api/products?${params}`)
    const data = await res.json()
    
    products.value = data.products || []
    totalCount.value = data.total || 0
    category.value = data.category || null
    categoryName.value = category.value?.name || categorySlug.value
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

watch(() => route.params.slug, (newSlug) => {
  categorySlug.value = newSlug || ''
  categoryName.value = '全部商品'
  currentPage.value = 1
  fetchProducts()
}, { immediate: true })

function changePage(page) {
  currentPage.value = page
  fetchProducts()
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
  fetchProducts()
}

const totalPages = ref(1)
watch([totalCount, pageSize], () => {
  totalPages.value = Math.ceil(totalCount.value / pageSize.value)
})
</script>

<template>
  <div class="category-page">
    <div class="container">
      <div class="category-layout">
        <CategorySidebar :activeSlug="categorySlug" />
        <div class="category-main">
          <div class="category-header">
            <h1 class="page-title">{{ categoryName }}</h1>
            <p class="product-count" v-if="!loading">共 {{ totalCount }} 个产品</p>
          </div>

          <div class="category-controls" v-if="products.length > 0">
            <div class="sort-options">
              <span class="sort-label">排序:</span>
              <button
                class="sort-btn"
                :class="{ active: sortBy === 'createdAt' }"
                @click="changeSort('createdAt')"
              >最新</button>
              <button
                class="sort-btn"
                :class="{ active: sortBy === 'price' }"
                @click="changeSort('price')"
              >价格</button>
              <button
                class="sort-btn"
                :class="{ active: sortBy === 'name' }"
                @click="changeSort('name')"
              >名称</button>
            </div>
          </div>

          <div v-if="loading" class="loading">加载中...</div>

          <div v-else-if="error" class="error">
            <p>{{ error }}</p>
          </div>

          <div v-else-if="products.length === 0" class="no-products">
            <p>该分类暂无产品</p>
          </div>

          <div v-else>
            <div class="products-grid">
              <ProductCard v-for="product in products" :key="product.id" :product="product" />
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
    </div>
  </div>
</template>

<style scoped>
.category-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.category-layout {
  display: flex;
  gap: 30px;
  align-items: flex-start;
}

.category-main {
  flex: 1;
  min-width: 0;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 15px;
}

.category-header {
  margin-bottom: 30px;
  border-bottom: 2px solid #d4a574;
  padding-bottom: 20px;
}

.page-title {
  color: #d4a574;
  font-size: 32px;
  margin: 0 0 10px;
}

.product-count {
  color: #888;
  font-size: 14px;
}

.category-controls {
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

.loading, .error, .no-products {
  text-align: center;
  padding: 80px 20px;
  color: #888;
  font-size: 18px;
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

@media (max-width: 992px) {
  .products-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .products-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .category-layout {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .products-grid {
    grid-template-columns: 1fr;
  }
}
</style>
