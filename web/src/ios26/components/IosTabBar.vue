<script setup lang="ts">
import { computed, ref, watch, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import IosIcon from './IosIcon.vue'
import { IOS_TABS, IOS_ROOT } from '../constants'

const route = useRoute()
const router = useRouter()

const tabRefs = ref<HTMLElement[]>([])
const trackRef = ref<HTMLElement | null>(null)

const indicatorStyle = ref({
  transform: 'translateX(0px)',
  width: '0px',
  opacity: '0',
})

function isActive(matchPaths: string[]): boolean {
  const p = route.path.replace(/\/$/, '') || '/ios26'
  if (matchPaths.includes('/ios26')) {
    if (p === '/ios26') return true
  }
  return matchPaths.some(m => m !== '/ios26' && p.startsWith(m))
}

const items = computed(() => IOS_TABS.map(t => ({ ...t, active: isActive(t.match) })))

const activeIndex = computed(() => items.value.findIndex(t => t.active))

function updateIndicator() {
  nextTick(() => {
    const idx = activeIndex.value
    const el = tabRefs.value[idx]
    const track = trackRef.value
    if (!el || !track) return
    const trackRect = track.getBoundingClientRect()
    const tabRect = el.getBoundingClientRect()
    indicatorStyle.value = {
      transform: `translateX(${tabRect.left - trackRect.left}px)`,
      width: `${tabRect.width}px`,
      opacity: '1',
    }
  })
}

watch(activeIndex, updateIndicator)
onMounted(updateIndicator)

function go(path: string) {
  if (route.path !== path) router.push(path)
}
</script>

<template>
  <header class="ios-topnav ios-no-select" role="banner">
    <div class="ios-topnav__inner">
      <!-- Brand (left) -->
      <div class="ios-topnav__brand" @click="router.push(IOS_ROOT)">
        <div class="ios-topnav__logo">
          <IosIcon
            path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"
            :size="18"
            :stroke="2.2"
          />
        </div>
        <span class="ios-topnav__brand-name">DevPrompt</span>
      </div>

      <!-- Nav tabs (center) with morphing glass indicator -->
      <nav class="ios-topnav__nav" aria-label="主导航">
        <div ref="trackRef" class="ios-topnav__track">
          <!-- Liquid Glass sliding indicator -->
          <div class="ios-topnav__indicator" :style="indicatorStyle" aria-hidden="true" />

          <button
            v-for="(t, i) in items"
            :key="t.name"
            :ref="el => { if (el) tabRefs[i] = el as HTMLElement }"
            class="ios-topnav__tab"
            :class="{ 'is-active': t.active }"
            :aria-current="t.active ? 'page' : undefined"
            type="button"
            @click="go(t.path)"
          >
            <IosIcon :path="t.icon" :size="17" :stroke="t.active ? 2.2 : 1.6" />
            <span class="ios-topnav__tab-label">{{ t.label }}</span>
          </button>
        </div>
      </nav>

      <!-- Actions (right) -->
      <div class="ios-topnav__actions">
        <router-link
          to="/ios26/dashboard"
          class="ios-topnav__action-btn"
          :class="{ 'is-active': route.path.startsWith('/ios26/dashboard') }"
          title="我的"
        >
          <IosIcon
            path="M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Zm0 2c-3.866 0-7 2.239-7 5v1h14v-1c0-2.761-3.134-5-7-5Z"
            :size="19"
            :stroke="1.8"
          />
        </router-link>
      </div>
    </div>
  </header>
</template>

<style scoped>
/* ===== Top Navigation Bar ===== */
.ios-topnav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 60;
  height: var(--ios-topbar-h);
  background: var(--ios-glass-bg-light-lg);
  backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(var(--ios-glass-saturate));
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(var(--ios-glass-saturate));
  border-bottom: 0.5px solid var(--ios-color-separator);
}

/* Top highlight layer */
.ios-topnav::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, var(--ios-glass-highlight-light) 0%, transparent 60%);
  opacity: 0.4;
  z-index: 0;
}

.ios-topnav__inner {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 20px;
  gap: 12px;
}

/* Brand */
.ios-topnav__brand {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  flex-shrink: 0;
  padding: 6px 8px 6px 0;
  border-radius: var(--ios-radius-md);
  transition: opacity 0.15s;
}

.ios-topnav__brand:hover {
  opacity: 0.8;
}

.ios-topnav__logo {
  width: 30px;
  height: 30px;
  border-radius: var(--ios-radius-sm);
  background: var(--ios-color-tint);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ios-topnav__brand-name {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: -0.3px;
  color: var(--ios-color-label-primary);
  white-space: nowrap;
}

/* Nav tabs */
.ios-topnav__nav {
  flex: 1;
  display: flex;
  justify-content: center;
}

.ios-topnav__track {
  position: relative;
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px;
  background: var(--ios-color-fill-quaternary);
  border-radius: var(--ios-radius-pill);
}

/* Morphing Liquid Glass indicator */
.ios-topnav__indicator {
  position: absolute;
  top: 4px;
  bottom: 4px;
  left: 0;
  background: var(--ios-glass-bg-light-md);
  backdrop-filter: blur(var(--ios-glass-blur-md)) saturate(var(--ios-glass-saturate));
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-md)) saturate(var(--ios-glass-saturate));
  border: 0.5px solid var(--ios-glass-border-light);
  border-radius: var(--ios-radius-pill);
  box-shadow: var(--ios-glass-shadow-press);
  transition:
    transform 0.45s cubic-bezier(0.34, 1.56, 0.64, 1),
    width 0.45s cubic-bezier(0.34, 1.56, 0.64, 1),
    opacity 0.2s;
  pointer-events: none;
}

.ios-topnav__tab {
  position: relative;
  z-index: 1;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border: none;
  background: transparent;
  color: var(--ios-color-label-secondary);
  cursor: pointer;
  border-radius: var(--ios-radius-pill);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: -0.15px;
  white-space: nowrap;
  transition: color 0.2s var(--ios-ease-gentle);
  min-height: 34px;
}

.ios-topnav__tab:hover:not(.is-active) {
  color: var(--ios-color-label-primary);
}

.ios-topnav__tab.is-active {
  color: var(--ios-color-label-primary);
  font-weight: 600;
}

.ios-topnav__tab-label {
  line-height: 1;
}

/* Actions */
.ios-topnav__actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.ios-topnav__action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: var(--ios-radius-pill);
  border: none;
  background: transparent;
  color: var(--ios-color-label-secondary);
  cursor: pointer;
  text-decoration: none;
  transition: background 0.15s var(--ios-ease-stiff),
    color 0.15s var(--ios-ease-stiff);
}

.ios-topnav__action-btn:hover {
  background: var(--ios-color-fill-tertiary);
  color: var(--ios-color-label-primary);
}

.ios-topnav__action-btn.is-active {
  color: var(--ios-color-tint);
  background: var(--ios-color-tint-soft);
}

/* Responsive: hide center nav on small screens (mobile bottom bar handles it) */
@media (max-width: 768px) {
  .ios-topnav__nav {
    display: none;
  }

  .ios-topnav__brand-name {
    display: none;
  }

  .ios-topnav__inner {
    justify-content: space-between;
    padding: 0 16px;
  }
}
</style>
