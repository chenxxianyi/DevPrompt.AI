<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { usePromptStore } from '@/store/prompt'
import { useUiStore } from '@/store/ui'
import PromptCard from '@/components/PromptCard.vue'

const router = useRouter()
const auth = useAuthStore()
const promptStore = usePromptStore()
const ui = useUiStore()

onMounted(async () => {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  try {
    await promptStore.fetchFavoriteTemplates()
  } catch {
    ui.showToast('加载收藏失败，请稍后重试')
  }
})
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6 py-8 pb-16">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-[28px] font-extrabold tracking-tight">我的收藏</h1>
          <p class="mt-1 text-[15px]" style="color:var(--text-secondary)">沉淀你高频使用的 Prompt 模板，随时复用</p>
        </div>
        <router-link to="/prompts" class="btn btn-ghost text-sm px-4 py-2 rounded-xl no-underline">浏览模板库</router-link>
      </div>

      <div v-if="promptStore.favoriteLoading" class="text-center py-16" style="color:var(--text-muted)">加载中...</div>

      <div v-else-if="promptStore.favoriteTemplates.length === 0" class="glass p-10 text-center">
        <h2 class="text-lg font-bold mb-2">还没有收藏模板</h2>
        <p class="text-sm mb-5" style="color:var(--text-secondary)">去模板库挑选几个适合你的 Prompt，沉淀成自己的效率资产。</p>
        <router-link to="/prompts" class="btn btn-primary text-sm px-5 py-2.5 rounded-xl no-underline">去收藏模板</router-link>
      </div>

      <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-4">
        <router-link
          v-for="t in promptStore.favoriteTemplates"
          :key="t.id"
          :to="`/prompts/${t.slug}`"
          class="no-underline"
        >
          <PromptCard :template="t" />
        </router-link>
      </div>
    </div>
  </div>
</template>
