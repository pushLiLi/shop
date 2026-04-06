import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path: '/products',
      name: 'all-products',
      component: () => import('../views/CategoryView.vue')
    },
    {
      path: '/category/:slug',
      name: 'category',
      component: () => import('../views/CategoryView.vue')
    },
    {
      path: '/products/:id',
      name: 'product-detail',
      component: () => import('../views/ProductDetailView.vue')
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/SearchView.vue')
    },
    {
      path: '/cart',
      name: 'cart',
      component: () => import('../views/CartView.vue')
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: () => import('../views/CheckoutView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/orders',
      name: 'orders',
      component: () => import('../views/OrdersView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/orders/:id',
      name: 'order-detail',
      component: () => import('../views/OrderDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/favorites',
      name: 'favorites',
      component: () => import('../views/FavoritesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/:slug(about|services|privacy-policy|statement)',
      name: 'page',
      component: () => import('../views/PageView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/notifications/:id',
      name: 'notification-detail',
      component: () => import('../views/NotificationDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      component: () => import('../views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/AdminDashboard.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'products',
          name: 'admin-products',
          component: () => import('../views/admin/AdminProducts.vue'),
          meta: { title: '商品管理' }
        },
        {
          path: 'orders',
          name: 'admin-orders',
          component: () => import('../views/admin/AdminOrders.vue'),
          meta: { title: '订单管理' }
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('../views/admin/AdminUsers.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'chat',
          name: 'admin-chat',
          component: () => import('../views/admin/AdminChat.vue'),
          meta: { title: '客服消息' }
        },
        {
          path: 'categories',
          name: 'admin-categories',
          component: () => import('../views/admin/AdminCategories.vue'),
          meta: { title: '分类管理' }
        },
        {
          path: 'banners',
          name: 'admin-banners',
          component: () => import('../views/admin/AdminBanners.vue'),
          meta: { title: '轮播图管理', requiresSuperAdmin: true }
        },
        {
          path: 'pages',
          name: 'admin-pages',
          component: () => import('../views/admin/AdminPages.vue'),
          meta: { title: '页面管理', requiresSuperAdmin: true }
        },
        {
          path: 'settings',
          name: 'admin-settings',
          component: () => import('../views/admin/AdminSettings.vue'),
          meta: { title: '站点设置', requiresSuperAdmin: true }
        },
        {
          path: 'payment-methods',
          name: 'admin-payment-methods',
          component: () => import('../views/admin/AdminPaymentMethods.vue'),
          meta: { title: '付款方式', requiresSuperAdmin: true }
        },
        {
          path: 'contact-methods',
          name: 'admin-contact-methods',
          component: () => import('../views/admin/AdminContactMethods.vue'),
          meta: { title: '联系方式', requiresSuperAdmin: true }
        },
        {
          path: 'cleanup',
          name: 'admin-cleanup',
          component: () => import('../views/admin/AdminCleanup.vue'),
          meta: { title: '数据清理', requiresSuperAdmin: true }
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')
    const userStr = localStorage.getItem('user')
    
    if (!token || !userStr) {
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }

    if (to.meta.requiresAdmin) {
      try {
        const user = JSON.parse(userStr)
        if (user.role !== 'admin' && user.role !== 'service') {
          next({ name: 'home' })
          return
        }
      } catch {
        next({ name: 'login', query: { redirect: to.fullPath } })
        return
      }
    }

    if (to.meta.requiresSuperAdmin) {
      try {
        const user = JSON.parse(userStr)
        if (user.role !== 'admin') {
          next({ name: 'home' })
          return
        }
      } catch {
        next({ name: 'login', query: { redirect: to.fullPath } })
        return
      }
    }
  }
  next()
})

export default router
