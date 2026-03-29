import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CategoryView from '../views/CategoryView.vue'
import ProductDetailView from '../views/ProductDetailView.vue'
import SearchView from '../views/SearchView.vue'
import CartView from '../views/CartView.vue'
import CheckoutView from '../views/CheckoutView.vue'
import OrdersView from '../views/OrdersView.vue'
import FavoritesView from '../views/FavoritesView.vue'
import PageView from '../views/PageView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import NotificationDetailView from '../views/NotificationDetailView.vue'
import AdminLayout from '../views/admin/AdminLayout.vue'
import AdminProducts from '../views/admin/AdminProducts.vue'
import AdminBanners from '../views/admin/AdminBanners.vue'
import AdminCategories from '../views/admin/AdminCategories.vue'
import AdminPages from '../views/admin/AdminPages.vue'
import AdminSettings from '../views/admin/AdminSettings.vue'
import AdminOrders from '../views/admin/AdminOrders.vue'
import AdminDashboard from '../views/admin/AdminDashboard.vue'
import AdminUsers from '../views/admin/AdminUsers.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/products',
      name: 'all-products',
      component: CategoryView
    },
    {
      path: '/category/:slug',
      name: 'category',
      component: CategoryView
    },
    {
      path: '/products/:id',
      name: 'product-detail',
      component: ProductDetailView
    },
    {
      path: '/search',
      name: 'search',
      component: SearchView
    },
    {
      path: '/cart',
      name: 'cart',
      component: CartView
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: CheckoutView,
      meta: { requiresAuth: true }
    },
    {
      path: '/orders',
      name: 'orders',
      component: OrdersView,
      meta: { requiresAuth: true }
    },
    {
      path: '/favorites',
      name: 'favorites',
      component: FavoritesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/:slug(about|services|privacy-policy|statement)',
      name: 'page',
      component: PageView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/notifications/:id',
      name: 'notification-detail',
      component: NotificationDetailView,
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: AdminDashboard,
          meta: { title: '仪表盘' }
        },
        {
          path: 'products',
          name: 'admin-products',
          component: AdminProducts,
          meta: { title: '商品管理' }
        },
        {
          path: 'orders',
          name: 'admin-orders',
          component: AdminOrders,
          meta: { title: '订单管理' }
        },
        {
          path: 'users',
          name: 'admin-users',
          component: AdminUsers,
          meta: { title: '用户管理' }
        },
        {
          path: 'categories',
          name: 'admin-categories',
          component: AdminCategories,
          meta: { title: '分类管理' }
        },
        {
          path: 'banners',
          name: 'admin-banners',
          component: AdminBanners,
          meta: { title: '轮播图管理', requiresSuperAdmin: true }
        },
        {
          path: 'pages',
          name: 'admin-pages',
          component: AdminPages,
          meta: { title: '页面管理', requiresSuperAdmin: true }
        },
        {
          path: 'settings',
          name: 'admin-settings',
          component: AdminSettings,
          meta: { title: '站点设置', requiresSuperAdmin: true }
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
