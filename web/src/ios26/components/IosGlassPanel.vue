<script setup lang="ts">
/**
 * IosGlassPanel — Liquid Glass 基础容器
 * size: 'sm' | 'md' | 'lg'
 * inset: 是否自带内边距（默认 true）
 */
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg'
  inset?: boolean
  /** 渲染元素：默认 div */
  as?: string
  /** 是否在 hover 时轻微抬升 */
  interactive?: boolean
}>(), {
  size: 'md',
  inset: true,
  as: 'div',
  interactive: false,
})

const klass = computed(() => {
  return [
    'ios-glass',
    props.size === 'sm' ? 'ios-glass-sm' : props.size === 'lg' ? 'ios-glass-lg' : '',
    props.inset ? 'ios-glass--inset' : '',
    props.interactive ? 'ios-glass--interactive' : '',
  ].filter(Boolean).join(' ')
})
</script>

<template>
  <component :is="as" :class="klass">
    <slot />
  </component>
</template>

<style scoped>
.ios-glass--inset {
  padding: var(--ios-space-5);
}

.ios-glass--interactive {
  transition: transform 0.25s var(--ios-ease-gentle),
    box-shadow 0.25s var(--ios-ease-gentle);
  cursor: pointer;
}

.ios-glass--interactive:hover {
  transform: translateY(-2px);
}
</style>
