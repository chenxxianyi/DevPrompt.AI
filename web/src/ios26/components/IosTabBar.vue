<script setup lang="ts">
import { computed, ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import IosIcon from './IosIcon.vue'
import { IOS_TABS, IOS_ROOT } from '../constants'

const route = useRoute()
const router = useRouter()

const tabRefs = ref<HTMLElement[]>([])
const trackRef = ref<HTMLElement | null>(null)
let resizeObserver: ResizeObserver | null = null

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

const isCompact = ref(false)
let lastScrollY = 0

watch(isCompact, () => {
  updateIndicator()
})

function handleScroll() {
  const currentY = window.scrollY
  if (currentY > lastScrollY && currentY > 60) {
    isCompact.value = true
  } else if (currentY < lastScrollY) {
    isCompact.value = false
  }
  lastScrollY = currentY
}

onMounted(() => {
  lastScrollY = window.scrollY
  updateIndicator()
  window.addEventListener('resize', updateIndicator)
  window.addEventListener('scroll', handleScroll, { passive: true })
  if (trackRef.value) {
    resizeObserver = new ResizeObserver(updateIndicator)
    resizeObserver.observe(trackRef.value)
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateIndicator)
  window.removeEventListener('scroll', handleScroll)
  resizeObserver?.disconnect()
})

function go(path: string) {
  if (route.path !== path) router.push(path)
}
</script>

<template>
  <header class="ios-topnav ios-no-select" :class="{ 'is-compact': isCompact }" role="banner">
    <div class="ios-topnav__inner">
      <!-- Brand (left) -->
      <router-link :to="IOS_ROOT" class="ios-topnav__brand">
        <div class="ios-topnav__logo">
          <IosIcon
            path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"
            :size="18"
            :stroke="2.2"
          />
        </div>
        <span class="ios-topnav__brand-name">DevPrompt</span>
      </router-link>

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
        <!-- 切换回经典版 -->
        <router-link
          to="/"
          class="ios-topnav__classic-btn"
          title="切换回经典版"
        >
          <IosIcon
            path="M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2v-4M9 21H5a2 2 0 0 1-2-2v-4m0 0h18"
            :size="15"
            :stroke="1.8"
          />
          <span class="ios-topnav__classic-label">经典版</span>
        </router-link>

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
/* ===== Top Navigation Bar — floating centered pill ===== */
.ios-topnav {
  position: fixed;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 60;
  height: 64px;
  width: 80vw;
  max-width: calc(100vw - 40px);
  clip-path: inset(0 round 32px);
  transition: width 0.55s var(--ios-ease-gentle);
  border-radius: 32px;
  background: var(--ios-glass-bg-light-lg);
  backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(2.4) brightness(1.04);
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(2.4) brightness(1.04);
  border: 0.5px solid var(--ios-glass-border-light);
  box-shadow:
    var(--ios-glass-shadow-layer),
    inset 0 1px 0 rgba(255, 255, 255, 0.80);
}

/* Top highlight layer */
.ios-topnav::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  pointer-events: none;
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.88) 0%,
    rgba(200, 235, 255, 0.20) 35%,
    transparent 65%
  );
  opacity: 0.7;
  z-index: 0;
}

.ios-topnav__inner {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 20px;
  gap: 14px;
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
  text-decoration: none;
  color: inherit;
}

.ios-topnav__brand:hover {
  opacity: 0.8;
}

.ios-topnav__logo {
  width: 36px;
  height: 36px;
  border-radius: var(--ios-radius-md);
  background: var(--ios-color-tint-gradient, var(--ios-color-tint));
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ios-topnav__brand-name {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.4px;
  color: var(--ios-color-label-primary);
  white-space: nowrap;
  overflow: hidden;
  max-width: 120px;
  transition: max-width 0.55s var(--ios-ease-gentle), opacity 0.3s var(--ios-ease-gentle);
}

.ios-topnav.is-compact .ios-topnav__brand-name {
  max-width: 0;
  opacity: 0;
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
  /* 不用 backdrop-filter，改用明确的白色半透明填充，避免嵌套模糊失效 */
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(125, 211, 252, 0.50);
  border-radius: var(--ios-radius-pill);
  box-shadow:
    0 2px 12px rgba(63, 164, 232, 0.20),
    0 1px 4px rgba(63, 164, 232, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 1);
  transition:
    transform 0.45s cubic-bezier(0.34, 1.56, 0.64, 1),
    width   0.45s cubic-bezier(0.34, 1.56, 0.64, 1),
    opacity 0.2s;
  pointer-events: none;
}

.ios-topnav__tab {
  position: relative;
  z-index: 1;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 18px;
  border: none;
  background: transparent;
  color: var(--ios-color-label-secondary);
  cursor: pointer;
  border-radius: var(--ios-radius-pill);
  font-size: 15px;
  font-weight: 500;
  letter-spacing: -0.2px;
  white-space: nowrap;
  transition: color 0.2s var(--ios-ease-gentle), padding 0.55s var(--ios-ease-gentle);
  min-height: 42px;
}

.ios-topnav.is-compact .ios-topnav__tab {
  padding: 9px 10px;
}

.ios-topnav__tab:hover:not(.is-active) {
  color: var(--ios-color-label-primary);
}

.ios-topnav__tab.is-active {
  color: var(--ios-color-tint);
  font-weight: 700;
}

.ios-topnav__tab-label {
  line-height: 1;
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  max-width: 80px;
  transition: max-width 0.55s var(--ios-ease-gentle), opacity 0.3s var(--ios-ease-gentle);
}

.ios-topnav.is-compact .ios-topnav__tab-label {
  max-width: 0;
  opacity: 0;
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

/* 经典版切换按钮 */
.ios-topnav__classic-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  border-radius: var(--ios-radius-pill);
  border: 1px solid var(--ios-color-separator);
  background: var(--ios-color-fill-quaternary);
  color: var(--ios-color-label-secondary);
  font-size: 12px;
  font-weight: 500;
  text-decoration: none;
  white-space: nowrap;
  transition: background 0.18s var(--ios-ease-stiff),
    color 0.18s var(--ios-ease-stiff),
    border-color 0.18s var(--ios-ease-stiff);
}

.ios-topnav__classic-btn:hover {
  background: var(--ios-color-fill-secondary);
  color: var(--ios-color-label-primary);
  border-color: var(--ios-color-separator-opaque);
}

/* Responsive: mobile — stretch to near full width */
@media (max-width: 768px) {
  .ios-topnav {
    left: 12px;
    right: 12px;
    top: 8px;
    transform: none;
    width: auto;
    min-width: unset;
    max-width: unset;
    border-radius: 24px;
    height: 56px;
  }

  .ios-topnav__nav {
    display: none;
  }

  .ios-topnav__brand-name {
    display: none;
  }

  .ios-topnav__inner {
    justify-content: space-between;
    padding: 0 14px;
  }

  .ios-topnav__classic-label {
    display: none;
  }

  .ios-topnav__classic-btn {
    padding: 5px 8px;
    border: none;
    background: transparent;
  }
}

/* Compact: hide classic label text on desktop too */
.ios-topnav__classic-label {
  overflow: hidden;
  max-width: 4em;
  white-space: nowrap;
  transition: max-width 0.55s var(--ios-ease-gentle), opacity 0.3s var(--ios-ease-gentle);
}

.ios-topnav.is-compact .ios-topnav__classic-label {
  max-width: 0;
  opacity: 0;
}

.ios-topnav.is-compact {
  width: 380px;
}

.ios-topnav.is-compact .ios-topnav__classic-btn {
  padding: 5px 8px;
  border: none;
  background: transparent;
}
</style>
