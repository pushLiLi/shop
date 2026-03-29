<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToastStore()

onMounted(() => {
  if (authStore.isLoggedIn) {
    router.push('/profile')
  }
})

const activeTab = ref('login')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const name = ref('')
const error = ref('')
const loading = ref(false)

const registerCaptchaId = ref('')
const registerCaptchaImage = ref('')
const registerCaptchaCode = ref('')

const loginCaptchaRequired = ref(false)
const loginCaptchaId = ref('')
const loginCaptchaImage = ref('')
const loginCaptchaCode = ref('')

async function refreshRegisterCaptcha() {
  try {
    const res = await fetch('/api/auth/captcha')
    const data = await res.json()
    registerCaptchaId.value = data.captchaId
    registerCaptchaImage.value = data.captchaImage
    registerCaptchaCode.value = ''
  } catch (e) {
    console.error(e)
  }
}

async function refreshLoginCaptcha() {
  try {
    const res = await fetch('/api/auth/captcha')
    const data = await res.json()
    loginCaptchaId.value = data.captchaId
    loginCaptchaImage.value = data.captchaImage
    loginCaptchaCode.value = ''
  } catch (e) {
    console.error(e)
  }
}

const handleLogin = async () => {
  error.value = ''

  if (!email.value || !password.value) {
    error.value = '请输入邮箱和密码'
    return
  }

  if (loginCaptchaRequired.value && !loginCaptchaCode.value) {
    error.value = '请输入验证码'
    return
  }

  loading.value = true

  try {
    const user = await authStore.login(
      email.value,
      password.value,
      loginCaptchaId.value,
      loginCaptchaCode.value
    )

    if (route.query.redirect) {
      router.push(route.query.redirect)
    } else if (user.role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/')
    }
  } catch (e) {
    error.value = e.message
    if (e.requireCaptcha) {
      loginCaptchaRequired.value = true
      refreshLoginCaptcha()
    }
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  error.value = ''

  if (!email.value || !password.value) {
    error.value = '请输入邮箱和密码'
    return
  }

  if (password.value.length < 6) {
    error.value = '密码至少需要6个字符'
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }

  if (!registerCaptchaCode.value) {
    error.value = '请输入验证码'
    return
  }

  loading.value = true

  try {
    await authStore.register(
      email.value,
      password.value,
      name.value,
      registerCaptchaId.value,
      registerCaptchaCode.value
    )
    toast.success('注册成功，已自动登录')
    router.push('/')
  } catch (e) {
    error.value = e.message
    refreshRegisterCaptcha()
  } finally {
    loading.value = false
  }
}

const switchTab = (tab) => {
  activeTab.value = tab
  error.value = ''
  if (tab === 'register' && !registerCaptchaImage.value) {
    refreshRegisterCaptcha()
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <h1>{{ activeTab === 'login' ? '用户登录' : '用户注册' }}</h1>
        <p class="subtitle">{{ activeTab === 'login' ? '登录您的账户' : '创建新账户' }}</p>
      </div>

      <div class="tabs">
        <button
          :class="['tab-btn', { active: activeTab === 'login' }]"
          @click="switchTab('login')"
        >
          登录
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'register' }]"
          @click="switchTab('register')"
        >
          注册
        </button>
      </div>

      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>邮箱</label>
          <input
            v-model="email"
            type="email"
            placeholder="请输入邮箱"
            autocomplete="email"
          >
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="password"
            type="password"
            placeholder="请输入密码"
            autocomplete="current-password"
          >
        </div>

        <div v-if="loginCaptchaRequired" class="form-group">
          <label>验证码</label>
          <div class="captcha-row">
            <input
              v-model="loginCaptchaCode"
              type="text"
              placeholder="请输入验证码"
              maxlength="4"
              autocomplete="off"
            >
            <img
              v-if="loginCaptchaImage"
              :src="loginCaptchaImage"
              class="captcha-img"
              @click="refreshLoginCaptcha"
              title="点击刷新"
            >
            <button type="button" class="btn-refresh-captcha" @click="refreshLoginCaptcha">
              刷新
            </button>
          </div>
        </div>

        <div v-if="error" class="error-msg">{{ error }}</div>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <form v-else @submit.prevent="handleRegister" class="login-form">
        <div class="form-group">
          <label>邮箱</label>
          <input
            v-model="email"
            type="email"
            placeholder="请输入邮箱"
            autocomplete="email"
          >
        </div>

        <div class="form-group">
          <label>用户名 (可选)</label>
          <input
            v-model="name"
            type="text"
            placeholder="请输入用户名"
            autocomplete="username"
          >
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="password"
            type="password"
            placeholder="请输入密码 (至少6位)"
            autocomplete="new-password"
          >
        </div>

        <div class="form-group">
          <label>确认密码</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            autocomplete="new-password"
          >
        </div>

        <div class="form-group">
          <label>验证码</label>
          <div class="captcha-row">
            <input
              v-model="registerCaptchaCode"
              type="text"
              placeholder="请输入验证码"
              maxlength="4"
              autocomplete="off"
            >
            <img
              v-if="registerCaptchaImage"
              :src="registerCaptchaImage"
              class="captcha-img"
              @click="refreshRegisterCaptcha"
              title="点击刷新"
            >
            <button type="button" class="btn-refresh-captcha" @click="refreshRegisterCaptcha">
              刷新
            </button>
          </div>
        </div>

        <div v-if="error" class="error-msg">{{ error }}</div>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 70vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
}

.login-container {
  width: 100%;
  max-width: 420px;
  background: #2d2d2d;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 28px;
  color: #d4a574;
  margin: 0 0 8px;
}

.subtitle {
  color: #888;
  margin: 0;
}

.tabs {
  display: flex;
  margin-bottom: 30px;
  border-bottom: 1px solid #444;
}

.tab-btn {
  flex: 1;
  padding: 12px;
  background: transparent;
  border: none;
  color: #888;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
}

.tab-btn:hover {
  color: #fff;
}

.tab-btn.active {
  color: #d4a574;
  border-bottom-color: #d4a574;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: #ccc;
  font-size: 14px;
}

.form-group input {
  padding: 14px 16px;
  background: #1a1a1a;
  border: 1px solid #444;
  border-radius: 8px;
  color: #fff;
  font-size: 15px;
  transition: border-color 0.3s;
}

.form-group input::placeholder {
  color: #666;
}

.form-group input:focus {
  outline: none;
  border-color: #d4a574;
}

.captcha-row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.captcha-row input {
  flex: 1;
  min-width: 0;
}

.captcha-img {
  height: 46px;
  border-radius: 6px;
  cursor: pointer;
  flex-shrink: 0;
}

.captcha-img:hover {
  opacity: 0.8;
}

.btn-refresh-captcha {
  padding: 12px 14px;
  background: #3a3a3a;
  border: 1px solid #555;
  border-radius: 8px;
  color: #ccc;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  flex-shrink: 0;
}

.btn-refresh-captcha:hover {
  background: #444;
  color: #fff;
}

.error-msg {
  color: #e74c3c;
  font-size: 14px;
  text-align: center;
  padding: 10px;
  background: rgba(231, 76, 60, 0.1);
  border-radius: 6px;
}

.submit-btn {
  padding: 14px;
  background: #d4a574;
  border: none;
  border-radius: 8px;
  color: #1a1a1a;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 10px;
}

.submit-btn:hover:not(:disabled) {
  background: #c49564;
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
