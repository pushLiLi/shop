<script setup>
import { computed, onMounted, ref } from 'vue'
import { useFavoritesStore } from '../stores/favorites'
import { useCartStore } from '../stores/cart'
import { useToastStore } from '../stores/toast'

const favoritesStore = useFavoritesStore()
const cartStore = useCartStore()
const toast = useToastStore()

const searchQuery = ref('')
const sortBy = ref('recent')
const selectedIds = ref(new Set())

const items = computed(() => favoritesStore.items)
const loading = computed(() => favoritesStore.loading)

const sortedAndFilteredItems = computed(() => {
  let result = [...items.value]
  
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(fav => 
      fav.product?.name?.toLowerCase().includes(query)
    )
  }
  
  switch (sortBy.value) {
    case 'price-asc':
      result.sort((a, b) => (a.product?.price || 0) - (b.product?.price || 0))
      break
    case 'price-desc':
      result.sort((a, b) => (b.product?.price || 0) - (a.product?.price || 0))
      break
    case 'name':
      result.sort((a, b) => (a.product?.name || '').localeCompare(b.product?.name || ''))
      break
    default:
      break
  }
  
  return result
})

const isAllSelected = computed(() => {
  if (sortedAndFilteredItems.value.length === 0) return false
  return sortedAndFilteredItems.value.every(fav => selectedIds.value.has(fav.productId))
})

const selectedCount = computed(() => selectedIds.value.size)

onMounted(() => {
  favoritesStore.fetchFavorites()
})

function formatPrice(price) {
  return `$${Number(price).toFixed(2)}`
}

function toggleSelect(productId) {
  if (selectedIds.value.has(productId)) {
    selectedIds.value.delete(productId)
  } else {
    selectedIds.value.add(productId)
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value.clear()
  } else {
    sortedAndFilteredItems.value.forEach(fav => {
      selectedIds.value.add(fav.productId)
    })
  }
}

async function addToCart(product) {
  await cartStore.addItem(product, 1)
}

async function batchAddToCart() {
  if (selectedCount.value === 0) {
    toast.warning('请先选择商品')
    return
  }
  
  const selectedItems = sortedAndFilteredItems.value.filter(fav => 
    selectedIds.value.has(fav.productId)
  )
  
  for (const fav of selectedItems) {
    if (fav.product) {
      await cartStore.addItem(fav.product, 1)
    }
  }
  
  toast.success(`已添加 ${selectedCount.value} 件商品到购物车`)
  selectedIds.value.clear()
}

async function removeFavorite(productId) {
  selectedIds.value.delete(productId)
  await favoritesStore.removeItem(productId)
  toast.success('已取消收藏')
}
</script>

<template>
  <div class="favorites-page">
    <div class="container">
      <h1 class="page-title">我的收藏</h1>
      
      <div v-if="loading" class="loading">加载中...</div>
      
      <div v-else-if="items.length === 0" class="empty-favorites">
        <p>暂无收藏</p>
        <router-link to="/" class="continue-shopping">去购物</router-link>
      </div>
      
      <template v-else>
        <div class="toolbar">
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="搜索收藏商品..." 
            class="search-input"
          >
          
          <div class="toolbar-actions">
            <select v-model="sortBy" class="sort-select">
              <option value="recent">最近添加</option>
              <option value="price-asc">价格从低到高</option>
              <option value="price-desc">价格从高到低</option>
              <option value="name">商品名称 A-Z</option>
            </select>
            
            <label class="select-all">
              <input 
                type="checkbox" 
                :checked="isAllSelected" 
                @change="toggleSelectAll"
              >
              <span>全选</span>
            </label>
            
            <button 
              class="batch-cart-btn" 
              @click="batchAddToCart"
              :disabled="selectedCount === 0"
            >
              加入购物车{{ selectedCount > 0 ? ` (${selectedCount})` : '' }}
            </button>
          </div>
        </div>
        
        <div v-if="sortedAndFilteredItems.length === 0" class="empty-search">
          <p>未找到相关收藏</p>
        </div>
        
        <div v-else class="favorites-list">
          <div 
            v-for="fav in sortedAndFilteredItems" 
            :key="fav.productId" 
            class="favorite-item"
            :class="{ selected: selectedIds.has(fav.productId) }"
            @click="toggleSelect(fav.productId)"
          >
            <div class="checkbox" @click.stop>
              <input 
                type="checkbox" 
                :checked="selectedIds.has(fav.productId)"
                @change="toggleSelect(fav.productId)"
              >
            </div>
            
            <router-link 
              :to="'/products/' + fav.product?.id" 
              class="item-image"
              @click.stop
            >
              <img :src="fav.product?.imageUrl" :alt="fav.product?.name">
            </router-link>
            
            <router-link 
              :to="'/products/' + fav.product?.id" 
              class="item-name"
              @click.stop
            >
              {{ fav.product?.name }}
            </router-link>
            
            <div class="item-price">{{ formatPrice(fav.product?.price) }}</div>
            
            <button 
              class="cart-btn" 
              @click.prevent.stop="addToCart(fav.product)" 
              title="加入购物车"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="9" cy="21" r="1"></circle>
                <circle cx="20" cy="21" r="1"></circle>
                <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path>
              </svg>
            </button>
            
            <button 
              class="remove-btn" 
              @click.prevent.stop="removeFavorite(fav.productId)" 
              title="取消收藏"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" class="heart-icon">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
              </svg>
            </button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.favorites-page {
  background: #0f0f0f;
  min-height: 100vh;
  padding: 40px 0 60px;
}

.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-title {
  color: #d4a574;
  font-size: 28px;
  margin-bottom: 20px;
  border-bottom: 2px solid #d4a574;
  padding-bottom: 10px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  max-width: 300px;
  padding: 10px 15px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  color: #fff;
  font-size: 14px;
}

.search-input:focus {
  outline: none;
  border-color: #d4a574;
}

.search-input::placeholder {
  color: #666;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.sort-select {
  padding: 8px 12px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  color: #ccc;
  font-size: 14px;
  cursor: pointer;
}

.sort-select:focus {
  outline: none;
  border-color: #d4a574;
}

.select-all {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #888;
  font-size: 14px;
  cursor: pointer;
}

.select-all input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.batch-cart-btn {
  padding: 8px 16px;
  background: #d4a574;
  border: none;
  border-radius: 6px;
  color: #1a1a1a;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.batch-cart-btn:hover:not(:disabled) {
  background: #e0b585;
}

.batch-cart-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading, .empty-favorites, .empty-search {
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

.favorites-list {
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
}

.favorite-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #2a2a2a;
  gap: 15px;
  cursor: pointer;
  transition: background 0.2s;
}

.favorite-item:last-child {
  border-bottom: none;
}

.favorite-item:hover {
  background: #252525;
}

.favorite-item.selected {
  background: #2a2520;
}

.checkbox {
  flex-shrink: 0;
}

.checkbox input {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.item-image {
  flex-shrink: 0;
  width: 60px;
  height: 60px;
  background: #fff;
  border-radius: 6px;
  overflow: hidden;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-name {
  flex: 1;
  color: #ccc;
  text-decoration: none;
  font-size: 14px;
  line-height: 1.4;
}

.item-name:hover {
  color: #d4a574;
}

.item-price {
  color: #d4a574;
  font-weight: bold;
  font-size: 14px;
  white-space: nowrap;
  min-width: 70px;
  text-align: right;
}

.cart-btn {
  background: transparent;
  border: 1px solid #555;
  padding: 8px;
  cursor: pointer;
  color: #888;
  border-radius: 6px;
  transition: all 0.2s;
}

.cart-btn:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.remove-btn {
  background: transparent;
  border: none;
  padding: 8px;
  cursor: pointer;
  color: #d4a574;
  transition: transform 0.2s;
}

.remove-btn:hover {
  transform: scale(1.15);
}

.heart-icon {
  fill: #d4a574;
  transition: fill 0.2s;
}

.remove-btn:hover .heart-icon {
  fill: none;
}

@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-input {
    max-width: none;
  }
  
  .toolbar-actions {
    justify-content: space-between;
  }
  
  .favorite-item {
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .checkbox {
    order: 1;
  }
  
  .item-image {
    order: 2;
  }
  
  .item-name {
    order: 3;
    width: calc(100% - 100px);
  }
  
  .item-price {
    order: 4;
    width: calc(100% - 100px);
    text-align: left;
    margin-left: 90px;
    margin-top: -5px;
  }
  
  .cart-btn {
    order: 5;
    margin-left: auto;
  }
  
  .remove-btn {
    order: 6;
  }
}
</style>
