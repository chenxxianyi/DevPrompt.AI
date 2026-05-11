<script setup lang="ts">
import { onMounted, watch, computed } from 'vue'
import { usePromptStore } from '@/store/prompt'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosSegmentedControl from '../components/IosSegmentedControl.vue'
import IosIcon from '../components/IosIcon.vue'

const prompt = usePromptStore()

onMounted(async () => {
  await prompt.fetchCategories()
  await prompt.fetchTemplates()
})

let searchTimer: ReturnType<typeof setTimeout> | null = null
watch(() => prompt.searchKeyword, (val: string) => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    prompt.setKeyword(val)
  }, 400)
})

const sortOptions = [
  { value: 'hot', label: '最热门' },
  { value: 'new', label: '最新' },
  { value: 'likes', label: '最多点赞' },
]

const pageButtons = computed(() => {
  const pages: number[] = []
  const total = prompt.totalPages
  const current = prompt.page
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

function goToPage(p: number) {
  if (p >= 1 && p <= prompt.totalPages) {
    prompt.setPage(p)
  }
}
</script>

<template>
  <div class="ios-prompts">
    <IosNavBar large>
      <template #subbar>
        <div class="ios-prompts__search">
          <span class="ios-prompts__search-icon">
            <IosIcon path="M11 19a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM21 21l-4.35-4.35" :size="17" :stroke="2" />
          </span>
          <input
            v-model="prompt.searchKeyword"
            class="ios-input ios-prompts__search-input"
            type="search"
            placeholder="搜索 Prompt 模板..."
          >
        </div>
      </template>
    </IosNavBar>

    <div class="ios-page">
      <!-- Sort -->
      <IosSegmentedControl
        :model-value="prompt.sortBy"
        :options="sortOptions"
        size="md"
        @update:model-value="(v: string) => prompt.setSort(v)"
      />

      <!-- Categories chips -->
      <div class="ios-prompts__categories ios-scroll-x">
        <button
          v-for="cat in prompt.allCategories"
          :key="cat.id"
          class="ios-chip ios-press"
          :class="{ 'is-active': prompt.selectedCategory === cat.name }"
          @click="prompt.setCategory(cat.name)"
        >
          {{ cat.name }}
        </button>
      </div>

      <!-- Loading -->
      <div v-if="prompt.loading" class="ios-prompts__loading">
        <span class="ios-spinner" /> 加载中...
      </div>

      <!-- Grid -->
      <div v-else class="ios-prompts__grid ios-prompts__grid--web">
        <router-link
          v-for="t in prompt.filteredTemplates"
          :key="t.id"
          :to="`/ios26/prompts/${t.slug}`"
          class="ios-prompts__link"
        >
          <IosGlassPanel interactive>
            <div class="ios-text-headline">{{ t.title }}</div>
            <p class="ios-text-subheadline ios-prompts__desc">{{ t.description }}</p>
            <div class="ios-prompts__tags">
              <span v-for="tag in t.tags.slice(0, 3)" :key="tag" class="ios-tag">{{ tag }}</span>
            </div>
            <div class="ios-prompts__meta">
              <span class="ios-prompts__meta-item">
                <IosIcon path="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" :size="12" :stroke="2" />
                {{ t.likeCount }}
              </span>
              <span class="ios-prompts__meta-item">
                <IosIcon path="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" :size="12" :stroke="2" />
                {{ t.favoriteCount }}
              </span>
              <span class="ios-prompts__meta-item">
                <IosIcon path="M3 12a9 9 0 1 0 18 0 9 9 0 0 0-18 0Zm6-2 6 2-6 2v-4Z" :size="12" :stroke="2" />
                {{ t.useCount }}
              </span>
            </div>
          </IosGlassPanel>
        </router-link>
      </div>

      <p v-if="!prompt.loading && prompt.filteredTemplates.length === 0" class="ios-prompts__empty">
        暂无匹配的模板
      </p>

      <!-- Pagination -->
      <nav v-if="prompt.totalPages > 1" class="ios-pagination">
        <button class="ios-pagination__btn ios-press" :disabled="prompt.page <= 1" @click="goToPage(prompt.page - 1)">
          <IosIcon path="M15 18l-6-6 6-6" :size="16" :stroke="2.2" />
        </button>
        <button
          v-for="p in pageButtons"
          :key="p"
          class="ios-pagination__btn ios-press"
          :class="{ 'is-active': p === prompt.page }"
          @click="goToPage(p)"
        >
          {{ p }}
        </button>
        <button class="ios-pagination__btn ios-press" :disabled="prompt.page >= prompt.totalPages" @click="goToPage(prompt.page + 1)">
          <IosIcon path="M9 18l6-6-6-6" :size="16" :stroke="2.2" />
        </button>
      </nav>
    </div>
  </div>
</template>

<style scoped>
.ios-prompts {
  display: flex;
  flex-direction: column;
}

.ios-prompts__search {
  position: relative;
  max-width: 320px;
}

.ios-prompts__search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--ios-color-label-tertiary);
  pointer-events: none;
}

.ios-prompts__search-input {
  padding-left: 42px;
  font-size: 14px;
  height: 36px;
  padding-top: 0;
  padding-bottom: 0;
  line-height: 36px;
}

.ios-prompts__categories {
  display: flex;
  gap: 8px;
  margin: 14px 0 18px;
  padding: 4px 0;
}

.ios-prompts__loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px;
  color: var(--ios-color-label-secondary);
}

.ios-prompts__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 14px;
}

.ios-prompts__grid--web {
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.ios-prompts__link {
  text-decoration: none;
  color: inherit;
  display: block;
}

.ios-prompts__desc {
  margin: 6px 0 10px;
  color: var(--ios-color-label-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ios-prompts__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
}

.ios-prompts__meta {
  display: flex;
  gap: 14px;
  font-size: 12px;
  color: var(--ios-color-label-tertiary);
}

.ios-prompts__meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.ios-prompts__empty {
  text-align: center;
  padding: 40px 0;
  color: var(--ios-color-label-tertiary);
}

.ios-pagination {
  display: flex;
  justify-content: center;
  gap: 6px;
  margin-top: 24px;
}

.ios-pagination__btn {
  min-width: 36px;
  height: 36px;
  padding: 0 12px;
  border-radius: var(--ios-radius-md);
  border: none;
  background: var(--ios-color-fill-quaternary);
  color: var(--ios-color-label-primary);
  font-weight: 500;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.ios-pagination__btn.is-active {
  background: var(--ios-color-tint);
  color: #fff;
}

.ios-pagination__btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
