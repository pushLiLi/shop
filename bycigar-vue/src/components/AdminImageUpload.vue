<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const API_BASE = 'http://localhost:3000/api'
const dragOver = ref(false)
const uploading = ref(false)
const error = ref('')

const imageUrl = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const handleDragOver = (e) => {
  e.preventDefault()
  dragOver.value = true
}

const handleDragLeave = () => {
  dragOver.value = false
}

const handleDrop = (e) => {
  e.preventDefault()
  dragOver.value = false
  const files = e.dataTransfer.files
  if (files.length > 0) {
    uploadFile(files[0])
  }
}

const handleFileSelect = (e) => {
  const files = e.target.files
  if (files.length > 0) {
    uploadFile(files[0])
  }
}

const uploadFile = async (file) => {
  if (!file.type.startsWith('image/')) {
    error.value = '请选择图片文件'
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    error.value = '图片大小不能超过 10MB'
    return
  }

  error.value = ''
  uploading.value = true

  try {
    const formData = new FormData()
    formData.append('file', file)

    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/admin/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    })

    const data = await res.json()

    if (!res.ok) {
      throw new Error(data.error || '上传失败')
    }

    imageUrl.value = data.url
  } catch (e) {
    error.value = e.message
  } finally {
    uploading.value = false
  }
}

const removeImage = () => {
  imageUrl.value = ''
}
</script>

<template>
  <div class="image-upload">
    <div 
      v-if="!imageUrl"
      class="upload-area"
      :class="{ 'drag-over': dragOver, uploading }"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
      @click="$refs.fileInput.click()"
    >
      <input 
        ref="fileInput"
        type="file" 
        accept="image/*" 
        class="file-input"
        @change="handleFileSelect"
      >
      <div v-if="uploading" class="uploading-text">
        <span class="spinner"></span>
        上传中...
      </div>
      <div v-else class="upload-hint">
        <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="17 8 12 3 7 8"></polyline>
          <line x1="12" y1="3" x2="12" y2="15"></line>
        </svg>
        <p>拖拽图片到此处或点击上传</p>
        <span class="hint">支持 JPG、PNG、GIF、WebP，最大 10MB</span>
      </div>
    </div>

    <div v-else class="preview-area">
      <div class="preview-image-wrap">
        <img :src="imageUrl" alt="Preview" class="preview-image">
      </div>
      <div class="preview-actions">
        <button class="btn-replace" @click="$refs.fileInput.click()">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="17 8 12 3 7 8"></polyline>
            <line x1="12" y1="3" x2="12" y2="15"></line>
          </svg>
          替换图片
        </button>
        <button class="btn-remove" @click="removeImage">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
          删除图片
        </button>
      </div>
      <input 
        ref="fileInput"
        type="file" 
        accept="image/*" 
        class="file-input"
        @change="handleFileSelect"
      >
    </div>

    <p v-if="error" class="error-text">{{ error }}</p>
  </div>
</template>

<style scoped>
.image-upload {
  width: 100%;
}

.upload-area {
  border: 2px dashed #ddd;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: #fafafa;
}

.upload-area:hover {
  border-color: #d4a574;
  background: #fff;
}

.upload-area.drag-over {
  border-color: #d4a574;
  background: #fff8f0;
}

.upload-area.uploading {
  pointer-events: none;
  opacity: 0.7;
}

.file-input {
  display: none;
}

.upload-hint {
  color: #666;
}

.upload-hint svg {
  color: #ccc;
  margin-bottom: 10px;
}

.upload-hint p {
  margin: 0 0 5px;
  font-size: 14px;
}

.upload-hint .hint {
  font-size: 12px;
  color: #999;
}

.uploading-text {
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #ddd;
  border-top-color: #d4a574;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.preview-area {
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
  border: 1px solid #eee;
}

.preview-image-wrap {
  min-height: 150px;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-image {
  width: 100%;
  max-height: 300px;
  object-fit: contain;
  display: block;
}

.preview-actions {
  display: flex;
  gap: 8px;
  padding: 10px;
  border-top: 1px solid #eee;
  background: #fafafa;
}

.btn-replace,
.btn-remove {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-replace {
  background: #fff;
  color: #333;
  border: 1px solid #ddd;
}

.btn-replace:hover {
  background: #f5f5f5;
  border-color: #d4a574;
}

.btn-remove {
  background: #fff;
  color: #dc3545;
  border: 1px solid #ddd;
}

.btn-remove:hover {
  background: #fff5f5;
  border-color: #dc3545;
}

.error-text {
  color: #dc3545;
  font-size: 12px;
  margin: 8px 0 0;
}
</style>
