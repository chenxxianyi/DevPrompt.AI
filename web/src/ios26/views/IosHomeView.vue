<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePromptStore } from '@/store/prompt'
import { useUiStore } from '@/store/ui'
import { aiTools } from '@/mock/data'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosIcon from '../components/IosIcon.vue'

const router = useRouter()
const prompt = usePromptStore()
const ui = useUiStore()

const quickInput = ref('')
const quickTool = ref('')

const generators = [
  {
    id: 'project',
    title: '项目 Prompt',
    desc: '输入项目信息，生成完整的项目开发 Prompt',
    icon: 'M3 7.5A1.5 1.5 0 0 1 4.5 6h4.879a1.5 1.5 0 0 1 1.06.44l1.122 1.12A1.5 1.5 0 0 0 12.62 8H19.5A1.5 1.5 0 0 1 21 9.5V18a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18V7.5Z',
    tint: '#A78BFA',
    shadow: 'rgba(167,139,250,0.35)',
    cardTint: 'purple' as const,
  },
  {
    id: 'cursor-rules',
    title: 'Cursor Rules',
    desc: '根据语言、框架，一键生成 .cursorrules 配置',
    icon: 'M5.636 4.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728ZM12 8.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V15a.75.75 0 0 1-1.5 0v-2.25H9a.75.75 0 0 1 0-1.5h2.25V9a.75.75 0 0 1 .75-.75Z',
    tint: '#7DD3FC',
    shadow: 'rgba(125,211,252,0.35)',
    cardTint: 'teal' as const,
  },
  {
    id: 'claude-code',
    title: 'Claude Code',
    desc: '针对 Claude Code 场景，生成结构化任务 Prompt',
    icon: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2 14 8 20 8',
    tint: '#FDB78A',
    shadow: 'rgba(253,183,138,0.35)',
    cardTint: 'orange' as const,
  },
  {
    id: 'optimize',
    title: 'Prompt 优化',
    desc: '将模糊 Prompt 优化为专业、精确的版本',
    icon: 'M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z',
    tint: '#6EE7B7',
    shadow: 'rgba(110,231,183,0.35)',
    cardTint: 'green' as const,
  },
]

const capabilities = [
  { title: '多工具适配', desc: 'Cursor / Claude / GPT / Gemini / DeepSeek / Qwen', icon: 'M13 10V3L4 14h7v7l9-11h-7z', color: '#A78BFA' },
  { title: '模板库', desc: '覆盖前端、后端、DevOps 等高质量模板', icon: 'M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20', color: '#7DD3FC' },
  { title: '一键优化', desc: '把粗糙 Prompt 优化为专业版本', icon: 'M12 20h9 M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z', color: '#F9A8D4' },
  { title: '即拿即用', desc: '支持复制和 Markdown 导出，零成本接入', icon: 'M9 9h13v13H9z M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1', color: '#6EE7B7' },
]

onMounted(async () => {
  await prompt.fetchTemplates()
})

function goGenerator(type?: string) {
  const query: Record<string, string> = {}
  if (type) query.type = type
  if (quickTool.value) query.targetTool = quickTool.value
  router.push({ path: '/ios26/generator', query })
}

function quickGenerate() {
  if (!quickInput.value.trim()) {
    ui.showToast('请先描述你的任务或 Prompt')
    return
  }
  const query: Record<string, string> = {
    rawPrompt: quickInput.value,
    type: 'optimize',
  }
  if (quickTool.value) query.targetTool = quickTool.value
  router.push({ path: '/ios26/generator', query })
}
</script>

<template>
  <div class="ios-home">
    <div class="ios-page ios-home__page">
      <!-- Hero quick input -->
      <section class="ios-home__hero">
        <p class="ios-home__hero-tagline ios-text-subheadline">
          <span class="ios-home__pulse" />
          支持 Cursor / Claude / GPT / Gemini
        </p>
        <h2 class="ios-home__hero-title">为开发者生成<br />可直接使用的 AI 编程 Prompt</h2>

        <IosGlassPanel class="ios-home__quick" :inset="false">
          <div class="ios-home__quick-row">
            <textarea
              v-model="quickInput"
              class="ios-home__quick-input"
              placeholder="描述你的项目、任务或想优化的 Prompt..."
              rows="2"
              @keydown.enter.meta.prevent="quickGenerate"
            />
            <IosButton variant="filled" size="md" class="ios-home__quick-btn" @click="quickGenerate">
              <IosIcon path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z" :size="18" :stroke="2" />
              生成
            </IosButton>
          </div>

          <div class="ios-home__tools">
            <span class="ios-home__tools-label">目标工具</span>
            <div class="ios-home__chips ios-scroll-x">
              <button
                v-for="t in aiTools"
                :key="t.name"
                class="ios-chip ios-press"
                :class="{ 'is-active': quickTool === t.name }"
                @click="quickTool = quickTool === t.name ? '' : t.name"
              >
                <span class="ios-chip__dot" :style="{ background: t.color }" />
                {{ t.name }}
              </button>
            </div>
          </div>
        </IosGlassPanel>

        <div class="ios-home__hero-cta">
          <IosButton variant="glass" size="md" to="/ios26/prompts">
            <IosIcon path="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20" :size="18" :stroke="2" />
            浏览模板库
          </IosButton>
        </div>
      </section>

      <!-- Generator entries -->
      <section class="ios-section">
        <header class="ios-section__header">
          <h3 class="ios-text-title-3">生成器</h3>
        </header>
        <div class="ios-home__cards">
          <IosGlassPanel
            v-for="g in generators"
            :key="g.id"
            interactive
            :tint="g.cardTint"
            class="ios-home__card"
            @click="goGenerator(g.id)"
          >
            <span class="ios-home__card-icon" :style="{ background: g.tint, boxShadow: `0 8px 24px ${g.shadow}` }">
              <IosIcon :path="g.icon" :size="22" :stroke="1.8" />
            </span>
            <div class="ios-home__card-title ios-text-headline">{{ g.title }}</div>
            <div class="ios-text-subheadline ios-home__card-desc">{{ g.desc }}</div>
            <div class="ios-home__card-arrow">
              <IosIcon path="M5 12h14M13 6l6 6-6 6" :size="15" :stroke="2" />
            </div>
          </IosGlassPanel>
        </div>
      </section>

      <!-- Hot templates -->
      <section class="ios-section">
        <header class="ios-section__header">
          <h3 class="ios-text-title-3">热门模板</h3>
          <router-link to="/ios26/prompts" class="ios-section__more">
            查看全部
            <IosIcon path="M9 18l6-6-6-6" :size="14" :stroke="2.4" />
          </router-link>
        </header>
        <div v-if="prompt.loading" class="ios-home__loading">加载中...</div>
        <div v-else class="ios-home__templates">
          <router-link
            v-for="t in prompt.filteredTemplates.slice(0, 4)"
            :key="t.id"
            :to="`/ios26/prompts/${t.slug}`"
            class="ios-home__template-link"
          >
            <IosGlassPanel interactive tint="blue" class="ios-home__template">
              <div class="ios-text-headline">{{ t.title }}</div>
              <p class="ios-text-subheadline ios-home__template-desc">{{ t.description }}</p>
              <div class="ios-home__tags">
                <span v-for="tag in t.tags.slice(0, 3)" :key="tag" class="ios-tag">{{ tag }}</span>
              </div>
              <div class="ios-home__meta">
                <span class="ios-home__meta-item">
                  <IosIcon path="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" :size="12" :stroke="2" />
                  {{ t.likeCount }}
                </span>
                <span class="ios-home__meta-item">
                  <IosIcon path="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" :size="12" :stroke="2" />
                  {{ t.favoriteCount }}
                </span>
              </div>
            </IosGlassPanel>
          </router-link>
        </div>
      </section>

      <!-- Capabilities -->
      <section class="ios-section">
        <header class="ios-section__header">
          <h3 class="ios-text-title-3">核心能力</h3>
        </header>
        <div class="ios-home__caps">
          <IosGlassPanel
            v-for="(cap, i) in capabilities"
            :key="cap.title"
            size="sm"
            :tint="(['purple','teal','pink','green'] as const)[i]"
            class="ios-home__cap"
          >
            <span class="ios-home__cap-icon" :style="{ color: cap.color, background: cap.color + '22' }">
              <IosIcon :path="cap.icon" :size="18" :stroke="2" />
            </span>
            <div class="ios-text-headline ios-home__cap-title">{{ cap.title }}</div>
            <p class="ios-text-subheadline ios-home__cap-desc">{{ cap.desc }}</p>
          </IosGlassPanel>
        </div>
      </section>

      <!-- CTA -->
      <section class="ios-section">
        <IosGlassPanel size="lg" class="ios-home__cta">
          <h3 class="ios-text-title-2 ios-home__cta-title">解锁无限生成次数</h3>
          <p class="ios-text-body ios-home__cta-sub">升级 Pro 会员，享受更多高级功能与额度</p>
          <IosButton variant="filled" size="lg" to="/ios26/pricing">查看会员方案</IosButton>
        </IosGlassPanel>
      </section>
    </div>
  </div>
</template>

<style scoped>
.ios-home {
  display: flex;
  flex-direction: column;
}

.ios-home__hero {
  padding: 48px 0 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.ios-home__hero-tagline {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--ios-color-tint);
  background: var(--ios-color-tint-soft);
  padding: 6px 12px;
  border-radius: var(--ios-radius-pill);
  font-weight: 500;
  margin-bottom: 16px;
}

.ios-home__pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--ios-color-systemGreen);
  box-shadow: 0 0 0 0 var(--ios-color-systemGreen);
  animation: ios-pulse 2s ease-in-out infinite;
}

.ios-home__hero-title {
  font-size: clamp(32px, 4vw, 52px);
  line-height: 1.08;
  letter-spacing: -0.8px;
  font-weight: 800;
  color: var(--ios-color-label-primary);
  margin: 0 0 28px 0;
  max-width: 680px;
}

.ios-home__quick {
  padding: 20px;
  width: 100%;
  max-width: 680px;
}

.ios-home__quick-row {
  display: flex;
  gap: 10px;
  align-items: stretch;
}

.ios-home__quick-input {
  flex: 1;
  font-size: 15px;
  line-height: 22px;
  border: none;
  background: transparent;
  color: var(--ios-color-label-primary);
  resize: none;
  outline: none;
  padding: 8px 6px;
  min-height: 44px;
  font-family: inherit;
}

.ios-home__quick-input::placeholder {
  color: var(--ios-color-label-tertiary);
}

.ios-home__quick-btn {
  align-self: stretch;
}

.ios-home__tools {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.ios-home__tools-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--ios-color-label-secondary);
  flex-shrink: 0;
  letter-spacing: 0.06px;
  text-transform: uppercase;
}

.ios-home__chips {
  display: flex;
  gap: 6px;
  flex: 1;
  min-width: 0;
  padding: 4px 0;
}

.ios-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-color-fill-quaternary);
  border: 1px solid transparent;
  color: var(--ios-color-label-primary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
}

.ios-chip.is-active {
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
  border-color: var(--ios-color-tint-soft-strong);
}

.ios-chip__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.ios-home__hero-cta {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  gap: 10px;
}

.ios-section {
  margin-top: 40px;
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

.ios-section__more {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 15px;
  color: var(--ios-color-tint);
  font-weight: 500;
}

.ios-home__cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

@media (max-width: 960px) {
  .ios-home__cards {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 520px) {
  .ios-home__cards {
    grid-template-columns: 1fr;
  }
}

.ios-home__card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  cursor: pointer;
}

.ios-home__card-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.ios-home__card-title {
  margin: 0;
  color: var(--ios-color-label-primary);
}

.ios-home__card-desc {
  color: var(--ios-color-label-secondary);
  margin: 0;
  flex: 1;
  font-size: 13px;
  line-height: 1.5;
}

.ios-home__card-arrow {
  display: flex;
  align-items: center;
  color: var(--ios-color-tint);
  margin-top: 4px;
}

.ios-home__loading {
  padding: 24px;
  text-align: center;
  color: var(--ios-color-label-secondary);
}

.ios-home__templates {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
}

.ios-home__template-link {
  text-decoration: none;
  color: inherit;
  display: block;
}

.ios-home__template {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ios-home__template-desc {
  margin: 0;
  color: var(--ios-color-label-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ios-home__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.ios-tag {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
  letter-spacing: 0.06px;
}

.ios-home__meta {
  display: flex;
  gap: 14px;
  font-size: 12px;
  color: var(--ios-color-label-tertiary);
  margin-top: 4px;
}

.ios-home__meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.ios-home__caps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}

.ios-home__cap {
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ios-home__cap-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--ios-radius-md);
  flex-shrink: 0;
}

.ios-home__cap-title {
  margin: 0;
  color: var(--ios-color-label-primary);
}

.ios-home__cap-desc {
  margin: 0;
  color: var(--ios-color-label-secondary);
  font-size: 13px;
  line-height: 1.5;
}

.ios-home__cta {
  text-align: center;
  padding: 40px 32px;
}

.ios-home__cta-title {
  margin: 0 0 8px;
  color: var(--ios-color-label-primary);
}

.ios-home__cta-sub {
  margin: 0 0 20px;
  color: var(--ios-color-label-secondary);
}
</style>
