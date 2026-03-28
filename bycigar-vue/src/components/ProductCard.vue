<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { useFavoritesStore } from '../stores/favorites'
import { useToastStore } from '../stores/toast'

const emit = defineEmits(['addToCart'])

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()
const favoritesStore = useFavoritesStore()
const toast = useToastStore()

const isFavorite = computed(() => {
  return favoritesStore.items.some(item => item.productId === props.product.id)
})

const formatPrice = (price) => {
  return `$${Number(price).toFixed(2)}`
}

async function addToCartHandler(e) {
  e.preventDefault()
  e.stopPropagation()
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }
  await cartStore.addItem(props.product, 1)
  toast.success('已添加到购物车')
}

async function toggleFavorite(e) {
  e.preventDefault()
  e.stopPropagation()
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }
  if (isFavorite.value) {
    await favoritesStore.removeItem(props.product.id)
  } else {
    await favoritesStore.addItem(props.product)
  }
}
</script>

<template>
  <div class="product-card">
    <router-link :to="'/products/' + product.id" class="product-image">
      <img :src="product.imageUrl" :alt="product.name">
      <button class="favorite-btn" @click="toggleFavorite" :class="{ active: isFavorite }">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" :fill="isFavorite ? '#d4a574' : 'none'" stroke="currentColor" stroke-width="2">
          <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
        </svg>
      </button>
    </router-link>
    <div class="product-info">
      <h3 class="product-name">
        <router-link :to="'/products/' + product.id">{{ product.name }}</router-link>
      </h3>
      <div class="product-bottom">
        <button class="add-cart-btn" @click="addToCartHandler">加入购物车</button>
        <div class="product-price">{{ formatPrice(product.price) }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
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
  position: relative;
}

.favorite-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(0,0,0,0.5);
  border: none;
  padding: 8px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s;
  color: #fff;
}

.favorite-btn:hover {
  background: rgba(0,0,0,0.8);
}

.favorite-btn.active {
  color: #d4a574;
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

@media (max-width: 768px) {
  .product-card {
    display: flex;
    flex-direction: row;
    background: #1a1a1a;
    border-radius: 8px;
    overflow: hidden;
  }

  .product-card:hover {
    transform: none;
    box-shadow: none;
  }

  .product-image {
    flex-shrink: 0;
    width: 140px;
    height: 140px;
    padding: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .product-image img {
    width: 100%;
    height: 100%;
    aspect-ratio: unset;
    object-fit: cover;
  }

  .favorite-btn {
    top: 6px;
    right: 6px;
    padding: 5px;
  }

  .product-info {
    flex: 1;
    padding: 12px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-width: 0;
  }

  .product-name {
    font-size: 14px;
    margin: 0 0 8px;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .product-bottom {
    flex-wrap: wrap;
    gap: 8px;
  }

  .add-cart-btn {
    padding: 5px 10px;
    font-size: 11px;
  }

  .product-price {
    font-size: 15px;
  }
}
</style>
