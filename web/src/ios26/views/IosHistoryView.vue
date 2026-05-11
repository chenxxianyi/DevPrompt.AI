<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useGeneratorStore, type GeneratorTab } from '@/store/generator'
import { getHistory } from '@/api/generator'
import type { GeneratedPrompt } from '@/types'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosSheet from '../components/IosSheet.vue'
import IosIcon from '../components/IosIcon.vue'

const router = useRouter()
const ui = useUiStore()
const gen = useGeneratorStore()

const loading = ref(false)
const list = ref<GeneratedPrompt[]>([])
const page = ref(1)
const total = ref(0)
const pageSize = 20
const activeType = ref('')
const selectedItem = ref<GeneratedPrompt | null>(null)
const showDetail = ref(false)

const typeOptions = [
  { value: '', label: '全部' },
  { value: 'project', label: '项目' },
  { value: 'cursor-rules', label: 'Cursor' },
  { value: 'claude-code', label: 'Claude' },
  { value: 'optimize', label: '优化' },
]

const typeLabelMap: Record<string, string> = {
  project: '项目 Prompt',
  'cursor-rules': 'Cursor 规则',
  'claude-code': 'Claude Code',
  optimize: 'Prompt 优化',
}

function typeLabel(t: string) {
  return typeLabelMap[t] || t
}

async function fetchData() {
  loading.value = true
  try {
    const params: any = { page: page.value, pageSize, type: activeType.value || undefined }
    const res = await getHistory(params)
    list.value = res.data.data.list
    total.value = res.data.data.total
  } catch {
    /* silent */
  } finally {
    loading.value = false
  }
}

function switchType(t: string) {
  activeType.value = t
  page.value = 1
  fetchData()
}

function openDetail(item: GeneratedPrompt) {
  selectedItem.value = item
  showDetail.value = true
}

function formatTime(t: string): string {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function formatInput(input: any): string {
  if (!input) return ''
  let parsed = input
  if (typeof input === 'string') {
    try { parsed = JSON.parse(input) } catch { return input }
  }
  if (typeof parsed !== 'object') return String(parsed)
  if (parsed.features) {
    const f = parsed.features
    return Array.isArray(f) ? f.join('\n') : String(f)
  }
  if (parsed.rawPrompt) return String(parsed.rawPrompt)
  if (parsed.task) return String(parsed.task)
  if (parsed.requirements) {
    const r = parsed.requirements
    return Array.isArray(r) ? r.join('\n') : String(r)
  }
  if (parsed.rules) {
    const r = parsed.rules
    return Array.isArray(r) ? r.join('\n') : String(r)
  }
  return String(Object.values(parsed).find((v: any) => typeof v === 'string' && v.length > 10) || '')
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  ui.showToast('已复制', 'success')
}

function regenerateFrom(item: GeneratedPrompt) {
  gen.loadFromHistory({
    type: item.type as GeneratorTab,
    input: item.input,
  })
  router.push({ path: '/ios26/generator', query: { type: item.type } })
  showDetail.value = false
}

const totalPages = computed(() => Math.ceil(total.value / pageSize) || 1)

onMounted(fetchData)
</script>

<template>
  <div class="ios-history">
    <IosNavBar :title="'历史记录'" :show-back="true" back-to="/ios26/generator" />

    <div class="ios-page">
      <!-- Filter chips -->
      <div class="ios-history__filters ios-scroll-x">
        <button
          v-for="opt in typeOptions"
          :key="opt.value"
          class="ios-chip ios-press"
          :class="{ 'is-active': activeType === opt.value }"
          @click="switchType(opt.value)"
        >
          {{ opt.label }}
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="ios-history__state">
        <span class="ios-spinner" /> 加载中...
      </div>

      <!-- Empty -->
      <div v-else-if="list.length === 0" class="ios-history__empty">
        <IosGlassPanel size="md">
          <IosIcon path="M3 12a9 9 0 1 0 18 0 9 9 0 0 0-18 0Zm9-5v5l3 2" :size="36" :stroke="1.6" />
          <p class="ios-text-headline ios-history__empty-title">暂无生成记录</p>
          <p class="ios-text-subheadline ios-history__empty-sub">前往生成器，开始你的第一次生成</p>
          <IosButton variant="filled" size="md" to="/ios26/generator">前往生成器</IosButton>
        </IosGlassPanel>
      </div>

      <!-- List -->
      <div v-else class="ios-history__list">
        <IosGlassPanel
          v-for="item in list"
          :key="item.id"
          interactive
          size="md"
          @click="openDetail(item)"
        >
          <div class="ios-history__row">
            <div class="ios-history__info">
              <div class="ios-history__head">
                <span class="ios-tag">{{ typeLabel(item.type) }}</span>
                <span v-if="item.title" class="ios-text-subheadline-emph ios-history__title">{{ item.title }}</span>
              </div>
              <p class="ios-text-footnote ios-history__excerpt">
                {{ item.output?.slice(0, 160) || '（无内容）' }}
              </p>
              <span class="ios-text-caption-1 ios-history__time">{{ formatTime(item.createdAt) }}</span>
            </div>
            <IosIcon path="M9 18l6-6-6-6" :size="18" :stroke="2.2" />
          </div>
        </IosGlassPanel>
      </div>

      <!-- Pagination -->
      <nav v-if="total > pageSize" class="ios-pagination">
        <button class="ios-pagination__btn ios-press" :disabled="page <= 1" @click="page > 1 && (page--, fetchData())">
          <IosIcon path="M15 18l-6-6 6-6" :size="16" :stroke="2.2" />
        </button>
        <span class="ios-pagination__info">{{ page }} / {{ totalPages }}</span>
        <button class="ios-pagination__btn ios-press" :disabled="page >= totalPages" @click="page < totalPages && (page++, fetchData())">
          <IosIcon path="M9 18l6-6-6-6" :size="16" :stroke="2.2" />
        </button>
      </nav>
    </div>

    <!-- Detail Sheet -->
    <IosSheet v-model="showDetail" :title="selectedItem?.title || '生成记录'">
      <template v-if="selectedItem">
        <div class="ios-history__sheet-meta">
          <span class="ios-tag">{{ typeLabel(selectedItem.type) }}</span>
          <span class="ios-text-caption-1 ios-history__time">{{ formatTime(selectedItem.createdAt) }}</span>
        </div>

        <section class="ios-history__sheet-section">
          <header class="ios-history__sheet-section-head">
            <span class="ios-text-subheadline-emph">输入文案</span>
            <IosButton variant="tinted" size="sm" @click="copy(formatInput(selectedItem.input))">复制</IosButton>
          </header>
          <pre class="ios-mono ios-history__sheet-pre">{{ formatInput(selectedItem.input) }}</pre>
        </section>

        <section class="ios-history__sheet-section">
          <header class="ios-history__sheet-section-head">
            <span class="ios-text-subheadline-emph">生成结果</span>
            <IosButton variant="tinted" size="sm" @click="copy(selectedItem.output)">复制</IosButton>
          </header>
          <pre class="ios-mono ios-history__sheet-pre">{{ selectedItem.output }}</pre>
        </section>
      </template>

      <template #footer>
        <IosButton variant="filled" size="md" block @click="selectedItem && regenerateFrom(selectedItem)">
          再次生成
        </IosButton>
      </template>
    </IosSheet>
  </div>
</template>

<style scoped>
.ios-history {
  display: flex;
  flex-direction: column;
}

.ios-history__filters {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  padding: 4px 0;
}

.ios-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
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

.ios-history__state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 60px 0;
  color: var(--ios-color-label-secondary);
}

.ios-history__empty {
  text-align: center;
  margin-top: 24px;
}

.ios-history__empty-title {
  margin: 14px 0 4px;
}

.ios-history__empty-sub {
  margin: 0 0 16px;
  color: var(--ios-color-label-secondary);
}

.ios-history__list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.ios-history__row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ios-history__info {
  flex: 1;
  min-width: 0;
}

.ios-history__head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.ios-history__title {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--ios-color-label-primary);
  flex: 1;
  min-width: 0;
}

.ios-history__excerpt {
  margin: 0 0 4px;
  color: var(--ios-color-label-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ios-history__time {
  color: var(--ios-color-label-tertiary);
}

.ios-tag {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
  letter-spacing: 0.06px;
  white-space: nowrap;
}

.ios-pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
  margin-top: 24px;
}

.ios-pagination__btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: var(--ios-color-fill-quaternary);
  color: var(--ios-color-label-primary);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.ios-pagination__btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.ios-pagination__info {
  font-size: 14px;
  color: var(--ios-color-label-secondary);
  font-weight: 500;
  letter-spacing: -0.2px;
}

.ios-history__sheet-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.ios-history__sheet-section {
  margin-bottom: 16px;
}

.ios-history__sheet-section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.ios-history__sheet-pre {
  margin: 0;
  padding: 12px 14px;
  background: var(--ios-color-fill-quaternary);
  border-radius: var(--ios-radius-md);
  font-size: 12.5px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 300px;
  overflow-y: auto;
  color: var(--ios-color-label-primary);
}
</style>
