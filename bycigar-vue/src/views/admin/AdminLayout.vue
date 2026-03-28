<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const sidebarCollapsed = ref(false)

const menuItems = [
  { path: '/admin/products', name: '商品管理', icon: 'box' },
  { path: '/admin/banners', name: '轮播图管理', icon: 'image' },
  { path: '/admin/categories', name: '分类管理', icon: 'folder' },
  { path: '/admin/pages', name: '页面管理', icon: 'file-text' },
  { path: '/admin/settings', name: '站点设置', icon: 'settings' }
]

const isActive = (path) => route.path === path

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const goHome = () => router.push('/')
const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="admin-layout" :class="{ collapsed: sidebarCollapsed }">
    <aside class="sidebar">
      <div class="sidebar-header">
        <router-link to="/admin" class="logo">
          <span v-if="!sidebarCollapsed">HUAUGE</span>
          <span v-else>H</span>
        </router-link>
      </div>
      
      <nav class="sidebar-nav">
        <router-link 
          v-for="item in menuItems" 
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
        >
          <svg v-if="item.icon === 'box'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
            <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
            <line x1="12" y1="22.08" x2="12" y2="12"></line>
          </svg>
          <svg v-else-if="item.icon === 'image'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <circle cx="8.5" cy="8.5" r="1.5"></circle>
            <polyline points="21 15 16 10 5 21"></polyline>
          </svg>
          <svg v-else-if="item.icon === 'folder'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
          </svg>
          <svg v-else-if="item.icon === 'file-text'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
            <polyline points="10 9 9 9 8 9"></polyline>
          </svg>
          <svg v-else-if="item.icon === 'settings'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
          <span v-if="!sidebarCollapsed">{{ item.name }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <button class="toggle-btn" @click="toggleSidebar">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline v-if="sidebarCollapsed" points="9 18 15 12 9 6"></polyline>
            <polyline v-else points="15 18 9 12 15 6"></polyline>
          </svg>
        </button>
      </div>
    </aside>

    <div class="main-content">
      <header class="admin-header">
        <div class="header-left">
          <h1 class="page-title">{{ route.meta.title || '后台管理' }}</h1>
        </div>
        <div class="header-right">
          <span class="user-info">{{ authStore.userName }}</span>
          <button class="btn-home" @click="goHome" title="返回前台">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
          </button>
          <button class="btn-logout" @click="handleLogout">退出</button>
        </div>
      </header>

      <main class="content-area">
        <router-view></router-view>
      </main>
    </div>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: #f5f5f5;
}

.sidebar {
  width: 240px;
  background: #1a1a1a;
  color: #fff;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 100;
}

.admin-layout.collapsed .sidebar {
  width: 60px;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #333;
}

.logo {
  font-family: 'Playfair Display', serif;
  font-size: 20px;
  font-weight: 600;
  color: #d4a574;
  text-decoration: none;
  letter-spacing: 3px;
  transition: all 0.3s ease;
}

.logo:hover {
  color: #e8c49a;
  letter-spacing: 4px;
}

.sidebar-nav {
  flex: 1;
  padding: 15px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  color: #ccc;
  text-decoration: none;
  transition: all 0.2s;
}

.nav-item:hover {
  background: #2d2d2d;
  color: #fff;
}

.nav-item.active {
  background: #d4a574;
  color: #1a1a1a;
}

.sidebar-footer {
  padding: 15px;
  border-top: 1px solid #333;
}

.toggle-btn {
  width: 100%;
  padding: 8px;
  background: #2d2d2d;
  border: none;
  border-radius: 4px;
  color: #ccc;
  cursor: pointer;
  display: flex;
  justify-content: center;
  transition: all 0.2s;
}

.toggle-btn:hover {
  background: #3d3d3d;
  color: #fff;
}

.main-content {
  flex: 1;
  margin-left: 240px;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s;
}

.admin-layout.collapsed .main-content {
  margin-left: 60px;
}

.admin-header {
  background: #fff;
  padding: 15px 25px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  flex-shrink: 0;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-info {
  color: #666;
  font-size: 14px;
}

.btn-home {
  background: transparent;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 5px;
  transition: color 0.2s;
}

.btn-home:hover {
  color: #d4a574;
}

.btn-logout {
  background: #f5f5f5;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  color: #666;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-logout:hover {
  background: #e0e0e0;
}

.content-area {
  padding: 25px;
  flex: 1;
  color: #333;
}
</style>
