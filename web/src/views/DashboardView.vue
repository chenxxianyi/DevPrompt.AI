<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useUiStore } from '@/store/ui'
import { useGeneratorStore } from '@/store/generator'
import { getHistory } from '@/api/generator'
import { getGenerateStats } from '@/api/user'
import type { GeneratedPrompt } from '@/types'
import type { GeneratorTab } from '@/store/generator'

const router = useRouter()
const auth = useAuthStore()
const ui = useUiStore()
const gen = useGeneratorStore()

const loading = ref(false)
const loadFailed = ref(false)
const stats = ref({ totalUsed: 0, dailyLimit: 5 })
const recentList = ref<GeneratedPrompt[]>([])

const typeLabelMap: Record<string, string> = {
  project: '项目 Prompt', 'cursor-rules': 'Cursor 规则',
  'claude-code': 'Claude Code', optimize: 'Prompt 优化',
}

const remainingQuota = computed(() => {
  const limit = stats.value.dailyLimit
  if (limit >= 999999) return Infinity
  return Math.max(0, limit - stats.value.totalUsed)
})

const quotaPercent = computed(() => {
  const limit = stats.value.dailyLimit
  if (limit >= 999999) return 0
  return Math.min(100, Math.round((stats.value.totalUsed / limit) * 100))
})

const showUpgradeTip = computed(() => {
  if (!auth.isLoggedIn) return false
  if (auth.user?.membershipLevel !== 'free') return false
  return remainingQuota.value <= 5
})

onMounted(async () => {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  loading.value = true
  try {
    // 从后端 API 获取真实额度数据
    const statsRes = await getGenerateStats()
    stats.value.totalUsed = statsRes.data.data.dailyUsed
    stats.value.dailyLimit = statsRes.data.data.dailyLimit

    // 获取最近生成记录
    const res = await getHistory({ page: 1, pageSize: 5 })
    recentList.value = res.data.data.list
  } catch {
    loadFailed.value = true
    stats.value.dailyLimit = 0 // 标记加载失败，隐藏错误额度
  } finally {
    loading.value = false
  }
})

function formatTime(t: string): string {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function regenerate(item: GeneratedPrompt) {
  gen.loadFromHistory({
    type: item.type as GeneratorTab,
    input: item.input,
  })
  router.push({ path: '/generator', query: { type: item.type } })
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  ui.showToast('已复制到剪贴板', 'success')
}

function goGenerator(type?: string) {
  const query: Record<string, string> = {}
  if (type) query.type = type
  router.push({ path: '/generator', query })
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="py-8 pb-16">
        <!-- 页面标题 -->
        <div class="flex items-center justify-between mb-8">
          <div>
            <h1 class="text-[28px] font-extrabold tracking-tight">我的工作台</h1>
            <p class="mt-1 text-[15px]" style="color:var(--text-secondary)">管理你的 Prompt 资产和生成记录</p>
          </div>
          <router-link
            to="/generator"
            class="hidden md:inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-[13px] font-medium no-underline transition-all"
            style="background:rgba(139,92,246,.1);color:var(--accent)"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
            去生成
          </router-link>
        </div>

        <div v-if="loading" class="flex items-center justify-center py-20">
          <div class="w-6 h-6 border-2 rounded-full animate-spin" style="border-color:var(--accent) transparent transparent" />
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Left Column: Quota + Member -->
          <div class="space-y-5">
            <!-- 今日额度 -->
            <div class="glass p-6">
              <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" style="color:var(--text-secondary)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                今日生成额度
              </h3>
              <div v-if="loadFailed" class="text-sm py-2" style="color:var(--text-muted)">加载失败</div>
              <template v-else>
              <div class="text-[32px] font-black mb-1">
                <span v-if="remainingQuota === Infinity">∞</span>
                <template v-else>
                  <span style="color:var(--accent-hover)">{{ remainingQuota }}</span>
                  <span class="text-base font-normal" style="color:var(--text-muted)"> / {{ stats.dailyLimit }}</span>
                </template>
              </div>
              <div class="text-xs" style="color:var(--text-muted)">剩余生成次数</div>
              <div v-if="remainingQuota !== Infinity" class="mt-4 h-2 rounded-full overflow-hidden" style="background:var(--bg-base);border:1px solid var(--border)">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :style="{
                    width: quotaPercent + '%',
                    background: quotaPercent > 80 ? 'var(--rose)' : 'var(--accent)',
                  }"
                />
              </div>
            </template>
            </div>

            <!-- 会员状态 -->
            <div class="glass p-6">
              <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" style="color:var(--text-secondary)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
                会员状态
              </h3>
              <div class="flex items-center gap-3">
                <span
                  class="text-sm font-bold px-3 py-1 rounded-lg"
                  :style="{
                    background: auth.user?.membershipLevel === 'free' ? 'rgba(139,92,246,.1)' : 'rgba(16,185,129,.15)',
                    color: auth.user?.membershipLevel === 'free' ? 'var(--accent)' : 'var(--green)',
                  }"
                >
                  {{ auth.user?.membershipLevel === 'free' ? '免费用户' : (auth.user?.membershipLevel === 'pro' ? 'Pro 会员' : auth.user?.membershipLevel) }}
                </span>
              </div>
              <div v-if="auth.user?.membershipLevel === 'free'" class="mt-4">
                <button class="w-full text-sm py-2 rounded-xl btn btn-primary justify-center" @click="router.push('/pricing')">
                  升级 Pro
                </button>
              </div>
            </div>

            <!-- 快捷入口 -->
            <div class="glass p-6">
              <h3 class="text-sm font-semibold mb-4" style="color:var(--text-secondary)">快捷生成</h3>
              <div class="flex flex-col gap-2">
                <button class="text-sm px-4 py-2.5 rounded-xl flex items-center gap-2 transition-all cursor-pointer border" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);text-align:left;font-family:inherit" @click="goGenerator('project')">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7.5A1.5 1.5 0 0 1 4.5 6h4.879a1.5 1.5 0 0 1 1.06.44l1.122 1.12A1.5 1.5 0 0 0 12.62 8H19.5A1.5 1.5 0 0 1 21 9.5V18a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18V7.5Z"/></svg>
                  项目 Prompt 生成
                </button>
                <button class="text-sm px-4 py-2.5 rounded-xl flex items-center gap-2 transition-all cursor-pointer border" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);text-align:left;font-family:inherit" @click="goGenerator('cursor-rules')">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M5.636 4.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728Z"/></svg>
                  Cursor Rules
                </button>
                <button class="text-sm px-4 py-2.5 rounded-xl flex items-center gap-2 transition-all cursor-pointer border" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);text-align:left;font-family:inherit" @click="goGenerator('optimize')">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
                  Prompt 优化
                </button>
              </div>
            </div>
          </div>

          <!-- Right Column: Recent History -->
          <div class="md:col-span-2 space-y-5">
            <!-- 升级提示 -->
            <div v-if="showUpgradeTip" class="p-5 rounded-xl flex items-start gap-3" style="background:rgba(245,158,11,.1);border:1px solid rgba(245,158,11,.2)">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="var(--amber)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
              <div class="flex-1">
                <p class="text-sm font-medium" style="color:var(--amber)">免费额度即将用尽</p>
                <p class="text-xs mt-1" style="color:var(--text-secondary)">升级 Pro 会员获得更多生成次数和高级功能</p>
                <button class="mt-2 text-xs px-4 py-1.5 rounded-lg font-medium" style="background:var(--amber);color:#fff" @click="router.push('/pricing')">查看方案</button>
              </div>
            </div>

            <!-- 最近生成 -->
            <div class="glass p-6">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold flex items-center gap-2" style="color:var(--text-secondary)">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                  最近生成记录
                </h3>
                <router-link to="/generator/history" class="text-xs no-underline" style="color:var(--accent)">查看全部</router-link>
              </div>

              <div v-if="recentList.length === 0" class="text-center py-8" style="color:var(--text-muted)">
                <p class="text-sm">暂无生成记录</p>
                <router-link to="/generator" class="inline-block mt-2 text-xs font-medium no-underline" style="color:var(--accent)">开始第一次生成 →</router-link>
              </div>

              <div v-else class="flex flex-col gap-3">
                <div
                  v-for="item in recentList" :key="item.id"
                  class="p-4 rounded-xl transition-all cursor-pointer border"
                  style="background:var(--bg-elevated);border-color:var(--border)"
                  @click="regenerate(item)"
                >
                  <div class="flex items-center justify-between gap-3">
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center gap-2 mb-1">
                        <span class="text-[11px] px-2 py-0.5 rounded font-medium" style="background:rgba(139,92,246,.1);color:var(--accent)">
                          {{ typeLabelMap[item.type] || item.type }}
                        </span>
                        <span v-if="item.title" class="text-sm truncate font-medium" style="color:var(--text-primary)">{{ item.title }}</span>
                      </div>
                      <p class="text-[12px] line-clamp-1" style="color:var(--text-secondary)">
                        {{ item.output?.slice(0, 120) || '（无内容）' }}
                      </p>
                      <span class="text-[11px]" style="color:var(--text-tertiary)">{{ formatTime(item.createdAt) }}</span>
                    </div>
                    <div class="flex items-center gap-1.5">
                      <button class="text-xs px-2.5 py-1.5 rounded-lg font-medium transition-all cursor-pointer" style="background:rgba(139,92,246,.1);color:var(--accent);border:none;white-space:nowrap" @click.stop="copy(item.output)">复制</button>
                      <button class="text-xs px-2.5 py-1.5 rounded-lg font-medium transition-all cursor-pointer" style="background:rgba(139,92,246,.1);color:var(--accent);border:none;white-space:nowrap" @click.stop="regenerate(item)">再次生成</button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 收藏模板 & 使用建议 -->
            <div class="glass p-6">
              <h3 class="text-sm font-semibold mb-2 flex items-center gap-2" style="color:var(--text-secondary)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                收藏的模板
              </h3>
              <p class="text-sm mb-4" style="color:var(--text-muted)">浏览模板库，收藏你常用的 Prompt 模板</p>
              <router-link to="/prompts" class="inline-flex items-center gap-1.5 text-sm px-4 py-2 rounded-lg no-underline font-medium transition-all" style="background:rgba(139,92,246,.1);color:var(--accent)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>
                浏览模板库
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
