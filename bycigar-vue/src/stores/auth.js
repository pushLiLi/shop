import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const API_BASE = '/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const userName = computed(() => user.value?.name || user.value?.email?.split('@')[0] || '用户')

  async function login(email, password, captchaId, captchaCode) {
    const body = { email, password }
    if (captchaId && captchaCode) {
      body.captchaId = captchaId
      body.captchaCode = captchaCode
    }

    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })
    
    const data = await response.json()
    
    if (!response.ok) {
      const err = new Error(data.error || '登录失败')
      err.requireCaptcha = !!data.requireCaptcha
      throw err
    }
    
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
    
    return data.user
  }

  async function register(email, password, name, captchaId, captchaCode) {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password, name, captchaId, captchaCode })
    })
    
    const data = await response.json()
    
    if (!response.ok) {
      throw new Error(data.error || '注册失败')
    }
    
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
    
    return data.user
  }

  async function validateToken() {
    if (!token.value) return false
    
    try {
      const res = await fetch(`${API_BASE}/auth/me`, {
        headers: { 'Authorization': `Bearer ${token.value}` }
      })
      
      if (!res.ok) {
        logout()
        return false
      }
      
      const data = await res.json()
      user.value = data.user
      localStorage.setItem('user', JSON.stringify(data.user))
      return true
    } catch (e) {
      logout()
      return false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    userName,
    login,
    register,
    validateToken,
    logout,
    getAuthHeaders: () => ({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token.value}`
    })
  }
})
