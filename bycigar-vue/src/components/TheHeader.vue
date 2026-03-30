<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useFavoritesStore } from '../stores/favorites'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { useNotificationsStore } from '../stores/notifications'
import NotificationPanel from './NotificationPanel.vue'

const router = useRouter()
const cartStore = useCartStore()
const favoritesStore = useFavoritesStore()
const authStore = useAuthStore()
const toast = useToastStore()
const notificationsStore = useNotificationsStore()
const isMenuOpen = ref(false)
const showMobileSearch = ref(false)
const showNotice = ref(!localStorage.getItem('notice_closed'))
const searchKeyword = ref('')
const mobileSearchKeyword = ref('')
const showUserMenu = ref(false)
const isAdmin = computed(() => authStore.isAdmin)
const showSearchDropdown = ref(false)
const searchSuggestions = ref([])
const suggestionsLoading = ref(false)
let suggestTimer = null

const HISTORY_KEY = 'search_history'
const MAX_HISTORY = 8

function getSearchHistory() {
  try {
    return JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]')
  } catch {
    return []
  }
}

function addToHistory(keyword) {
  const history = getSearchHistory().filter(h => h !== keyword)
  history.unshift(keyword)
  localStorage.setItem(HISTORY_KEY, JSON.stringify(history.slice(0, MAX_HISTORY)))
}

function removeHistoryItem(keyword) {
  const history = getSearchHistory().filter(h => h !== keyword)
  localStorage.setItem(HISTORY_KEY, JSON.stringify(history))
}

function clearSearchHistory() {
  localStorage.removeItem(HISTORY_KEY)
}

const searchHistory = ref(getSearchHistory())

function onSearchInputFocus() {
  if (!searchKeyword.value.trim()) {
    searchHistory.value = getSearchHistory()
  }
  showSearchDropdown.value = true
}

async function fetchSuggestions(keyword) {
  if (!keyword || keyword.length < 1) {
    searchSuggestions.value = []
    return
  }
  suggestionsLoading.value = true
  try {
    const res = await fetch(`/api/products/suggest?q=${encodeURIComponent(keyword)}&limit=6`)
    const data = await res.json()
    searchSuggestions.value = data.suggestions || []
  } catch {
    searchSuggestions.value = []
  } finally {
    suggestionsLoading.value = false
  }
}

function onSearchInputChange() {
  clearTimeout(suggestTimer)
  if (!searchKeyword.value.trim()) {
    searchSuggestions.value = []
    searchHistory.value = getSearchHistory()
    return
  }
  suggestTimer = setTimeout(() => {
    fetchSuggestions(searchKeyword.value.trim())
  }, 300)
}

function closeSearchDropdown() {
  showSearchDropdown.value = false
}

function handleSuggestionClick(product) {
  closeSearchDropdown()
  router.push('/product/' + product.id)
}

function handleHistoryClick(keyword) {
  searchKeyword.value = keyword
  addToHistory(keyword)
  closeSearchDropdown()
  router.push('/search?q=' + encodeURIComponent(keyword))
}

function removeHistory(keyword) {
  removeHistoryItem(keyword)
  searchHistory.value = getSearchHistory()
}

function clearHistory() {
  clearSearchHistory()
  searchHistory.value = []
}

const menuItems = [
  { name: '首页', path: '/', children: [] },
  { name: '全部商品', path: '/products', children: [] },
  { name: '关于我们', path: '/about', children: [] }
]

watch(isMenuOpen, (val) => {
  document.body.style.overflow = val ? 'hidden' : ''
})

watch(showMobileSearch, (val) => {
  if (val) {
    isMenuOpen.value = false
  }
})

watch(() => authStore.isLoggedIn, (loggedIn) => {
  if (loggedIn) {
    notificationsStore.initWSListener()
  } else {
    notificationsStore.cleanupWSListener()
  }
})

onMounted(() => {
  if (authStore.isLoggedIn) {
    cartStore.fetchCart()
    favoritesStore.fetchFavorites()
    notificationsStore.initWSListener()
  }
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('visibilitychange', handleVisibilityChange)
})

const handleLogout = () => {
  notificationsStore.cleanupWSListener()
  authStore.logout()
  showUserMenu.value = false
  router.push('/')
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const closeMenu = () => {
  isMenuOpen.value = false
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    addToHistory(searchKeyword.value.trim())
    closeSearchDropdown()
    router.push('/search?q=' + encodeURIComponent(searchKeyword.value.trim()))
  }
}

const handleMobileSearch = () => {
  if (mobileSearchKeyword.value.trim()) {
    addToHistory(mobileSearchKeyword.value.trim())
    router.push('/search?q=' + encodeURIComponent(mobileSearchKeyword.value.trim()))
    showMobileSearch.value = false
    mobileSearchKeyword.value = ''
  }
}

const handleCartClick = () => {
  if (!authStore.isLoggedIn) {
    toast.error('请先登录')
    return
  }
  cartStore.openCart()
}

const handleOverlayClick = (e) => {
  if (e.target === e.currentTarget) {
    closeMenu()
  }
}

const handleClickOutside = (e) => {
  const wrapper = document.querySelector('.user-menu-wrapper')
  if (wrapper && !wrapper.contains(e.target)) {
    showUserMenu.value = false
  }
  const searchWrapper = document.querySelector('.search-form-wrapper')
  if (searchWrapper && !searchWrapper.contains(e.target)) {
    closeSearchDropdown()
  }
}

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  clearTimeout(suggestTimer)
  notificationsStore.cleanupWSListener()
})

const handleVisibilityChange = () => {
  if (!authStore.isLoggedIn) return
  if (!document.hidden) {
    notificationsStore.fetchUnreadCount()
  }
}
</script>

<template>
  <header class="site-header">
    <div class="header-top" v-if="showNotice">
      <div class="container">
        <div class="top-notice">
          尊敬的客户，为确保您的购物体验，下单前请仔细阅读我们<a href="/services">服务条款</a>。
          <button class="notice-close" @click="showNotice = false; localStorage.setItem('notice_closed', '1')">&times;</button>
        </div>
      </div>
    </div>

    <div class="header-main">
      <div class="container">
        <div class="header-content">
          <div class="header-left">
            <button class="mobile-menu-btn" @click="toggleMenu" aria-label="打开菜单">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="3" y1="12" x2="21" y2="12"></line>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <line x1="3" y1="18" x2="21" y2="18"></line>
              </svg>
            </button>
            <router-link to="/" class="logo">
              HUAUGE
            </router-link>
            <nav class="header-nav" :class="{ 'is-open': isMenuOpen }">
              <ul class="nav-list">
                <li v-for="item in menuItems" :key="item.path" class="nav-item">
                  <router-link :to="item.path" class="nav-link" @click="closeMenu">
                    {{ item.name }}
                  </router-link>
                </li>
              </ul>
            </nav>
            <Transition name="overlay-fade">
              <div v-if="isMenuOpen" class="menu-overlay" @click="closeMenu"></div>
            </Transition>
          </div>

          <div class="header-center">
            <div class="search-form-wrapper">
              <form class="search-form" @submit.prevent="handleSearch">
                <input 
                  v-model="searchKeyword" 
                  type="text" 
                  class="search-input" 
                  placeholder="搜索"
                  autocomplete="off"
                  @focus="onSearchInputFocus"
                  @input="onSearchInputChange"
                >
                <button type="submit" class="search-btn">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="m21 21-4.35-4.35"></path>
                  </svg>
                </button>
              </form>
              <Transition name="dropdown">
                <div v-if="showSearchDropdown" class="search-dropdown">
                  <div v-if="suggestionsLoading" class="dropdown-loading">搜索中...</div>
                  <div v-else-if="searchSuggestions.length > 0" class="dropdown-suggestions">
                    <div 
                      v-for="item in searchSuggestions" 
                      :key="item.id" 
                      class="suggestion-item"
                      @mousedown.prevent="handleSuggestionClick(item)"
                    >
                      <img v-if="item.thumbnailUrl" :src="item.thumbnailUrl" class="suggestion-thumb" alt="">
                      <div class="suggestion-thumb-placeholder" v-else></div>
                      <div class="suggestion-info">
                        <span class="suggestion-name">{{ item.name }}</span>
                        <span class="suggestion-price">¥{{ item.price }}</span>
                      </div>
                    </div>
                  </div>
                  <div v-else-if="searchKeyword.trim() && !suggestionsLoading" class="dropdown-empty">
                    无搜索建议
                  </div>
                  <div v-if="!searchKeyword.trim() && searchHistory.length > 0" class="dropdown-history">
                    <div class="history-header">
                      <span class="history-title">搜索历史</span>
                      <button class="history-clear" @mousedown.prevent="clearHistory">清空</button>
                    </div>
                    <div class="history-list">
                      <div 
                        v-for="h in searchHistory" 
                        :key="h" 
                        class="history-item"
                        @mousedown.prevent="handleHistoryClick(h)"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg>
                        <span class="history-text">{{ h }}</span>
                        <button class="history-remove" @mousedown.prevent.stop="removeHistory(h)">&times;</button>
                      </div>
                    </div>
                  </div>
                </div>
              </Transition>
            </div>
          </div>

          <div class="header-right">
            <div class="header-icons">
              <button class="icon-item mobile-search-btn" @click="showMobileSearch = !showMobileSearch" aria-label="搜索">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="11" cy="11" r="8"></circle>
                  <path d="m21 21-4.35-4.35"></path>
                </svg>
              </button>
              <template v-if="authStore.isLoggedIn">
                <div class="user-menu-wrapper">
                  <button class="icon-item user-btn" @click="showUserMenu = !showUserMenu">
                    <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                      <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                    <span class="user-name">{{ authStore.userName }}</span>
                  </button>
                  <Transition name="dropdown">
                    <div v-if="showUserMenu" class="user-dropdown">
                      <router-link to="/profile" class="dropdown-item" @click="showUserMenu = false">个人信息</router-link>
                      <router-link to="/orders" class="dropdown-item" @click="showUserMenu = false">我的订单</router-link>
                      <router-link to="/favorites" class="dropdown-item" @click="showUserMenu = false">我的收藏</router-link>
                      <router-link v-if="isAdmin" to="/admin" class="dropdown-item admin-link" @click="showUserMenu = false">后台管理</router-link>
                      <button class="dropdown-item logout-btn" @click="handleLogout">退出登录</button>
                    </div>
                  </Transition>
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
              <button v-if="authStore.isLoggedIn" @click="notificationsStore.openPanel()" class="icon-item" title="通知">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                  <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
                </svg>
                <span class="icon-badge" v-if="notificationsStore.unreadCount">{{ notificationsStore.unreadCount }}</span>
              </button>
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

        <Transition name="search-slide">
          <div v-if="showMobileSearch" class="mobile-search">
            <form class="mobile-search-form" @submit.prevent="handleMobileSearch">
              <input 
                v-model="mobileSearchKeyword" 
                type="text" 
                class="mobile-search-input" 
                placeholder="搜索商品..."
                ref="mobileSearchInput"
                autofocus
              >
              <button type="submit" class="mobile-search-submit">搜索</button>
            </form>
            <div v-if="!mobileSearchKeyword && searchHistory.length" class="mobile-search-history">
              <div class="mobile-history-header">
                <span>搜索历史</span>
                <button @click="clearHistory">清空</button>
              </div>
              <div class="mobile-history-tags">
                <button 
                  v-for="h in searchHistory" 
                  :key="h" 
                  class="mobile-history-tag"
                  @click="mobileSearchKeyword = h; handleMobileSearch()"
                >{{ h }}</button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </div>
    <NotificationPanel />
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

.top-notice {
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.notice-close {
  position: absolute;
  right: 0;
  background: transparent;
  border: none;
  color: #999;
  font-size: 18px;
  cursor: pointer;
  padding: 0 5px;
  line-height: 1;
  transition: color 0.2s;
}

.notice-close:hover {
  color: #fff;
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

.header-right {
  display: flex;
  align-items: center;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.search-form-wrapper {
  position: relative;
  width: 280px;
}

.search-form {
  display: flex;
  align-items: center;
  background: #2d2d2d;
  border-radius: 4px;
  overflow: hidden;
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

.search-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #2d2d2d;
  border-radius: 0 0 8px 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
}

.search-dropdown::-webkit-scrollbar {
  width: 4px;
}

.search-dropdown::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 2px;
}

.dropdown-loading,
.dropdown-empty {
  padding: 20px;
  text-align: center;
  color: #888;
  font-size: 13px;
}

.dropdown-suggestions {
  padding: 4px 0;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  cursor: pointer;
  transition: background 0.15s;
}

.suggestion-item:hover {
  background: rgba(212, 165, 116, 0.1);
}

.suggestion-thumb {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.suggestion-thumb-placeholder {
  width: 40px;
  height: 40px;
  background: #3a3a3a;
  border-radius: 4px;
  flex-shrink: 0;
}

.suggestion-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.suggestion-name {
  color: #ddd;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.suggestion-price {
  color: #d4a574;
  font-size: 12px;
}

.dropdown-history {
  padding: 8px 0;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 14px;
}

.history-title {
  color: #888;
  font-size: 12px;
}

.history-clear {
  background: transparent;
  border: none;
  color: #666;
  font-size: 12px;
  cursor: pointer;
  transition: color 0.2s;
}

.history-clear:hover {
  color: #d4a574;
}

.history-list {
  padding: 0;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  cursor: pointer;
  transition: background 0.15s;
  color: #ccc;
}

.history-item:hover {
  background: rgba(212, 165, 116, 0.1);
}

.history-item svg {
  color: #666;
  flex-shrink: 0;
}

.history-text {
  flex: 1;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.history-remove {
  background: transparent;
  border: none;
  color: #555;
  font-size: 16px;
  cursor: pointer;
  padding: 0 2px;
  line-height: 1;
  transition: color 0.2s;
}

.history-remove:hover {
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
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 44px;
  min-height: 44px;
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
  padding: 10px;
  min-width: 44px;
  min-height: 44px;
  align-items: center;
  justify-content: center;
}

@media (max-width: 992px) {
  .header-left {
    gap: 10px;
  }

  .mobile-menu-btn {
    display: flex;
  }

  .header-nav {
    position: fixed;
    top: 0;
    left: -100%;
    width: 280px;
    height: 100vh;
    background: #1a1a1a;
    flex-direction: column;
    justify-content: flex-start;
    padding: 20px 0;
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
    padding: 10px 20px;
  }

  .nav-link {
    height: auto;
    padding: 15px 0;
    border-bottom: 1px solid #333;
  }

  .nav-link::after {
    display: none;
  }

  .menu-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1000;
  }

  .header-center {
    display: none;
  }

  .mobile-search-btn {
    display: flex;
  }

  .mobile-search {
    padding: 0 0 15px;
  }

  .mobile-search-form {
    display: flex;
    gap: 8px;
  }

  .mobile-search-input {
    flex: 1;
    background: #2d2d2d;
    border: 1px solid #444;
    border-radius: 4px;
    padding: 10px 15px;
    color: #fff;
    font-size: 14px;
    outline: none;
  }

  .mobile-search-input::placeholder {
    color: #888;
  }

  .mobile-search-input:focus {
    border-color: #d4a574;
  }

  .mobile-search-submit {
    background: #d4a574;
    color: #1a1a1a;
    border: none;
    padding: 10px 20px;
    border-radius: 4px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
  }

  .mobile-search-history {
    margin-top: 10px;
  }

  .mobile-history-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .mobile-history-header span {
    color: #888;
    font-size: 12px;
  }

  .mobile-history-header button {
    background: transparent;
    border: none;
    color: #666;
    font-size: 12px;
    cursor: pointer;
  }

  .mobile-history-header button:hover {
    color: #d4a574;
  }

  .mobile-history-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .mobile-history-tag {
    background: #2d2d2d;
    border: 1px solid #333;
    color: #ccc;
    padding: 5px 12px;
    border-radius: 16px;
    font-size: 12px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .mobile-history-tag:hover {
    border-color: #d4a574;
    color: #d4a574;
  }

  .user-name {
    display: none;
  }

  .header-right {
    flex-shrink: 0;
  }

  .header-icons {
    gap: 5px;
  }

  .icon-item {
    min-width: 36px;
    min-height: 36px;
  }

  .icon-item svg {
    width: 20px;
    height: 20px;
  }

  .icon-badge {
    top: -2px;
    right: -6px;
    font-size: 9px;
    padding: 1px 4px;
    min-width: 14px;
  }
}

.mobile-search-btn {
  display: none;
}

.search-slide-enter-active,
.search-slide-leave-active {
  transition: all 0.3s ease;
}

.search-slide-enter-from,
.search-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.overlay-fade-enter-active,
.overlay-fade-leave-active {
  transition: opacity 0.3s ease;
}

.overlay-fade-enter-from,
.overlay-fade-leave-to {
  opacity: 0;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}
</style>
