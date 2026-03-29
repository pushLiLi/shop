import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSettingsStore } from '../stores/useSettingsStore'

describe('settings store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('initial state has empty settings', () => {
    const store = useSettingsStore()
    expect(store.settings).toEqual({})
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('fetchSettings updates settings from API', async () => {
    const mockSettings = { footer_description: '关于我们', footer_service_time: '9:00-18:00' }
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true, data: mockSettings })
    })

    const store = useSettingsStore()
    await store.fetchSettings()

    expect(store.settings).toEqual(mockSettings)
    expect(store.loading).toBe(false)
  })

  it('getSetting returns value by key', () => {
    const store = useSettingsStore()
    store.settings = { footer_description: '测试描述' }

    expect(store.getSetting('footer_description')).toBe('测试描述')
  })

  it('getSetting returns defaultValue for missing key', () => {
    const store = useSettingsStore()
    expect(store.getSetting('nonexistent', '默认')).toBe('默认')
  })

  it('footerDescription getter', () => {
    const store = useSettingsStore()
    store.settings = { footer_description: '雪茄描述' }
    expect(store.footerDescription).toBe('雪茄描述')
  })

  it('footerDescription defaults to empty string', () => {
    const store = useSettingsStore()
    expect(store.footerDescription).toBe('')
  })

  it('updateSetting calls API and updates local state', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true })
    })
    localStorage.setItem('token', 'test')

    const store = useSettingsStore()
    const result = await store.updateSetting('footer_description', '新描述')

    expect(result).toBe(true)
    expect(store.settings.footer_description).toBe('新描述')
    expect(global.fetch).toHaveBeenCalledWith(
      '/api/admin/settings/footer_description',
      expect.objectContaining({ method: 'PUT' })
    )
  })

  it('updateSetting failure returns false and sets error', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: false, message: '权限不足' })
    })
    localStorage.setItem('token', 'test')

    const store = useSettingsStore()
    const result = await store.updateSetting('key', 'val')

    expect(result).toBe(false)
    expect(store.error).toBe('权限不足')
  })
})
