import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const message = ref('')
  const type = ref('success')
  const visible = ref(false)
  let timer = null

  function show(msg, t = 'success') {
    message.value = msg
    type.value = t
    visible.value = true

    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      visible.value = false
    }, 2000)
  }

  function success(msg) {
    show(msg, 'success')
  }

  function error(msg) {
    show(msg, 'error')
  }

  function hide() {
    visible.value = false
    if (timer) clearTimeout(timer)
  }

  return { message, type, visible, show, success, error, hide }
})
