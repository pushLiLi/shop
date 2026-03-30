import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useContactMethodsStore = defineStore('contactMethods', () => {
  const methods = ref([])
  const loading = ref(false)

  async function fetchMethods() {
    loading.value = true
    try {
      const res = await fetch('/api/contact-methods')
      const data = await res.json()
      methods.value = data.contactMethods || []
    } catch (e) {
      console.error('Failed to fetch contact methods:', e)
    } finally {
      loading.value = false
    }
  }

  return { methods, loading, fetchMethods }
})
