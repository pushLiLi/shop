import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCartStore } from '../stores/cart'

describe('cart store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('initial state is empty', () => {
    const store = useCartStore()
    expect(store.items).toEqual([])
    expect(store.total).toBe(0)
    expect(store.count).toBe(0)
    expect(store.isOpen).toBe(false)
  })

  it('total computes sum of price * quantity', () => {
    const store = useCartStore()
    store.setItems([
      { id: 1, quantity: 2, product: { price: 100 } },
      { id: 2, quantity: 3, product: { price: 50 } }
    ])
    expect(store.total).toBe(350)
  })

  it('count computes sum of quantities', () => {
    const store = useCartStore()
    store.setItems([
      { id: 1, quantity: 2, product: { price: 100 } },
      { id: 2, quantity: 3, product: { price: 50 } }
    ])
    expect(store.count).toBe(5)
  })

  it('fetchCart updates items', async () => {
    const mockItems = [
      { id: 1, quantity: 1, product: { price: 100 } }
    ]
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ items: mockItems })
    })
    localStorage.setItem('token', 'test')

    const store = useCartStore()
    await store.fetchCart()

    expect(store.items).toEqual(mockItems)
    expect(store.loading).toBe(false)
  })

  it('fetchCart with empty response', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({})
    })

    const store = useCartStore()
    await store.fetchCart()

    expect(store.items).toEqual([])
  })

  it('addItem calls API then fetchCart', async () => {
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ json: () => Promise.resolve({ success: true }) })
      .mockResolvedValueOnce({ json: () => Promise.resolve({ items: [{ id: 1, quantity: 1, product: { price: 100 } }] }) })

    localStorage.setItem('token', 'test')
    const store = useCartStore()
    await store.addItem({ id: 5 }, 2)

    expect(global.fetch).toHaveBeenCalledTimes(2)
    expect(global.fetch).toHaveBeenNthCalledWith(1, '/api/cart', expect.objectContaining({ method: 'POST' }))
  })

  it('updateQuantity optimistically updates local state', () => {
    const store = useCartStore()
    store.setItems([{ id: 1, quantity: 1, product: { price: 100 } }])

    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true })
    })
    localStorage.setItem('token', 'test')

    store.updateQuantity(1, 5)

    expect(store.items[0].quantity).toBe(5)
  })

  it('updateQuantity with quantity < 1 calls removeItem', () => {
    const store = useCartStore()
    store.setItems([{ id: 1, quantity: 2, product: { price: 100 } }])

    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true })
    })
    localStorage.setItem('token', 'test')

    store.updateQuantity(1, 0)

    expect(store.items).toEqual([])
  })

  it('removeItem optimistically removes from state', () => {
    const store = useCartStore()
    store.setItems([{ id: 1, quantity: 1, product: { price: 100 } }])

    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true })
    })
    localStorage.setItem('token', 'test')

    store.removeItem(1)

    expect(store.items).toEqual([])
  })

  it('clear empties items', () => {
    const store = useCartStore()
    store.setItems([{ id: 1, quantity: 1, product: { price: 100 } }])
    store.clear()
    expect(store.items).toEqual([])
  })

  it('openCart/closeCart toggles isOpen', () => {
    const store = useCartStore()
    expect(store.isOpen).toBe(false)
    store.openCart()
    expect(store.isOpen).toBe(true)
    store.closeCart()
    expect(store.isOpen).toBe(false)
  })

  it('toggleCart flips isOpen', () => {
    const store = useCartStore()
    expect(store.isOpen).toBe(false)
    store.toggleCart()
    expect(store.isOpen).toBe(true)
    store.toggleCart()
    expect(store.isOpen).toBe(false)
  })
})
