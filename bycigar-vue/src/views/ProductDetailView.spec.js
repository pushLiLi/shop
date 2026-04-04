import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { flushPromises } from '@vue/test-utils'
import ProductDetailView from './ProductDetailView.vue'

const mockProduct = {
  id: 1,
  name: 'Premium Cigar',
  price: 49.99,
  imageUrl: '/cigar.jpg',
  brand: 'Cuba Brand',
  description: 'A fine cigar',
  category: { id: 1, name: '雪茄', slug: 'cigars' },
  stock: 10
}

const mockRelatedProducts = {
  products: [
    { id: 2, name: 'Related Cigar', price: 29.99, imageUrl: '/related.jpg', stock: 5 },
    { id: 3, name: 'Another Cigar', price: 19.99, imageUrl: '/another.jpg', stock: 3, category: { id: 1 } }
  ]
}

function setupFetch(options = {}) {
  return vi.fn((url) => {
    if (url.includes('/api/products/') && !url.includes('?')) {
      const response = options.productResponse || mockProduct
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(response)
      })
    }
    if (url.includes('/api/products?')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(options.relatedResponse || mockRelatedProducts)
      })
    }
    return Promise.resolve({ ok: false, json: () => Promise.resolve({ error: 'Not found' }) })
  })
}

async function mountDetail(options = {}) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/products/:id', component: { template: '<div />' } },
      { path: '/', component: { template: '<div />' } },
      { path: '/category/:slug', component: { template: '<div />' } },
      { path: '/login', component: { template: '<div />' } }
    ]
  })

  const pinia = createPinia()
  setActivePinia(pinia)

  if (options.loggedIn) {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({
      role: 'customer', name: 'Test User', email: 'test@test.com'
    }))
  }

  global.fetch = setupFetch(options)

  router.push('/products/1')
  await router.isReady()
  await flushPromises()

  const wrapper = mount(ProductDetailView, {
    global: {
      plugins: [router, pinia],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        },
        ProductCard: {
          template: '<div class="product-card-stub">{{ product.name }}</div>',
          props: ['product']
        }
      }
    }
  })

  await flushPromises()
  return { wrapper, router, pinia }
}

describe('ProductDetailView', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows loading state initially', async () => {
    global.fetch = vi.fn(() => new Promise(() => {}))

    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/products/:id', component: { template: '<div />' } },
        { path: '/', component: { template: '<div />' } }
      ]
    })
    const pinia = createPinia()
    setActivePinia(pinia)

    router.push('/products/1')
    await router.isReady()
    await flushPromises()

    const wrapper = mount(ProductDetailView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product'] }
        }
      }
    })

    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.find('.loading').text()).toBe('加载中...')
  })

  it('renders product name', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.product-title').text()).toBe('Premium Cigar')
  })

  it('renders product price formatted', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.product-price-main').text()).toBe('¥49.99')
  })

  it('renders product brand', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.product-brand').text()).toBe('Cuba Brand')
  })

  it('renders product description', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.product-description').text()).toContain('A fine cigar')
  })

  it('renders product image', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.main-image img').attributes('src')).toBe('/cigar.jpg')
  })

  it('renders breadcrumb with category link', async () => {
    const { wrapper } = await mountDetail()
    const breadcrumb = wrapper.find('.breadcrumb')
    expect(breadcrumb.text()).toContain('首页')
    expect(breadcrumb.text()).toContain('雪茄')
    expect(breadcrumb.text()).toContain('Premium Cigar')
  })

  it('renders quantity selector with default 1', async () => {
    const { wrapper } = await mountDetail()
    const qtyInput = wrapper.find('.qty-input')
    expect(qtyInput.element.value).toBe('1')
  })

  it('increases quantity on plus click', async () => {
    const { wrapper } = await mountDetail()
    const btns = wrapper.findAll('.qty-btn')
    await btns[1].trigger('click')
    const qtyInput = wrapper.find('.qty-input')
    expect(qtyInput.element.value).toBe('2')
  })

  it('decreases quantity on minus click (min 1)', async () => {
    const { wrapper } = await mountDetail()
    const btns = wrapper.findAll('.qty-btn')
    await btns[1].trigger('click')
    await btns[0].trigger('click')
    const qtyInput = wrapper.find('.qty-input')
    expect(qtyInput.element.value).toBe('1')
  })

  it('does not decrease below 1', async () => {
    const { wrapper } = await mountDetail()
    const btns = wrapper.findAll('.qty-btn')
    await btns[0].trigger('click')
    expect(wrapper.find('.qty-input').element.value).toBe('1')
  })

  it('renders add to cart button', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.buy-btn').text()).toBe('加入购物车')
  })

  it('renders favorite button', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.favorite-btn').exists()).toBe(true)
  })

  it('renders related products', async () => {
    const { wrapper } = await mountDetail()
    expect(wrapper.find('.related-section').exists()).toBe(true)
    expect(wrapper.text()).toContain('相关产品')
    const cards = wrapper.findAll('.product-card-stub')
    expect(cards.length).toBeGreaterThan(0)
  })

  it('shows error state on fetch failure', async () => {
    global.fetch = vi.fn(() => Promise.resolve({
      ok: false,
      json: () => Promise.resolve({ error: 'Not found' })
    }))

    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/products/:id', component: { template: '<div />' } },
        { path: '/', component: { template: '<div />' } }
      ]
    })
    const pinia = createPinia()
    setActivePinia(pinia)

    router.push('/products/999')
    await router.isReady()
    await flushPromises()

    const wrapper = mount(ProductDetailView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product'] }
        }
      }
    })
    await flushPromises()

    expect(wrapper.find('.error').exists()).toBe(true)
  })
})
