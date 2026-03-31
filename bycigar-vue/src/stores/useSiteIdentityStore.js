import { defineStore } from 'pinia'
import { ref } from 'vue'

const API_BASE = '/api'

export const useSiteIdentityStore = defineStore('siteIdentity', () => {
  const identity = ref({
    title: 'BYCIGAR | 权威正品雪茄在线购买商城',
    metaDescription: '',
    faviconUrl: '/favicon.png'
  })

  const applied = ref(false)

  function applyToDocument(data) {
    if (data.title) {
      document.title = data.title
    }
    if (data.faviconUrl) {
      const favicon = document.querySelector("link[rel='icon']")
      if (favicon) {
        favicon.href = data.faviconUrl
      }
    }
    applied.value = true
  }

  async function fetchSiteIdentity() {
    try {
      const res = await fetch(`${API_BASE}/site-identity`)
      if (res.ok) {
        const data = await res.json()
        identity.value = data
        applyToDocument(data)
      }
    } catch (e) {
      console.warn('Failed to fetch site identity:', e)
    }
  }

  return {
    identity,
    applied,
    fetchSiteIdentity,
    applyToDocument
  }
})
