import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const toasts = ref<{ id: number; message: string; type: string }[]>([])
  let nextId = 0

  function showToast(message: string, type = '') {
    const id = nextId++
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, 2500)
  }

  return { toasts, showToast }
})
