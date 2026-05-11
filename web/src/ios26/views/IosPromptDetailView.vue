<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePromptStore } from '@/store/prompt'
import { useUiStore } from '@/store/ui'
import type { PromptTemplate } from '@/types'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosIcon from '../components/IosIcon.vue'

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
  return promptStore.allTemplates.filter((item: PromptTemplate) => item.categoryId === t.categoryId && item.id !== t.id).slice(0, 3)
})

function copyContent() {
  if (template.value) {
    navigator.clipboard.writeText(template.value.content)
    ui.showToast('已复制', 'success')
  }
}

function useTemplate() {
  if (!template.value) return
  router.push({
    path: '/ios26/generator',
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
  <div class="ios-detail">
    <IosNavBar :title="template?.title || '模板详情'" :show-back="true" back-to="/ios26/prompts">
      <template #trailing>
        <button v-if="template" class="ios-detail__nav-btn ios-press" type="button" :aria-pressed="template.isFavorited" @click="toggleFav">
          <IosIcon
            path="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
            :size="20"
            :stroke="2"
            :filled="!!template.isFavorited"
          />
        </button>
      </template>
    </IosNavBar>

    <div class="ios-page" v-if="template">
      <!-- Title block -->
      <div class="ios-detail__head">
        <h1 class="ios-text-large-title ios-detail__title">{{ template.title }}</h1>
        <p class="ios-text-callout ios-detail__desc">{{ template.description }}</p>
        <div class="ios-detail__tags">
          <span v-for="tag in template.tags" :key="tag" class="ios-tag">{{ tag }}</span>
        </div>
      </div>

      <!-- Stat actions -->
      <IosGlassPanel size="md" class="ios-detail__actions">
        <button class="ios-detail__action ios-press" @click="toggleLike">
          <IosIcon
            path="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"
            :size="20"
            :stroke="2"
            :filled="!!template.isLiked"
          />
          <span class="ios-text-footnote-emph">{{ template.likeCount }}</span>
          <span class="ios-text-caption-1 ios-detail__action-label">点赞</span>
        </button>
        <span class="ios-detail__action-divider" />
        <button class="ios-detail__action ios-press" @click="toggleFav">
          <IosIcon
            path="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
            :size="20"
            :stroke="2"
            :filled="!!template.isFavorited"
          />
          <span class="ios-text-footnote-emph">{{ template.favoriteCount }}</span>
          <span class="ios-text-caption-1 ios-detail__action-label">收藏</span>
        </button>
        <span class="ios-detail__action-divider" />
        <button class="ios-detail__action ios-press" @click="copyContent">
          <IosIcon path="M9 9h13v13H9z M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" :size="20" :stroke="2" />
          <span class="ios-text-footnote-emph">复制</span>
          <span class="ios-text-caption-1 ios-detail__action-label">内容</span>
        </button>
      </IosGlassPanel>

      <!-- Content block -->
      <IosGlassPanel size="md" class="ios-detail__content-card">
        <div class="ios-detail__content-head">
          <span class="ios-text-subheadline-emph">Prompt 正文</span>
          <button class="ios-detail__copy ios-press" @click="copyContent">
            <IosIcon path="M9 9h13v13H9z M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" :size="14" :stroke="2" />
            复制
          </button>
        </div>
        <pre class="ios-mono ios-detail__pre">{{ template.content }}</pre>
      </IosGlassPanel>

      <!-- Use template CTA -->
      <IosGlassPanel size="md" class="ios-detail__cta">
        <div class="ios-text-headline">把模板带入生成器</div>
        <p class="ios-text-subheadline ios-detail__cta-sub">将此模板内容预填到生成器，快速生成定制版 Prompt</p>
        <IosButton variant="filled" size="lg" block @click="useTemplate">
          <IosIcon path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z" :size="18" :stroke="2" />
          使用此模板生成
        </IosButton>
      </IosGlassPanel>

      <!-- Related -->
      <section v-if="relatedTemplates.length" class="ios-detail__section">
        <header class="ios-section__header">
          <h3 class="ios-text-title-3">相关模板</h3>
        </header>
        <div class="ios-detail__related">
          <router-link
            v-for="t in relatedTemplates"
            :key="t.id"
            :to="`/ios26/prompts/${t.slug}`"
            class="ios-detail__related-link"
          >
            <IosGlassPanel interactive size="sm">
              <div class="ios-text-headline">{{ t.title }}</div>
              <p class="ios-text-subheadline ios-detail__related-desc">{{ t.description }}</p>
              <div class="ios-detail__tags">
                <span v-for="tag in t.tags.slice(0, 2)" :key="tag" class="ios-tag">{{ tag }}</span>
              </div>
            </IosGlassPanel>
          </router-link>
        </div>
      </section>
    </div>

    <div v-else-if="loading" class="ios-detail__state">
      <span class="ios-spinner" />
      <span class="ios-text-callout">加载模板...</span>
    </div>

    <div v-else class="ios-detail__state">
      <span class="ios-text-headline">模板未找到</span>
      <IosButton variant="tinted" to="/ios26/prompts" size="md">返回模板列表</IosButton>
    </div>
  </div>
</template>

<style scoped>
.ios-detail {
  display: flex;
  flex-direction: column;
}


.ios-detail__nav-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--ios-color-fill-quaternary);
  border: none;
  color: var(--ios-color-systemYellow);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.ios-detail__head {
  padding: 4px 0 16px;
}

.ios-detail__title {
  margin: 0 0 8px;
  color: var(--ios-color-label-primary);
}

.ios-detail__desc {
  margin: 0 0 12px;
  color: var(--ios-color-label-secondary);
}

.ios-detail__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.ios-detail__actions {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 12px;
  margin-bottom: 16px;
}

.ios-detail__action {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  flex: 1;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--ios-color-label-primary);
  padding: 6px 4px;
  min-height: 44px;
}

.ios-detail__action-label {
  color: var(--ios-color-label-tertiary);
}

.ios-detail__action-divider {
  width: 1px;
  height: 28px;
  background: var(--ios-color-separator);
}

.ios-detail__content-card {
  margin-bottom: 16px;
}

.ios-detail__content-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.ios-detail__copy {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border: none;
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
  font-size: 13px;
  font-weight: 600;
  border-radius: var(--ios-radius-pill);
  cursor: pointer;
}

.ios-detail__pre {
  margin: 0;
  padding: 14px;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  background: var(--ios-color-fill-quaternary);
  color: var(--ios-color-label-primary);
  border-radius: var(--ios-radius-md);
  max-height: 520px;
  overflow-y: auto;
}

.ios-detail__cta {
  text-align: center;
  margin-bottom: 24px;
}

.ios-detail__cta-sub {
  margin: 6px 0 16px;
  color: var(--ios-color-label-secondary);
}

.ios-detail__section {
  margin-top: 16px;
}

.ios-section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px 12px;
}

.ios-section__header h3 {
  margin: 0;
  color: var(--ios-color-label-primary);
}

.ios-detail__related {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.ios-detail__related-link {
  text-decoration: none;
  color: inherit;
}

.ios-detail__related-desc {
  margin: 6px 0 8px;
  color: var(--ios-color-label-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ios-detail__state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 80px 0;
  color: var(--ios-color-label-secondary);
}
</style>
