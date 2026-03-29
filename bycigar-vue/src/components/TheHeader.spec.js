import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import TheHeader from './TheHeader.vue'

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/login', component: { template: '<div />' } },
      { path: '/search', component: { template: '<div />' } },
      { path: '/profile', component: { template: '<div />' } },
      { path: '/orders', component: { template: '<div />' } },
      { path: '/favorites', component: { template: '<div />' } },
      { path: '/admin', component: { template: '<div />' } }
    ]
  })
}

async function mountHeader(options = {}) {
  const pinia = options.pinia || createPinia()
  setActivePinia(pinia)
  const router = options.router || createTestRouter()

  const wrapper = mount(TheHeader, {
    global: {
      plugins: [router, pinia],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        },
        Transition: {
          template: '<slot />'
        }
      }
    }
  })
  await router.isReady()
  return wrapper
}

describe('TheHeader', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows logo text', async () => {
    const wrapper = await mountHeader()
    expect(wrapper.text()).toContain('HUAUGE')
  })

  it('shows login link when not logged in', async () => {
    const wrapper = await mountHeader()
    const loginLinks = wrapper.findAll('a[href="/login"]')
    expect(loginLinks.length).toBeGreaterThanOrEqual(1)
  })

  it('shows user name when logged in', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'TestUser', role: 'customer' }))
    const wrapper = await mountHeader()
    expect(wrapper.text()).toContain('TestUser')
  })

  it('shows admin link when user is admin', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'Admin', role: 'admin' }))
    const wrapper = await mountHeader()

    const userBtn = wrapper.find('.user-btn')
    await userBtn.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.admin-link').exists()).toBe(true)
    expect(wrapper.text()).toContain('后台管理')
  })

  it('does not show admin link for regular user', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'User', role: 'customer' }))
    const wrapper = await mountHeader()

    const userBtn = wrapper.find('.user-btn')
    await userBtn.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.admin-link').exists()).toBe(false)
  })

  it('shows navigation menu items', async () => {
    const wrapper = await mountHeader()
    expect(wrapper.text()).toContain('首页')
    expect(wrapper.text()).toContain('全部商品')
    expect(wrapper.text()).toContain('关于我们')
  })

  it('has search form', async () => {
    const wrapper = await mountHeader()
    expect(wrapper.find('.search-form').exists()).toBe(true)
    expect(wrapper.find('.search-input').exists()).toBe(true)
  })

  it('navigates to search on form submit', async () => {
    const router = createTestRouter()
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = await mountHeader({ router })

    await wrapper.find('.search-input').setValue('cigar')
    await wrapper.find('.search-form').trigger('submit.prevent')

    expect(pushSpy).toHaveBeenCalledWith('/search?q=cigar')
  })

  it('does not navigate on empty search', async () => {
    const router = createTestRouter()
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = await mountHeader({ router })

    await wrapper.find('.search-input').setValue('   ')
    await wrapper.find('.search-form').trigger('submit.prevent')

    expect(pushSpy).not.toHaveBeenCalled()
  })

  it('shows cart badge with item count', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'User', role: 'customer' }))
    const pinia = createPinia()
    setActivePinia(pinia)

    const wrapper = await mountHeader({ pinia })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = [
      { id: 1 }, { id: 2 }, { id: 3 }
    ]
    await wrapper.vm.$nextTick()

    const badges = wrapper.findAll('.icon-badge')
    const cartBadge = badges.find(b => b.text() === '3')
    expect(cartBadge).toBeTruthy()
  })

  it('calls logout and navigates to home', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'User', role: 'customer' }))
    const router = createTestRouter()
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = await mountHeader({ router })

    await wrapper.find('.user-btn').trigger('click')
    await wrapper.vm.$nextTick()

    await wrapper.find('.logout-btn').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/')
  })

  it('toggles mobile menu', async () => {
    const wrapper = await mountHeader()
    expect(wrapper.find('.header-nav').classes()).not.toContain('is-open')

    await wrapper.find('.mobile-menu-btn').trigger('click')
    expect(wrapper.find('.header-nav').classes()).toContain('is-open')

    await wrapper.find('.mobile-menu-btn').trigger('click')
    expect(wrapper.find('.header-nav').classes()).not.toContain('is-open')
  })
})
