<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import IosIcon from './IosIcon.vue'
import { IOS_TITLES, IOS_ROOT } from '../constants'

const props = withDefaults(defineProps<{
  title?: string
  large?: boolean
  showBack?: boolean
  backTo?: string
}>(), {
  large: false,
})

const route = useRoute()
const router = useRouter()

const computedTitle = computed(() => {
  if (props.title) return props.title
  const path = route.path.replace(/\/$/, '') || '/ios26'
  if (IOS_TITLES[path]) return IOS_TITLES[path]
  const matched = Object.keys(IOS_TITLES)
    .filter(k => path.startsWith(k))
    .sort((a, b) => b.length - a.length)[0]
  return matched ? IOS_TITLES[matched] : 'DevPrompt'
})

const canGoBack = computed(() => {
  if (props.showBack !== undefined) return props.showBack
  const rest = route.path.replace(IOS_ROOT, '').replace(/^\//, '')
  return rest.length > 0 && !['', 'prompts', 'generator', 'dashboard'].includes(rest)
})

function goBack() {
  if (props.backTo) router.push(props.backTo)
  else if (window.history.length > 1) router.back()
  else router.push(IOS_ROOT)
}
</script>

<template>
  <header class="ios-navbar ios-no-select">
    <div class="ios-navbar__inner">
      <!-- Left: back button + title -->
      <div class="ios-navbar__left">
        <button
          v-if="canGoBack"
          class="ios-navbar__back"
          type="button"
          aria-label="返回"
          @click="goBack"
        >
          <IosIcon path="M15 18l-6-6 6-6" :size="18" :stroke="2" />
        </button>
        <div class="ios-navbar__title-group">
          <h1 class="ios-navbar__title" :class="{ 'ios-navbar__title--large': large }">
            {{ computedTitle }}
          </h1>
          <slot name="subtitle" />
        </div>
        <slot name="leading" />
      </div>

      <!-- Right: actions -->
      <div class="ios-navbar__right">
        <slot name="trailing" />
      </div>
    </div>

    <!-- Optional sub-bar (search, filters, etc.) -->
    <div v-if="$slots.subbar" class="ios-navbar__subbar">
      <slot name="subbar" />
    </div>
  </header>
</template>

<style scoped>
.ios-navbar {
  position: sticky;
  top: calc(var(--ios-topbar-h) + 8px);
  z-index: 50;
  width: 100%;
  max-width: var(--ios-content-max-web);
  margin: 0 auto;
  border-radius: var(--ios-radius-xxl);
  background: var(--ios-glass-bg-light-md);
  backdrop-filter: blur(var(--ios-glass-blur-md)) saturate(2.2) brightness(1.03);
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-md)) saturate(2.2) brightness(1.03);
  border: 0.5px solid var(--ios-glass-border-light);
  box-shadow: var(--ios-glass-shadow-layer), inset 0 1px 0 rgba(255,255,255,0.75);
}

.ios-navbar__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: var(--ios-topbar-h);
  padding: 0 var(--ios-web-page-padding);
  gap: 16px;
}

.ios-navbar__left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.ios-navbar__back {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--ios-radius-sm);
  border: none;
  background: var(--ios-color-fill-quaternary);
  color: var(--ios-color-label-secondary);
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s var(--ios-ease-stiff),
    color 0.15s var(--ios-ease-stiff);
}

.ios-navbar__back:hover {
  background: var(--ios-color-fill-tertiary);
  color: var(--ios-color-label-primary);
}

.ios-navbar__title-group {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.ios-navbar__title {
  font-size: 17px;
  line-height: 1.3;
  letter-spacing: -0.35px;
  font-weight: 600;
  color: var(--ios-color-label-primary);
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ios-navbar__title--large {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.ios-navbar__right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.ios-navbar__subbar {
  padding: 0 var(--ios-web-page-padding) 14px;
}

/* On mobile, use smaller padding and reset sticky top */
@media (max-width: 768px) {
  .ios-navbar {
    top: calc(var(--ios-topbar-h) + 6px);
    max-width: calc(100vw - 24px);
    border-radius: var(--ios-radius-xl);
  }

  .ios-navbar__inner {
    padding: 0 var(--ios-page-padding);
    min-height: var(--ios-navbar-h);
  }

  .ios-navbar__subbar {
    padding: 0 var(--ios-page-padding) 10px;
  }

  .ios-navbar__title--large {
    font-size: 18px;
  }
}
</style>
