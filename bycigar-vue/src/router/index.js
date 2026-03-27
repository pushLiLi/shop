import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CategoryView from '../views/CategoryView.vue'
import ProductDetailView from '../views/ProductDetailView.vue'
import SearchView from '../views/SearchView.vue'
import CartView from '../views/CartView.vue'
import CheckoutView from '../views/CheckoutView.vue'
import OrdersView from '../views/OrdersView.vue'
import FavoritesView from '../views/FavoritesView.vue'
import AboutView from '../views/AboutView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import ServicesView from '../views/ServicesView.vue'
import PrivacyPolicyView from '../views/PrivacyPolicyView.vue'
import ReturnsPolicyView from '../views/ReturnsPolicyView.vue'
import AdminLayout from '../views/admin/AdminLayout.vue'
import AdminProducts from '../views/admin/AdminProducts.vue'
import AdminBanners from '../views/admin/AdminBanners.vue'
import AdminCategories from '../views/admin/AdminCategories.vue'

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
      path: '/about',
      name: 'about',
      component: AboutView
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
      path: '/services',
      name: 'services',
      component: ServicesView
    },
    {
      path: '/privacy-policy',
      name: 'privacy-policy',
      component: PrivacyPolicyView
    },
    {
      path: '/returns-policy',
      name: 'returns-policy',
      component: ReturnsPolicyView
    },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/products'
        },
        {
          path: 'products',
          name: 'admin-products',
          component: AdminProducts,
          meta: { title: '商品管理' }
        },
        {
          path: 'banners',
          name: 'admin-banners',
          component: AdminBanners,
          meta: { title: '轮播图管理' }
        },
        {
          path: 'categories',
          name: 'admin-categories',
          component: AdminCategories,
          meta: { title: '分类管理' }
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
