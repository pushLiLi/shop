import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import AdminImageUpload from './AdminImageUpload.vue'

vi.mock('vue-advanced-cropper', () => ({
  Cropper: {
    template: '<div class="mock-cropper"><slot /></div>'
  }
}))

vi.mock('vue-advanced-cropper/dist/style.css', () => ({}))

describe('AdminImageUpload', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows upload area when no image', () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    expect(wrapper.find('.upload-area').exists()).toBe(true)
    expect(wrapper.text()).toContain('拖拽图片到此处或点击上传')
  })

  it('shows preview when image URL exists', () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '/media/test.jpg' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    expect(wrapper.find('.preview-area').exists()).toBe(true)
    expect(wrapper.find('.preview-image').attributes('src')).toBe('/media/test.jpg')
  })

  it('shows replace and remove buttons in preview mode', () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '/media/test.jpg' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    expect(wrapper.find('.btn-replace').exists()).toBe(true)
    expect(wrapper.find('.btn-remove').exists()).toBe(true)
  })

  it('clears image on remove click', async () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '/media/test.jpg' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    await wrapper.find('.btn-remove').trigger('click')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([''])
  })

  it('shows error for non-image file', async () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })

    const file = new File(['text'], 'test.txt', { type: 'text/plain' })
    const fileInput = wrapper.findAll('input[type="file"]')[0]

    Object.defineProperty(fileInput.element, 'files', {
      value: [file],
      configurable: true
    })
    await fileInput.trigger('change')

    expect(wrapper.find('.error-text').exists()).toBe(true)
    expect(wrapper.text()).toContain('请选择图片文件')
  })

  it('shows error for file over 10MB', async () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })

    const bigFile = new File(['x'.repeat(11 * 1024 * 1024)], 'big.jpg', { type: 'image/jpeg' })
    const fileInput = wrapper.findAll('input[type="file"]')[0]

    Object.defineProperty(fileInput.element, 'files', {
      value: [bigFile],
      configurable: true
    })
    await fileInput.trigger('change')

    expect(wrapper.find('.error-text').exists()).toBe(true)
    expect(wrapper.text()).toContain('10MB')
  })

  it('has file input element', () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    const fileInputs = wrapper.findAll('input[type="file"]')
    expect(fileInputs.length).toBeGreaterThanOrEqual(1)
  })

  it('has drag and drop support', () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    const uploadArea = wrapper.find('.upload-area')
    expect(uploadArea.exists()).toBe(true)
  })

  it('enters drag-over state on dragover', async () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    await wrapper.find('.upload-area').trigger('dragover', { preventDefault: () => {} })
    expect(wrapper.find('.upload-area').classes()).toContain('drag-over')
  })

  it('leaves drag-over state on dragleave', async () => {
    const wrapper = mount(AdminImageUpload, {
      props: { modelValue: '' },
      global: {
        stubs: {
          Teleport: { template: '<div><slot /></div>' }
        }
      }
    })
    await wrapper.find('.upload-area').trigger('dragover', { preventDefault: () => {} })
    expect(wrapper.find('.upload-area').classes()).toContain('drag-over')

    await wrapper.find('.upload-area').trigger('dragleave')
    expect(wrapper.find('.upload-area').classes()).not.toContain('drag-over')
  })
})
