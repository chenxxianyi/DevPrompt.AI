<script setup lang="ts">
/**
 * IosGlassPanel — Liquid Glass 基础容器
 * size: 'sm' | 'md' | 'lg'
 * inset: 是否自带内边距（默认 true）
 * tint: 马卡龙色调，叠加彩色渐变光晕
 */
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg'
  inset?: boolean
  as?: string
  interactive?: boolean
  tint?: 'blue' | 'purple' | 'pink' | 'green' | 'orange' | 'teal' | 'yellow' | 'indigo'
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
    props.tint ? `ios-glass--tint-${props.tint}` : '',
  ].filter(Boolean).join(' ')
})
</script>

<template>
  <component :is="as" :class="klass">
    <slot />
  </component>
</template>

<style scoped>
/* Light sweep: a diagonal white stripe moves across background-position */
@keyframes glass-shine {
  0%   { background-position: -200% center; }
  100% { background-position:  200% center; }
}

.ios-glass--inset {
  padding: var(--ios-space-5);
}

/* ─── interactive base ─── */
.ios-glass--interactive {
  overflow: hidden;
  background-image: none;
  background-size: 300% 100%;
  transition:
    transform        0.28s var(--ios-ease-gentle),
    box-shadow       0.28s var(--ios-ease-gentle),
    border-color     0.28s var(--ios-ease-gentle),
    filter           0.28s var(--ios-ease-gentle),
    background-color 0.28s var(--ios-ease-gentle);
  cursor: pointer;
}

/* ─── On hover: boost ::before specular highlight ─── */
.ios-glass--interactive::before {
  transition: opacity 0.28s var(--ios-ease-gentle), background 0.28s var(--ios-ease-gentle);
}

.ios-glass--interactive:hover::before {
  opacity: 1 !important;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 1.0)  0%,
    rgba(200, 235, 255, 0.45) 20%,
    transparent               50%,
    rgba(255, 255, 255, 0.05) 100%
  ) !important;
}

/* ─── Hover: lift + layered glow + inner edge highlight ─── */
.ios-glass--interactive:hover {
  transform: translateY(-5px) scale(1.012);
  filter: saturate(1.28) brightness(1.05);
  /* 光扫条纹叠加在背景上 */
  background-image: linear-gradient(
    115deg,
    transparent 20%,
    rgba(255, 255, 255, 0.55) 45%,
    rgba(255, 255, 255, 0.70) 50%,
    rgba(255, 255, 255, 0.55) 55%,
    transparent 80%
  );
  animation: glass-shine 1.4s var(--ios-ease-gentle) forwards;

  box-shadow:
    /* 远端大光晕 */
    0 32px 72px rgba(63,  164, 232, 0.32),
    /* 中层彩色漫射 */
    0 12px 32px rgba(167, 139, 250, 0.22),
    /* 近端锐利阴影 */
    0  4px 10px rgba(63,  164, 232, 0.16),
    /* 顶部内发光边缘（最亮） */
    inset 0  1.5px 0 rgba(255, 255, 255, 1.00),
    /* 底部内边蓝色折射 */
    inset 0 -1px  0 rgba(147, 197, 253, 0.20),
    /* 左侧内边高光 */
    inset 1.5px 0 0 rgba(255, 255, 255, 0.65),
    /* 外边框光环 */
    0 0 0 1.5px rgba(147, 197, 253, 0.60);
}

/* ─── Tint overlay also becomes more vivid on hover ─── */
.ios-glass--interactive:hover::after {
  opacity: 1;
  filter: saturate(1.7) brightness(1.15);
}

/* ===== Macaron Tint Overlays ===== */
.ios-glass--tint-blue::after,
.ios-glass--tint-purple::after,
.ios-glass--tint-pink::after,
.ios-glass--tint-green::after,
.ios-glass--tint-orange::after,
.ios-glass--tint-teal::after,
.ios-glass--tint-yellow::after,
.ios-glass--tint-indigo::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  pointer-events: none;
  z-index: 0;
  opacity: 0.9;
  transition: opacity 0.28s var(--ios-ease-gentle), filter 0.28s var(--ios-ease-gentle);
}

.ios-glass--tint-blue::after {
  background: linear-gradient(145deg, rgba(147, 197, 253, 0.30) 0%, rgba(96, 165, 250, 0.12) 50%, transparent 100%);
}
.ios-glass--tint-purple::after {
  background: linear-gradient(145deg, rgba(196, 181, 253, 0.32) 0%, rgba(167, 139, 250, 0.12) 50%, transparent 100%);
}
.ios-glass--tint-pink::after {
  background: linear-gradient(145deg, rgba(249, 168, 212, 0.34) 0%, rgba(244, 114, 182, 0.12) 50%, transparent 100%);
}
.ios-glass--tint-green::after {
  background: linear-gradient(145deg, rgba(110, 231, 183, 0.32) 0%, rgba(52, 211, 153, 0.12) 50%, transparent 100%);
}
.ios-glass--tint-orange::after {
  background: linear-gradient(145deg, rgba(253, 186, 116, 0.34) 0%, rgba(251, 146, 60, 0.12) 50%, transparent 100%);
}
.ios-glass--tint-teal::after {
  background: linear-gradient(145deg, rgba(125, 211, 252, 0.32) 0%, rgba(56, 189, 248, 0.12) 50%, transparent 100%);
}
.ios-glass--tint-yellow::after {
  background: linear-gradient(145deg, rgba(253, 224, 71, 0.30) 0%, rgba(250, 204, 21, 0.10) 50%, transparent 100%);
}
.ios-glass--tint-indigo::after {
  background: linear-gradient(145deg, rgba(165, 180, 252, 0.32) 0%, rgba(129, 140, 248, 0.12) 50%, transparent 100%);
}
</style>
