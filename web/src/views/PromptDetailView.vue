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
      <button class="inline-flex items-center gap-1.5 text-sm mb-5 no-underline" style="color:var(--text-secondary)" @click="router.push('/prompts')">← 返回模板列表</button>

      <div class="mb-6">
        <h1 class="text-[28px] font-extrabold mb-3">{{ template.title }}</h1>
        <div class="flex flex-wrap gap-2 mb-4">
          <span class="tag">{{ template.tags[0] }}</span>
          <span v-for="(tag, i) in template.tags.slice(1)" :key="tag" class="tag" :class="i % 2 ? 'tag-cyan' : ''">{{ tag }}</span>
        </div>
        <div class="flex gap-2.5">
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
            {{ template.isLiked ? '♥' : '♡' }} {{ template.likeCount }} 点赞
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
            {{ template.isFavorited ? '★' : '☆' }} 收藏
          </button>
          <button class="action-btn" @click="copyContent">📋 一键复制</button>
        </div>
      </div>

      <!-- Content -->
      <div class="p-7 rounded-2xl mb-8 relative border" style="background:var(--bg-elevated);border-color:var(--border)">
        <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg absolute top-4 right-4" @click="copyContent">📋 复制</button>
        <pre class="text-[13.5px] leading-relaxed whitespace-pre-wrap break-words" style="font-family:'JetBrains Mono',monospace;color:var(--text-primary)">{{ template.content }}</pre>
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
