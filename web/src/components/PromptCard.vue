<script setup lang="ts">
import type { PromptTemplate } from '@/types'

defineProps<{ template: PromptTemplate }>()
</script>

<template>
  <div class="glass p-5 cursor-pointer hover:-translate-y-0.5 transition-transform duration-300">
    <h4 class="text-[15px] font-semibold mb-1.5">{{ template.title }}</h4>
    <p class="text-[13px] text-[var(--text-secondary)] mb-3 line-clamp-2 leading-relaxed">{{ template.description }}</p>
    <div class="flex flex-wrap gap-1.5 mb-3">
      <span v-for="(tag, i) in template.tags" :key="tag" class="tag" :class="i % 2 ? 'tag-cyan' : ''">{{ tag }}</span>
    </div>
    <div class="flex gap-4 text-[13px] text-[var(--text-muted)]">
      <span>♡ {{ template.likeCount }}</span>
      <span>☆ {{ template.favoriteCount }}</span>
    </div>
  </div>
</template>
