import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import ProductCard from './ProductCard.vue'

const mockProduct = {
  id: 1,
  name: 'Test Cigar',
  price: 29.99,
  imageUrl: '/media/cigar.jpg',
  stock: 10
}

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/products/:id', component: { template: '<div />' } },
      { path: '/login', component: { template: '<div />' } }
    ]
  })
}

async function mountCard(props = {}, options = {}) {
  const router = options.router || createTestRouter()
  const pinia = options.pinia || createPinia()
  setActivePinia(pinia)

  if (options.setupStores) {
    options.setupStores(pinia)
  }

  const wrapper = mount(ProductCard, {
    props: { product: mockProduct, ...props },
    global: {
      plugins: [router, pinia],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        }
      }
    }
  })
  await router.isReady()
  return wrapper
}

describe('ProductCard', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders product name and price', async () => {
    const wrapper = await mountCard()
    expect(wrapper.text()).toContain('Test Cigar')
    expect(wrapper.text()).toContain('¥29.99')
  })

  it('renders product image', async () => {
    const wrapper = await mountCard()
    const img = wrapper.find('img')
    expect(img.attributes('src')).toBe('/media/cigar.jpg')
  })

  it('links to product detail page', async () => {
    const wrapper = await mountCard()
    const links = wrapper.findAll('a[href="/products/1"]')
    expect(links.length).toBeGreaterThanOrEqual(1)
  })

  it('has add to cart button', async () => {
    const wrapper = await mountCard()
    expect(wrapper.find('.add-cart-btn').text()).toBe('加入购物车')
  })

  it('redirects to login when adding to cart while logged out', async () => {
    const router = createTestRouter()
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = await mountCard({}, { router })

    await wrapper.find('.add-cart-btn').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/login')
  })

  it('adds to cart when logged in', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'User', role: 'customer' }))

    const wrapper = await mountCard()

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    const addItemSpy = vi.spyOn(cartStore, 'addItem').mockResolvedValue()

    await wrapper.find('.add-cart-btn').trigger('click')
    expect(addItemSpy).toHaveBeenCalledWith(mockProduct, 1)
  })

  it('renders in vertical mode by default', async () => {
    const wrapper = await mountCard()
    expect(wrapper.find('.product-card').classes()).not.toContain('horizontal')
  })

  it('renders in horizontal mode when prop is true', async () => {
    const wrapper = await mountCard({ horizontal: true })
    expect(wrapper.find('.product-card').classes()).toContain('horizontal')
  })

  it('shows favorite button', async () => {
    const wrapper = await mountCard()
    expect(wrapper.find('.favorite-btn').exists()).toBe(true)
  })

  it('redirects to login when toggling favorite while logged out', async () => {
    const router = createTestRouter()
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = await mountCard({}, { router })

    await wrapper.find('.favorite-btn').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/login')
  })

  it('adds to favorites when logged in and not favorite', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'User', role: 'customer' }))

    const wrapper = await mountCard()
    const { useFavoritesStore } = await import('../stores/favorites')
    const favStore = useFavoritesStore()
    const addSpy = vi.spyOn(favStore, 'addItem').mockResolvedValue()

    await wrapper.find('.favorite-btn').trigger('click')
    expect(addSpy).toHaveBeenCalledWith(mockProduct)
  })

  it('removes from favorites when already favorite', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'User', role: 'customer' }))

    const pinia = createPinia()
    const wrapper = await mountCard({}, {
      pinia,
      setupStores: (p) => {
        setActivePinia(p)
      }
    })

    const { useFavoritesStore } = await import('../stores/favorites')
    const favStore = useFavoritesStore()
    favStore.items = [{ productId: 1 }]
    const removeSpy = vi.spyOn(favStore, 'removeItem').mockResolvedValue()

    await wrapper.vm.$nextTick()
    await wrapper.find('.favorite-btn').trigger('click')
    expect(removeSpy).toHaveBeenCalledWith(1)
  })

  it('shows sold out overlay when stock is 0', async () => {
    const soldOutProduct = { ...mockProduct, stock: 0 }
    const wrapper = await mountCard({ product: soldOutProduct })
    expect(wrapper.find('.sold-out-overlay').exists()).toBe(true)
    expect(wrapper.find('.sold-out-text').text()).toBe('已售罄')
  })

  it('does not show sold out overlay when stock > 0', async () => {
    const wrapper = await mountCard()
    expect(wrapper.find('.sold-out-overlay').exists()).toBe(false)
  })

  it('hides add to cart button and shows sold out tag when stock is 0', async () => {
    const soldOutProduct = { ...mockProduct, stock: 0 }
    const wrapper = await mountCard({ product: soldOutProduct })
    expect(wrapper.find('.add-cart-btn').exists()).toBe(false)
    expect(wrapper.find('.sold-out-tag').exists()).toBe(true)
    expect(wrapper.find('.sold-out-tag').text()).toBe('已售罄')
  })
})
