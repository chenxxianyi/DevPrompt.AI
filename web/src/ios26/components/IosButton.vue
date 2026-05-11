<script setup lang="ts">
/**
 * IosButton — iOS 26 风格按钮
 * variant:
 *   filled   实心填充（主操作）
 *   tinted   语义淡色填充
 *   plain    无背景仅文字
 *   glass    Liquid Glass 玻璃按钮
 *   destructive 红色破坏性操作
 * size: sm(36) | md(44 default) | lg(52)
 */
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'filled' | 'tinted' | 'plain' | 'glass' | 'destructive'
  size?: 'sm' | 'md' | 'lg'
  block?: boolean
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
  to?: string
}>(), {
  variant: 'filled',
  size: 'md',
  block: false,
  loading: false,
  disabled: false,
  type: 'button',
})

const emit = defineEmits<{ (e: 'click', ev: MouseEvent): void }>()

const klass = computed(() => [
  'ios-btn',
  `ios-btn--${props.variant}`,
  `ios-btn--${props.size}`,
  props.block ? 'ios-btn--block' : '',
  props.disabled || props.loading ? 'is-disabled' : '',
  'ios-press',
].filter(Boolean).join(' '))

function onClick(ev: MouseEvent) {
  if (props.disabled || props.loading) {
    ev.preventDefault()
    ev.stopPropagation()
    return
  }
  emit('click', ev)
}
</script>

<template>
  <router-link
    v-if="to"
    :to="to"
    :class="klass"
    role="button"
    @click="onClick"
  >
    <span v-if="loading" class="ios-spinner mr-2" />
    <slot />
  </router-link>
  <button
    v-else
    :type="type"
    :class="klass"
    :disabled="disabled || loading"
    @click="onClick"
  >
    <span v-if="loading" class="ios-spinner mr-2" />
    <slot />
  </button>
</template>

<style scoped>
.ios-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-weight: 600;
  letter-spacing: -0.2px;
  border-radius: var(--ios-radius-button);
  border: 1px solid transparent;
  cursor: pointer;
  outline: none;
  text-decoration: none;
  white-space: nowrap;
  min-height: var(--ios-touch-min);
  padding: 0 20px;
  font-size: 17px;
  line-height: 1;
}

.ios-btn--sm {
  min-height: 36px;
  padding: 0 14px;
  font-size: 15px;
  border-radius: var(--ios-radius-md);
}

.ios-btn--md {
  min-height: 44px;
  padding: 0 20px;
  font-size: 17px;
}

.ios-btn--lg {
  min-height: 52px;
  padding: 0 24px;
  font-size: 17px;
  font-weight: 600;
  border-radius: var(--ios-radius-lg);
}

.ios-btn--block {
  width: 100%;
}

/* Filled — macaron gradient */
.ios-btn--filled {
  background: var(--ios-color-tint-gradient, var(--ios-color-tint));
  background-size: 200% 200%;
  background-position: 0% 50%;
  color: #fff;
  box-shadow: 0 6px 20px rgba(63, 164, 232, 0.35), 0 2px 6px rgba(167, 139, 250, 0.18);
  transition: background-position 0.4s var(--ios-ease-gentle),
    box-shadow 0.2s var(--ios-ease-stiff),
    transform 0.18s var(--ios-ease-snappy);
}
.ios-btn--filled:hover:not(.is-disabled) {
  background-position: 100% 50%;
  box-shadow: 0 8px 28px rgba(63, 164, 232, 0.45), 0 3px 10px rgba(167, 139, 250, 0.24);
}

/* Tinted */
.ios-btn--tinted {
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
}
.ios-btn--tinted:hover:not(.is-disabled) {
  background: var(--ios-color-tint-soft-strong);
}

/* Plain */
.ios-btn--plain {
  background: transparent;
  color: var(--ios-color-tint);
}
.ios-btn--plain:hover:not(.is-disabled) {
  background: var(--ios-color-fill-quaternary);
}

/* Glass */
.ios-btn--glass {
  background: var(--ios-glass-bg-light-md);
  backdrop-filter: blur(var(--ios-glass-blur-md)) saturate(var(--ios-glass-saturate));
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-md)) saturate(var(--ios-glass-saturate));
  border-color: var(--ios-glass-border-light);
  color: var(--ios-color-label-primary);
  box-shadow: var(--ios-glass-shadow-press);
}
.ios-btn--glass:hover:not(.is-disabled) {
  background: var(--ios-glass-bg-light-lg);
}

/* Destructive */
.ios-btn--destructive {
  background: var(--ios-color-systemRed);
  color: #fff;
}
.ios-btn--destructive:hover:not(.is-disabled) {
  filter: brightness(1.05);
}

.is-disabled,
.ios-btn[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
  box-shadow: none;
}

.mr-2 {
  margin-right: 6px;
}
</style>
