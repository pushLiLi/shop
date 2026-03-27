<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useFavoritesStore } from '../stores/favorites'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const router = useRouter()
const cartStore = useCartStore()
const favoritesStore = useFavoritesStore()
const authStore = useAuthStore()
const toast = useToastStore()
const isMenuOpen = ref(false)
const searchKeyword = ref('')
const showUserMenu = ref(false)
const isAdmin = computed(() => authStore.isAdmin)

const menuItems = [
  { name: '首页', path: '/', children: [] },
  {
    name: '全部商品',
    path: '/products',
    children: [
      { name: '古巴雪茄', path: '/category/cuban' },
      { name: '尼加拉瓜雪茄', path: '/category/nicaraguan' },
      { name: '多米尼加雪茄', path: '/category/dominican' },
      { name: '雪茄配件', path: '/category/accessories' }
    ]
  },
  { name: '关于我们', path: '/about', children: [] }
]

onMounted(() => {
  if (authStore.isLoggedIn) {
    cartStore.fetchCart()
    favoritesStore.fetchFavorites()
  }
})

const handleLogout = () => {
  authStore.logout()
  showUserMenu.value = false
  router.push('/')
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push('/search?q=' + encodeURIComponent(searchKeyword.value.trim()))
  } else {
    alert('请输入搜索关键词')
  }
}

const handleCartClick = () => {
  if (!authStore.isLoggedIn) {
    toast.error('请先登录')
    return
  }
  cartStore.openCart()
}
</script>

<template>
  <header class="site-header">
    <div class="header-top">
      <div class="container">
        <div class="top-notice">
          尊敬的客户，为确保您的购物体验，下单前请仔细阅读我们<a href="/services">服务条款</a>。
        </div>
      </div>
    </div>

    <div class="header-main">
      <div class="container">
        <div class="header-content">
          <div class="header-left">
            <router-link to="/" class="logo">
              HUAUGE
            </router-link>
            <nav class="header-nav" :class="{ 'is-open': isMenuOpen }">
              <ul class="nav-list">
                <li v-for="item in menuItems" :key="item.path" class="nav-item">
                  <router-link v-if="item.children.length === 0" :to="item.path" class="nav-link">
                    {{ item.name }}
                  </router-link>
                  <div v-else class="nav-dropdown">
                    <span class="nav-link">
                      {{ item.name }}
                      <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-left: 4px;">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </span>
                    <ul class="dropdown-menu">
                      <li v-for="child in item.children" :key="child.path">
                        <router-link :to="child.path" class="dropdown-link">
                          {{ child.name }}
                        </router-link>
                      </li>
                    </ul>
                  </div>
                </li>
              </ul>
            </nav>
            <button class="mobile-menu-btn" @click="toggleMenu">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="3" y1="12" x2="21" y2="12"></line>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <line x1="3" y1="18" x2="21" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="header-center">
            <form class="search-form" @submit.prevent="handleSearch">
              <input 
                v-model="searchKeyword" 
                type="text" 
                class="search-input" 
                placeholder="搜索"
              >
              <button type="submit" class="search-btn">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="11" cy="11" r="8"></circle>
                  <path d="m21 21-4.35-4.35"></path>
                </svg>
              </button>
            </form>
          </div>

          <div class="header-right">
            <div class="header-icons">
              <template v-if="authStore.isLoggedIn">
                <div class="user-menu-wrapper" @mouseleave="showUserMenu = false">
                  <button class="icon-item user-btn" @mouseenter="showUserMenu = true">
                    <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                      <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                    <span class="user-name">{{ authStore.userName }}</span>
                  </button>
                  <div v-if="showUserMenu" class="user-dropdown">
                    <router-link to="/profile" class="dropdown-item">个人信息</router-link>
                    <router-link to="/orders" class="dropdown-item">我的订单</router-link>
                    <router-link to="/favorites" class="dropdown-item">我的收藏</router-link>
                    <router-link v-if="isAdmin" to="/admin" class="dropdown-item admin-link">后台管理</router-link>
                    <button class="dropdown-item logout-btn" @click="handleLogout">退出登录</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <router-link to="/login" class="icon-item" title="登录">
                  <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                    <circle cx="12" cy="7" r="4"></circle>
                  </svg>
                </router-link>
              </template>
              <router-link to="/favorites" class="icon-item" title="收藏">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
                </svg>
                <span class="icon-badge" v-if="favoritesStore.items.length">{{ favoritesStore.items.length }}</span>
              </router-link>
            <button @click="handleCartClick" class="icon-item" title="购物车">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="9" cy="21" r="1"></circle>
                  <circle cx="20" cy="21" r="1"></circle>
                  <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path>
                </svg>
                <span class="icon-badge" v-if="cartStore.items.length">{{ cartStore.items.length }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  background: #1a1a1a;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-top {
  background: #2d2d2d;
  padding: 8px 0;
  font-size: 13px;
  color: #ccc;
}

.top-notice a {
  color: #d4a574;
  text-decoration: none;
}

.top-notice a:hover {
  text-decoration: underline;
}

.header-main {
  padding: 15px 0;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 15px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 25px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 25px;
}

.logo {
  font-family: 'Playfair Display', serif;
  font-size: 28px;
  font-weight: 600;
  color: #d4a574;
  letter-spacing: 3px;
  text-decoration: none;
  transition: all 0.3s ease;
}

.logo:hover {
  color: #e8c49a;
  letter-spacing: 4px;
}

.header-nav {
  display: flex;
}

.nav-list {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
  gap: 15px;
}

.nav-item {
  position: relative;
}

.nav-link {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 44px;
  padding: 0 18px;
  color: #c9a87c;
  text-decoration: none;
  font-family: 'Playfair Display', serif;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 1.5px;
  white-space: nowrap;
  position: relative;
  transition: color 0.3s ease;
}

.nav-link:hover {
  color: #d4a574;
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 8px;
  left: 18px;
  right: 18px;
  height: 1px;
  background: #d4a574;
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.nav-link:hover::after {
  transform: scaleX(1);
}

.nav-dropdown {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 180px;
  background: #2d2d2d;
  border-radius: 4px;
  padding: 10px 0;
  list-style: none;
  margin: 0;
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}

.nav-dropdown:hover .dropdown-menu {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

.dropdown-link {
  display: block;
  padding: 10px 18px;
  color: #c9a87c;
  text-decoration: none;
  font-family: 'Playfair Display', serif;
  font-size: 13px;
  letter-spacing: 0.5px;
  transition: all 0.3s;
  white-space: nowrap;
}

.dropdown-link:hover {
  color: #d4a574;
  background: rgba(212, 165, 116, 0.1);
}

.header-right {
  display: flex;
  align-items: center;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.search-form {
  display: flex;
  align-items: center;
  background: #2d2d2d;
  border-radius: 4px;
  overflow: hidden;
  width: 280px;
}

.search-input {
  background: transparent;
  border: none;
  padding: 10px 15px;
  color: #fff;
  font-size: 14px;
  width: 240px;
  outline: none;
}

.search-input::placeholder {
  color: #888;
}

.search-btn {
  background: transparent;
  border: none;
  padding: 8px 12px;
  color: #888;
  cursor: pointer;
  transition: color 0.3s;
}

.search-btn:hover {
  color: #d4a574;
}

.header-icons {
  display: flex;
  align-items: center;
  gap: 15px;
}

.icon-item {
  position: relative;
  color: #fff;
  text-decoration: none;
  transition: color 0.3s;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
}

.icon-item:hover {
  color: #d4a574;
}

.icon-badge {
  position: absolute;
  top: -5px;
  right: -8px;
  background: #d4a574;
  color: #1a1a1a;
  font-size: 10px;
  font-weight: bold;
  padding: 2px 5px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
}

.user-menu-wrapper {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-btn:hover {
  background: rgba(255,255,255,0.1);
}

.user-name {
  color: #fff;
  font-size: 13px;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  min-width: 140px;
  background: #2d2d2d;
  border-radius: 6px;
  padding: 8px 0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  z-index: 100;
}

.dropdown-item {
  display: block;
  padding: 10px 16px;
  color: #ccc;
  text-decoration: none;
  font-size: 13px;
  transition: all 0.2s;
  background: transparent;
  border: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
}

.dropdown-item:hover {
  color: #d4a574;
  background: rgba(255,255,255,0.05);
}

.logout-btn {
  border-top: 1px solid #444;
  margin-top: 4px;
  padding-top: 12px;
}

.admin-link {
  color: #d4a574 !important;
}

.mobile-menu-btn {
  display: none;
  background: transparent;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 5px;
}

@media (max-width: 992px) {
  .header-nav {
    position: fixed;
    top: 0;
    left: -100%;
    width: 280px;
    height: 100vh;
    background: #1a1a1a;
    flex-direction: column;
    justify-content: flex-start;
    padding: 60px 20px 20px;
    transition: left 0.3s;
    z-index: 1001;
  }

  .header-nav.is-open {
    left: 0;
  }

  .nav-list {
    flex-direction: column;
    width: 100%;
    gap: 0;
  }

  .nav-link {
    height: auto;
    padding: 15px 0;
    border-bottom: 1px solid #333;
  }

  .nav-link::after {
    display: none;
  }

  .dropdown-menu {
    position: static;
    opacity: 1;
    visibility: visible;
    transform: none;
    background: transparent;
    padding-left: 15px;
    box-shadow: none;
  }

  .mobile-menu-btn {
    display: block;
  }

  .header-center {
    display: none;
  }
}
</style>
