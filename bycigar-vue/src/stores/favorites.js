import { defineStore } from 'pinia'

const API_BASE = '/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useFavoritesStore = defineStore('favorites', {
  state: () => ({
    items: [],
    loading: false
  }),
  
  getters: {
    count: (state) => state.items.length
  },
  
  actions: {
    setItems(items) {
      this.items = items
    },
    
    async fetchFavorites() {
      try {
        this.loading = true
        const res = await fetch(`${API_BASE}/favorites`, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        this.items = data.items || []
      } catch (e) {
        console.error('Fetch favorites failed:', e)
      } finally {
        this.loading = false
      }
    },
    
    async addItem(product) {
      try {
        const res = await fetch(`${API_BASE}/favorites`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({ productId: product.id })
        })
        const data = await res.json()
        if (data.success) {
          await this.fetchFavorites()
        }
      } catch (e) {
        console.error('Add favorite failed:', e)
      }
    },
    
    async removeItem(productId) {
      try {
        const res = await fetch(`${API_BASE}/favorites/${productId}`, {
          method: 'DELETE',
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (data.success) {
          await this.fetchFavorites()
        }
      } catch (e) {
        console.error('Remove favorite failed:', e)
      }
    },
    
    clear() {
      this.items = []
    }
  }
})
