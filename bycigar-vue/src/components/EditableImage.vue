<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const props = defineProps({
  src: { type: String, required: true },
  alt: { type: String, default: '' },
  class: { type: String, default: '' }
})

const emit = defineEmits(['update'])

const authStore = useAuthStore()
const isHovering = ref(false)

const canEdit = computed(() => authStore.isAdmin)

async function handleUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await fetch('/api/admin/upload', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: formData
    })
    if (res.ok) {
      const data = await res.json()
      emit('update', data.url)
    }
  } catch (e) {
    console.error('Upload failed:', e)
  }
}
</script>

<template>
  <div 
    class="editable-image-wrapper"
    :class="props.class"
    @mouseenter="isHovering = true"
    @mouseleave="isHovering=false"
  >
    <img :src="src" :alt="alt" class="editable-image">
    <div v-if="canEdit && isHovering" class="image-overlay">
      <label class="upload-btn">
        更换图片
        <input type="file" accept="image/*" @change="handleUpload" hidden>
      </label>
    </div>
  </div>
</template>

<style scoped>
.editable-image-wrapper {
  position: relative;
  display: inline-block;
}

.editable-image {
  display: block;
  max-width: 100%;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-btn {
  background: #d4a574;
  color: #1a1a1a;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.upload-btn:hover {
  background: #e5b484;
}
</style>
