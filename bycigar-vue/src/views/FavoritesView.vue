<script setup>
import { computed, onMounted } from 'vue'
import { useFavoritesStore } from '../stores/favorites'
import ProductCard from '../components/ProductCard.vue'

const favoritesStore = useFavoritesStore()

const items = computed(() => favoritesStore.items)
const loading = computed(() => favoritesStore.loading)

onMounted(() => {
  favoritesStore.fetchFavorites()
})
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
      
      <div v-else>
        <div class="favorites-grid">
          <div v-for="fav in items" :key="fav.productId" class="favorite-item">
            <ProductCard :product="fav.product" />
            <button class="remove-btn" @click="favoritesStore.removeItem(fav.productId)">
              移除
            </button>
          </div>
        </div>
      </div>
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
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-title {
  color: #d4a574;
  font-size: 28px;
  margin-bottom: 30px;
  border-bottom: 2px solid #d4a574;
  padding-bottom: 10px;
}

.loading, .empty-favorites {
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

.favorites-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.favorite-item {
  position: relative;
}

.remove-btn {
  position: absolute;
  top: 50px;
  right: 10px;
  background: rgba(0,0,0,0.7);
  color: #e74;
  border: 1px solid #e74;
  padding: 5px 10px;
  font-size: 12px;
  cursor: pointer;
  border-radius: 4px;
  opacity: 0;
  transition: opacity 0.3s;
  z-index: 10;
}

.favorite-item:hover .remove-btn {
  opacity: 1;
}

.remove-btn:hover {
  background: #e74;
  color: #1a1a1a;
}

@media (max-width: 992px) {
  .favorites-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .favorites-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .favorites-grid {
    grid-template-columns: 1fr;
  }
}
</style>
