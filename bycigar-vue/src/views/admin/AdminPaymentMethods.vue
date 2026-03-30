<script setup>
import { ref, onMounted } from 'vue'
import { useToastStore } from '../../stores/toast'
import AdminImageUpload from '../../components/AdminImageUpload.vue'

const API_BASE = '/api'
const toast = useToastStore()

const methods = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingMethod = ref(null)
const saving = ref(false)

const form = ref({
  name: '',
  qrCodeUrl: '',
  instructions: '',
  isActive: true,
  sortOrder: 0
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

async function fetchMethods() {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/payment-methods`, { headers: authHeaders() })
    const data = await res.json()
    methods.value = data.paymentMethods || []
  } catch (e) {
    toast.error('获取付款方式失败')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingMethod.value = null
  form.value = { name: '', qrCodeUrl: '', instructions: '', isActive: true, sortOrder: 0 }
  showModal.value = true
}

function openEditModal(method) {
  editingMethod.value = method
  form.value = {
    name: method.name,
    qrCodeUrl: method.qrCodeUrl,
    instructions: method.instructions,
    isActive: method.isActive,
    sortOrder: method.sortOrder
  }
  showModal.value = true
}

async function saveMethod() {
  if (!form.value.name.trim()) {
    toast.error('请输入付款方式名称')
    return
  }
  saving.value = true
  try {
    const url = editingMethod.value
      ? `${API_BASE}/admin/payment-methods/${editingMethod.value.id}`
      : `${API_BASE}/admin/payment-methods`
    const res = await fetch(url, {
      method: editingMethod.value ? 'PUT' : 'POST',
      headers: authHeaders(),
      body: JSON.stringify(form.value)
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
  if (!confirm(`确定删除付款方式「${method.name}」？`)) return
  try {
    const res = await fetch(`${API_BASE}/admin/payment-methods/${method.id}`, {
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
    const res = await fetch(`${API_BASE}/admin/payment-methods/${method.id}`, {
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
  <div class="admin-payment-methods">
    <div class="toolbar">
      <div class="toolbar-left">
        <h2>付款方式管理</h2>
      </div>
      <button class="btn-add" @click="openCreateModal">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        添加付款方式
      </button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="methods.length === 0" class="empty">
      <p>暂无付款方式</p>
      <button class="btn-add" @click="openCreateModal">添加第一个付款方式</button>
    </div>

    <div v-else class="methods-list">
      <div v-for="method in methods" :key="method.id" class="method-card">
        <div class="method-info">
          <div class="method-header">
            <span class="method-name">{{ method.name }}</span>
            <span :class="['status-tag', method.isActive ? 'active' : 'inactive']">
              {{ method.isActive ? '启用' : '停用' }}
            </span>
          </div>
          <div v-if="method.instructions" class="method-instructions">{{ method.instructions }}</div>
          <div class="method-meta">
            <span>排序: {{ method.sortOrder }}</span>
          </div>
        </div>
        <div class="method-qr" v-if="method.qrCodeUrl">
          <img :src="method.qrCodeUrl" alt="收款码" class="qr-preview" />
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
          <h3>{{ editingMethod ? '编辑付款方式' : '添加付款方式' }}</h3>
          <button class="modal-close" @click="showModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>名称</label>
            <input v-model="form.name" type="text" placeholder="如：微信支付、支付宝、PayPal" />
          </div>
          <div class="form-group">
            <label>收款码/收款信息图片</label>
            <AdminImageUpload v-model="form.qrCodeUrl" />
          </div>
          <div class="form-group">
            <label>付款说明</label>
            <textarea v-model="form.instructions" placeholder="如：请扫码支付后上传付款截图" rows="3"></textarea>
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
.admin-payment-methods {
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

.method-info {
  flex: 1;
  min-width: 0;
}

.method-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.method-name {
  font-weight: 600;
  color: #333;
  font-size: 15px;
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

.method-instructions {
  color: #666;
  font-size: 13px;
  margin-bottom: 4px;
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

.form-group textarea {
  resize: vertical;
  min-height: 60px;
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
