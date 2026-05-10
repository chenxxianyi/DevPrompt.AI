<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { usePromptStore } from '@/store/prompt'
import PromptCard from '@/components/PromptCard.vue'

const prompt = usePromptStore()

onMounted(async () => {
  await prompt.fetchCategories()
  await prompt.fetchTemplates()
})

// 搜索防抖：输入变化后触发搜索
let searchTimer: ReturnType<typeof setTimeout> | null = null
watch(() => prompt.searchKeyword, (val) => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    prompt.setKeyword(val)
  }, 400)
})

function goToPage(p: number) {
  if (p >= 1 && p <= prompt.totalPages) {
    prompt.setPage(p)
  }
}

// 生成分页按钮
const pageButtons = () => {
  const pages: number[] = []
  const total = prompt.totalPages
  const current = prompt.page
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="flex gap-6 py-8">
        <!-- Sidebar -->
        <aside class="w-[220px] flex-shrink-0 hidden lg:block">
          <div class="glass p-5 sticky top-24">
            <h3 class="text-xs font-semibold uppercase tracking-wider mb-3" style="color:var(--text-muted)">分类</h3>
            <ul class="list-none">
              <li
                v-for="cat in prompt.allCategories" :key="cat.id"
                class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm cursor-pointer transition-all duration-300 mb-0.5"
                :class="prompt.selectedCategory === cat.name ? 'bg-[rgba(139,92,246,.1)] text-[var(--accent)] font-semibold' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[rgba(139,92,246,.06)]'"
                @click="prompt.setCategory(cat.name)"
              >
                {{ cat.name }}
              </li>
            </ul>
          </div>
        </aside>

        <!-- Main -->
        <div class="flex-1 min-w-0">
          <!-- Toolbar -->
          <div class="flex items-center gap-3 mb-5 flex-wrap">
            <div class="relative flex-1 min-w-[200px]">
              <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-sm" style="color:var(--text-muted)">🔍</span>
              <input
                v-model="prompt.searchKeyword"
                class="form-input pl-10"
                type="text"
                placeholder="搜索 Prompt 模板..."
              >
            </div>
            <select v-model="prompt.sortBy" @change="prompt.setSort(prompt.sortBy)" class="px-3.5 py-2.5 rounded-xl outline-none text-sm cursor-pointer" style="background:var(--bg-elevated);border:1px solid var(--border);color:var(--text-secondary);font-family:inherit">
              <option value="hot">最热门</option>
              <option value="new">最新</option>
              <option value="likes">最多点赞</option>
            </select>
          </div>

          <!-- Loading -->
          <div v-if="prompt.loading" class="text-center py-10" style="color:var(--text-muted)">加载中...</div>

          <!-- Grid -->
          <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-4">
            <router-link
              v-for="t in prompt.filteredTemplates" :key="t.id"
              :to="`/prompts/${t.slug}`"
              class="no-underline"
            >
              <PromptCard :template="t" />
            </router-link>
          </div>
          <p v-if="!prompt.loading && prompt.filteredTemplates.length === 0" class="text-center py-10" style="color:var(--text-muted)">暂无匹配的模板</p>

          <!-- Pagination -->
          <div v-if="prompt.totalPages > 1" class="flex items-center justify-center gap-1.5 mt-8 pb-10">
            <button
              class="w-9 h-9 rounded-lg flex items-center justify-center text-sm cursor-pointer border transition-all"
              :disabled="prompt.page <= 1"
              style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
              @click="goToPage(prompt.page - 1)"
            >‹</button>
            <button
              v-for="p in pageButtons()" :key="p"
              class="w-9 h-9 rounded-lg flex items-center justify-center text-sm cursor-pointer border transition-all"
              :style="p === prompt.page
                ? { background:'var(--accent)', borderColor:'var(--accent)', color:'#fff' }
                : { background:'var(--bg-elevated)', borderColor:'var(--border)', color:'var(--text-secondary)' }"
              @click="goToPage(p)"
            >{{ p }}</button>
            <button
              class="w-9 h-9 rounded-lg flex items-center justify-center text-sm cursor-pointer border transition-all"
              :disabled="prompt.page >= prompt.totalPages"
              style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
              @click="goToPage(prompt.page + 1)"
            >›</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
