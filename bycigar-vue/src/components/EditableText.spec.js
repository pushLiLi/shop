import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import EditableText from './EditableText.vue'

describe('EditableText', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  const defaultProps = {
    configKey: 'site_name',
    modelValue: 'HUAUGE',
    tag: 'span'
  }

  it('displays the modelValue text', () => {
    const wrapper = mount(EditableText, { props: defaultProps })
    expect(wrapper.text()).toContain('HUAUGE')
  })

  it('does not show input when user is not admin', () => {
    const wrapper = mount(EditableText, { props: defaultProps })
    expect(wrapper.find('input').exists()).toBe(false)
    expect(wrapper.find('span').exists()).toBe(true)
  })

  it('shows input when admin clicks on text', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    const wrapper = mount(EditableText, { props: defaultProps })
    await wrapper.find('span').trigger('click')
    expect(wrapper.find('input').exists()).toBe(true)
    expect(wrapper.find('input').element.value).toBe('HUAUGE')
  })

  it('emits update on blur when value changed', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    const mockFetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({}) })
    global.fetch = mockFetch

    const wrapper = mount(EditableText, { props: defaultProps })
    await wrapper.find('span').trigger('click')

    const input = wrapper.find('input')
    await input.setValue('NEW NAME')
    await input.trigger('blur')

    expect(mockFetch).toHaveBeenCalledWith(
      '/api/admin/config/site_name',
      expect.objectContaining({ method: 'PUT' })
    )
    expect(wrapper.emitted('update')).toBeTruthy()
    expect(wrapper.emitted('update')[0]).toEqual(['NEW NAME'])
  })

  it('does not call API if value unchanged on blur', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    const mockFetch = vi.fn()
    global.fetch = mockFetch

    const wrapper = mount(EditableText, { props: defaultProps })
    await wrapper.find('span').trigger('click')

    await wrapper.find('input').trigger('blur')
    expect(mockFetch).not.toHaveBeenCalled()
    expect(wrapper.emitted('update')).toBeFalsy()
  })

  it('returns to span display after blur', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    global.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({}) })

    const wrapper = mount(EditableText, { props: defaultProps })
    await wrapper.find('span').trigger('click')
    expect(wrapper.find('input').exists()).toBe(true)

    await wrapper.find('input').trigger('blur')
    expect(wrapper.find('input').exists()).toBe(false)
    expect(wrapper.find('span').exists()).toBe(true)
  })
})
