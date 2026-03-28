import { defineStore } from 'pinia'
import { useToastStore } from './toast'

const API_BASE = '/api'
const pendingUpdates = new Map()

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
    loading: false,
    isOpen: false
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
        this.removeItem(cartItemId)
        return
      }
      
      this.items = this.items.map(item =>
        item.id === cartItemId ? { ...item, quantity } : item
      )
      
      if (pendingUpdates.has(cartItemId)) {
        clearTimeout(pendingUpdates.get(cartItemId))
      }
      
      pendingUpdates.set(cartItemId, setTimeout(async () => {
        pendingUpdates.delete(cartItemId)
        try {
          const res = await fetch(`${API_BASE}/cart/${cartItemId}`, {
            method: 'PUT',
            headers: getAuthHeaders(),
            body: JSON.stringify({ quantity })
          })
          const data = await res.json()
          if (!data.success) {
            useToastStore().error('更新失败')
          }
        } catch (e) {
          useToastStore().error('更新失败')
        }
      }, 300))
    },
    
    removeItem(cartItemId) {
      if (pendingUpdates.has(cartItemId)) {
        clearTimeout(pendingUpdates.get(cartItemId))
        pendingUpdates.delete(cartItemId)
      }
      
      const oldItem = this.items.find(item => item.id === cartItemId)
      this.items = this.items.filter(item => item.id !== cartItemId)
      
      fetch(`${API_BASE}/cart/${cartItemId}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
      })
        .then(res => res.json())
        .then(data => {
          if (!data.success && oldItem) {
            this.items.push(oldItem)
            useToastStore().error('删除失败')
          }
        })
        .catch(() => {
          if (oldItem) {
            this.items.push(oldItem)
            useToastStore().error('删除失败')
          }
        })
    },
    
    clear() {
      this.items = []
    },
    
    openCart() {
      this.isOpen = true
    },
    
    closeCart() {
      this.isOpen = false
    },
    
    toggleCart() {
      this.isOpen = !this.isOpen
    }
  }
})
