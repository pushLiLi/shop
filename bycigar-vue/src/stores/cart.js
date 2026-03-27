import { defineStore } from 'pinia'

const API_BASE = 'http://localhost:3000/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [],
    loading: false
  }),
  
  getters: {
    total: (state) => {
      return state.items.reduce((sum, item) => {
        return sum + (item.product?.price || 0) * item.quantity
      }, 0)
    },
    
    count: (state) => {
      return state.items.reduce((sum, item) => sum + item.quantity, 0)
    }
  },
  
  actions: {
    setItems(items) {
      this.items = items
    },
    
    async fetchCart() {
      try {
        this.loading = true
        const res = await fetch(`${API_BASE}/cart`, {
          headers: getAuthHeaders()
        })
        const data = await res.json()
        this.items = data.items || []
      } catch (e) {
        console.error('Fetch cart failed:', e)
      } finally {
        this.loading = false
      }
    },
    
    async addItem(product, quantity = 1) {
      try {
        const res = await fetch(`${API_BASE}/cart`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({ productId: product.id, quantity })
        })
        const data = await res.json()
        if (data.success) {
          await this.fetchCart()
        }
      } catch (e) {
        console.error('Add to cart failed:', e)
      }
    },
    
    async updateQuantity(cartItemId, quantity) {
      if (quantity < 1) {
        await this.removeItem(cartItemId)
        return
      }
      try {
        const res = await fetch(`${API_BASE}/cart/${cartItemId}`, {
          method: 'PUT',
          headers: getAuthHeaders(),
          body: JSON.stringify({ quantity })
        })
        const data = await res.json()
        if (data.success) {
          await this.fetchCart()
        }
      } catch (e) {
        console.error('Update quantity failed:', e)
      }
    },
    
    async removeItem(cartItemId) {
      try {
        const res = await fetch(`${API_BASE}/cart/${cartItemId}`, {
          method: 'DELETE',
          headers: getAuthHeaders()
        })
        const data = await res.json()
        if (data.success) {
          await this.fetchCart()
        }
      } catch (e) {
        console.error('Remove item failed:', e)
      }
    },
    
    clear() {
      this.items = []
    }
  }
})
