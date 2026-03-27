<script setup>
import { ref, onMounted } from 'vue'
import AdminImageUpload from '../../components/AdminImageUpload.vue'

const API_BASE = 'http://localhost:3000/api'

const banners = ref([])
const loading = ref(false)
const showModal = ref(false)
const modalMode = ref('create')
const saving = ref(false)

const form = ref({
  id: null,
  title: '',
  imageUrl: '',
  link: '',
  sortOrder: 0,
  isActive: true
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const fetchBanners = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/banners`, {
      headers: authHeaders()
    })
    banners.value = await res.json()
  } catch (e) {
    console.error('Error fetching banners:', e)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  modalMode.value = 'create'
  form.value = {
    id: null,
    title: '',
    imageUrl: '',
    link: '',
    sortOrder: 0,
    isActive: true
  }
  showModal.value = true
}

const openEditModal = (banner) => {
  modalMode.value = 'edit'
  form.value = {
    id: banner.id,
    title: banner.title || '',
    imageUrl: banner.imageUrl || '',
    link: banner.link || '',
    sortOrder: banner.sortOrder || 0,
    isActive: banner.isActive
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveBanner = async () => {
  if (!form.value.imageUrl) {
    alert('请上传轮播图片')
    return
  }

  saving.value = true
  try {
    const url = modalMode.value === 'create' 
      ? `${API_BASE}/admin/banners`
      : `${API_BASE}/admin/banners/${form.value.id}`
    
    const body = {
      title: form.value.title,
      image: form.value.imageUrl,
      link: form.value.link || '#',
      sortOrder: parseInt(form.value.sortOrder) || 0,
      isActive: form.value.isActive
    }

    const res = await fetch(url, {
      method: modalMode.value === 'create' ? 'POST' : 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(body)
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '保存失败')
    }

    closeModal()
    fetchBanners()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

const deleteBanner = async (banner) => {
  if (!confirm(`确定要删除该轮播图吗？`)) return

  try {
    const res = await fetch(`${API_BASE}/admin/banners/${banner.id}`, {
      method: 'DELETE',
      headers: authHeaders()
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '删除失败')
    }

    fetchBanners()
  } catch (e) {
    alert(e.message)
  }
}

const toggleActive = async (banner) => {
  try {
    const res = await fetch(`${API_BASE}/admin/banners/${banner.id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        title: banner.title,
        image: banner.imageUrl,
        link: banner.link,
        sortOrder: banner.sortOrder,
        isActive: !banner.isActive
      })
    })

    if (res.ok) {
      banner.isActive = !banner.isActive
    }
  } catch (e) {
    console.error('Error toggling active:', e)
  }
}

onMounted(() => {
  fetchBanners()
})
</script>

<template>
  <div class="admin-banners">
    <div class="toolbar">
      <div class="toolbar-left">
        <h3>轮播图列表</h3>
        <span class="count">共 {{ banners.length }} 个</span>
      </div>
      <div class="toolbar-right">
        <button class="btn-add" @click="openCreateModal">+ 添加轮播图</button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="banner-list">
      <div v-for="banner in banners" :key="banner.id" class="banner-item">
        <div class="banner-image">
          <img :src="banner.imageUrl" :alt="banner.title">
        </div>
        <div class="banner-info">
          <div class="banner-title">{{ banner.title || '无标题' }}</div>
          <div class="banner-link">{{ banner.link || '-' }}</div>
          <div class="banner-meta">
            <span>排序: {{ banner.sortOrder || 0 }}</span>
            <span 
              class="badge" 
              :class="banner.isActive ? 'badge-success' : 'badge-danger'"
              @click="toggleActive(banner)"
              style="cursor: pointer"
            >
              {{ banner.isActive ? '启用' : '禁用' }}
            </span>
          </div>
        </div>
        <div class="banner-actions">
          <button class="btn-edit" @click="openEditModal(banner)">编辑</button>
          <button class="btn-delete" @click="deleteBanner(banner)">删除</button>
        </div>
      </div>

      <div v-if="banners.length === 0" class="empty-text">
        暂无轮播图数据
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ modalMode === 'create' ? '添加轮播图' : '编辑轮播图' }}</h3>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标题（可选）</label>
            <input v-model="form.title" type="text" placeholder="如: 新品上市">
          </div>
          <div class="form-group">
            <label>轮播图片 <span class="required">*</span></label>
            <AdminImageUpload v-model="form.imageUrl" />
          </div>
          <div class="form-group">
            <label>跳转链接</label>
            <input v-model="form.link" type="text" placeholder="如: /category/cohiba 或 https://...">
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>排序（数字越小越靠前）</label>
              <input v-model.number="form.sortOrder" type="number" min="0" placeholder="0">
            </div>
            <div class="form-group checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="form.isActive">
                启用
              </label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button class="btn-save" :disabled="saving" @click="saveBanner">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-banners {
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

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.toolbar-left h3 {
  margin: 0;
  font-size: 16px;
}

.count {
  color: #999;
  font-size: 14px;
}

.btn-add {
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

.banner-list {
  padding: 20px;
}

.banner-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
  transition: all 0.2s;
}

.banner-item:hover {
  border-color: #ddd;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.banner-image {
  width: 200px;
  flex-shrink: 0;
}

.banner-image img {
  width: 100%;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
}

.banner-info {
  flex: 1;
}

.banner-title {
  font-weight: 600;
  margin-bottom: 5px;
}

.banner-link {
  color: #666;
  font-size: 13px;
  margin-bottom: 8px;
}

.banner-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  font-size: 13px;
  color: #999;
}

.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.badge-success {
  background: #e8f5e9;
  color: #2e7d32;
}

.badge-danger {
  background: #ffebee;
  color: #c62828;
}

.banner-actions {
  display: flex;
  gap: 8px;
}

.btn-edit,
.btn-delete {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.btn-edit {
  background: #e3f2fd;
  color: #1565c0;
}

.btn-edit:hover {
  background: #bbdefb;
}

.btn-delete {
  background: #ffebee;
  color: #c62828;
}

.btn-delete:hover {
  background: #ffcdd2;
}

.empty-text {
  text-align: center;
  color: #999;
  padding: 40px;
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
  max-width: 550px;
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
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #333;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #d4a574;
}

.form-row {
  display: flex;
  gap: 15px;
}

.form-row .form-group {
  flex: 1;
}

.checkbox-group {
  display: flex;
  align-items: flex-end;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  padding: 10px 0;
}

.checkbox-label input {
  width: auto;
}

.required {
  color: #dc3545;
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
