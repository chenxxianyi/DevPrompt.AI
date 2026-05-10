<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUiStore } from '@/store/ui'
import { getHistory } from '@/api/generator'
import type { GeneratedPrompt } from '@/types'

const ui = useUiStore()

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

const typeLabel = computed(() => {
  const map: Record<string, string> = {
    project: '项目 Prompt', 'cursor-rules': 'Cursor 规则',
    'claude-code': 'Claude Code', optimize: 'Prompt 优化',
  }
  return (t: string) => map[t] || t
})

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
  // 后端返回的 input 是 JSON 字符串，需要先解析
  let parsed = input
  if (typeof input === 'string') {
    try { parsed = JSON.parse(input) } catch { return input }
  }
  if (typeof parsed !== 'object') return String(parsed)
  // 优先显示 features，这是用户输入的核心提示词
  if (parsed.features) {
    const f = parsed.features
    return Array.isArray(f) ? f.join('\n') : String(f)
  }
  // 其他类型生成器的回退
  if (parsed.rawPrompt) return String(parsed.rawPrompt)
  if (parsed.task) return String(parsed.task)
  return String(Object.values(parsed).find((v: any) => typeof v === 'string' && v.length > 10) || '')
}

function formatTime(t: string): string {
  if (!t) return ''
  // 将 ISO 时间或 Go 时间转为 年月日 时分秒
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  ui.showToast('已复制到剪贴板', 'success')
}

onMounted(fetchData)
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="py-8 pb-16">
        <h1 class="text-[28px] font-extrabold mb-2">📋 生成历史</h1>
        <p class="mb-8 text-[15px]" style="color:var(--text-secondary)">查看所有 AI 生成记录</p>

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

        <!-- 历史列表 -->
        <div v-if="loading" class="flex items-center justify-center py-20">
          <div class="w-6 h-6 border-2 rounded-full animate-spin" style="border-color:var(--accent) transparent transparent" />
        </div>

        <div v-else-if="list.length === 0" class="glass p-10 text-center">
          <p class="text-[15px]" style="color:var(--text-secondary)">暂无生成记录</p>
          <router-link to="/generator" class="inline-block mt-3 text-sm font-medium no-underline" style="color:var(--accent)">去生成第一个 Prompt →</router-link>
        </div>

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
                  <span v-if="item.title" class="text-sm truncate" style="color:var(--text-primary)">{{ item.title }}</span>
                </div>
                <p class="text-[13px] line-clamp-2 leading-relaxed" style="color:var(--text-secondary)">
                  {{ item.output?.slice(0, 200) || '（无内容）' }}
                </p>
                <div class="flex items-center gap-4 mt-3 text-[12px]" style="color:var(--text-tertiary)">
                  <span>{{ formatTime(item.createdAt) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div v-if="total > pageSize" class="flex items-center justify-center gap-2 mt-4">
            <button
              class="px-3 py-1.5 rounded-lg text-sm border transition-all"
              :disabled="page <= 1"
              style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
              @click="page > 1 && (page--, fetchData())"
            >
              上一页
            </button>
            <span class="text-sm px-3" style="color:var(--text-secondary)">{{ page }} / {{ Math.ceil(total / pageSize) }}</span>
            <button
              class="px-3 py-1.5 rounded-lg text-sm border transition-all"
              :disabled="page >= Math.ceil(total / pageSize)"
              style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary)"
              @click="page < Math.ceil(total / pageSize) && (page++, fetchData())"
            >
              下一页
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 详情弹窗 -->
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
            <button class="text-xl leading-none cursor-pointer" style="color:var(--text-secondary)" @click="showDetail = false">✕</button>
          </div>

          <div class="mb-5 text-[13px]" style="color:var(--text-secondary)">
            <span>🕐 {{ formatTime(selectedItem.createdAt) }}</span>
          </div>
          <!-- 输入文案 -->
          <div class="mb-5">
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-sm font-semibold" style="color:var(--text-primary)">📝 输入文案</h3>
              <button class="text-xs px-3 py-1 rounded-lg font-medium transition-all cursor-pointer" style="background:rgba(139,92,246,.1);color:var(--accent)" @click="copy(formatInput(selectedItem.input))">复制文案</button>
            </div>
            <pre class="text-[13px] p-4 rounded-xl leading-relaxed whitespace-pre-wrap overflow-x-auto" style="background:var(--bg-base);color:var(--text-primary);border:1px solid var(--border);max-height:300px;overflow-y:auto">{{ formatInput(selectedItem.input) }}</pre>
          </div>

          <!-- 生成内容 -->
          <div>
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-sm font-semibold" style="color:var(--text-primary)">📤 生成结果</h3>
              <button class="text-xs px-3 py-1 rounded-lg font-medium transition-all cursor-pointer" style="background:rgba(139,92,246,.1);color:var(--accent)" @click="copy(selectedItem.output)">复制内容</button>
            </div>
            <pre class="text-[13px] p-4 rounded-xl leading-relaxed whitespace-pre-wrap overflow-x-auto" style="background:var(--bg-base);color:var(--text-primary);border:1px solid var(--border);max-height:500px;overflow-y:auto">{{ selectedItem.output }}</pre>
          </div>

          <div class="mt-6 flex justify-end">
            <button class="text-sm px-4 py-2 rounded-lg font-medium transition-all cursor-pointer" style="background:rgba(139,92,246,.1);color:var(--accent)" @click="showDetail = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
