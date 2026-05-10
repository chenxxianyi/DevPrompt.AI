<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useGeneratorStore } from '@/store/generator'
import { getHistory } from '@/api/generator'
import type { GeneratedPrompt } from '@/types'
import type { GeneratorTab } from '@/store/generator'

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
  { value: 'project', label: '项目 Prompt' },
  { value: 'cursor-rules', label: 'Cursor 规则' },
  { value: 'claude-code', label: 'Claude Code' },
  { value: 'optimize', label: 'Prompt 优化' },
]

const typeLabelMap: Record<string, string> = {
  project: '项目 Prompt', 'cursor-rules': 'Cursor 规则',
  'claude-code': 'Claude Code', optimize: 'Prompt 优化',
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
    // ignore
  } finally {
    loading.value = false
  }
}

function switchType(t: string) {
  activeType.value = t
  page.value = 1
  fetchData()
}

function viewDetail(item: GeneratedPrompt) {
  selectedItem.value = item
  showDetail.value = true
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

function formatTime(t: string): string {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  ui.showToast('已复制到剪贴板', 'success')
}

function regenerateFrom(item: GeneratedPrompt) {
  // 将历史记录加载到生成器并跳转
  gen.loadFromHistory({
    type: item.type as GeneratorTab,
    input: item.input,
  })
  router.push({ path: '/generator', query: { type: item.type } })
  showDetail.value = false
}

const totalPages = computed(() => Math.ceil(total.value / pageSize) || 1)

onMounted(fetchData)
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="py-8 pb-16">
        <!-- 页面标题 -->
        <div class="flex items-center justify-between mb-6">
          <div>
            <h1 class="text-[28px] font-extrabold tracking-tight flex items-center gap-2">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              生成历史
            </h1>
            <p class="text-[15px]" style="color:var(--text-secondary)">查看所有 AI 生成记录</p>
          </div>
          <router-link
            to="/generator"
            class="hidden md:inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-[13px] font-medium no-underline transition-all"
            style="background:rgba(139,92,246,.1);color:var(--accent)"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
            返回生成器
          </router-link>
        </div>

        <!-- 类型筛选 -->
        <div class="flex flex-wrap gap-2 mb-6">
          <button
            v-for="opt in typeOptions" :key="opt.value"
            class="px-4 py-1.5 rounded-lg text-[13px] font-medium transition-all cursor-pointer border"
            :class="activeType === opt.value ? 'selected' : ''"
            :style="activeType === opt.value
              ? { background: 'rgba(139,92,246,.12)', borderColor: 'var(--accent)', color: 'var(--accent-hover)' }
              : { background: 'var(--bg-elevated)', borderColor: 'var(--border)', color: 'var(--text-secondary)' }"
            @click="switchType(opt.value)"
          >
            {{ opt.label }}
          </button>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center py-20">
          <div class="w-6 h-6 border-2 rounded-full animate-spin" style="border-color:var(--accent) transparent transparent" />
        </div>

        <!-- Empty -->
        <div v-else-if="list.length === 0" class="glass p-10 text-center">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="color:var(--text-muted);margin:0 auto 12px"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          <p class="text-[15px]" style="color:var(--text-secondary)">暂无生成记录</p>
          <router-link to="/generator" class="inline-block mt-3 text-sm font-medium no-underline" style="color:var(--accent)">去生成第一个 Prompt →</router-link>
        </div>

        <!-- List -->
        <div v-else class="flex flex-col gap-4">
          <div
            v-for="item in list" :key="item.id"
            class="glass p-5 rounded-xl cursor-pointer transition-all hover:translate-y-[-1px]"
            style="border-color:var(--border)"
            @click="viewDetail(item)"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-2">
                  <span class="text-xs px-2 py-0.5 rounded font-medium" style="background:rgba(139,92,246,.1);color:var(--accent)">
                    {{ typeLabel(item.type) }}
                  </span>
                  <span v-if="item.title" class="text-sm font-medium truncate" style="color:var(--text-primary)">{{ item.title }}</span>
                </div>
                <p class="text-[13px] line-clamp-2 leading-relaxed" style="color:var(--text-secondary)">
                  {{ item.output?.slice(0, 200) || '（无内容）' }}
                </p>
                <div class="flex items-center gap-4 mt-3 text-[12px]" style="color:var(--text-tertiary)">
                  <span class="flex items-center gap-1">
                    <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                    {{ formatTime(item.createdAt) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div v-if="total > pageSize" class="flex items-center justify-center gap-2 mt-4">
            <button
              class="px-3 py-1.5 rounded-lg text-sm border transition-all cursor-pointer"
              :disabled="page <= 1"
              style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
              @click="page > 1 && (page--, fetchData())"
            >
              上一页
            </button>
            <span class="text-sm px-3" style="color:var(--text-secondary)">{{ page }} / {{ totalPages }}</span>
            <button
              class="px-3 py-1.5 rounded-lg text-sm border transition-all cursor-pointer"
              :disabled="page >= totalPages"
              style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
              @click="page < totalPages && (page++, fetchData())"
            >
              下一页
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Detail Modal -->
    <Teleport to="body">
      <div v-if="showDetail && selectedItem" class="fixed inset-0 z-[100] flex items-start justify-center pt-12 pb-12 overflow-y-auto" style="background:rgba(0,0,0,.6);backdrop-filter:blur(4px)" @click.self="showDetail = false">
        <div class="w-full max-w-3xl mx-4 glass p-7 rounded-2xl" style="border-color:var(--border)">
          <div class="flex items-center justify-between mb-5">
            <div class="flex items-center gap-2">
              <span class="text-xs px-2 py-0.5 rounded font-medium" style="background:rgba(139,92,246,.1);color:var(--accent)">
                {{ typeLabel(selectedItem.type) }}
              </span>
              <span class="text-lg font-bold" style="color:var(--text-primary)">{{ selectedItem.title || '生成记录' }}</span>
            </div>
            <button class="text-xl leading-none cursor-pointer border-none" style="color:var(--text-secondary);background:transparent" @click="showDetail = false">✕</button>
          </div>

          <div class="mb-5 flex items-center gap-2 text-[13px]" style="color:var(--text-secondary)">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            {{ formatTime(selectedItem.createdAt) }}
          </div>

          <!-- Input -->
          <div class="mb-5">
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-sm font-semibold flex items-center gap-1.5" style="color:var(--text-primary)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                输入文案
              </h3>
              <button class="text-xs px-3 py-1 rounded-lg font-medium transition-all cursor-pointer border" style="background:rgba(139,92,246,.1);color:var(--accent);border-color:transparent" @click="copy(formatInput(selectedItem.input))">复制文案</button>
            </div>
            <pre class="text-[13px] p-4 rounded-xl leading-relaxed whitespace-pre-wrap overflow-x-auto" style="background:var(--bg-base);color:var(--text-primary);border:1px solid var(--border);max-height:300px;overflow-y:auto">{{ formatInput(selectedItem.input) }}</pre>
          </div>

          <!-- Output -->
          <div>
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-sm font-semibold flex items-center gap-1.5" style="color:var(--text-primary)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                生成结果
              </h3>
              <button class="text-xs px-3 py-1 rounded-lg font-medium transition-all cursor-pointer border" style="background:rgba(139,92,246,.1);color:var(--accent);border-color:transparent" @click="copy(selectedItem.output)">复制内容</button>
            </div>
            <pre class="text-[13px] p-4 rounded-xl leading-relaxed whitespace-pre-wrap overflow-x-auto" style="background:var(--bg-base);color:var(--text-primary);border:1px solid var(--border);max-height:500px;overflow-y:auto">{{ selectedItem.output }}</pre>
          </div>

          <!-- Action buttons -->
          <div class="mt-6 flex items-center justify-between">
            <button
              class="inline-flex items-center gap-1.5 px-5 py-2.5 rounded-xl text-sm font-medium transition-all cursor-pointer"
              style="background:rgba(139,92,246,.12);color:var(--accent-hover);border:1px solid rgba(139,92,246,.25)"
              @click="regenerateFrom(selectedItem)"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
              再次生成
            </button>
            <button
              class="text-sm px-4 py-2.5 rounded-xl font-medium transition-all cursor-pointer"
              style="background:rgba(139,92,246,.1);color:var(--accent)"
              @click="showDetail = false"
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
