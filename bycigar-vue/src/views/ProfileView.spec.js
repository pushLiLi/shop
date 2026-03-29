import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { flushPromises } from '@vue/test-utils'
import ProfileView from './ProfileView.vue'

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/profile', component: ProfileView },
      { path: '/login', component: { template: '<div>Login</div>' } },
      { path: '/', component: { template: '<div>Home</div>' } }
    ]
  })
}

const mockUser = {
  role: 'customer',
  name: 'Test User',
  email: 'test@test.com'
}

const mockOrders = {
  orders: [
    {
      id: 1,
      orderNo: 'ORD-001',
      status: 'paid',
      total: 99.99,
      createdAt: '2026-01-15T10:00:00Z',
      items: [
        { id: 1, quantity: 2, price: 49.99, product: { name: 'Cigar', imageUrl: '/cigar.jpg' } }
      ]
    }
  ]
}

const mockAddresses = {
  addresses: [
    {
      id: 1,
      fullName: 'John Doe',
      addressLine1: '123 Main St',
      addressLine2: 'Apt 4',
      city: 'New York',
      state: 'NY',
      zipCode: '10001',
      phone: '212-555-1234',
      isDefault: true
    }
  ]
}

async function mountProfile(options = {}) {
  const router = createTestRouter()
  const pinia = createPinia()
  setActivePinia(pinia)

  localStorage.setItem('token', 'test-token')
  localStorage.setItem('user', JSON.stringify(options.user || mockUser))

  let callCount = 0
  global.fetch = vi.fn((url) => {
    callCount++
    if (url.includes('/api/auth/me')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ user: options.user || mockUser })
      })
    }
    if (url.includes('/api/orders')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(options.ordersResponse || mockOrders)
      })
    }
    if (url.includes('/api/addresses')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(options.addressesResponse || mockAddresses)
      })
    }
    if (url.includes('/api/auth/captcha')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ captchaId: 'cap-123', captchaImage: 'data:image/png;base64,abc' })
      })
    }
    return Promise.resolve({
      ok: true,
      json: () => Promise.resolve({})
    })
  })

  const wrapper = mount(ProfileView, {
    global: {
      plugins: [router, pinia],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        },
        AddressForm: {
          template: '<div class="address-form-stub" />',
          props: ['address', 'mode', 'saveFunction']
        }
      }
    }
  })

  await router.push('/profile')
  await router.isReady()
  await flushPromises()
  return { wrapper, router, pinia }
}

describe('ProfileView', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders 4 tab buttons', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    expect(navItems.length).toBe(5)
    expect(navItems[0].text()).toBe('个人信息')
    expect(navItems[1].text()).toBe('地址管理')
    expect(navItems[2].text()).toBe('我的订单')
    expect(navItems[3].text()).toBe('账户安全')
    expect(navItems[4].text()).toBe('退出登录')
  })

  it('shows user info tab by default', async () => {
    const { wrapper } = await mountProfile()
    expect(wrapper.find('.content-section').text()).toContain('个人信息')
    expect(wrapper.text()).toContain('Test User')
    expect(wrapper.text()).toContain('test@test.com')
  })

  it('shows user role', async () => {
    const { wrapper } = await mountProfile()
    expect(wrapper.text()).toContain('普通用户')
  })

  it('shows admin role for admin user', async () => {
    const { wrapper } = await mountProfile({
      user: { role: 'admin', name: 'Admin', email: 'admin@test.com' }
    })
    expect(wrapper.text()).toContain('管理员')
  })

  it('shows edit button in info tab', async () => {
    const { wrapper } = await mountProfile()
    expect(wrapper.find('.btn-edit').text()).toBe('编辑信息')
  })

  it('shows edit form when edit button clicked', async () => {
    const { wrapper } = await mountProfile()
    await wrapper.find('.btn-edit').trigger('click')
    await flushPromises()
    expect(wrapper.find('.edit-form').exists()).toBe(true)
    expect(wrapper.find('.btn-cancel').exists()).toBe(true)
    expect(wrapper.find('.btn-save').exists()).toBe(true)
  })

  it('returns to info view on cancel', async () => {
    const { wrapper } = await mountProfile()
    await wrapper.find('.btn-edit').trigger('click')
    await flushPromises()
    await wrapper.find('.btn-cancel').trigger('click')
    await flushPromises()
    expect(wrapper.find('.edit-form').exists()).toBe(false)
    expect(wrapper.find('.info-list').exists()).toBe(true)
  })

  it('switches to orders tab', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[2].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('我的订单')
    expect(wrapper.text()).toContain('ORD-001')
  })

  it('shows order items in orders tab', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[2].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('Cigar')
    expect(wrapper.text()).toContain('$49.99')
  })

  it('shows empty orders message', async () => {
    const { wrapper } = await mountProfile({
      ordersResponse: { orders: [] }
    })
    const navItems = wrapper.findAll('.nav-item')
    await navItems[2].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('暂无订单记录')
  })

  it('switches to addresses tab', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[1].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('地址管理')
    expect(wrapper.text()).toContain('John Doe')
  })

  it('shows address details', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[1].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('123 Main St')
    expect(wrapper.text()).toContain('New York')
  })

  it('shows default badge on default address', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[1].trigger('click')
    await flushPromises()
    expect(wrapper.find('.default-badge').exists()).toBe(true)
    expect(wrapper.find('.default-badge').text()).toBe('默认')
  })

  it('shows add address button when under limit', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[1].trigger('click')
    await flushPromises()
    expect(wrapper.find('.btn-add-address').exists()).toBe(true)
    expect(wrapper.find('.btn-add-address').text()).toContain('新增地址')
  })

  it('shows empty address state', async () => {
    const { wrapper } = await mountProfile({
      addressesResponse: { addresses: [] }
    })
    const navItems = wrapper.findAll('.nav-item')
    await navItems[1].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('您还没有保存任何收货地址')
  })

  it('switches to security tab', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[3].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('账户安全')
    expect(wrapper.text()).toContain('原密码')
    expect(wrapper.text()).toContain('新密码')
  })

  it('shows captcha in security tab', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[3].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('验证码')
    expect(global.fetch).toHaveBeenCalledWith('/api/auth/captcha')
  })

  it('validates password mismatch in security tab', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[3].trigger('click')
    await flushPromises()

    const inputs = wrapper.findAll('.password-form input')
    await inputs[1].setValue('newpass1')
    await inputs[2].setValue('newpass2')
    await wrapper.find('.password-form .btn-save').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('两次输入的新密码不一致')
  })

  it('shows logout button', async () => {
    const { wrapper } = await mountProfile()
    const logoutBtn = wrapper.findAll('.nav-item')[4]
    expect(logoutBtn.text()).toBe('退出登录')
    expect(logoutBtn.classes()).toContain('logout')
  })

  it('shows user name in sidebar', async () => {
    const { wrapper } = await mountProfile()
    expect(wrapper.find('.user-name').text()).toBe('Test User')
    expect(wrapper.find('.user-email').text()).toBe('test@test.com')
  })

  it('uses email prefix as fallback name', async () => {
    const { wrapper } = await mountProfile({
      user: { role: 'customer', name: '', email: 'john@example.com' }
    })
    expect(wrapper.find('.user-name').text()).toBe('john')
  })

  it('shows address count', async () => {
    const { wrapper } = await mountProfile()
    const navItems = wrapper.findAll('.nav-item')
    await navItems[1].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('1 / 5 个地址')
  })
})
