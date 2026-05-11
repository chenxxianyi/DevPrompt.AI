<script setup lang="ts">
/**
 * IosSegmentedControl — iOS 风格分段控件
 * v-model 绑定当前 value
 */
defineProps<{
  modelValue: string
  options: Array<{ value: string; label: string }>
  /** 紧凑模式 */
  size?: 'sm' | 'md'
}>()

defineEmits<{ (e: 'update:modelValue', v: string): void }>()
</script>

<template>
  <div class="ios-segmented" :class="`ios-segmented--${size || 'md'}`">
    <button
      v-for="opt in options"
      :key="opt.value"
      type="button"
      class="ios-segmented__item ios-press"
      :class="{ 'is-active': modelValue === opt.value }"
      @click="$emit('update:modelValue', opt.value)"
    >
      <span class="ios-segmented__label">{{ opt.label }}</span>
    </button>
  </div>
</template>

<style scoped>
.ios-segmented {
  display: inline-flex;
  padding: 2px;
  background: var(--ios-color-fill-tertiary);
  border-radius: var(--ios-radius-md);
  width: 100%;
  max-width: 480px;
  gap: 0;
}

.ios-segmented__item {
  flex: 1;
  border: none;
  background: transparent;
  color: var(--ios-color-label-primary);
  font-size: 13px;
  font-weight: 500;
  letter-spacing: -0.08px;
  padding: 8px 10px;
  border-radius: 7px;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s var(--ios-ease-gentle),
    box-shadow 0.2s var(--ios-ease-gentle),
    color 0.2s var(--ios-ease-gentle);
  min-height: 32px;
}

.ios-segmented--md .ios-segmented__item {
  font-size: 14px;
  padding: 10px 12px;
  min-height: 36px;
}

.ios-segmented__item.is-active {
  background: var(--ios-color-bg-grouped-secondary);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.06);
  color: var(--ios-color-label-primary);
  font-weight: 600;
}

@media (prefers-color-scheme: dark) {
  .ios-segmented__item.is-active {
    background: var(--ios-color-bg-tertiary);
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.45), 0 1px 2px rgba(0, 0, 0, 0.3);
  }
}

.ios26-app[data-theme='dark'] .ios-segmented__item.is-active {
  background: var(--ios-color-bg-tertiary);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.45), 0 1px 2px rgba(0, 0, 0, 0.3);
}
</style>
