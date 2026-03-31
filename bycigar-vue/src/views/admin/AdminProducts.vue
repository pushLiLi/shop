<script setup>
import { ref, onMounted, computed } from 'vue'
import AdminImageUpload from '../../components/AdminImageUpload.vue'
import { useToastStore } from '../../stores/toast'

const API_BASE = '/api'
const toast = useToastStore()

const products = ref([])
const categories = ref([])
const loading = ref(false)
const showModal = ref(false)
const modalMode = ref('create')
const saving = ref(false)

const search = ref('')
const filterCategory = ref('')
const filterFeatured = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 20
const sortBy = ref('id')
const sortOrder = ref('desc')

const form = ref({
  id: null,
  name: '',
  slug: '',
  description: '',
  price: 0,
  imageUrl: '',
  thumbnailUrl: '',
  categoryId: '',
  stock: 0,
  isActive: true,
  isFeatured: false,
  images: ''
})

const selectedIds = ref([])
const selectAll = ref(false)
const imageList = ref([])

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const fetchCategories = async () => {
  try {
    const res = await fetch(`${API_BASE}/categories`)
    categories.value = await res.json()
  } catch (e) {
    console.error('Error fetching categories:', e)
  }
}

const fetchProducts = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: currentPage.value,
      limit
    })
    if (search.value) params.append('search', search.value)
    if (filterCategory.value) params.append('categoryId', filterCategory.value)
    if (filterFeatured.value) params.append('featured', filterFeatured.value)
    if (sortBy.value) params.append('sortBy', sortBy.value)
    if (sortOrder.value) params.append('sortOrder', sortOrder.value)

    const res = await fetch(`${API_BASE}/admin/products?${params}`, {
      headers: authHeaders()
    })
    const data = await res.json()
    products.value = data.products || []
    totalPages.value = data.totalPages || 1
  } catch (e) {
    console.error('Error fetching products:', e)
  } finally {
    loading.value = false
  }
}

const handleSort = (field) => {
  if (sortBy.value === field) {
    if (sortOrder.value === 'desc') {
      sortOrder.value = 'asc'
    } else {
      sortBy.value = ''
      sortOrder.value = 'desc'
    }
  } else {
    sortBy.value = field
    sortOrder.value = 'desc'
  }
  currentPage.value = 1
  fetchProducts()
}

const sortIcon = (field) => {
  if (sortBy.value !== field) return ''
  return sortOrder.value === 'desc' ? ' ↓' : ' ↑'
}

const handleSearch = () => {
  currentPage.value = 1
  fetchProducts()
}

const resetFilters = () => {
  search.value = ''
  filterCategory.value = ''
  filterFeatured.value = ''
  currentPage.value = 1
  fetchProducts()
}

const openCreateModal = () => {
  modalMode.value = 'create'
  form.value = {
    id: null,
    name: '',
    slug: '',
    description: '',
    price: 0,
    imageUrl: '',
    thumbnailUrl: '',
    categoryId: '',
    stock: 0,
    isActive: true,
    isFeatured: false,
    images: ''
  }
  imageList.value = []
  showModal.value = true
}

const openEditModal = (product) => {
  modalMode.value = 'edit'
  form.value = {
    id: product.id,
    name: product.name,
    slug: product.slug,
    description: product.description || '',
    price: product.price,
    imageUrl: product.imageUrl || '',
    thumbnailUrl: product.thumbnailUrl || '',
    categoryId: product.categoryId || '',
    stock: product.stock || 0,
    isActive: product.isActive,
    isFeatured: product.isFeatured
  }
  const extra = product.images ? product.images.split(',').filter(Boolean) : []
  imageList.value = extra
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveProduct = async () => {
  if (!form.value.name || !form.value.price) {
    toast.error('请填写商品名称和价格')
    return
  }

  saving.value = true
  try {
    const url = modalMode.value === 'create' 
      ? `${API_BASE}/admin/products`
      : `${API_BASE}/admin/products/${form.value.id}`
    
    const body = {
      name: form.value.name,
      slug: form.value.slug,
      description: form.value.description,
      price: parseFloat(form.value.price),
      imageUrl: form.value.imageUrl,
      thumbnailUrl: form.value.thumbnailUrl || '',
      categoryId: form.value.categoryId ? parseInt(form.value.categoryId) : 0,
      stock: parseInt(form.value.stock) || 0,
      isActive: form.value.isActive,
      isFeatured: form.value.isFeatured,
      images: imageList.value.join(',')
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
    fetchProducts()
  } catch (e) {
    toast.error(e.message)
  } finally {
    saving.value = false
  }
}

const deleteProduct = async (product) => {
  if (!confirm(`确定要删除商品「${product.name}」吗？`)) return

  try {
    const res = await fetch(`${API_BASE}/admin/products/${product.id}`, {
      method: 'DELETE',
      headers: authHeaders()
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '删除失败')
    }

    fetchProducts()
  } catch (e) {
    toast.error(e.message)
  }
}

const toggleSelectAll = () => {
  if (selectAll.value) {
    selectedIds.value = products.value.map(p => p.id)
  } else {
    selectedIds.value = []
  }
}

const batchUpdateStatus = async (isActive) => {
  if (selectedIds.value.length === 0) {
    toast.error('请选择商品')
    return
  }
  try {
    const res = await fetch(`${API_BASE}/admin/products/batch/status`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ ids: selectedIds.value, isActive })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '操作失败')
    toast.success(data.message)
    selectedIds.value = []
    selectAll.value = false
    fetchProducts()
  } catch (e) {
    toast.error(e.message)
  }
}

const batchDelete = async () => {
  if (selectedIds.value.length === 0) {
    toast.error('请选择商品')
    return
  }
  if (!confirm(`确定要删除选中的 ${selectedIds.value.length} 个商品吗？`)) return
  try {
    const res = await fetch(`${API_BASE}/admin/products/batch`, {
      method: 'DELETE',
      headers: authHeaders(),
      body: JSON.stringify({ ids: selectedIds.value })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '删除失败')
    toast.success(data.message)
    selectedIds.value = []
    selectAll.value = false
    fetchProducts()
  } catch (e) {
    toast.error(e.message)
  }
}

const addExtraImage = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  e.target.value = ''
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await fetch(`${API_BASE}/admin/upload`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` },
      body: formData
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '上传失败')
    imageList.value.push(data.url)
  } catch (e) {
    toast.error(e.message)
  }
}

const removeExtraImage = (index) => {
  imageList.value.splice(index, 1)
}

const toggleFeatured = async (product) => {
  try {
    const res = await fetch(`${API_BASE}/admin/products/${product.id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        name: product.name,
        slug: product.slug,
        description: product.description,
        price: product.price,
        imageUrl: product.imageUrl,
        categoryId: product.categoryId,
        stock: product.stock,
        isActive: product.isActive,
        isFeatured: !product.isFeatured
      })
    })

    if (res.ok) {
      product.isFeatured = !product.isFeatured
    }
  } catch (e) {
    console.error('Error toggling featured:', e)
  }
}

const toggleActive = async (product) => {
  try {
    const res = await fetch(`${API_BASE}/admin/products/${product.id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        name: product.name,
        slug: product.slug,
        description: product.description,
        price: product.price,
        imageUrl: product.imageUrl,
        categoryId: product.categoryId,
        stock: product.stock,
        isActive: !product.isActive,
        isFeatured: product.isFeatured
      })
    })

    if (res.ok) {
      product.isActive = !product.isActive
    }
  } catch (e) {
    console.error('Error toggling active:', e)
  }
}

const formatPrice = (price) => `¥${parseFloat(price).toFixed(2)}`

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchProducts()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchProducts()
  }
}

onMounted(() => {
  fetchCategories()
  fetchProducts()
})
</script>

<template>
  <div class="admin-products">
    <div class="toolbar">
      <div class="toolbar-left">
        <button class="btn-refresh" :class="{ spinning: loading }" @click="fetchProducts" title="刷新商品列表">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
        </button>
        <div class="search-box">
          <input 
            v-model="search" 
            type="text" 
            placeholder="搜索商品名称..."
            @keyup.enter="handleSearch"
          >
          <button class="btn-search" @click="handleSearch">搜索</button>
        </div>
        <select v-model="filterCategory" @change="handleSearch">
          <option value="">所有分类</option>
          <template v-for="cat in categories" :key="cat.id">
            <optgroup :label="cat.name">
              <option :value="cat.id">{{ cat.name }} - 全部</option>
              <option v-for="child in cat.children" :key="child.id" :value="child.id">{{ child.name }}</option>
            </optgroup>
          </template>
        </select>
        <select v-model="filterFeatured" @change="handleSearch">
          <option value="">全部商品</option>
          <option value="true">推荐商品</option>
          <option value="false">非推荐</option>
        </select>
        <button class="btn-reset" @click="resetFilters">重置</button>
      </div>
      <div class="toolbar-right">
        <button class="btn-add" @click="openCreateModal">+ 添加商品</button>
        <template v-if="selectedIds.length > 0">
          <button class="btn-batch" @click="batchUpdateStatus(true)">批量上架</button>
          <button class="btn-batch" @click="batchUpdateStatus(false)">批量下架</button>
          <button class="btn-batch btn-batch-danger" @click="batchDelete">批量删除</button>
          <span class="selected-count">已选 {{ selectedIds.length }} 件</span>
        </template>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th style="width: 40px"><input type="checkbox" v-model="selectAll" @change="toggleSelectAll"></th>
            <th style="width: 60px" class="sortable-th" @click="handleSort('id')">ID{{ sortIcon('id') }}</th>
            <th style="width: 80px">图片</th>
            <th class="sortable-th" @click="handleSort('name')">商品名称{{ sortIcon('name') }}</th>
            <th style="width: 120px">分类</th>
            <th style="width: 100px" class="sortable-th" @click="handleSort('price')">价格{{ sortIcon('price') }}</th>
            <th style="width: 80px" class="sortable-th" @click="handleSort('stock')">库存{{ sortIcon('stock') }}</th>
            <th style="width: 80px">推荐</th>
            <th style="width: 80px">状态</th>
            <th style="width: 140px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="product in products" :key="product.id">
            <td><input type="checkbox" :value="product.id" v-model="selectedIds"></td>
            <td>{{ product.id }}</td>
            <td>
              <img 
                v-if="product.imageUrl" 
                :src="product.imageUrl" 
                class="product-thumb"
                @error="$event.target.style.display='none'"
              >
            </td>
            <td>
              <div class="product-name">{{ product.name }}</div>
            </td>
            <td>{{ product.category?.name || '-' }}</td>
            <td>{{ formatPrice(product.price) }}</td>
            <td>{{ product.stock }} 件</td>
            <td>
              <span 
                class="badge" 
                :class="product.isFeatured ? 'badge-warning' : 'badge-default'"
                @click="toggleFeatured(product)"
                style="cursor: pointer"
              >
                {{ product.isFeatured ? '推荐' : '-' }}
              </span>
            </td>
            <td>
              <span 
                class="badge" 
                :class="product.isActive ? 'badge-success' : 'badge-danger'"
                @click="toggleActive(product)"
                style="cursor: pointer"
              >
                {{ product.isActive ? '上架' : '下架' }}
              </span>
            </td>
            <td>
              <button class="btn-edit" @click="openEditModal(product)">编辑</button>
              <button class="btn-delete" @click="deleteProduct(product)">删除</button>
            </td>
          </tr>
          <tr v-if="products.length === 0">
            <td colspan="10" class="empty-text">暂无商品数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="totalPages > 1" class="pagination">
      <button :disabled="currentPage === 1" @click="prevPage">上一页</button>
      <span>{{ currentPage }} / {{ totalPages }}</span>
      <button :disabled="currentPage === totalPages" @click="nextPage">下一页</button>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ modalMode === 'create' ? '添加商品' : '编辑商品' }}</h3>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-section">
            <div class="section-title">商品状态</div>
            <div class="status-switches">
              <div class="switch-item">
                <span class="switch-label">上架销售</span>
                <label class="switch">
                  <input type="checkbox" v-model="form.isActive">
                  <span class="slider"></span>
                </label>
                <span class="switch-value">{{ form.isActive ? '是' : '否' }}</span>
              </div>
              <div class="switch-item">
                <span class="switch-label">首页推荐</span>
                <label class="switch">
                  <input type="checkbox" v-model="form.isFeatured">
                  <span class="slider"></span>
                </label>
                <span class="switch-value">{{ form.isFeatured ? '是' : '否' }}</span>
              </div>
            </div>
          </div>

          <div class="form-section">
            <div class="section-title">基本信息</div>
            <div class="form-group">
              <label>商品名称 <span class="required">*</span></label>
              <input v-model="form.name" type="text" placeholder="请输入商品名称">
            </div>
            <div class="form-group">
              <label>Slug（留空自动生成）</label>
              <input v-model="form.slug" type="text" placeholder="如: cohiba-esplendidos">
            </div>
            <div class="form-group">
              <label>分类</label>
              <select v-model="form.categoryId">
                <option value="">请选择分类</option>
                <template v-for="cat in categories" :key="cat.id">
                  <optgroup :label="cat.name">
                    <option v-for="child in cat.children" :key="child.id" :value="child.id">{{ child.name }}</option>
                  </optgroup>
                </template>
              </select>
            </div>
          </div>

          <div class="form-section">
            <div class="section-title">价格与库存</div>
            <div class="form-row">
              <div class="form-group">
                <label>价格 <span class="required">*</span></label>
                <div class="input-group">
                  <span class="input-prefix">¥</span>
                  <input v-model.number="form.price" type="number" step="0.01" min="0" placeholder="0.00">
                </div>
              </div>
              <div class="form-group">
                <label>库存</label>
                <div class="input-group">
                  <input v-model.number="form.stock" type="number" min="0" placeholder="0">
                  <span class="input-suffix">件</span>
                </div>
              </div>
            </div>
          </div>

          <div class="form-section">
            <div class="section-title">商品描述</div>
            <div class="form-group">
              <textarea v-model="form.description" rows="3" placeholder="请输入商品描述"></textarea>
            </div>
          </div>

          <div class="form-section">
            <div class="section-title">商品图片</div>
            <div class="form-group">
              <AdminImageUpload v-model="form.imageUrl" v-model:thumbnail="form.thumbnailUrl" :aspect-ratio="1" />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button class="btn-save" :disabled="saving" @click="saveProduct">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-products {
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
  flex-wrap: wrap;
  gap: 10px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  color: #999;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
  flex-shrink: 0;
}

.btn-refresh:hover {
  background: #f0f0f0;
  color: #333;
}

.btn-refresh.spinning svg {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.search-box {
  display: flex;
  gap: 5px;
}

.search-box input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 200px;
}

select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
}

.btn-search,
.btn-reset,
.btn-add {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-search {
  background: #1a1a1a;
  color: #fff;
}

.btn-search:hover {
  background: #333;
}

.btn-reset {
  background: #f5f5f5;
  color: #666;
}

.btn-reset:hover {
  background: #e0e0e0;
}

.btn-add {
  background: #d4a574;
  color: #fff;
}

.btn-add:hover {
  background: #c49464;
}

.loading {
  padding: 40px;
  text-align: center;
  color: #666;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f9f9f9;
  font-weight: 600;
  color: #333;
  white-space: nowrap;
}

.sortable-th {
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.sortable-th:hover {
  background: #f0f0f0;
}

.data-table td {
  color: #666;
}

.product-thumb {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 4px;
}

.product-name {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.badge-warning {
  background: #fff8e1;
  color: #f57c00;
}

.badge-default {
  background: #f5f5f5;
  color: #999;
}

.btn-edit,
.btn-delete {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  margin-right: 5px;
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
  padding: 40px !important;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  padding: 20px;
  border-top: 1px solid #eee;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: #fff;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  border-color: #d4a574;
  color: #d4a574;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
  max-width: 600px;
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
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #d4a574;
}

.form-section {
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.form-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
  padding-left: 10px;
  border-left: 3px solid #d4a574;
}

.status-switches {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}

.switch-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.switch-label {
  font-size: 14px;
  color: #666;
  min-width: 70px;
}

.switch-value {
  font-size: 13px;
  color: #999;
  min-width: 24px;
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

.input-group {
  display: flex;
  align-items: stretch;
}

.input-group input {
  flex: 1;
  border-radius: 4px 0 0 4px;
}

.input-prefix,
.input-suffix {
  display: flex;
  align-items: center;
  padding: 0 12px;
  background: #f5f5f5;
  border: 1px solid #ddd;
  color: #666;
  font-size: 14px;
}

.input-prefix {
  border-right: none;
  border-radius: 4px 0 0 4px;
}

.input-suffix {
  border-left: none;
  border-radius: 0 4px 4px 0;
}

.input-group .input-prefix + input {
  border-radius: 0 4px 4px 0;
}

.input-group input:not(:last-child) {
  border-radius: 4px 0 0 4px;
}

.form-row {
  display: flex;
  gap: 15px;
}

.form-row .form-group {
  flex: 1;
}

.checkbox-row {
  display: flex;
  gap: 20px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
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
