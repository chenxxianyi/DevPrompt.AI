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
      <span class="flex items-center gap-1">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/></svg>
        {{ template.likeCount }}
      </span>
      <span class="flex items-center gap-1">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
        {{ template.favoriteCount }}
      </span>
    </div>
  </div>
</template>
