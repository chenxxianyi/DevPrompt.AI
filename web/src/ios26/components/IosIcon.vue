<script setup lang="ts">
defineProps<{
  /** SVG path d 属性内容 */
  path: string
  /** 尺寸（px），默认 20 */
  size?: number | string
  /** 线宽，默认 1.6 */
  stroke?: number | string
  /** 是否填充：true 用 fill，false 用 stroke（默认 stroke） */
  filled?: boolean
}>()
</script>

<template>
  <svg
    :width="size || 20"
    :height="size || 20"
    viewBox="0 0 24 24"
    :fill="filled ? 'currentColor' : 'none'"
    :stroke="filled ? 'none' : 'currentColor'"
    :stroke-width="stroke || 1.7"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <path :d="path" />
  </svg>
</template>
