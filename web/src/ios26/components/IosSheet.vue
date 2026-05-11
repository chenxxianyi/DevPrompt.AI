<script setup lang="ts">
import { watch } from 'vue'

const props = defineProps<{
  modelValue: boolean
  title?: string
  /** 关闭按钮文本 */
  closeText?: string
}>()

const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()

function close() {
  emit('update:modelValue', false)
}

// 防止滚动穿透
watch(() => props.modelValue, (open: boolean) => {
  if (typeof document === 'undefined') return
  if (open) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition name="ios-sheet-fade">
      <div
        v-if="modelValue"
        class="ios-sheet-backdrop"
        @click.self="close"
      >
        <Transition name="ios-sheet" appear>
          <div v-if="modelValue" class="ios-sheet" role="dialog" aria-modal="true">
            <div class="ios-sheet__grabber" />
            <div v-if="title || $slots.header" class="ios-sheet__header">
              <slot name="header">
                <h2 class="ios-text-headline ios-sheet__title">{{ title }}</h2>
              </slot>
              <button
                class="ios-sheet__close ios-press"
                type="button"
                :aria-label="closeText || '关闭'"
                @click="close"
              >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                  <path d="M18 6 6 18M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div class="ios-sheet__body">
              <slot />
            </div>
            <div v-if="$slots.footer" class="ios-sheet__footer">
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.ios-sheet-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  z-index: 200;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: env(safe-area-inset-bottom, 0);
}

.ios-sheet {
  width: 100%;
  max-width: 720px;
  max-height: 90vh;
  background: var(--ios-color-bg-grouped-secondary);
  border-top-left-radius: var(--ios-radius-sheet);
  border-top-right-radius: var(--ios-radius-sheet);
  border: 1px solid var(--ios-color-separator);
  border-bottom: none;
  box-shadow: var(--ios-glass-shadow-background);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding-top: 8px;
}

/* Desktop: center it */
@media (min-width: 768px) {
  .ios-sheet-backdrop {
    align-items: center;
  }
  .ios-sheet {
    border-radius: var(--ios-radius-sheet);
    border-bottom: 1px solid var(--ios-color-separator);
    max-height: 80vh;
  }
}

.ios-sheet__grabber {
  width: 36px;
  height: 5px;
  border-radius: 9999px;
  background: var(--ios-color-label-tertiary);
  margin: 0 auto;
}

.ios-sheet__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px 12px;
}

.ios-sheet__title {
  margin: 0;
  color: var(--ios-color-label-primary);
}

.ios-sheet__close {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: var(--ios-color-fill-tertiary);
  color: var(--ios-color-label-secondary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
}

.ios-sheet__body {
  padding: 0 20px 20px;
  overflow-y: auto;
  flex: 1;
}

.ios-sheet__footer {
  padding: 12px 20px calc(env(safe-area-inset-bottom, 0) + 16px);
  border-top: 0.5px solid var(--ios-color-separator);
  background: var(--ios-color-bg-grouped-secondary);
}

/* Transitions */
.ios-sheet-fade-enter-active,
.ios-sheet-fade-leave-active {
  transition: opacity 0.3s var(--ios-ease-gentle);
}
.ios-sheet-fade-enter-from,
.ios-sheet-fade-leave-to {
  opacity: 0;
}

.ios-sheet-enter-active {
  transition: transform 0.5s var(--ios-ease-gentle);
}
.ios-sheet-leave-active {
  transition: transform 0.3s var(--ios-ease-gentle);
}
.ios-sheet-enter-from,
.ios-sheet-leave-to {
  transform: translateY(100%);
}
</style>
