import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../stores/auth'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('isLoggedIn returns false when no token/user', () => {
    const store = useAuthStore()
    expect(store.isLoggedIn).toBe(false)
  })

  it('isLoggedIn returns true when token and user exist', () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, email: 'a@b.com', role: 'customer' }))
    const store = useAuthStore()
    expect(store.isLoggedIn).toBe(true)
  })

  it('isAdmin returns true for admin role', () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, role: 'admin' }))
    const store = useAuthStore()
    expect(store.isAdmin).toBe(true)
  })

  it('isAdmin returns false for customer role', () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ id: 1, role: 'customer' }))
    const store = useAuthStore()
    expect(store.isAdmin).toBe(false)
  })

  it('login success writes token/user to state and localStorage', async () => {
    const mockUser = { id: 1, email: 'test@test.com', name: 'Test', role: 'customer' }
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ token: 'jwt-token', user: mockUser })
    })

    const store = useAuthStore()
    await store.login('test@test.com', 'password123')

    expect(store.token).toBe('jwt-token')
    expect(store.user).toEqual(mockUser)
    expect(localStorage.getItem('token')).toBe('jwt-token')
    expect(JSON.parse(localStorage.getItem('user'))).toEqual(mockUser)
  })

  it('login failure throws error', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: '邮箱或密码错误' })
    })

    const store = useAuthStore()
    await expect(store.login('test@test.com', 'wrong')).rejects.toThrow('邮箱或密码错误')
  })

  it('login failure with requireCaptcha', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: '密码错误', requireCaptcha: true })
    })

    const store = useAuthStore()
    try {
      await store.login('test@test.com', 'wrong')
    } catch (err) {
      expect(err.requireCaptcha).toBe(true)
      expect(err.message).toBe('密码错误')
    }
  })

  it('register success auto-logs in', async () => {
    const mockUser = { id: 2, email: 'new@test.com', name: 'New', role: 'customer' }
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: () => Promise.resolve({ token: 'reg-token', user: mockUser })
    })

    const store = useAuthStore()
    await store.register('new@test.com', 'password123', 'New', 'cap-id', '1234')

    expect(store.token).toBe('reg-token')
    expect(store.user.name).toBe('New')
  })

  it('register failure throws error', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: '验证码错误' })
    })

    const store = useAuthStore()
    await expect(store.register('new@test.com', '123456', 'New', 'cap-id', 'wrong'))
      .rejects.toThrow('验证码错误')
  })

  it('logout clears token/user and localStorage', () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ id: 1 }))

    const store = useAuthStore()
    store.logout()

    expect(store.token).toBe('')
    expect(store.user).toBeNull()
    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('user')).toBeNull()
  })

  it('validateToken success updates user', async () => {
    const mockUser = { id: 1, email: 'test@test.com', name: 'Updated', role: 'customer' }
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ user: mockUser })
    })

    localStorage.setItem('token', 'valid-token')
    const store = useAuthStore()
    const result = await store.validateToken()

    expect(result).toBe(true)
    expect(store.user.name).toBe('Updated')
  })

  it('validateToken failure calls logout', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: '无效的token' })
    })

    localStorage.setItem('token', 'bad-token')
    const store = useAuthStore()
    const result = await store.validateToken()

    expect(result).toBe(false)
    expect(store.token).toBe('')
    expect(store.user).toBeNull()
  })

  it('getAuthHeaders returns Bearer token', () => {
    localStorage.setItem('token', 'my-token')
    const store = useAuthStore()
    const headers = store.getAuthHeaders()

    expect(headers.Authorization).toBe('Bearer my-token')
    expect(headers['Content-Type']).toBe('application/json')
  })

  it('userName computed from user name', () => {
    localStorage.setItem('token', 't')
    localStorage.setItem('user', JSON.stringify({ name: '张三', role: 'customer' }))
    const store = useAuthStore()
    expect(store.userName).toBe('张三')
  })

  it('userName falls back to email prefix', () => {
    localStorage.setItem('token', 't')
    localStorage.setItem('user', JSON.stringify({ email: 'user@test.com', role: 'customer' }))
    const store = useAuthStore()
    expect(store.userName).toBe('user')
  })

  it('userName falls back to default', () => {
    const store = useAuthStore()
    expect(store.userName).toBe('用户')
  })
})
