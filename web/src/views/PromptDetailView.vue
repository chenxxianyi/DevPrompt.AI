<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePromptStore } from '@/store/prompt'
import { useUiStore } from '@/store/ui'
import PromptCard from '@/components/PromptCard.vue'
import type { PromptTemplate } from '@/types'

const route = useRoute()
const router = useRouter()
const promptStore = usePromptStore()
const ui = useUiStore()

const template = ref<PromptTemplate | null>(null)
const loading = ref(true)

onMounted(async () => {
  const slug = route.params.slug as string
  template.value = await promptStore.getTemplateBySlug(slug)
  loading.value = false
})

const relatedTemplates = computed(() => {
  const t = template.value
  if (!t) return []
  return promptStore.allTemplates.filter(item => item.categoryId === t.categoryId && item.id !== t.id).slice(0, 3)
})

function copyContent() {
  if (template.value) {
    navigator.clipboard.writeText(template.value.content)
    ui.showToast('已复制到剪贴板', 'success')
  }
}

function useTemplate() {
  if (!template.value) return
  router.push({
    path: '/generator',
    query: {
      type: 'project',
      templateSlug: template.value.slug,
    },
  })
}

async function toggleLike() {
  if (!template.value) return
  const liked = await promptStore.toggleLike(template.value.id)
  if (liked !== null) {
    template.value.isLiked = liked
    template.value.likeCount += liked ? 1 : -1
    ui.showToast(liked ? '已点赞' : '已取消点赞', 'success')
  } else {
    ui.showToast('操作失败，请稍后重试')
  }
}

async function toggleFav() {
  if (!template.value) return
  const favorited = await promptStore.toggleFavorite(template.value.id)
  if (favorited !== null) {
    template.value.isFavorited = favorited
    template.value.favoriteCount += favorited ? 1 : -1
    ui.showToast(favorited ? '已收藏' : '已取消收藏', 'success')
  } else {
    ui.showToast('操作失败，请稍后重试')
  }
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6 py-8 pb-16" v-if="template">
      <button class="inline-flex items-center gap-1.5 text-sm mb-5 no-underline cursor-pointer" style="color:var(--text-secondary)" @click="router.push('/prompts')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5"/><path d="m12 19-7-7 7-7"/></svg>
        返回模板列表
      </button>

      <div class="mb-6">
        <h1 class="text-[28px] font-extrabold mb-3">{{ template.title }}</h1>
        <div class="flex flex-wrap gap-2 mb-4">
          <span class="tag">{{ template.tags[0] }}</span>
          <span v-for="(tag, i) in template.tags.slice(1)" :key="tag" class="tag" :class="i % 2 ? 'tag-cyan' : ''">{{ tag }}</span>
        </div>
        <div class="flex flex-wrap gap-2.5">
          <button
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium cursor-pointer transition-all border"
            :class="template.isLiked ? 'liked' : ''"
            :style="{
              background: template.isLiked ? 'rgba(244,63,94,.08)' : 'var(--bg-elevated)',
              borderColor: template.isLiked ? 'rgba(244,63,94,.3)' : 'var(--border)',
              color: template.isLiked ? 'var(--rose)' : 'var(--text-secondary)',
            }"
            @click="toggleLike"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" :fill="template.isLiked ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/></svg>
            {{ template.likeCount }} 点赞
          </button>
          <button
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium cursor-pointer transition-all border"
            :class="template.isFavorited ? 'faved' : ''"
            :style="{
              background: template.isFavorited ? 'rgba(245,158,11,.08)' : 'var(--bg-elevated)',
              borderColor: template.isFavorited ? 'rgba(245,158,11,.3)' : 'var(--border)',
              color: template.isFavorited ? 'var(--amber)' : 'var(--text-secondary)',
            }"
            @click="toggleFav"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" :fill="template.isFavorited ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
            收藏
          </button>
          <button
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium cursor-pointer transition-all border"
            style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
            @click="copyContent"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
            一键复制
          </button>
        </div>
      </div>

      <!-- Content -->
      <div class="p-7 rounded-2xl mb-6 relative border" style="background:var(--bg-elevated);border-color:var(--border)">
        <div class="absolute top-4 right-4 flex gap-2">
          <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg flex items-center gap-1" @click="copyContent">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
            复制
          </button>
        </div>
        <pre class="text-[13.5px] leading-relaxed whitespace-pre-wrap break-words" style="font-family:'JetBrains Mono',monospace;color:var(--text-primary)">{{ template.content }}</pre>
      </div>

      <!-- Use this template -->
      <div class="glass p-6 rounded-2xl mb-8 text-center" style="border-color:var(--border)">
        <p class="text-sm mb-4" style="color:var(--text-secondary)">将此模板内容带入生成器，快速生成定制版 Prompt</p>
        <button class="btn btn-primary px-6 py-2.5 rounded-xl text-sm" @click="useTemplate">
          <svg class="inline-block mr-1.5" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
          使用此模板生成
        </button>
      </div>

      <!-- Related -->
      <div v-if="relatedTemplates.length">
        <h3 class="text-lg font-bold mb-4">相关模板</h3>
        <div class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-3.5">
          <router-link
            v-for="t in relatedTemplates" :key="t.id"
            :to="`/prompts/${t.slug}`"
            class="no-underline"
          >
            <PromptCard :template="t" />
          </router-link>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-else-if="loading" class="flex items-center justify-center py-32" style="color:var(--text-muted)">
      <span class="spinner"></span>
    </div>

    <!-- Not found -->
    <div v-else class="flex flex-col items-center justify-center py-32 text-[var(--text-muted)]">
      <p class="text-lg mb-4">模板未找到</p>
      <router-link to="/prompts" class="btn btn-ghost px-4 py-2 rounded-lg text-sm no-underline">返回模板列表</router-link>
    </div>
  </div>
</template>

<style scoped>
.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 18px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border);
  background: var(--bg-elevated);
  color: var(--text-secondary);
}
.action-btn:hover {
  border-color: var(--border-active);
  color: var(--text-primary);
}
</style>
