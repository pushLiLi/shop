<script setup>
import { ref, onMounted, computed } from 'vue'

const API_BASE = '/api'

const categories = ref([])
const loading = ref(false)
const showModal = ref(false)
const modalMode = ref('create')
const saving = ref(false)

const form = ref({
  id: null,
  name: '',
  slug: '',
  parentId: ''
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const parentCategories = computed(() => {
  return categories.value.filter(c => !c.parentId)
})

const fetchCategories = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/categories`, {
      headers: authHeaders()
    })
    categories.value = await res.json()
  } catch (e) {
    console.error('Error fetching categories:', e)
  } finally {
    loading.value = false
  }
}

const openCreateModal = (parentId = null) => {
  modalMode.value = 'create'
  form.value = {
    id: null,
    name: '',
    slug: '',
    parentId: parentId || ''
  }
  showModal.value = true
}

const openEditModal = (category) => {
  modalMode.value = 'edit'
  form.value = {
    id: category.id,
    name: category.name,
    slug: category.slug,
    parentId: category.parentId || ''
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveCategory = async () => {
  if (!form.value.name) {
    alert('请填写分类名称')
    return
  }

  saving.value = true
  try {
    const url = modalMode.value === 'create' 
      ? `${API_BASE}/admin/categories`
      : `${API_BASE}/admin/categories/${form.value.id}`
    
    const body = {
      name: form.value.name,
      slug: form.value.slug,
      parentId: form.value.parentId ? parseInt(form.value.parentId) : null
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
    fetchCategories()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

const deleteCategory = async (category) => {
  if (!confirm(`确定要删除分类「${category.name}」吗？`)) return

  try {
    const res = await fetch(`${API_BASE}/admin/categories/${category.id}`, {
      method: 'DELETE',
      headers: authHeaders()
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '删除失败')
    }

    fetchCategories()
  } catch (e) {
    alert(e.message)
  }
}

const getChildren = (parentId) => {
  return categories.value.filter(c => c.parentId === parentId)
}

onMounted(() => {
  fetchCategories()
})
</script>

<template>
  <div class="admin-categories">
    <div class="toolbar">
      <div class="toolbar-left">
        <h3>分类列表</h3>
        <span class="count">共 {{ categories.length }} 个</span>
      </div>
      <div class="toolbar-right">
        <button class="btn-add" @click="openCreateModal()">+ 添加顶级分类</button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="category-list">
      <div v-for="category in parentCategories" :key="category.id" class="category-group">
        <div class="category-item parent">
          <div class="category-info">
            <span class="category-name">{{ category.name }}</span>
            <span class="category-slug">{{ category.slug }}</span>
          </div>
          <div class="category-actions">
            <button class="btn-add-child" @click="openCreateModal(category.id)">+ 添加子分类</button>
            <button class="btn-edit" @click="openEditModal(category)">编辑</button>
            <button class="btn-delete" @click="deleteCategory(category)">删除</button>
          </div>
        </div>
        
        <div v-for="child in getChildren(category.id)" :key="child.id" class="category-item child">
          <div class="category-info">
            <span class="indent">└</span>
            <span class="category-name">{{ child.name }}</span>
            <span class="category-slug">{{ child.slug }}</span>
          </div>
          <div class="category-actions">
            <button class="btn-edit" @click="openEditModal(child)">编辑</button>
            <button class="btn-delete" @click="deleteCategory(child)">删除</button>
          </div>
        </div>
      </div>

      <div v-if="categories.length === 0" class="empty-text">
        暂无分类数据
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ modalMode === 'create' ? '添加分类' : '编辑分类' }}</h3>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>分类名称 <span class="required">*</span></label>
            <input v-model="form.name" type="text" placeholder="如: 高希霸">
          </div>
          <div class="form-group">
            <label>Slug（留空自动生成）</label>
            <input v-model="form.slug" type="text" placeholder="如: cohiba">
          </div>
          <div class="form-group">
            <label>父级分类</label>
            <select v-model="form.parentId">
              <option value="">无（顶级分类）</option>
              <option 
                v-for="cat in parentCategories" 
                :key="cat.id" 
                :value="cat.id"
                :disabled="modalMode === 'edit' && cat.id === form.id"
              >
                {{ cat.name }}
              </option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button class="btn-save" :disabled="saving" @click="saveCategory">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-categories {
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

.category-list {
  padding: 20px;
}

.category-group {
  margin-bottom: 10px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  border: 1px solid #eee;
  border-radius: 6px;
  margin-bottom: 5px;
}

.category-item.parent {
  background: #fafafa;
}

.category-item.child {
  margin-left: 30px;
  background: #fff;
}

.category-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.indent {
  color: #ccc;
  font-size: 14px;
}

.category-name {
  font-weight: 500;
}

.category-slug {
  color: #999;
  font-size: 13px;
}

.category-actions {
  display: flex;
  gap: 8px;
}

.btn-add-child,
.btn-edit,
.btn-delete {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.btn-add-child {
  background: #e8f5e9;
  color: #2e7d32;
}

.btn-add-child:hover {
  background: #c8e6c9;
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
  max-width: 450px;
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
