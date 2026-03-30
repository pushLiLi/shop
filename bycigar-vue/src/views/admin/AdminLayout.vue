<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const sidebarCollapsed = ref(false)

const allMenuGroups = [
  {
    title: '概览',
    items: [
      { path: '/admin', name: '仪表盘', icon: 'dashboard' }
    ]
  },
  {
    title: '商品',
    items: [
      { path: '/admin/products', name: '商品管理', icon: 'box' },
      { path: '/admin/categories', name: '分类管理', icon: 'folder' }
    ]
  },
  {
    title: '交易',
    items: [
      { path: '/admin/orders', name: '订单管理', icon: 'shopping-bag' },
      { path: '/admin/payment-methods', name: '付款方式', icon: 'credit-card' }
    ]
  },
  {
    title: '客户',
    items: [
      { path: '/admin/users', name: '用户管理', icon: 'users' },
      { path: '/admin/chat', name: '客服消息', icon: 'message-circle' }
    ]
  },
  {
    title: '系统',
    superAdminOnly: true,
    items: [
      { path: '/admin/banners', name: '轮播图管理', icon: 'image' },
      { path: '/admin/pages', name: '页面管理', icon: 'file-text' },
      { path: '/admin/settings', name: '站点设置', icon: 'settings' },
      { path: '/admin/contact-methods', name: '联系方式', icon: 'phone' }
    ]
  }
]

const menuGroups = computed(() => {
  if (authStore.isSuperAdmin) return allMenuGroups
  return allMenuGroups.filter(group => !group.superAdminOnly)
})

const isActive = (path) => {
  if (path === '/admin') return route.path === '/admin'
  return route.path === path
}

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
        <div v-for="group in menuGroups" :key="group.title" class="nav-group">
          <div v-if="!sidebarCollapsed" class="nav-group-title">{{ group.title }}</div>
          <router-link 
            v-for="item in group.items" 
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
          >
          <svg v-if="item.icon === 'dashboard'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"></rect>
            <rect x="14" y="3" width="7" height="7"></rect>
            <rect x="14" y="14" width="7" height="7"></rect>
            <rect x="3" y="14" width="7" height="7"></rect>
          </svg>
          <svg v-else-if="item.icon === 'box'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
            <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
            <line x1="12" y1="22.08" x2="12" y2="12"></line>
          </svg>
          <svg v-else-if="item.icon === 'shopping-bag'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <path d="M16 10a4 4 0 0 1-8 0"></path>
          </svg>
          <svg v-else-if="item.icon === 'users'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
          <svg v-else-if="item.icon === 'message-circle'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
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
          <svg v-else-if="item.icon === 'credit-card'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="1" y="4" width="22" height="16" rx="2" ry="2"></rect>
            <line x1="1" y1="10" x2="23" y2="10"></line>
          </svg>
          <svg v-else-if="item.icon === 'phone'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"></path>
          </svg>
          <span v-if="!sidebarCollapsed">{{ item.name }}</span>
        </router-link>
        </div>
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
  padding: 10px 0;
  overflow-y: auto;
}

.nav-group {
  margin-bottom: 6px;
}

.nav-group-title {
  padding: 10px 20px 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #777;
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

.admin-layout.collapsed .nav-group + .nav-group {
  border-top: 1px solid #333;
  padding-top: 6px;
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
