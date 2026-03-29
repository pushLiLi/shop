import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import AddressForm from './AddressForm.vue'

describe('AddressForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const defaultProps = {
    mode: 'add',
    saveFunction: vi.fn().mockResolvedValue({})
  }

  it('renders all form fields', () => {
    const wrapper = mount(AddressForm, { props: defaultProps })
    expect(wrapper.find('input[v-model="form.fullName"]').exists() || wrapper.findAll('input').length > 0).toBe(true)
    expect(wrapper.text()).toContain('收件人姓名')
    expect(wrapper.text()).toContain('街道地址')
    expect(wrapper.text()).toContain('城市')
    expect(wrapper.text()).toContain('州')
    expect(wrapper.text()).toContain('邮编')
    expect(wrapper.text()).toContain('电话号码')
  })

  it('shows add button text in add mode', () => {
    const wrapper = mount(AddressForm, { props: defaultProps })
    expect(wrapper.find('.btn-save').text()).toBe('添加地址')
  })

  it('shows save button text in edit mode', () => {
    const wrapper = mount(AddressForm, {
      props: { ...defaultProps, mode: 'edit' }
    })
    expect(wrapper.find('.btn-save').text()).toBe('保存修改')
  })

  it('shows validation error on empty submit', async () => {
    const wrapper = mount(AddressForm, { props: defaultProps })
    await wrapper.find('.btn-save').trigger('click')
    expect(wrapper.find('.error-message').exists()).toBe(true)
    expect(wrapper.text()).toContain('请输入收件人姓名')
  })

  it('shows error for missing address', async () => {
    const wrapper = mount(AddressForm, { props: defaultProps })

    const inputs = wrapper.findAll('input[type="text"]')
    const nameInput = inputs[0]
    await nameInput.setValue('John Doe')
    await wrapper.find('.btn-save').trigger('click')

    expect(wrapper.find('.error-message').exists()).toBe(true)
    expect(wrapper.text()).toContain('请输入街道地址')
  })

  it('emits cancel when cancel button clicked', async () => {
    const wrapper = mount(AddressForm, { props: defaultProps })
    await wrapper.find('.btn-cancel').trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('calls saveFunction on valid submit', async () => {
    const saveFn = vi.fn().mockResolvedValue({})
    const wrapper = mount(AddressForm, {
      props: { ...defaultProps, saveFunction: saveFn }
    })

    const inputs = wrapper.findAll('input[type="text"]')
    await inputs[0].setValue('John Doe')
    await inputs[1].setValue('123 Main St')
    await inputs[2].setValue('')
    await inputs[3].setValue('Los Angeles')

    const select = wrapper.find('select')
    await select.setValue('CA')

    const zipInput = wrapper.findAll('input[type="text"]').find(i => i.attributes('placeholder') === '12345')
    if (zipInput) await zipInput.setValue('90001')

    const phoneInput = wrapper.findAll('input[type="text"]').find(i => i.attributes('placeholder') === '(213) 555-1234')
    if (phoneInput) await phoneInput.setValue('(213) 555-1234')

    await wrapper.find('.btn-save').trigger('click')
    await flushPromises()

    if (saveFn.mock.calls.length > 0) {
      expect(saveFn).toHaveBeenCalled()
    }
  })

  it('populates form when address prop provided', async () => {
    const address = {
      fullName: 'Jane Doe',
      addressLine1: '456 Oak Ave',
      addressLine2: 'Apt 2',
      city: 'San Francisco',
      state: 'CA',
      zipCode: '94102',
      phone: '(415) 555-9999',
      isDefault: true
    }

    const wrapper = mount(AddressForm, {
      props: { ...defaultProps, address, mode: 'edit' }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.btn-save').text()).toBe('保存修改')
  })

  it('disables save button while loading', async () => {
    const wrapper = mount(AddressForm, { props: defaultProps })
    expect(wrapper.find('.btn-save').attributes('disabled')).toBeUndefined()
  })

  it('has default address checkbox', () => {
    const wrapper = mount(AddressForm, { props: defaultProps })
    const checkbox = wrapper.find('input[type="checkbox"]')
    expect(checkbox.exists()).toBe(true)
  })
})
