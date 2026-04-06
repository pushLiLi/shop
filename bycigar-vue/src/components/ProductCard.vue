<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { useFavoritesStore } from '../stores/favorites'
import { useToastStore } from '../stores/toast'
import { formatPrice } from '../composables/useFormatPrice'
import { useShare } from '../composables/useShare'

const emit = defineEmits(['addToCart'])

const props = defineProps({
  product: {
    type: Object,
    required: true
  },
  horizontal: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()
const favoritesStore = useFavoritesStore()
const toast = useToastStore()
const { shareProduct } = useShare()

const isFavorite = computed(() => {
  return favoritesStore.items.some(item => item.productId === props.product.id)
})

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

async function shareHandler(e) {
  e.preventDefault()
  e.stopPropagation()
  await shareProduct(props.product)
}
</script>

<template>
  <div class="product-card" :class="{ horizontal: horizontal }" @click="router.push('/products/' + product.id)">
    <router-link :to="'/products/' + product.id" class="product-image">
      <img :src="product.thumbnailUrl || product.imageUrl" :alt="product.name" loading="lazy">
      <div v-if="product.stock === 0" class="sold-out-overlay">
        <span class="sold-out-text">已售罄</span>
      </div>
      <button class="favorite-btn" @click="toggleFavorite" :class="{ active: isFavorite }">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" :fill="isFavorite ? '#d4a574' : 'none'" stroke="currentColor" stroke-width="2">
          <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
        </svg>
      </button>
      <button class="share-btn" @click="shareHandler">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="18" cy="5" r="3"></circle>
          <circle cx="6" cy="12" r="3"></circle>
          <circle cx="18" cy="19" r="3"></circle>
          <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
          <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
        </svg>
      </button>
    </router-link>
    <div class="product-info">
      <h3 class="product-name">
        <router-link :to="'/products/' + product.id">{{ product.name }}</router-link>
      </h3>
      <div class="product-bottom">
        <button v-if="product.stock > 0" class="add-cart-btn" @click="addToCartHandler">
          <svg class="cart-icon" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>
          <span class="cart-text">加入购物车</span>
        </button>
        <span v-else class="sold-out-tag">已售罄</span>
        <div class="product-price">{{ formatPrice(product) }}</div>
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
  cursor: pointer;
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

.share-btn {
  position: absolute;
  top: 10px;
  left: 10px;
  background: rgba(0,0,0,0.5);
  border: none;
  padding: 8px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s;
  color: #fff;
}

.share-btn:hover {
  background: rgba(0,0,0,0.8);
  color: #d4a574;
}

.product-image img {
  width: 100%;
  height: auto;
  aspect-ratio: 1;
  object-fit: cover;
}

.sold-out-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
  pointer-events: none;
}

.sold-out-text {
  display: inline-block;
  padding: 8px 24px;
  background: #e53e3e;
  border-radius: 4px;
  color: #fff;
  font-size: 22px;
  font-weight: 800;
  letter-spacing: 6px;
  transform: rotate(-15deg);
  box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.4);
  user-select: none;
}

.sold-out-tag {
  color: #e53e3e;
  font-size: 12px;
  font-weight: 600;
  padding: 6px 12px;
  border: 1px solid #e53e3e;
  border-radius: 4px;
  white-space: nowrap;
}

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 13px;
  margin: 0 0 10px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
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
  gap: 8px;
  flex-wrap: wrap;
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
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.add-cart-btn .cart-icon {
  display: none;
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

.product-card.horizontal {
  display: flex;
  flex-direction: row;
  width: 100%;
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
}

.product-card.horizontal:hover {
  transform: none;
  box-shadow: none;
}

.product-card.horizontal .product-image {
  flex-shrink: 0;
  width: 180px;
  height: 180px;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-card.horizontal .product-image img {
  width: 100%;
  height: 100%;
  aspect-ratio: unset;
  object-fit: cover;
}

.product-card.horizontal .favorite-btn {
  top: 8px;
  right: 8px;
  padding: 6px;
}

.product-card.horizontal .product-info {
  flex: 1;
  padding: 14px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-width: 0;
}

.product-card.horizontal .product-name {
  font-size: 15px;
  margin: 0 0 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-card.horizontal .product-bottom {
  flex-wrap: wrap;
  gap: 8px;
}

.product-card.horizontal .add-cart-btn {
  padding: 6px 12px;
  font-size: 12px;
}

.product-card.horizontal .product-price {
  font-size: 16px;
}

.product-card.horizontal .share-btn {
  top: 8px;
  left: 8px;
  padding: 6px;
}

.product-card.horizontal .sold-out-text {
  font-size: 18px;
  padding: 6px 18px;
  letter-spacing: 4px;
}

@media (max-width: 768px) {
  .add-cart-btn {
    padding: 8px 14px;
    font-size: 13px;
  }

  .favorite-btn {
    padding: 8px;
  }

  .product-info {
    padding: 12px;
  }

  .product-name {
    font-size: 12px;
    margin-bottom: 8px;
  }

  .product-bottom {
    gap: 6px;
  }

  .product-price {
    font-size: 13px;
  }

  .product-card.horizontal .product-image {
    width: 120px;
    height: 120px;
    padding: 8px;
  }

  .product-card.horizontal .favorite-btn {
    top: 6px;
    right: 6px;
    padding: 5px;
  }

  .product-card.horizontal .share-btn {
    top: 6px;
    left: 6px;
    padding: 5px;
  }

  .product-card.horizontal .product-info {
    padding: 10px;
  }

  .product-card.horizontal .product-name {
    font-size: 14px;
    margin: 0 0 6px;
  }

  .product-card.horizontal .add-cart-btn {
    padding: 8px 12px;
    font-size: 12px;
  }

  .product-card.horizontal .product-price {
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .product-info {
    padding: 10px;
  }

  .product-name {
    font-size: 12px;
    margin-bottom: 6px;
  }

  .add-cart-btn {
    padding: 6px;
    font-size: 0;
  }

  .add-cart-btn .cart-icon {
    display: block;
    width: 16px;
    height: 16px;
  }

  .add-cart-btn .cart-text {
    display: none;
  }

  .product-price {
    font-size: 13px;
  }

  .product-image {
    padding: 6px;
  }

  .favorite-btn {
    top: 6px;
    right: 6px;
    padding: 5px;
  }

  .favorite-btn svg {
    width: 14px;
    height: 14px;
  }

  .share-btn {
    top: 6px;
    left: 6px;
    padding: 5px;
  }

  .share-btn svg {
    width: 14px;
    height: 14px;
  }

  .product-card.horizontal .product-image {
    width: 100px;
    height: 100px;
    padding: 6px;
  }

  .product-card.horizontal .product-name {
    font-size: 13px;
  }

  .product-card.horizontal .add-cart-btn .cart-text {
    display: inline;
  }

  .product-card.horizontal .add-cart-btn {
    font-size: 12px;
    padding: 6px 10px;
  }

  .product-card.horizontal .add-cart-btn .cart-icon {
    display: none;
  }
}
</style>
