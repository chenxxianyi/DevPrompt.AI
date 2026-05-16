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
  <Transition name="btt">
    <button
      v-if="visible"
      class="btt"
      type="button"
      aria-label="回到顶部"
      title="回到顶部"
      @click="scrollToTop"
    >
      <svg
        width="18"
        height="18"
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
.btt {
  position: fixed;
  right: 24px;
  bottom: 32px;
  z-index: 70;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border);
  background: var(--bg-card);
  backdrop-filter: blur(16px) saturate(1.3);
  -webkit-backdrop-filter: blur(16px) saturate(1.3);
  color: var(--accent-hover);
  cursor: pointer;
  box-shadow:
    0 4px 20px rgba(0, 0, 0, 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.04);
  transition:
    transform 0.22s cubic-bezier(0.34, 1.56, 0.64, 1),
    box-shadow 0.25s ease,
    border-color 0.25s ease,
    color 0.25s ease;
  -webkit-tap-highlight-color: transparent;
  user-select: none;
  -webkit-user-select: none;
}

@media (hover: hover) {
  .btt:hover {
    transform: translateY(-3px);
    border-color: var(--accent);
    box-shadow:
      0 6px 24px var(--accent-glow),
      0 2px 8px rgba(0, 0, 0, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.06);
    color: var(--accent-hover);
  }
}

.btt:active {
  transform: scale(0.92);
  opacity: 0.75;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .btt {
    right: 16px;
    bottom: 24px;
    width: 40px;
    height: 40px;
    border-radius: 12px;
  }
}

/* Transition */
.btt-enter-active {
  transition:
    opacity 0.25s ease,
    transform 0.28s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.btt-leave-active {
  transition:
    opacity 0.18s ease,
    transform 0.18s ease;
}

.btt-enter-from {
  opacity: 0;
  transform: translateY(12px) scale(0.85);
}

.btt-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.9);
}
</style>
