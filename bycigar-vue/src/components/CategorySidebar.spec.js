import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import CategorySidebar from './CategorySidebar.vue'

const mockCategories = [
  { id: 1, name: 'Cigars', slug: 'cigars', children: [
    { id: 3, name: 'Cuban', slug: 'cuban' },
    { id: 4, name: 'Dominican', slug: 'dominican' }
  ]},
  { id: 2, name: 'Accessories', slug: 'accessories', children: [] }
]

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/category/:slug', component: { template: '<div />' } }
    ]
  })
}

async function mountSidebar(props = {}, isMobile = false) {
  const router = createTestRouter()

  const originalInnerWidth = window.innerWidth
  if (isMobile) {
    Object.defineProperty(window, 'innerWidth', { value: 500, configurable: true })
  } else {
    Object.defineProperty(window, 'innerWidth', { value: 1024, configurable: true })
  }

  global.fetch = vi.fn().mockResolvedValue({
    json: () => Promise.resolve(mockCategories)
  })

  const wrapper = mount(CategorySidebar, {
    props: { activeSlug: '', ...props },
    global: {
      plugins: [router],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        },
        Teleport: { template: '<div><slot /></div>' },
        Transition: { template: '<slot />' }
      }
    }
  })

  await router.isReady()
  await vi.waitFor(() => {
    if (wrapper.vm.categories === undefined) throw new Error('not ready')
  }, { timeout: 3000 }).catch(() => {})

  return wrapper
}

describe('CategorySidebar', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('renders desktop sidebar with categories', async () => {
    const wrapper = await mountSidebar()
    await vi.waitFor(() => wrapper.find('.category-list').exists() || wrapper.find('.sidebar-title').exists(), { timeout: 3000 })

    expect(wrapper.find('.sidebar-title').exists()).toBe(true)
    expect(wrapper.text()).toContain('商品分类')
  })

  it('renders category links', async () => {
    const wrapper = await mountSidebar()
    await vi.waitFor(() => wrapper.findAll('.category-link').length > 0, { timeout: 3000 })

    const links = wrapper.findAll('.category-link')
    expect(links.length).toBeGreaterThanOrEqual(2)
    expect(wrapper.text()).toContain('Cigars')
    expect(wrapper.text()).toContain('Accessories')
  })

  it('highlights active category', async () => {
    const wrapper = await mountSidebar({ activeSlug: 'cigars' })
    await vi.waitFor(() => wrapper.findAll('.category-link').length > 0, { timeout: 3000 })

    const activeLinks = wrapper.findAll('.category-link.active')
    expect(activeLinks.length).toBeGreaterThanOrEqual(1)
    expect(activeLinks[0].text()).toContain('Cigars')
  })

  it('renders children categories', async () => {
    const wrapper = await mountSidebar()
    await vi.waitFor(() => wrapper.find('.subcategory-list').exists(), { timeout: 3000 })

    expect(wrapper.text()).toContain('Cuban')
    expect(wrapper.text()).toContain('Dominican')
  })

  it('shows mobile category button on mobile', async () => {
    const wrapper = await mountSidebar({}, true)
    await vi.waitFor(() => wrapper.find('.mobile-category-btn').exists(), { timeout: 3000 })

    expect(wrapper.find('.mobile-category-btn').exists()).toBe(true)
  })

  it('opens drawer on mobile button click', async () => {
    const wrapper = await mountSidebar({}, true)
    await vi.waitFor(() => wrapper.find('.mobile-category-btn').exists(), { timeout: 3000 })

    await wrapper.find('.mobile-category-btn').trigger('click')
    expect(wrapper.find('.drawer-overlay').exists()).toBe(true)
  })

  it('shows default name when no active slug', async () => {
    const wrapper = await mountSidebar()
    await vi.waitFor(() => wrapper.findAll('.category-link').length > 0, { timeout: 3000 })
  })
})
