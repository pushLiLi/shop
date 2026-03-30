<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToastStore } from '../../stores/toast'
import AdminImageUpload from '../../components/AdminImageUpload.vue'

const API_BASE = '/api'
const toast = useToastStore()

const methods = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingMethod = ref(null)
const saving = ref(false)

const channelTypes = [
  { value: 'phone', label: '电话' },
  { value: 'email', label: '邮箱' },
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'wechat', label: '微信' },
  { value: 'qq', label: 'QQ' },
  { value: 'telegram', label: 'Telegram' },
  { value: 'custom', label: '自定义' }
]

const valuePlaceholders = {
  phone: '如：400-888-9999',
  email: '如：support@example.com',
  whatsapp: '含国际区号，如：8613800138000',
  wechat: '微信号',
  qq: 'QQ号码',
  telegram: 'Telegram 用户名',
  custom: '链接或文本'
}

const valueLabels = {
  phone: '电话号码',
  email: '邮箱地址',
  whatsapp: 'WhatsApp 号码',
  wechat: '微信号',
  qq: 'QQ 号',
  telegram: 'Telegram 用户名',
  custom: '链接或文本'
}

const form = ref({
  type: 'phone',
  label: '',
  value: '',
  qrCodeUrl: '',
  isActive: true,
  sortOrder: 0
})

const showQRUpload = computed(() => {
  return form.value.type === 'wechat' || form.value.type === 'qq' || form.value.qrCodeUrl
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const getTypeLabel = (type) => {
  return channelTypes.find(t => t.value === type)?.label || type
}

async function fetchMethods() {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/contact-methods`, { headers: authHeaders() })
    const data = await res.json()
    methods.value = data.contactMethods || []
  } catch (e) {
    toast.error('获取联系方式失败')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingMethod.value = null
  form.value = { type: 'phone', label: '', value: '', qrCodeUrl: '', isActive: true, sortOrder: 0 }
  showModal.value = true
}

function openEditModal(method) {
  editingMethod.value = method
  form.value = {
    type: method.type,
    label: method.label,
    value: method.value,
    qrCodeUrl: method.qrCodeUrl || '',
    isActive: method.isActive,
    sortOrder: method.sortOrder
  }
  showModal.value = true
}

async function saveMethod() {
  if (!form.value.label.trim()) {
    toast.error('请输入显示名称')
    return
  }
  if (!form.value.value.trim()) {
    toast.error('请输入联系方式')
    return
  }
  saving.value = true
  try {
    const url = editingMethod.value
      ? `${API_BASE}/admin/contact-methods/${editingMethod.value.id}`
      : `${API_BASE}/admin/contact-methods`
    const payload = { ...form.value }
    if (!showQRUpload.value) {
      payload.qrCodeUrl = ''
    }
    const res = await fetch(url, {
      method: editingMethod.value ? 'PUT' : 'POST',
      headers: authHeaders(),
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '操作失败')
    toast.success(editingMethod.value ? '更新成功' : '创建成功')
    showModal.value = false
    fetchMethods()
  } catch (e) {
    toast.error(e.message)
  } finally {
    saving.value = false
  }
}

async function deleteMethod(method) {
  if (!confirm(`确定删除联系方式「${method.label}」？`)) return
  try {
    const res = await fetch(`${API_BASE}/admin/contact-methods/${method.id}`, {
      method: 'DELETE',
      headers: authHeaders()
    })
    if (!res.ok) throw new Error('删除失败')
    toast.success('删除成功')
    fetchMethods()
  } catch (e) {
    toast.error(e.message)
  }
}

async function toggleActive(method) {
  try {
    const res = await fetch(`${API_BASE}/admin/contact-methods/${method.id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ isActive: !method.isActive })
    })
    if (!res.ok) throw new Error('操作失败')
    fetchMethods()
  } catch (e) {
    toast.error(e.message)
  }
}

onMounted(fetchMethods)
</script>

<template>
  <div class="admin-contact-methods">
    <div class="toolbar">
      <div class="toolbar-left">
        <h2>联系方式管理</h2>
      </div>
      <button class="btn-add" @click="openCreateModal">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        添加联系方式
      </button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="methods.length === 0" class="empty">
      <p>暂无联系方式</p>
      <button class="btn-add" @click="openCreateModal">添加第一个联系方式</button>
    </div>

    <div v-else class="methods-list">
      <div v-for="method in methods" :key="method.id" class="method-card">
        <div class="method-icon" :class="method.type">
          <svg v-if="method.type === 'phone'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
          <svg v-else-if="method.type === 'email'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
          <svg v-else-if="method.type === 'whatsapp'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
          <svg v-else-if="method.type === 'wechat'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8.5 2C5.46 2 3 4.46 3 7.5c0 1.68.75 3.18 1.94 4.2L4 14l2.6-1.3c.6.2 1.24.3 1.9.3.34 0 .67-.03 1-.08"></path><path d="M15.5 8c-3.04 0-5.5 2.46-5.5 5.5 0 3.04 2.46 5.5 5.5 5.5.66 0 1.3-.1 1.9-.3L20 20l-.94-2.3A5.47 5.47 0 0 0 21 13.5C21 10.46 18.54 8 15.5 8z"></path></svg>
          <svg v-else-if="method.type === 'qq'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><path d="M8 14s1.5 2 4 2 4-2 4-2"></path><line x1="9" y1="9" x2="9.01" y2="9"></line><line x1="15" y1="9" x2="15.01" y2="9"></line></svg>
          <svg v-else-if="method.type === 'telegram'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
          <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path></svg>
        </div>
        <div class="method-info">
          <div class="method-header">
            <span class="method-name">{{ method.label }}</span>
            <span class="type-tag">{{ getTypeLabel(method.type) }}</span>
            <span :class="['status-tag', method.isActive ? 'active' : 'inactive']">
              {{ method.isActive ? '启用' : '停用' }}
            </span>
          </div>
          <div class="method-value">{{ method.value }}</div>
          <div class="method-meta">
            <span>排序: {{ method.sortOrder }}</span>
          </div>
        </div>
        <div class="method-qr" v-if="method.qrCodeUrl">
          <img :src="method.qrCodeUrl" alt="二维码" class="qr-preview" />
        </div>
        <div class="method-actions">
          <button class="btn-toggle" @click="toggleActive(method)">
            {{ method.isActive ? '停用' : '启用' }}
          </button>
          <button class="btn-edit" @click="openEditModal(method)">编辑</button>
          <button class="btn-delete" @click="deleteMethod(method)">删除</button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingMethod ? '编辑联系方式' : '添加联系方式' }}</h3>
          <button class="modal-close" @click="showModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>渠道类型</label>
            <select v-model="form.type">
              <option v-for="ct in channelTypes" :key="ct.value" :value="ct.value">{{ ct.label }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>显示名称</label>
            <input v-model="form.label" type="text" :placeholder="'如：' + getTypeLabel(form.type) + '客服'" />
          </div>
          <div class="form-group">
            <label>{{ valueLabels[form.type] || '联系方式' }}</label>
            <input v-model="form.value" type="text" :placeholder="valuePlaceholders[form.type]" />
          </div>
          <div class="form-group" v-if="showQRUpload">
            <label>二维码图片</label>
            <AdminImageUpload v-model="form.qrCodeUrl" :aspect-ratio="1" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>排序（数字越小越靠前）</label>
              <input v-model.number="form.sortOrder" type="number" min="0" />
            </div>
            <div class="form-group">
              <label>状态</label>
              <select v-model="form.isActive">
                <option :value="true">启用</option>
                <option :value="false">停用</option>
              </select>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showModal = false">取消</button>
          <button class="btn-save" @click="saveMethod" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-contact-methods {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
}

.toolbar-left h2 {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.btn-add {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #d4a574;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-add:hover {
  background: #c49464;
}

.loading {
  padding: 40px;
  text-align: center;
  color: #666;
}

.empty {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.methods-list {
  padding: 15px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.method-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #eee;
}

.method-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #fff;
}

.method-icon.phone { background: #4caf50; }
.method-icon.email { background: #2196f3; }
.method-icon.whatsapp { background: #25d366; }
.method-icon.wechat { background: #07c160; }
.method-icon.qq { background: #12b7f5; }
.method-icon.telegram { background: #0088cc; }
.method-icon.custom { background: #ff9800; }

.method-info {
  flex: 1;
  min-width: 0;
}

.method-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.method-name {
  font-weight: 600;
  color: #333;
  font-size: 15px;
}

.type-tag {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
  background: #f0f0f0;
  color: #666;
}

.status-tag {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}

.status-tag.active {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-tag.inactive {
  background: #f5f5f5;
  color: #999;
}

.method-value {
  color: #666;
  font-size: 13px;
  margin-bottom: 4px;
  word-break: break-all;
}

.method-meta {
  font-size: 12px;
  color: #aaa;
}

.method-qr {
  flex-shrink: 0;
}

.qr-preview {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid #eee;
}

.method-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.btn-toggle,
.btn-edit,
.btn-delete {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.btn-toggle {
  background: #e3f2fd;
  color: #1565c0;
}

.btn-toggle:hover {
  background: #bbdefb;
}

.btn-edit {
  background: #fff8e1;
  color: #f57c00;
}

.btn-edit:hover {
  background: #ffecb3;
}

.btn-delete {
  background: #ffebee;
  color: #c62828;
}

.btn-delete:hover {
  background: #ffcdd2;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  border-radius: 8px;
  width: 90%;
  max-width: 520px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  font-size: 16px;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #d4a574;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 15px 20px;
  border-top: 1px solid #eee;
}

.btn-cancel,
.btn-save {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-cancel {
  background: #f5f5f5;
  color: #666;
}

.btn-cancel:hover {
  background: #e0e0e0;
}

.btn-save {
  background: #d4a574;
  color: #fff;
}

.btn-save:hover {
  background: #c49464;
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
