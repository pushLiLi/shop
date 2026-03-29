<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { marked } from 'marked'
import { useToastStore } from '../../stores/toast'

const API_BASE = '/api'
const toast = useToastStore()

const pages = ref([])
const loading = ref(false)
const saving = ref(false)
const selectedSlug = ref('')

const form = ref({
  title: '',
  content: ''
})

const authHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token')}`
})

const currentPage = computed(() => {
  return pages.value.find(p => p.slug === selectedSlug.value)
})

const renderedPreview = computed(() => {
  if (!form.value.content) return ''
  return marked(form.value.content)
})

const fetchPages = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/pages`, {
      headers: authHeaders()
    })
    pages.value = await res.json()
    if (pages.value.length > 0 && !selectedSlug.value) {
      selectPage(pages.value[0].slug)
    }
  } catch (e) {
    toast.error('获取页面列表失败')
  } finally {
    loading.value = false
  }
}

const selectPage = (slug) => {
  selectedSlug.value = slug
  const page = pages.value.find(p => p.slug === slug)
  if (page) {
    form.value = {
      title: page.title,
      content: page.content
    }
  }
}

const savePage = async () => {
  if (!form.value.title) {
    toast.error('请填写标题')
    return
  }

  saving.value = true
  try {
    const res = await fetch(`${API_BASE}/admin/pages/${selectedSlug.value}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(form.value)
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || '保存失败')
    }

    const index = pages.value.findIndex(p => p.slug === selectedSlug.value)
    if (index !== -1) {
      pages.value[index] = data
    }

    toast.success('保存成功')
  } catch (e) {
    toast.error(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(fetchPages)
</script>

<template>
  <div class="admin-pages">
    <div class="page-list">
      <h3>页面列表
        <button class="btn-refresh" :class="{ spinning: loading }" @click="fetchPages" title="刷新页面列表">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
        </button>
      </h3>
      <div class="list">
        <button
          v-for="page in pages"
          :key="page.slug"
          :class="['page-item', { active: selectedSlug === page.slug }]"
          @click="selectPage(page.slug)"
        >
          <span class="page-title-text">{{ page.title }}</span>
          <span class="page-slug">{{ page.slug }}</span>
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="selectedSlug" class="editor-container">
      <div class="editor-panels">
        <div class="editor-panel">
          <div class="panel-header">
            <span>Markdown 编辑</span>
          </div>
          <div class="form-group">
            <label>标题</label>
            <input v-model="form.title" type="text" placeholder="页面标题" />
          </div>
          <div class="form-group content-group">
            <label>内容</label>
            <textarea v-model="form.content" placeholder="使用 Markdown 格式编写内容..."></textarea>
          </div>
        </div>

        <div class="preview-panel">
          <div class="panel-header">
            <span>实时预览</span>
          </div>
          <div class="preview-content">
            <h1 class="preview-title">{{ form.title || '页面标题' }}</h1>
            <div class="markdown-body" v-html="renderedPreview"></div>
          </div>
        </div>
      </div>

      <div class="editor-footer">
        <button class="btn-save" :disabled="saving" @click="savePage">
          {{ saving ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-pages {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  min-height: calc(100vh - 200px);
}

.page-list {
  width: 200px;
  border-right: 1px solid #eee;
  flex-shrink: 0;
}

.page-list h3 {
  padding: 15px;
  margin: 0;
  font-size: 14px;
  color: #333;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
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

.list {
  padding: 10px;
}

.page-item {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 12px;
  border: none;
  background: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 4px;
}

.page-item:hover {
  background: #f5f5f5;
}

.page-item.active {
  background: #fef6ee;
}

.page-title-text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.page-item.active .page-title-text {
  color: #d4a574;
}

.page-slug {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.editor-panels {
  flex: 1;
  display: flex;
  min-height: 0;
}

.editor-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #eee;
  min-width: 0;
}

.preview-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.panel-header {
  padding: 12px 15px;
  background: #fafafa;
  border-bottom: 1px solid #eee;
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

.form-group {
  padding: 15px;
  border-bottom: 1px solid #eee;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  color: #333;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input:focus {
  outline: none;
  border-color: #d4a574;
}

.content-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 15px;
  border-bottom: none;
}

.content-group textarea {
  flex: 1;
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  line-height: 1.6;
  resize: none;
  min-height: 300px;
}

.content-group textarea:focus {
  outline: none;
  border-color: #d4a574;
}

.preview-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: #fafafa;
}

.preview-title {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.markdown-body {
  font-size: 14px;
  line-height: 1.8;
  color: #444;
}

.markdown-body :deep(h1) {
  font-size: 22px;
  margin: 20px 0 15px;
  color: #333;
}

.markdown-body :deep(h2) {
  font-size: 18px;
  margin: 18px 0 12px;
  color: #333;
}

.markdown-body :deep(h3) {
  font-size: 16px;
  margin: 15px 0 10px;
  color: #444;
}

.markdown-body :deep(p) {
  margin-bottom: 12px;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 20px;
  margin-bottom: 12px;
}

.markdown-body :deep(li) {
  margin-bottom: 6px;
}

.markdown-body :deep(a) {
  color: #d4a574;
}

.markdown-body :deep(blockquote) {
  border-left: 3px solid #d4a574;
  padding-left: 12px;
  margin: 12px 0;
  color: #666;
}

.markdown-body :deep(code) {
  background: #f0f0f0;
  padding: 2px 5px;
  border-radius: 3px;
  font-size: 13px;
}

.markdown-body :deep(pre) {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 12px 0;
}

.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
}

.editor-footer {
  padding: 15px 20px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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

@media (max-width: 900px) {
  .admin-pages {
    flex-direction: column;
  }

  .page-list {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #eee;
  }

  .list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding: 15px;
  }

  .page-item {
    width: auto;
    flex-direction: row;
    gap: 8px;
    padding: 8px 12px;
  }

  .page-slug {
    margin-top: 0;
  }

  .editor-panels {
    flex-direction: column;
  }

  .editor-panel {
    border-right: none;
    border-bottom: 1px solid #eee;
    min-height: 350px;
  }

  .preview-panel {
    min-height: 300px;
  }
}
</style>
