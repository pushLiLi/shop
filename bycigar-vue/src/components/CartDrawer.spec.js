import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import CartDrawer from './CartDrawer.vue'

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/products/:id', component: { template: '<div />' } },
      { path: '/checkout', component: { template: '<div />' } }
    ]
  })
}

function mountDrawer(cartState = {}) {
  const pinia = createPinia()
  setActivePinia(pinia)
  const router = createTestRouter()

  const wrapper = mount(CartDrawer, {
    global: {
      plugins: [router, pinia],
      stubs: {
        Teleport: {
          template: '<div><slot /></div>'
        },
        Transition: {
          template: '<slot />'
        }
      }
    }
  })

  if (cartState.items) {
    const { useCartStore } = require('../stores/cart')
    const store = useCartStore()
    store.items = cartState.items
    store.isOpen = cartState.isOpen !== undefined ? cartState.isOpen : true
    store.loading = cartState.loading || false
  }

  return { wrapper, pinia, router }
}

describe('CartDrawer', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('is hidden when cart is closed', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.isOpen = false
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.cart-drawer-overlay').exists()).toBe(false)
  })

  it('shows empty cart message when no items', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = []
    cartStore.isOpen = true
    cartStore.loading = false
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.empty-cart').exists()).toBe(true)
    expect(wrapper.text()).toContain('购物车是空的')
  })

  it('shows cart items when items exist', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = [{
      id: 1,
      productId: 10,
      quantity: 2,
      product: { name: 'Cigar A', price: 25.00, imageUrl: '/media/a.jpg' }
    }]
    cartStore.isOpen = true
    cartStore.loading = false
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.cart-items').exists()).toBe(true)
    expect(wrapper.text()).toContain('Cigar A')
    expect(wrapper.text()).toContain('¥25.00')
  })

  it('shows subtotal per item', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = [{
      id: 1,
      productId: 10,
      quantity: 3,
      product: { name: 'Cigar A', price: 10.00, imageUrl: '/media/a.jpg' }
    }]
    cartStore.isOpen = true
    cartStore.loading = false
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('¥30.00')
  })

  it('shows total in footer', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = [
      { id: 1, productId: 10, quantity: 2, product: { name: 'A', price: 20.00, imageUrl: '/a.jpg' } },
      { id: 2, productId: 11, quantity: 1, product: { name: 'B', price: 15.00, imageUrl: '/b.jpg' } }
    ]
    cartStore.isOpen = true
    cartStore.loading = false
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.drawer-footer').exists()).toBe(true)
    expect(wrapper.text()).toContain('¥55.00')
  })

  it('has close button', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = []
    cartStore.isOpen = true
    cartStore.loading = false
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.close-btn').exists()).toBe(true)
  })

  it('calls closeCart when close button clicked', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = []
    cartStore.isOpen = true
    cartStore.loading = false
    const closeSpy = vi.spyOn(cartStore, 'closeCart')
    await wrapper.vm.$nextTick()

    await wrapper.find('.close-btn').trigger('click')
    expect(closeSpy).toHaveBeenCalled()
  })

  it('shows loading state', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = []
    cartStore.isOpen = true
    cartStore.loading = true
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.text()).toContain('加载中')
  })

  it('navigates to checkout on checkout button click', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createTestRouter()
    const pushSpy = vi.spyOn(router, 'push')

    const wrapper = mount(CartDrawer, {
      global: {
        plugins: [router, pinia],
        stubs: {
          Teleport: { template: '<div><slot /></div>' },
          Transition: { template: '<slot />' }
        }
      }
    })

    const { useCartStore } = await import('../stores/cart')
    const cartStore = useCartStore()
    cartStore.items = [
      { id: 1, productId: 10, quantity: 1, product: { name: 'A', price: 10.00, imageUrl: '/a.jpg' } }
    ]
    cartStore.isOpen = true
    cartStore.loading = false
    await wrapper.vm.$nextTick()

    await wrapper.find('.checkout-btn').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/checkout')
  })
})
