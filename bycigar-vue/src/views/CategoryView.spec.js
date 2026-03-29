import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { flushPromises } from '@vue/test-utils'
import CategoryView from './CategoryView.vue'

function createTestRouter(initialPath = '/category/cigars') {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/category/:slug', component: CategoryView },
      { path: '/', component: { template: '<div />' } }
    ]
  })
}

const mockProducts = {
  products: [
    { id: 1, name: 'Cigar A', price: 29.99, imageUrl: '/a.jpg', stock: 10 },
    { id: 2, name: 'Cigar B', price: 39.99, imageUrl: '/b.jpg', stock: 5 }
  ],
  total: 2,
  category: { id: 1, name: '雪茄', slug: 'cigars' }
}

async function mountCategory(options = {}) {
  const router = createTestRouter(options.initialPath)
  const pinia = createPinia()
  setActivePinia(pinia)

  global.fetch = vi.fn(() => Promise.resolve({
    ok: true,
    json: () => Promise.resolve(options.fetchResponse || mockProducts)
  }))

  const wrapper = mount(CategoryView, {
    global: {
      plugins: [router, pinia],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        },
        ProductCard: {
          template: '<div class="product-card-stub">{{ product.name }}</div>',
          props: ['product', 'horizontal']
        },
        CategorySidebar: {
          template: '<div class="sidebar-stub" />',
          props: ['activeSlug']
        }
      }
    }
  })

  await router.push(options.initialPath || '/category/cigars')
  await router.isReady()
  await flushPromises()
  return { wrapper, router, pinia }
}

describe('CategoryView', () => {
  beforeEach(() => {
    localStorage.clear()
    global.innerWidth = 1024
  })

  it('shows loading state on initial load', async () => {
    let resolvePromise
    global.fetch = vi.fn(() => new Promise(resolve => { resolvePromise = resolve }))

    const router = createTestRouter()
    const pinia = createPinia()
    setActivePinia(pinia)

    const wrapper = mount(CategoryView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product', 'horizontal'] },
          CategorySidebar: { template: '<div />', props: ['activeSlug'] }
        }
      }
    })
    await router.push('/category/cigars')
    await router.isReady()

    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.find('.loading').text()).toBe('加载中...')
  })

  it('renders category name from API', async () => {
    const { wrapper } = await mountCategory()
    expect(wrapper.find('.page-title').text()).toBe('雪茄')
  })

  it('shows product count', async () => {
    const { wrapper } = await mountCategory()
    expect(wrapper.find('.product-count').text()).toContain('共 2 个产品')
  })

  it('renders product cards', async () => {
    const { wrapper } = await mountCategory()
    expect(wrapper.findAll('.product-card-stub')).toHaveLength(2)
    expect(wrapper.text()).toContain('Cigar A')
    expect(wrapper.text()).toContain('Cigar B')
  })

  it('renders sort buttons', async () => {
    const { wrapper } = await mountCategory()
    const sortBtns = wrapper.findAll('.sort-btn')
    expect(sortBtns).toHaveLength(3)
    expect(sortBtns[0].text()).toContain('最新')
    expect(sortBtns[1].text()).toContain('价格')
    expect(sortBtns[2].text()).toContain('名称')
  })

  it('highlights active sort button', async () => {
    const { wrapper } = await mountCategory()
    const sortBtns = wrapper.findAll('.sort-btn')
    expect(sortBtns[0].classes()).toContain('active')
    expect(sortBtns[1].classes()).not.toContain('active')
  })

  it('toggles sort order on same sort click', async () => {
    const { wrapper } = await mountCategory()
    const sortBtns = wrapper.findAll('.sort-btn')
    expect(sortBtns[0].text()).toContain('↓')

    await sortBtns[0].trigger('click')
    await flushPromises()
    expect(sortBtns[0].text()).toContain('↑')
  })

  it('shows no products message when empty', async () => {
    const { wrapper } = await mountCategory({
      fetchResponse: { products: [], total: 0, category: { id: 1, name: '空分类', slug: 'empty' } }
    })
    expect(wrapper.find('.no-products').exists()).toBe(true)
    expect(wrapper.find('.no-products').text()).toBe('该分类暂无产品')
  })

  it('shows error message on fetch failure', async () => {
    global.fetch = vi.fn(() => Promise.reject(new Error('Network error')))

    const router = createTestRouter()
    const pinia = createPinia()
    setActivePinia(pinia)

    const wrapper = mount(CategoryView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product', 'horizontal'] },
          CategorySidebar: { template: '<div />', props: ['activeSlug'] }
        }
      }
    })
    await router.push('/category/cigars')
    await router.isReady()
    await flushPromises()

    expect(wrapper.find('.error').exists()).toBe(true)
    expect(wrapper.find('.error').text()).toContain('Network error')
  })

  it('renders pagination when total pages > 1', async () => {
    const { wrapper } = await mountCategory({
      fetchResponse: {
        products: [{ id: 1, name: 'P1', price: 10, imageUrl: '/1.jpg', stock: 1 }],
        total: 25,
        category: { id: 1, name: '雪茄', slug: 'cigars' }
      }
    })
    expect(wrapper.find('.pagination').exists()).toBe(true)
    expect(wrapper.findAll('.page-btn').length).toBeGreaterThan(0)
  })

  it('disables prev button on first page', async () => {
    const { wrapper } = await mountCategory({
      fetchResponse: {
        products: [{ id: 1, name: 'P1', price: 10, imageUrl: '/1.jpg', stock: 1 }],
        total: 25,
        category: { id: 1, name: '雪茄', slug: 'cigars' }
      }
    })
    const prevBtn = wrapper.findAll('.page-btn')[0]
    expect(prevBtn.attributes('disabled')).toBeDefined()
  })

  it('renders CategorySidebar with activeSlug', async () => {
    const { wrapper } = await mountCategory()
    const sidebar = wrapper.find('.sidebar-stub')
    expect(sidebar.exists()).toBe(true)
  })

  it('uses slug as name when category is null', async () => {
    const { wrapper } = await mountCategory({
      fetchResponse: { products: [], total: 0, category: null }
    })
    expect(wrapper.find('.page-title').text()).toBe('cigars')
  })
})
