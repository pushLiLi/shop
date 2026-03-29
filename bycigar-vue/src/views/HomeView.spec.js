import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { flushPromises } from '@vue/test-utils'
import HomeView from './HomeView.vue'

function mockFetch(responses) {
  const urls = Object.keys(responses)
  return vi.fn((url) => {
    for (const key of urls) {
      if (url.includes(key)) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(responses[key])
        })
      }
    }
    return Promise.resolve({
      ok: true,
      json: () => Promise.resolve({})
    })
  })
}

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/category/:slug', component: { template: '<div />' } }
    ]
  })
}

const defaultResponses = {
  '/api/config': {},
  '/api/products?featured=true&limit=12': {
    products: [
      { id: 1, name: 'Featured Cigar', price: 29.99, imageUrl: '/img1.jpg', stock: 10 },
      { id: 2, name: 'Featured Cigar 2', price: 39.99, imageUrl: '/img2.jpg', stock: 5 }
    ]
  },
  '/api/banners': [],
  '/api/categories': []
}

async function mountHome(responses = {}) {
  const allResponses = { ...defaultResponses, ...responses }
  global.fetch = mockFetch(allResponses)

  const router = createTestRouter()
  const pinia = createPinia()
  setActivePinia(pinia)

  const wrapper = mount(HomeView, {
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
  await router.isReady()
  await flushPromises()
  return wrapper
}

describe('HomeView', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    localStorage.clear()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('shows loading state initially', async () => {
    let resolvePromise
    global.fetch = vi.fn(() => new Promise(resolve => { resolvePromise = resolve }))

    const router = createTestRouter()
    const pinia = createPinia()
    setActivePinia(pinia)

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product'] }
        }
      }
    })
    await router.isReady()

    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.find('.loading').text()).toBe('加载中...')
  })

  it('renders featured products after loading', async () => {
    const wrapper = await mountHome()
    expect(wrapper.findAll('.product-card-stub')).toHaveLength(2)
    expect(wrapper.text()).toContain('Featured Cigar')
    expect(wrapper.text()).toContain('Featured Cigar 2')
  })

  it('shows default section title when no config', async () => {
    const wrapper = await mountHome()
    expect(wrapper.text()).toContain('特别推荐')
  })

  it('shows custom section title from config', async () => {
    const wrapper = await mountHome({
      '/api/config': { home_featured_title: '热销商品' }
    })
    expect(wrapper.text()).toContain('热销商品')
  })

  it('renders banner slider with dots', async () => {
    const wrapper = await mountHome()
    const dots = wrapper.findAll('.dot')
    expect(dots.length).toBe(3)
  })

  it('renders slider prev/next buttons', async () => {
    const wrapper = await mountHome()
    expect(wrapper.find('.slider-btn.prev').exists()).toBe(true)
    expect(wrapper.find('.slider-btn.next').exists()).toBe(true)
  })

  it('renders API banners when returned', async () => {
    const wrapper = await mountHome({
      '/api/banners': [
        { imageUrl: '/banner1.jpg', title: 'Banner 1' },
        { imageUrl: '/banner2.jpg', title: 'Banner 2' }
      ]
    })
    expect(wrapper.findAll('.dot')).toHaveLength(2)
  })

  it('renders category sections with products', async () => {
    const wrapper = await mountHome({
      '/api/categories': [
        { id: 1, name: '雪茄', slug: 'cigars', _count: 5 }
      ],
      '/api/products?categoryId=1&limit=8': {
        products: [
          { id: 10, name: 'Cigar A', price: 15, imageUrl: '/a.jpg', stock: 3 }
        ]
      }
    })
    expect(wrapper.text()).toContain('雪茄')
    expect(wrapper.text()).toContain('查看更多')
    expect(wrapper.text()).toContain('Cigar A')
  })

  it('renders category section with link to category page', async () => {
    const wrapper = await mountHome({
      '/api/categories': [
        { id: 1, name: '雪茄', slug: 'cigars', _count: 5 }
      ],
      '/api/products?categoryId=1&limit=8': {
        products: [{ id: 10, name: 'Cigar A', price: 15, imageUrl: '/a.jpg', stock: 3 }]
      }
    })
    const links = wrapper.findAll('a[href="/category/cigars"]')
    expect(links.length).toBeGreaterThan(0)
  })

  it('does not render optional banner when no config', async () => {
    const wrapper = await mountHome()
    expect(wrapper.find('.banner-section').exists()).toBe(false)
  })

  it('renders optional banner when config has home_banner_1', async () => {
    const wrapper = await mountHome({
      '/api/config': { home_banner_1: '/promo.jpg' }
    })
    expect(wrapper.find('.banner-section').exists()).toBe(true)
    expect(wrapper.find('.full-width-banner').attributes('src')).toBe('/promo.jpg')
  })

  it('does not render category section with zero products', async () => {
    const wrapper = await mountHome({
      '/api/categories': [
        { id: 1, name: 'Empty', slug: 'empty', _count: 0 }
      ]
    })
    expect(wrapper.text()).not.toContain('Empty')
  })

  it('hides loading after data loads', async () => {
    const wrapper = await mountHome()
    expect(wrapper.find('.loading').exists()).toBe(false)
  })
})
