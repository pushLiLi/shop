<script setup>
import { ref, onMounted } from 'vue'
import AdminImageUpload from '../../components/AdminImageUpload.vue'
import { useToastStore } from '../../stores/toast'

const API_BASE = 'http://localhost:3000/api'
const toast = useToastStore()

const banners = ref([])
const loading = ref(false)
const showModal = ref(false)
const showDeleteConfirm = ref(false)
const deleteTarget = ref(null)
const modalMode = ref('create')
const saving = ref(false)

const dragIndex = ref(null)
const dragOverIndex = ref(null)

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
    const res = await fetch(`${API_BASE}/admin/banners`, { headers: authHeaders() })
    if (!res.ok) throw new Error('获取轮播图失败')
    banners.value = await res.json()
  } catch (e) {
    toast.error(e.message)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  modalMode.value = 'create'
  form.value = { id: null, title: '', imageUrl: '', link: '', sortOrder: 0, isActive: true }
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

const closeModal = () => { showModal.value = false }

const confirmDelete = (banner) => {
  deleteTarget.value = banner
  showDeleteConfirm.value = true
}

const cancelDelete = () => {
  showDeleteConfirm.value = false
  deleteTarget.value = null
}

const moveSortUp = () => {
  form.value.sortOrder = Math.max(0, form.value.sortOrder - 1)
}

const moveSortDown = () => {
  form.value.sortOrder++
}

const validateLink = (link) => {
  if (!link) return true
  return link.startsWith('/') || /^https?:\/\//.test(link)
}

const saveBanner = async () => {
  if (!form.value.imageUrl) {
    toast.error('请上传轮播图片')
    return
  }

  if (form.value.link && !validateLink(form.value.link)) {
    toast.error('跳转链接格式不正确，请输入站内路径或完整 URL')
    return
  }

  saving.value = true
  try {
    const url = modalMode.value === 'create'
      ? `${API_BASE}/admin/banners`
      : `${API_BASE}/admin/banners/${form.value.id}`

    const body = {
      title: form.value.title,
      imageUrl: form.value.imageUrl,
      link: form.value.link || '',
      sortOrder: parseInt(form.value.sortOrder) || 0,
      isActive: form.value.isActive
    }

    const res = await fetch(url, {
      method: modalMode.value === 'create' ? 'POST' : 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(body)
    })

    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '保存失败')

    toast.success(modalMode.value === 'create' ? '轮播图已添加' : '轮播图已更新')
    closeModal()
    fetchBanners()
  } catch (e) {
    toast.error(e.message)
  } finally {
    saving.value = false
  }
}

const executeDelete = async () => {
  const banner = deleteTarget.value
  showDeleteConfirm.value = false

  try {
    const res = await fetch(`${API_BASE}/admin/banners/${banner.id}`, {
      method: 'DELETE',
      headers: authHeaders()
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '删除失败')
    toast.success('轮播图已删除')
    fetchBanners()
  } catch (e) {
    toast.error(e.message)
  }
}

const toggleActive = async (banner) => {
  try {
    const res = await fetch(`${API_BASE}/admin/banners/${banner.id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        title: banner.title,
        imageUrl: banner.imageUrl,
        link: banner.link,
        sortOrder: banner.sortOrder,
        isActive: !banner.isActive
      })
    })
    if (res.ok) {
      banner.isActive = !banner.isActive
      toast.success(banner.isActive ? '已启用' : '已禁用')
    }
  } catch (e) {
    toast.error('状态切换失败')
  }
}

const onDragStart = (index) => { dragIndex.value = index }
const onDragOver = (e, index) => {
  e.preventDefault()
  dragOverIndex.value = index
}
const onDragLeave = () => { dragOverIndex.value = null }
const onDrop = async (dropIndex) => {
  if (dragIndex.value === null || dragIndex.value === dropIndex) return
  const item = banners.value.splice(dragIndex.value, 1)[0]
  banners.value.splice(dropIndex, 0, item)

  const updates = banners.value.map((b, i) =>
    fetch(`${API_BASE}/admin/banners/${b.id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ title: b.title, imageUrl: b.imageUrl, link: b.link, sortOrder: i, isActive: b.isActive })
    })
  )
  try {
    await Promise.all(updates)
    toast.success('排序已更新')
  } catch (e) {
    toast.error('排序更新失败')
    fetchBanners()
  }
  dragIndex.value = null
  dragOverIndex.value = null
}
const onDragEnd = () => { dragIndex.value = null; dragOverIndex.value = null }

onMounted(() => { fetchBanners() })
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
      <div
        v-for="(banner, index) in banners"
        :key="banner.id"
        class="banner-item"
        :class="{ 'drag-over': dragOverIndex === index, 'dragging': dragIndex === index }"
        draggable="true"
        @dragstart="onDragStart(index)"
        @dragover="onDragOver($event, index)"
        @dragleave="onDragLeave"
        @drop="onDrop(index)"
        @dragend="onDragEnd"
      >
        <div class="drag-handle" title="拖拽排序">⠿</div>
        <div class="banner-image">
          <img :src="banner.imageUrl" :alt="banner.title" loading="lazy">
        </div>
        <div class="banner-info">
          <div class="banner-title">{{ banner.title || '无标题' }}</div>
          <div class="banner-meta">
            <span class="meta-sort">排序: {{ banner.sortOrder || 0 }}</span>
          </div>
        </div>
        <div class="banner-controls">
          <label class="switch">
            <input type="checkbox" :checked="banner.isActive" @change="toggleActive(banner)">
            <span class="slider"></span>
          </label>
          <span class="switch-value">{{ banner.isActive ? '启用' : '禁用' }}</span>
        </div>
        <div class="banner-actions">
          <button class="btn-edit" @click="openEditModal(banner)">编辑</button>
          <button class="btn-delete" @click="confirmDelete(banner)">删除</button>
        </div>
      </div>

      <div v-if="banners.length === 0" class="empty-text">暂无轮播图数据</div>
    </div>

    <Transition name="modal">
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
              <AdminImageUpload v-model="form.imageUrl" :aspect-ratio="7/3" />
            </div>
            <div class="form-group">
              <label>跳转链接</label>
              <div class="input-group">
                <span class="input-prefix">链接</span>
                <input v-model="form.link" type="text" placeholder="留空或输入 /category/cohiba、https://...">
              </div>
            </div>
            <div class="form-group">
              <label>排序权重</label>
              <div class="sort-control">
                <button type="button" class="sort-btn" @click="moveSortUp">-</button>
                <input v-model.number="form.sortOrder" type="number" min="0" placeholder="0">
                <button type="button" class="sort-btn" @click="moveSortDown">+</button>
              </div>
              <div class="field-hint">数字越小排序越靠前</div>
            </div>
            <div class="form-section">
              <div class="section-title">显示状态</div>
              <div class="switch-item">
                <span class="switch-label">启用轮播</span>
                <label class="switch">
                  <input type="checkbox" v-model="form.isActive">
                  <span class="slider"></span>
                </label>
                <span class="switch-value">{{ form.isActive ? '启用' : '禁用' }}</span>
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
    </Transition>

    <Transition name="modal">
      <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="cancelDelete">
        <div class="modal modal-sm">
          <div class="modal-header">
            <h3>确认删除</h3>
            <button class="modal-close" @click="cancelDelete">&times;</button>
          </div>
          <div class="modal-body">
            <p class="confirm-text">确定要删除该轮播图吗？此操作不可恢复。</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="cancelDelete">取消</button>
            <button class="btn-delete-confirm" @click="executeDelete">确认删除</button>
          </div>
        </div>
      </div>
    </Transition>
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
  gap: 16px;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.25s ease;
  cursor: grab;
  background: #fff;
}

.banner-item:active {
  cursor: grabbing;
}

.banner-item:hover {
  border-color: #ddd;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.banner-item.dragging {
  opacity: 0.4;
}

.banner-item.drag-over {
  border-color: #d4a574;
  box-shadow: 0 0 0 2px rgba(212, 165, 116, 0.2);
}

.drag-handle {
  color: #ccc;
  font-size: 20px;
  line-height: 1;
  user-select: none;
  flex-shrink: 0;
  transition: color 0.2s;
}

.banner-item:hover .drag-handle {
  color: #999;
}

.banner-image {
  width: 240px;
  flex-shrink: 0;
}

.banner-image img {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border-radius: 6px;
  transition: transform 0.2s;
  background: #f0f0f0;
}

.banner-image img:hover {
  transform: scale(1.03);
}

.banner-info {
  flex: 1;
  min-width: 0;
}

.banner-title {
  font-weight: 600;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.banner-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 13px;
  color: #999;
}

.meta-sort {
  font-size: 13px;
  color: #999;
}

.banner-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.switch-value {
  font-size: 13px;
  color: #999;
  min-width: 24px;
}

.banner-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
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

.modal-sm {
  max-width: 400px;
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

.confirm-text {
  margin: 0;
  color: #333;
  font-size: 15px;
  line-height: 1.6;
}

.form-group {
  margin-bottom: 18px;
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
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #d4a574;
}

.required {
  color: #dc3545;
}

.input-group {
  display: flex;
  align-items: stretch;
}

.input-group input {
  flex: 1;
  border-radius: 0 4px 4px 0;
}

.input-prefix {
  display: flex;
  align-items: center;
  padding: 0 12px;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-right: none;
  border-radius: 4px 0 0 4px;
  color: #666;
  font-size: 14px;
}

.sort-control {
  display: flex;
  align-items: center;
  gap: 4px;
}

.sort-control input {
  width: 48px;
  height: 32px;
  text-align: center;
  padding: 0;
  border-radius: 4px;
  background: #f5f5f5;
  border: 1px solid #ddd;
  font-size: 14px;
  color: #333;
  -moz-appearance: textfield;
}

.sort-control input::-webkit-outer-spin-button,
.sort-control input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.sort-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #ddd;
  background: #f5f5f5;
  cursor: pointer;
  font-size: 16px;
  color: #555;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  padding: 0;
}

.sort-btn:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.field-hint {
  margin-top: 6px;
  font-size: 12px;
  color: #999;
}

.form-section {
  margin-bottom: 18px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.form-section:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  padding-left: 10px;
  border-left: 3px solid #d4a574;
  margin-bottom: 15px;
}

.switch-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.switch-label {
  font-size: 14px;
  color: #666;
  min-width: 60px;
}

.switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.3s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

.switch input:checked + .slider {
  background-color: #d4a574;
}

.switch input:checked + .slider:before {
  transform: translateX(20px);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 15px 20px;
  border-top: 1px solid #eee;
}

.btn-cancel,
.btn-save,
.btn-delete-confirm {
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

.btn-delete-confirm {
  background: #dc3545;
  color: #fff;
}

.btn-delete-confirm:hover {
  background: #c82333;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.25s ease;
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: transform 0.25s ease, opacity 0.25s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: scale(0.95) translateY(10px);
  opacity: 0;
}

@media (max-width: 768px) {
  .banner-item {
    flex-wrap: wrap;
    gap: 12px;
  }

  .banner-image {
    width: 100%;
  }

  .banner-image img {
    height: 160px;
  }

  .banner-info {
    flex: 1 1 calc(100% - 60px);
  }

  .banner-controls {
    flex: 0 0 auto;
  }

  .banner-actions {
    flex: 1;
    justify-content: flex-end;
  }
}
</style>
