<script setup>
import { ref, onMounted } from 'vue'
import { useToastStore } from '../../stores/toast'
import AdminImageUpload from '../../components/AdminImageUpload.vue'

const API_BASE = '/api'
const toast = useToastStore()

const loading = ref(false)
const saving = ref(false)

const form = ref({
  footer_description: '',
  footer_service_time: ''
})

const banners = ref({
  home_banner_1: ''
})

const siteIdentity = ref({
  site_title: '',
  site_meta_description: '',
  favicon_url: ''
})

const emailConfig = ref({
  email_enabled: 'false',
  email_smtp_host: 'smtp.qq.com',
  email_smtp_port: '465',
  email_smtp_username: '',
  email_smtp_password: '',
  email_from_name: ''
})

const testEmail = ref('')
const testingEmail = ref(false)

const bannerLabels = {
  home_banner_1: '横幅图 1（特别推荐下方）'
}

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const fetchSettings = async () => {
  loading.value = true
  try {
    const [settingsRes, configRes] = await Promise.all([
      fetch(`${API_BASE}/settings`),
      fetch(`${API_BASE}/config`)
    ])

    const settingsData = await settingsRes.json()
    if (settingsData.success) {
      form.value.footer_description = settingsData.data.footer_description || ''
      form.value.footer_service_time = settingsData.data.footer_service_time || ''
      emailConfig.value.email_enabled = settingsData.data.email_enabled || 'false'
      emailConfig.value.email_smtp_host = settingsData.data.email_smtp_host || 'smtp.qq.com'
      emailConfig.value.email_smtp_port = settingsData.data.email_smtp_port || '465'
      emailConfig.value.email_smtp_username = settingsData.data.email_smtp_username || ''
      emailConfig.value.email_smtp_password = settingsData.data.email_smtp_password || ''
      emailConfig.value.email_from_name = settingsData.data.email_from_name || ''
    }

    const configData = await configRes.json()
    banners.value.home_banner_1 = configData.home_banner_1 || ''

    const identityRes = await fetch(`${API_BASE}/site-identity`)
    if (identityRes.ok) {
      const identityData = await identityRes.json()
      siteIdentity.value.site_title = identityData.title || ''
      siteIdentity.value.site_meta_description = identityData.metaDescription || ''
      siteIdentity.value.favicon_url = identityData.faviconUrl || ''
    }
  } catch (e) {
    toast.error('获取设置失败')
  } finally {
    loading.value = false
  }
}

const sendTestEmail = async () => {
  if (!testEmail.value) {
    toast.error('请输入测试收件邮箱')
    return
  }
  testingEmail.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/email/test`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ to: testEmail.value })
    })
    const data = await res.json()
    if (res.ok) {
      toast.success('测试邮件发送成功，请检查收件箱')
    } else {
      toast.error(data.error || '测试邮件发送失败')
    }
  } catch (e) {
    toast.error('测试邮件发送失败')
  } finally {
    testingEmail.value = false
  }
}

const saveSettings = async () => {
  saving.value = true
  try {
    const updates = [
      fetch(`${API_BASE}/admin/settings/footer_description`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: form.value.footer_description })
      }),
      fetch(`${API_BASE}/admin/settings/footer_service_time`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: form.value.footer_service_time })
      }),
      fetch(`${API_BASE}/admin/config/home_banner_1`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: banners.value.home_banner_1 })
      }),
      fetch(`${API_BASE}/admin/config/site_title`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: siteIdentity.value.site_title })
      }),
      fetch(`${API_BASE}/admin/config/site_meta_description`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: siteIdentity.value.site_meta_description })
      }),
      fetch(`${API_BASE}/admin/config/favicon_url`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: siteIdentity.value.favicon_url })
      }),
      fetch(`${API_BASE}/admin/settings/email_enabled`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: emailConfig.value.email_enabled })
      }),
      fetch(`${API_BASE}/admin/settings/email_smtp_host`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: emailConfig.value.email_smtp_host })
      }),
      fetch(`${API_BASE}/admin/settings/email_smtp_port`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: emailConfig.value.email_smtp_port })
      }),
      fetch(`${API_BASE}/admin/settings/email_smtp_username`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: emailConfig.value.email_smtp_username })
      })
    ]

    if (emailConfig.value.email_smtp_password && emailConfig.value.email_smtp_password !== '****') {
      updates.push(fetch(`${API_BASE}/admin/settings/email_smtp_password`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: emailConfig.value.email_smtp_password })
      }))
    }

    updates.push(
      fetch(`${API_BASE}/admin/settings/email_from_name`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: emailConfig.value.email_from_name })
      })
    )

    await Promise.all(updates)
    toast.success('保存成功')
  } catch (e) {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchSettings)
</script>

<template>
  <div class="admin-settings">
    <div class="page-header">
      <h2>站点设置</h2>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <form v-else class="settings-form" @submit.prevent="saveSettings">
      <div class="form-section">
        <h3>网站基本信息</h3>

        <div class="content-tips">
          <div class="tip-item">
            <strong>网站标题：</strong>显示在浏览器标签页和搜索引擎结果中的标题
          </div>
          <div class="tip-item">
            <strong>META 描述：</strong>显示在搜索引擎结果中的网站简介，建议 50-160 字
          </div>
          <div class="tip-item">
            <strong>Favicon：</strong>显示在浏览器标签页的小图标，推荐 512×512 PNG 格式
          </div>
        </div>

        <div class="form-group">
          <label>网站标题</label>
          <input
            v-model="siteIdentity.site_title"
            type="text"
            placeholder="例如：BYCIGAR | 权威正品雪茄在线购买商城"
          />
          <span class="hint">显示在浏览器标签页，建议包含品牌名和简短描述</span>
        </div>

        <div class="form-group">
          <label>META 描述</label>
          <textarea
            v-model="siteIdentity.site_meta_description"
            rows="3"
            placeholder="输入网站 meta 描述..."
          ></textarea>
          <span class="hint">显示在搜索引擎结果中，建议包含关键词和品牌简介</span>
        </div>

        <div class="form-group">
          <label>Favicon 图标</label>
          <AdminImageUpload v-model="siteIdentity.favicon_url" :aspect-ratio="1" />
          <span class="hint">建议使用正方形图片，最佳尺寸 512×512 像素</span>
        </div>
      </div>

      <div class="form-section">
        <h3>首页横幅图</h3>

        <div class="content-tips">
          <div class="tip-item">
            <strong>推荐尺寸：</strong>1400 × 500px（宽高比 7:3）
          </div>
          <div class="tip-item">
            <strong>推荐内容：</strong>品牌故事、新品推荐、促销活动、节日专题等
          </div>
        </div>

        <div v-for="(label, key) in bannerLabels" :key="key" class="form-group">
          <label>{{ label }}</label>
          <AdminImageUpload v-model="banners[key]" :aspect-ratio="7/3" />
        </div>
      </div>

      <div class="form-section">
        <h3>页脚设置</h3>

        <div class="content-tips">
          <div class="tip-item">
            <strong>页脚描述：</strong>推荐包含品牌简介、经营理念、核心优势等，建议 50-150 字
          </div>
          <div class="tip-item">
            <strong>客服时间：</strong>推荐格式如"客服在线时间每周一至周六 9:00 到 18:00"
          </div>
        </div>

        <div class="form-group">
          <label>页脚描述</label>
          <textarea
            v-model="form.footer_description"
            rows="4"
            placeholder="输入页脚描述文字..."
          ></textarea>
          <span class="hint">显示在网站底部的主要描述文字，推荐包含品牌简介和核心优势</span>
        </div>

        <div class="form-group">
          <label>客服在线时间</label>
          <input
            v-model="form.footer_service_time"
            type="text"
            placeholder="例如：客服在线时间每周一至周六 9:00到18:00"
          />
          <span class="hint">显示在页脚描述下方的客服时间信息</span>
        </div>
      </div>

      <div class="form-section">
        <h3>邮件通知设置</h3>

        <div class="content-tips">
          <div class="tip-item">
            <strong>适用场景：</strong>订单发货后自动发送邮件通知客户，告知物流信息
          </div>
          <div class="tip-item">
            <strong>QQ邮箱：</strong>需在QQ邮箱设置中开启SMTP服务并获取授权码（非登录密码）
          </div>
          <div class="tip-item">
            <strong>通用性：</strong>支持任意SMTP服务（QQ、163、Gmail等），请根据服务商填写对应配置
          </div>
        </div>

        <div class="form-group">
          <label>启用邮件通知</label>
          <div class="toggle-group">
            <label class="toggle">
              <input type="radio" v-model="emailConfig.email_enabled" value="true" />
              <span>开启</span>
            </label>
            <label class="toggle">
              <input type="radio" v-model="emailConfig.email_enabled" value="false" />
              <span>关闭</span>
            </label>
          </div>
        </div>

        <div v-if="emailConfig.email_enabled === 'true'" class="email-fields">
          <div class="form-group">
            <label>SMTP 服务器地址</label>
            <input
              v-model="emailConfig.email_smtp_host"
              type="text"
              placeholder="例如：smtp.qq.com"
            />
            <span class="hint">QQ邮箱：smtp.qq.com，163邮箱：smtp.163.com，Gmail：smtp.gmail.com</span>
          </div>

          <div class="form-group">
            <label>SMTP 端口</label>
            <input
              v-model="emailConfig.email_smtp_port"
              type="text"
              placeholder="例如：465"
            />
            <span class="hint">QQ邮箱SSL端口：465，STARTTLS端口：587</span>
          </div>

          <div class="form-group">
            <label>发件邮箱账号</label>
            <input
              v-model="emailConfig.email_smtp_username"
              type="text"
              placeholder="例如：your_email@qq.com"
            />
            <span class="hint">用于发送邮件的邮箱地址</span>
          </div>

          <div class="form-group">
            <label>SMTP 授权码</label>
            <input
              v-model="emailConfig.email_smtp_password"
              type="password"
              placeholder="输入SMTP授权码"
            />
            <span class="hint">QQ邮箱：设置 → 账户 → POP3/SMTP服务 → 生成授权码；163邮箱：设置 → POP3/SMTP/IMAP → 开启并设置授权码</span>
          </div>

          <div class="form-group">
            <label>发件人显示名称</label>
            <input
              v-model="emailConfig.email_from_name"
              type="text"
              placeholder="例如：BYCIGAR官方商城"
            />
            <span class="hint">收件人看到的发件人名称</span>
          </div>

          <div class="form-group">
            <label>发送测试邮件</label>
            <div class="test-email-row">
              <input
                v-model="testEmail"
                type="email"
                placeholder="输入收件邮箱地址"
              />
              <button
                type="button"
                class="btn-test"
                :disabled="testingEmail"
                @click="sendTestEmail"
              >
                {{ testingEmail ? '发送中...' : '发送测试' }}
              </button>
            </div>
            <span class="hint">请先保存SMTP配置，再发送测试邮件验证配置是否正确</span>
          </div>
        </div>
      </div>

      <div class="form-actions">
        <button type="submit" class="btn-save" :disabled="saving">
          {{ saving ? '保存中...' : '保存设置' }}
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.admin-settings {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.settings-form {
  max-width: 600px;
}

.form-section {
  margin-bottom: 32px;
}

.form-section h3 {
  font-size: 15px;
  color: #333;
  margin: 0 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.content-tips {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px;
  background: #fdf8f3;
  border-radius: 6px;
  margin-bottom: 16px;
  font-size: 13px;
  color: #8a7560;
}

.content-tips .tip-item strong {
  color: #b08968;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #d4a574;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.form-group .hint {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: #999;
}

.toggle-group {
  display: flex;
  gap: 16px;
}

.toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #555;
}

.toggle input[type="radio"] {
  width: auto;
  margin: 0;
}

.email-fields {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.test-email-row {
  display: flex;
  gap: 10px;
}

.test-email-row input {
  flex: 1;
}

.btn-test {
  padding: 10px 20px;
  background: #fff;
  color: #d4a574;
  border: 1px solid #d4a574;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-test:hover {
  background: #d4a574;
  color: #fff;
}

.btn-test:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-actions {
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.btn-save {
  padding: 10px 24px;
  background: #d4a574;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-save:hover {
  background: #c49464;
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
