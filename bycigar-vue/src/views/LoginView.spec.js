import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { flushPromises } from '@vue/test-utils'
import LoginView from './LoginView.vue'

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/login', component: LoginView },
      { path: '/profile', component: { template: '<div>Profile</div>' } },
      { path: '/admin', component: { template: '<div>Admin</div>' } },
      { path: '/', component: { template: '<div>Home</div>' } }
    ]
  })
}

async function mountLogin(options = {}) {
  const router = createTestRouter()
  const pinia = createPinia()
  setActivePinia(pinia)

  if (options.loggedIn) {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({ role: 'customer', name: 'User', email: 'user@test.com' }))
  }

  global.fetch = vi.fn(() => Promise.resolve({
    ok: true,
    json: () => Promise.resolve({ captchaId: 'cap-123', captchaImage: 'data:image/png;base64,abc' })
  }))

  const wrapper = mount(LoginView, {
    global: {
      plugins: [router, pinia],
      stubs: {
        'router-link': {
          template: '<a :href="$attrs.to"><slot /></a>',
          props: []
        }
      }
    }
  })

  await router.isReady()
  await flushPromises()
  return { wrapper, router, pinia }
}

describe('LoginView', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows login form by default', async () => {
    const { wrapper } = await mountLogin()
    expect(wrapper.find('.login-form').exists()).toBe(true)
    expect(wrapper.text()).toContain('用户登录')
    expect(wrapper.find('.submit-btn').text()).toBe('登录')
  })

  it('shows login tab as active by default', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    expect(tabs[0].classes()).toContain('active')
    expect(tabs[1].classes()).not.toContain('active')
  })

  it('switches to register form on tab click', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('用户注册')
    expect(wrapper.find('input[autocomplete="new-password"]').exists()).toBe(true)
  })

  it('shows email and password fields in login form', async () => {
    const { wrapper } = await mountLogin()
    const inputs = wrapper.findAll('.form-group input')
    expect(inputs.length).toBeGreaterThanOrEqual(2)
  })

  it('shows register captcha field when on register tab', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('验证码')
    expect(wrapper.find('.captcha-row').exists()).toBe(true)
  })

  it('does not show login captcha by default', async () => {
    const { wrapper } = await mountLogin()
    expect(wrapper.find('.captcha-row').exists()).toBe(false)
  })

  it('shows error on empty login submission', async () => {
    const { wrapper } = await mountLogin()
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
    expect(wrapper.find('.error-msg').exists()).toBe(true)
    expect(wrapper.find('.error-msg').text()).toContain('请输入邮箱和密码')
  })

  it('shows error on empty register submission', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
    expect(wrapper.find('.error-msg').text()).toContain('请输入邮箱和密码')
  })

  it('shows error when register passwords do not match', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()

    const inputs = wrapper.findAll('form input')
    await inputs[0].setValue('test@test.com')
    await inputs[2].setValue('123456')
    await inputs[3].setValue('654321')
    await inputs[4].setValue('1234')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.find('.error-msg').text()).toContain('两次输入的密码不一致')
  })

  it('shows error when register password is too short', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()

    const inputs = wrapper.findAll('form input')
    await inputs[0].setValue('test@test.com')
    await inputs[2].setValue('12345')
    await inputs[3].setValue('12345')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.find('.error-msg').text()).toContain('密码至少需要6个字符')
  })

  it('disables submit button during loading', async () => {
    const { wrapper } = await mountLogin()
    expect(wrapper.find('.submit-btn').attributes('disabled')).toBeUndefined()
  })

  it('has tab buttons for login and register', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    expect(tabs).toHaveLength(2)
    expect(tabs[0].text()).toBe('登录')
    expect(tabs[1].text()).toBe('注册')
  })

  it('clears error when switching tabs', async () => {
    const { wrapper } = await mountLogin()
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
    expect(wrapper.find('.error-msg').exists()).toBe(true)

    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()
    expect(wrapper.find('.error-msg').exists()).toBe(false)
  })

  it('fetches captcha when switching to register tab', async () => {
    const { wrapper } = await mountLogin()
    const tabs = wrapper.findAll('.tab-btn')
    await tabs[1].trigger('click')
    await flushPromises()
    expect(global.fetch).toHaveBeenCalledWith('/api/auth/captcha')
  })
})
