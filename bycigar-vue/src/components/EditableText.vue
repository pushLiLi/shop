<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const props = defineProps({
  configKey: { type: String, required: true },
  modelValue: { type: String, required: true },
  tag: { type: String, default: 'span' },
  class: { type: String, default: '' }
})

const emit = defineEmits(['update'])

const authStore = useAuthStore()
const isEditing = ref(false)
const editValue = ref(props.modelValue)

const canEdit = computed(() => authStore.isAdmin)

async function saveEdit() {
  if (editValue.value !== props.modelValue) {
    try {
      const res = await fetch(`http://localhost:3000/api/admin/config/${props.configKey}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${authStore.token}`
        },
        body: JSON.stringify({ value: editValue.value })
      })
      if (res.ok) {
        emit('update', editValue.value)
      }
    } catch (e) {
      console.error('Save failed:', e)
    }
  }
  isEditing.value = false
}

function startEdit() {
  editValue.value = props.modelValue
  isEditing.value = true
}
</script>

<template>
  <span v-if="!canEdit || !isEditing" 
        :class="props.class" 
        @click="canEdit && startEdit()">
    {{ modelValue }}
  </span>
  <input v-else
         v-model="editValue"
         :class="props.class"
         @blur="saveEdit"
         @keyup.enter="saveEdit"
         autofocus
  />
</template>
