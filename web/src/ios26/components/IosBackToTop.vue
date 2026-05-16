<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const visible = ref(false)
const SCROLL_THRESHOLD = 400

function onScroll() {
  visible.value = window.scrollY > SCROLL_THRESHOLD
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<template>
  <Transition name="ios-btt">
    <button
      v-if="visible"
      class="ios-btt"
      type="button"
      aria-label="回到顶部"
      title="回到顶部"
      @click="scrollToTop"
    >
      <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.4"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M5 15l7-7 7 7" />
      </svg>
    </button>
  </Transition>
</template>

<style scoped>
.ios-btt {
  position: fixed;
  right: 24px;
  bottom: 32px;
  z-index: 70;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--ios-glass-border-light);
  background: var(--ios-glass-bg-light-lg);
  backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(2.4) brightness(1.04);
  -webkit-backdrop-filter: blur(var(--ios-glass-blur-lg)) saturate(2.4) brightness(1.04);
  color: var(--ios-color-tint);
  cursor: pointer;
  box-shadow:
    var(--ios-glass-shadow-layer),
    inset 0 1px 0 rgba(255, 255, 255, 0.65);
  transition:
    transform 0.22s var(--ios-ease-snappy),
    box-shadow 0.22s var(--ios-ease-stiff),
    background 0.22s var(--ios-ease-stiff),
    border-color 0.22s var(--ios-ease-stiff);
  -webkit-tap-highlight-color: transparent;
  user-select: none;
  -webkit-user-select: none;
}

@media (hover: hover) {
  .ios-btt:hover {
    transform: translateY(-2px) scale(1.06);
    box-shadow:
      var(--ios-glass-shadow-background),
      inset 0 1px 0 rgba(255, 255, 255, 0.75);
    border-color: var(--ios-color-tint);
  }
}

.ios-btt:active {
  transform: scale(0.92);
  opacity: 0.8;
}

/* 移动端：在底部 tab bar 上方留出空间 */
@media (max-width: 768px) {
  .ios-btt {
    right: 16px;
    bottom: calc(88px + env(safe-area-inset-bottom, 0px));
    width: 40px;
    height: 40px;
  }
}

/* Transition */
.ios-btt-enter-active {
  transition:
    opacity 0.25s var(--ios-ease-stiff),
    transform 0.28s var(--ios-ease-snappy);
}

.ios-btt-leave-active {
  transition:
    opacity 0.18s var(--ios-ease-stiff),
    transform 0.18s var(--ios-ease-stiff);
}

.ios-btt-enter-from {
  opacity: 0;
  transform: translateY(12px) scale(0.85);
}

.ios-btt-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.9);
}
</style>
