<script setup lang="ts">
import { useUiStore } from '@/store/ui'

const ui = useUiStore()
</script>

<template>
  <Teleport to="body">
    <div class="ios-toast-host ios-no-select">
      <TransitionGroup name="ios-toast">
        <div
          v-for="t in ui.toasts"
          :key="t.id"
          class="ios-toast"
          :class="t.type === 'success' ? 'ios-toast--success' : ''"
        >
          <span v-if="t.type === 'success'" class="ios-toast__dot" />
          {{ t.message }}
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.ios-toast-host {
  position: fixed;
  top: calc(env(safe-area-inset-top, 0) + 16px);
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  pointer-events: none;
}

.ios-toast {
  pointer-events: auto;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-glass-bg-light-lg);
  backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(var(--ios-glass-saturate));
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(var(--ios-glass-saturate));
  border: 1px solid var(--ios-glass-border-light);
  color: var(--ios-color-label-primary);
  font-size: 15px;
  font-weight: 500;
  letter-spacing: -0.2px;
  box-shadow: var(--ios-glass-shadow-layer);
  max-width: calc(100vw - 32px);
}

.ios-toast--success {
  color: var(--ios-color-systemGreen);
}

.ios-toast__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--ios-color-systemGreen);
}

.ios-toast-enter-active,
.ios-toast-leave-active {
  transition: transform 0.3s var(--ios-ease-gentle), opacity 0.3s var(--ios-ease-gentle);
}
.ios-toast-enter-from {
  transform: translateY(-20px) scale(0.96);
  opacity: 0;
}
.ios-toast-leave-to {
  transform: translateY(-12px) scale(0.96);
  opacity: 0;
}
</style>
