import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useFavoritesStore } from '../stores/favorites'

describe('favorites store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('initial state is empty', () => {
    const store = useFavoritesStore()
    expect(store.items).toEqual([])
    expect(store.count).toBe(0)
  })

  it('count returns items length', () => {
    const store = useFavoritesStore()
    store.setItems([{ id: 1, product: { id: 1 } }, { id: 2, product: { id: 2 } }])
    expect(store.count).toBe(2)
  })

  it('fetchFavorites updates items', async () => {
    const mockItems = [{ id: 1, product: { id: 1 } }]
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ items: mockItems })
    })
    localStorage.setItem('token', 'test')

    const store = useFavoritesStore()
    await store.fetchFavorites()

    expect(store.items).toEqual(mockItems)
    expect(store.loading).toBe(false)
  })

  it('addItem calls API then fetchFavorites', async () => {
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ json: () => Promise.resolve({ success: true }) })
      .mockResolvedValueOnce({ json: () => Promise.resolve({ items: [{ id: 1, product: { id: 5 } }] }) })

    localStorage.setItem('token', 'test')
    const store = useFavoritesStore()
    await store.addItem({ id: 5 })

    expect(global.fetch).toHaveBeenCalledTimes(2)
    expect(global.fetch).toHaveBeenNthCalledWith(1, '/api/favorites', expect.objectContaining({ method: 'POST' }))
  })

  it('removeItem calls API then fetchFavorites', async () => {
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ json: () => Promise.resolve({ success: true }) })
      .mockResolvedValueOnce({ json: () => Promise.resolve({ items: [] }) })

    localStorage.setItem('token', 'test')
    const store = useFavoritesStore()
    await store.removeItem(5)

    expect(global.fetch).toHaveBeenCalledTimes(2)
    expect(global.fetch).toHaveBeenNthCalledWith(1, '/api/favorites/5', expect.objectContaining({ method: 'DELETE' }))
  })

  it('clear empties items', () => {
    const store = useFavoritesStore()
    store.setItems([{ id: 1 }])
    store.clear()
    expect(store.items).toEqual([])
  })
})
