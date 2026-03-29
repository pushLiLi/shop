import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import EditableImage from './EditableImage.vue'

describe('EditableImage', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  const defaultProps = {
    src: '/media/test.jpg',
    alt: 'Test Image'
  }

  it('renders image with correct src and alt', () => {
    const wrapper = mount(EditableImage, { props: defaultProps })
    const img = wrapper.find('img')
    expect(img.attributes('src')).toBe('/media/test.jpg')
    expect(img.attributes('alt')).toBe('Test Image')
  })

  it('does not show overlay for non-admin', () => {
    const wrapper = mount(EditableImage, { props: defaultProps })
    wrapper.find('.editable-image-wrapper').trigger('mouseenter')
    expect(wrapper.find('.image-overlay').exists()).toBe(false)
  })

  it('shows overlay on hover when admin', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    const wrapper = mount(EditableImage, { props: defaultProps })
    await wrapper.find('.editable-image-wrapper').trigger('mouseenter')
    expect(wrapper.find('.image-overlay').exists()).toBe(true)
    expect(wrapper.find('.upload-btn').text()).toBe('更换图片')
  })

  it('hides overlay on mouseleave', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    const wrapper = mount(EditableImage, { props: defaultProps })
    await wrapper.find('.editable-image-wrapper').trigger('mouseenter')
    expect(wrapper.find('.image-overlay').exists()).toBe(true)

    await wrapper.find('.editable-image-wrapper').trigger('mouseleave')
    expect(wrapper.find('.image-overlay').exists()).toBe(false)
  })

  it('emits update on successful file upload', async () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'admin', name: 'Admin' }))

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ url: '/media/new.jpg' })
    })

    const wrapper = mount(EditableImage, { props: defaultProps })
    await wrapper.find('.editable-image-wrapper').trigger('mouseenter')

    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    const fileInput = wrapper.find('input[type="file"]')

    Object.defineProperty(fileInput.element, 'files', {
      value: [file],
      configurable: true
    })
    await fileInput.trigger('change')

    await vi.waitFor(() => {
      expect(wrapper.emitted('update')).toBeTruthy()
      expect(wrapper.emitted('update')[0]).toEqual(['/media/new.jpg'])
    })
  })
})
