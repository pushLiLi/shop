import { defineStore } from 'pinia'

const API_BASE = 'http://localhost:3000/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: {},
    loading: false,
    error: null
  }),

  getters: {
    getSetting: (state) => (key, defaultValue = '') => {
      return state.settings[key] || defaultValue
    },
    footerDescription: (state) => {
      return state.settings['footer_description'] || ''
    },
    footerServiceTime: (state) => {
      return state.settings['footer_service_time'] || ''
    }
  },

  actions: {
    async fetchSettings() {
      this.loading = true
      this.error = null
      try {
        const response = await fetch(`${API_BASE}/settings`)
        const data = await response.json()
        if (data.success) {
          this.settings = data.data || {}
        } else {
          this.error = data.message || '获取设置失败'
        }
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async updateSetting(key, value) {
      this.error = null
      try {
        const response = await fetch(`${API_BASE}/admin/settings/${key}`, {
          method: 'PUT',
          headers: getAuthHeaders(),
          body: JSON.stringify({ value })
        })
        const data = await response.json()
        if (data.success) {
          this.settings[key] = value
          return true
        } else {
          this.error = data.message || '更新设置失败'
          return false
        }
      } catch (err) {
        this.error = err.message
        return false
      }
    }
  }
})
