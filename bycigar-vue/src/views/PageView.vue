<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'

const API_BASE = 'http://localhost:3000/api'

const route = useRoute()
const page = ref(null)
const loading = ref(true)
const error = ref(false)

const fetchPage = async () => {
  loading.value = true
  error.value = false
  
  try {
    const res = await fetch(`${API_BASE}/pages/${route.params.slug}`)
    if (!res.ok) {
      error.value = true
      return
    }
    page.value = await res.json()
  } catch (e) {
    error.value = true
  } finally {
    loading.value = false
  }
}

const renderedContent = () => {
  if (!page.value?.content) return ''
  return marked(page.value.content)
}

watch(() => route.params.slug, fetchPage)
onMounted(fetchPage)
</script>

<template>
  <div class="page-view">
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      
      <div v-else-if="error" class="error">
        <h1>页面不存在</h1>
        <p>您访问的页面不存在或已被删除。</p>
        <router-link to="/" class="back-link">返回首页</router-link>
      </div>
      
      <article v-else class="page-content">
        <h1 class="page-title">{{ page?.title }}</h1>
        <div class="markdown-body" v-html="renderedContent()"></div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.page-view {
  min-height: 60vh;
  padding: 60px 20px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
}

.container {
  max-width: 800px;
  margin: 0 auto;
}

.loading,
.error {
  text-align: center;
  padding: 60px 20px;
}

.loading {
  color: #888;
  font-size: 18px;
}

.error h1 {
  color: #d4a574;
  font-size: 28px;
  margin-bottom: 16px;
}

.error p {
  color: #888;
  margin-bottom: 24px;
}

.back-link {
  color: #d4a574;
  text-decoration: none;
  padding: 10px 24px;
  border: 1px solid #d4a574;
  border-radius: 4px;
  transition: all 0.3s;
}

.back-link:hover {
  background: #d4a574;
  color: #1a1a1a;
}

.page-content {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  padding: 40px;
}

.page-title {
  color: #d4a574;
  font-size: 32px;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(212, 165, 116, 0.3);
}

.markdown-body {
  color: #ccc;
  line-height: 1.8;
  font-size: 16px;
}

.markdown-body :deep(h1) {
  color: #d4a574;
  font-size: 28px;
  margin: 30px 0 20px;
}

.markdown-body :deep(h2) {
  color: #d4a574;
  font-size: 24px;
  margin: 25px 0 15px;
}

.markdown-body :deep(h3) {
  color: #e0c4a8;
  font-size: 20px;
  margin: 20px 0 12px;
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 24px;
  margin-bottom: 16px;
}

.markdown-body :deep(li) {
  margin-bottom: 8px;
}

.markdown-body :deep(a) {
  color: #d4a574;
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.markdown-body :deep(blockquote) {
  border-left: 4px solid #d4a574;
  padding-left: 16px;
  margin: 20px 0;
  color: #999;
}

.markdown-body :deep(code) {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 14px;
}

.markdown-body :deep(pre) {
  background: rgba(0, 0, 0, 0.3);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 20px 0;
}

.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
}

.markdown-body :deep(th),
.markdown-body :deep(td) {
  border: 1px solid #444;
  padding: 12px;
  text-align: left;
}

.markdown-body :deep(th) {
  background: rgba(0, 0, 0, 0.2);
  color: #d4a574;
}

.markdown-body :deep(hr) {
  border: none;
  border-top: 1px solid #444;
  margin: 30px 0;
}

@media (max-width: 768px) {
  .page-view {
    padding: 40px 15px;
  }

  .page-content {
    padding: 24px;
  }

  .page-title {
    font-size: 24px;
  }
}
</style>
