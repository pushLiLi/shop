import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useToastStore } from '../stores/toast'

describe('toast store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('success sets message, type and visible', () => {
    const store = useToastStore()
    store.success('操作成功')

    expect(store.message).toBe('操作成功')
    expect(store.type).toBe('success')
    expect(store.visible).toBe(true)
  })

  it('error sets message, type and visible', () => {
    const store = useToastStore()
    store.error('出错了')

    expect(store.message).toBe('出错了')
    expect(store.type).toBe('error')
    expect(store.visible).toBe(true)
  })

  it('auto hides after 2 seconds', () => {
    const store = useToastStore()
    store.success('测试')

    expect(store.visible).toBe(true)

    vi.advanceTimersByTime(2000)

    expect(store.visible).toBe(false)
  })

  it('hide immediately sets visible to false', () => {
    const store = useToastStore()
    store.success('测试')
    expect(store.visible).toBe(true)

    store.hide()

    expect(store.visible).toBe(false)
  })

  it('show with custom type', () => {
    const store = useToastStore()
    store.show('自定义', 'warning')

    expect(store.type).toBe('warning')
    expect(store.message).toBe('自定义')
  })

  it('new toast resets timer', () => {
    const store = useToastStore()
    store.success('第一条')
    vi.advanceTimersByTime(1500)

    store.success('第二条')
    vi.advanceTimersByTime(1500)

    expect(store.visible).toBe(true)
    expect(store.message).toBe('第二条')
  })
})
