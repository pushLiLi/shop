import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { flushPromises } from '@vue/test-utils'
import SearchView from './SearchView.vue'

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/search', component: SearchView },
      { path: '/', component: { template: '<div />' } }
    ]
  })
}

const mockSearchResults = {
  products: [
    { id: 1, name: 'Cigar A', price: 29.99, imageUrl: '/a.jpg', stock: 10 },
    { id: 2, name: 'Cigar B', price: 39.99, imageUrl: '/b.jpg', stock: 5 }
  ],
  total: 2
}

async function mountSearch(options = {}) {
  const router = createTestRouter()
  const pinia = createPinia()
  setActivePinia(pinia)

  global.innerWidth = 1024

  global.fetch = vi.fn(() => Promise.resolve({
    ok: true,
    json: () => Promise.resolve(options.fetchResponse || mockSearchResults)
  }))

  const wrapper = mount(SearchView, {
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
        }
      }
    }
  })

  await router.push({ path: '/search', query: { q: options.query || 'cigar' } })
  await router.isReady()
  await flushPromises()
  return { wrapper, router }
}

describe('SearchView', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows search keyword in header', async () => {
    const { wrapper } = await mountSearch({ query: 'cigar' })
    expect(wrapper.find('.keyword').text()).toBe('cigar')
  })

  it('shows total result count', async () => {
    const { wrapper } = await mountSearch()
    expect(wrapper.text()).toContain('共 2 个结果')
  })

  it('renders search result products', async () => {
    const { wrapper } = await mountSearch()
    expect(wrapper.findAll('.product-card-stub')).toHaveLength(2)
    expect(wrapper.text()).toContain('Cigar A')
    expect(wrapper.text()).toContain('Cigar B')
  })

  it('renders sort buttons', async () => {
    const { wrapper } = await mountSearch()
    const sortBtns = wrapper.findAll('.sort-btn')
    expect(sortBtns).toHaveLength(3)
    expect(sortBtns[0].text()).toContain('最新')
    expect(sortBtns[1].text()).toContain('价格')
    expect(sortBtns[2].text()).toContain('名称')
  })

  it('shows no results message when empty', async () => {
    const { wrapper } = await mountSearch({
      fetchResponse: { products: [], total: 0 }
    })
    expect(wrapper.find('.no-results').exists()).toBe(true)
    expect(wrapper.find('.no-results').text()).toContain('未找到相关产品')
    expect(wrapper.text()).toContain('请尝试其他关键词')
  })

  it('shows empty state when no query', async () => {
    const router = createTestRouter()
    const pinia = createPinia()
    setActivePinia(pinia)

    global.fetch = vi.fn(() => Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ products: [], total: 0 })
    }))

    const wrapper = mount(SearchView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product', 'horizontal'] }
        }
      }
    })
    await router.push('/search')
    await router.isReady()
    await flushPromises()

    expect(wrapper.find('.empty-state').exists()).toBe(true)
    expect(wrapper.find('.empty-state').text()).toContain('请输入搜索关键词')
  })

  it('renders pagination when multiple pages', async () => {
    const { wrapper } = await mountSearch({
      fetchResponse: {
        products: [{ id: 1, name: 'P1', price: 10, imageUrl: '/1.jpg', stock: 1 }],
        total: 25
      }
    })
    expect(wrapper.find('.pagination').exists()).toBe(true)
  })

  it('does not show pagination for single page', async () => {
    const { wrapper } = await mountSearch()
    expect(wrapper.find('.pagination').exists()).toBe(false)
  })

  it('shows loading during search', async () => {
    let resolvePromise
    global.fetch = vi.fn(() => new Promise(resolve => { resolvePromise = resolve }))

    const router = createTestRouter()
    const pinia = createPinia()
    setActivePinia(pinia)

    const wrapper = mount(SearchView, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': { template: '<a><slot /></a>', props: [] },
          ProductCard: { template: '<div />', props: ['product', 'horizontal'] }
        }
      }
    })
    await router.push({ path: '/search', query: { q: 'test' } })
    await router.isReady()

    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.find('.loading').text()).toBe('搜索中...')
  })

  it('highlights active sort button', async () => {
    const { wrapper } = await mountSearch()
    const sortBtns = wrapper.findAll('.sort-btn')
    expect(sortBtns[0].classes()).toContain('active')
  })

  it('changes sort on button click', async () => {
    const { wrapper } = await mountSearch()
    const sortBtns = wrapper.findAll('.sort-btn')
    await sortBtns[1].trigger('click')
    await flushPromises()
    expect(sortBtns[1].classes()).toContain('active')
    expect(global.fetch).toHaveBeenCalled()
  })
})
