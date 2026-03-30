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
    }

    const configData = await configRes.json()
    banners.value.home_banner_1 = configData.home_banner_1 || ''
  } catch (e) {
    toast.error('获取设置失败')
  } finally {
    loading.value = false
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
      })
    ]

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
