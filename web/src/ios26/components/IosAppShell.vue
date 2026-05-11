<script setup lang="ts">
import { onMounted, watch, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import IosTabBar from './IosTabBar.vue'
import IosIcon from './IosIcon.vue'
import IosToast from './IosToast.vue'
import { IOS_THEME_STORAGE_KEY, IOS_TABS, type IosTheme } from '../constants'
import '../styles/index.css'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const theme = ref<IosTheme>('system')

function applyTheme(value: IosTheme) {
  theme.value = value
  if (typeof window !== 'undefined') {
    try {
      localStorage.setItem(IOS_THEME_STORAGE_KEY, value)
    } catch {
      /* noop */
    }
  }
}

onMounted(() => {
  auth.initAuth()
  if (auth.isLoggedIn) auth.fetchProfile()
  try {
    const saved = localStorage.getItem(IOS_THEME_STORAGE_KEY) as IosTheme | null
    if (saved === 'light' || saved === 'dark' || saved === 'system') {
      theme.value = saved
    }
  } catch {
    /* noop */
  }
})

watch(() => route.path, () => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = ''
  }
})

function isMobileActive(matchPaths: string[]): boolean {
  const p = route.path.replace(/\/$/, '') || '/ios26'
  if (matchPaths.includes('/ios26')) {
    if (p === '/ios26') return true
  }
  return matchPaths.some(m => m !== '/ios26' && p.startsWith(m))
}

const mobileTabs = computed(() => IOS_TABS.map(t => ({ ...t, active: isMobileActive(t.match) })))

defineExpose({ applyTheme })
</script>

<template>
  <div class="ios26-app" :data-theme="theme === 'system' ? undefined : theme">
    <div class="ios-aurora" aria-hidden="true" />

    <!-- Desktop: left sidebar -->
    <IosTabBar />

    <!-- Main content area -->
    <main class="ios-app__main">
      <router-view v-slot="{ Component }">
        <transition name="ios-page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Mobile bottom tab bar (visible only on small screens) -->
    <nav class="ios-mobile-tabbar ios-no-select" aria-label="移动端导航">
      <div class="ios-mobile-tabbar__chrome">
        <ul class="ios-mobile-tabbar__list">
          <li v-for="t in mobileTabs" :key="t.name">
            <button
              class="ios-mobile-tabbar__btn"
              :class="{ 'is-active': t.active }"
              type="button"
              @click="route.path !== t.path && router.push(t.path)"
            >
              <IosIcon :path="t.icon" :size="22" :stroke="t.active ? 2.3 : 1.7" />
              <span class="ios-mobile-tabbar__label">{{ t.label }}</span>
            </button>
          </li>
        </ul>
      </div>
    </nav>

    <IosToast />
  </div>
</template>

<style scoped>
/* Content pushed below top nav bar */
.ios-app__main {
  position: relative;
  z-index: 1;
  padding-top: var(--ios-topbar-h);
  min-height: 100dvh;
}

/* Mobile: add bottom padding for tab bar */
@media (max-width: 768px) {
  .ios-app__main {
    padding-bottom: calc(72px + env(safe-area-inset-bottom, 0px));
  }
}

/* Mobile bottom tab bar */
.ios-mobile-tabbar {
  display: none;
}

@media (max-width: 768px) {
  .ios-mobile-tabbar {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 60;
    padding: 6px 12px calc(env(safe-area-inset-bottom, 0px) + 6px);
    background: var(--ios-glass-bg-light-lg);
    backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(var(--ios-glass-saturate));
    -webkit-backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(var(--ios-glass-saturate));
    border-top: 0.5px solid var(--ios-color-separator);
  }

  .ios-mobile-tabbar__chrome {
    width: 100%;
  }

  .ios-mobile-tabbar__list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 2px;
  }

  .ios-mobile-tabbar__btn {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    min-height: 44px;
    padding: 4px 2px;
    border: none;
    background: transparent;
    color: var(--ios-color-label-secondary);
    cursor: pointer;
    border-radius: var(--ios-radius-md);
    transition: color 0.2s var(--ios-ease-gentle);
  }

  .ios-mobile-tabbar__btn.is-active {
    color: var(--ios-color-tint);
  }

  .ios-mobile-tabbar__label {
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.04px;
  }
}
</style>
