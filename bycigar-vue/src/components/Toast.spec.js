import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Toast from './Toast.vue'

describe('Toast', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  function mountToast() {
    return mount(Toast, {
      global: {
        stubs: {
          Transition: {
            template: '<slot />'
          }
        }
      }
    })
  }

  it('is hidden when toast store visible is false', () => {
    const wrapper = mountToast()
    expect(wrapper.find('.toast').exists()).toBe(false)
  })

  it('shows success toast with correct icon and message', async () => {
    const wrapper = mountToast()
    const { useToastStore } = await import('../stores/toast')
    const toast = useToastStore()
    toast.success('操作成功')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const toastEl = wrapper.find('.toast')
    expect(toastEl.exists()).toBe(true)
    expect(toastEl.classes()).toContain('success')
    expect(toastEl.find('.toast-icon').text()).toBe('✓')
    expect(toastEl.find('.toast-message').text()).toBe('操作成功')
  })

  it('shows error toast with correct icon and message', async () => {
    const wrapper = mountToast()
    const { useToastStore } = await import('../stores/toast')
    const toast = useToastStore()
    toast.error('出错了')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const toastEl = wrapper.find('.toast')
    expect(toastEl.exists()).toBe(true)
    expect(toastEl.classes()).toContain('error')
    expect(toastEl.find('.toast-icon').text()).toBe('✕')
    expect(toastEl.find('.toast-message').text()).toBe('出错了')
  })

  it('hides when toast store visible becomes false', async () => {
    const wrapper = mountToast()
    const { useToastStore } = await import('../stores/toast')
    const toast = useToastStore()
    toast.success('显示')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.toast').exists()).toBe(true)

    toast.hide()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.toast').exists()).toBe(false)
  })
})
