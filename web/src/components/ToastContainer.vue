<script setup lang="ts">
import { useUiStore } from '@/store/ui'

const ui = useUiStore()
</script>

<template>
  <TransitionGroup name="toast">
    <div
      v-for="t in ui.toasts" :key="t.id"
      class="fixed top-20 right-6 z-[200] px-5 py-3 rounded-xl text-sm font-medium flex items-center gap-2"
      :class="t.type === 'success' ? 'border-[rgba(16,185,129,.3)] bg-[rgba(16,185,129,.1)]' : ''"
      style="background:var(--bg-elevated);border:1px solid var(--border);color:var(--text-primary);box-shadow:0 4px 24px rgba(0,0,0,.3)"
    >
      {{ t.message }}
    </div>
  </TransitionGroup>
</template>

<style scoped>
.toast-enter-active { transition: all 0.3s ease; }
.toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from { transform: translateX(120%); opacity: 0; }
.toast-leave-to { transform: translateX(120%); opacity: 0; }
</style>
