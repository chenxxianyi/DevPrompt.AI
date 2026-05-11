<script setup lang="ts">
/**
 * IosListRow — iOS 风格列表行
 * size: sm(36) | md(44 default) | lg(58)
 * 提供 leading / default / value / trailing 4 个具名插槽
 * disclosure: 显示 chevron-right
 */
withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg'
  disclosure?: boolean
  /** 点击 -> 跳转 */
  to?: string
  /** 强调标题 */
  emph?: boolean
  /** 红色破坏样式（destructive） */
  destructive?: boolean
}>(), {
  size: 'md',
  disclosure: false,
  emph: false,
  destructive: false,
})

defineEmits<{ (e: 'click', ev: MouseEvent): void }>()
</script>

<template>
  <component
    :is="to ? 'router-link' : 'div'"
    :to="to"
    class="ios-list-row ios-press-soft"
    :class="[
      `ios-list-row--${size}`,
      to ? 'ios-list-row--interactive' : '',
      destructive ? 'ios-list-row--destructive' : '',
    ]"
    @click="(ev: MouseEvent) => $emit('click', ev)"
  >
    <span v-if="$slots.leading" class="ios-list-row__leading">
      <slot name="leading" />
    </span>
    <span class="ios-list-row__main">
      <span class="ios-list-row__title" :class="{ 'ios-text-body-emph': emph }">
        <slot />
      </span>
      <span v-if="$slots.subtitle" class="ios-list-row__subtitle">
        <slot name="subtitle" />
      </span>
    </span>
    <span v-if="$slots.value" class="ios-list-row__value">
      <slot name="value" />
    </span>
    <span v-if="$slots.trailing" class="ios-list-row__trailing">
      <slot name="trailing" />
    </span>
    <svg
      v-if="disclosure"
      class="ios-list-row__chevron"
      width="9"
      height="14"
      viewBox="0 0 9 14"
      fill="none"
      aria-hidden="true"
    >
      <path
        d="M1.5 1l6 6-6 6"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </component>
</template>

<style scoped>
.ios-list-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background: var(--ios-color-bg-grouped-secondary);
  color: var(--ios-color-label-primary);
  text-decoration: none;
  cursor: default;
  position: relative;
}

.ios-list-row--sm {
  min-height: 36px;
}
.ios-list-row--md {
  min-height: 44px;
}
.ios-list-row--lg {
  min-height: 58px;
  padding: 12px 16px;
}

.ios-list-row--interactive {
  cursor: pointer;
}
.ios-list-row--interactive:hover {
  background: var(--ios-color-fill-quaternary);
}

.ios-list-row--destructive .ios-list-row__title {
  color: var(--ios-color-systemRed);
}

.ios-list-row__leading {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  color: var(--ios-color-tint);
}

.ios-list-row__main {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.ios-list-row__title {
  font-size: 17px;
  line-height: 22px;
  letter-spacing: -0.43px;
  color: var(--ios-color-label-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ios-list-row__subtitle {
  font-size: 13px;
  line-height: 18px;
  letter-spacing: -0.08px;
  color: var(--ios-color-label-secondary);
  margin-top: 2px;
}

.ios-list-row__value {
  color: var(--ios-color-label-secondary);
  font-size: 17px;
  letter-spacing: -0.43px;
  white-space: nowrap;
}

.ios-list-row__trailing {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.ios-list-row__chevron {
  color: var(--ios-color-label-tertiary);
}
</style>
