<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePromptStore } from '@/store/prompt'
import { useUiStore } from '@/store/ui'
import PromptCard from '@/components/PromptCard.vue'
import { aiTools } from '@/mock/data'

const router = useRouter()
const prompt = usePromptStore()
const ui = useUiStore()

const quickInput = ref('')
const quickTool = ref('')

const generators = [
  { id: 'project', title: '项目 Prompt 生成', desc: '输入项目信息，自动生成完整的项目开发 Prompt', gradient: 'from-[var(--accent)] to-[var(--cyan)]', bg: 'rgba(139,92,246,.15)' },
  { id: 'cursor-rules', title: 'Cursor Rules 生成', desc: '根据语言、框架和代码风格，一键生成 .cursorrules 配置', gradient: 'from-[var(--cyan)] to-[var(--teal)]', bg: 'rgba(6,182,212,.15)' },
  { id: 'claude-code', title: 'Claude Code Prompt', desc: '针对 Claude Code 场景，生成结构化的任务 Prompt', gradient: 'from-[var(--amber)] to-[var(--orange)]', bg: 'rgba(245,158,11,.15)' },
  { id: 'optimize', title: 'Prompt 优化器', desc: '将模糊的 Prompt 优化为专业、精确的版本', gradient: 'from-[var(--green)] to-[var(--emerald)]', bg: 'rgba(16,185,129,.15)' },
]

const capabilities = [
  { title: '多工具适配', desc: '支持 Cursor、Claude Code、GPT、Gemini、DeepSeek、Qwen' },
  { title: '模板库', desc: '精选高质量 Prompt 模板，覆盖前端、后端、DevOps 等多个领域' },
  { title: '一键优化', desc: '将粗糙的 Prompt 优化为专业版本，提升 AI 输出质量' },
  { title: '一键复制', desc: '生成结果支持一键复制和 Markdown 渲染，即拿即用' },
]

onMounted(async () => {
  await prompt.fetchTemplates()
})

function goGenerator(type?: string) {
  const query: Record<string, string> = {}
  if (type) query.type = type
  if (quickTool.value) query.targetTool = quickTool.value
  // 如果快速输入框有内容，传给优化器
  if (quickInput.value && !type) {
    query.rawPrompt = quickInput.value
    query.type = 'optimize'
  }
  router.push({ path: '/generator', query })
}

function goGeneratorWithInput() {
  if (!quickInput.value.trim()) {
    ui.showToast('请先描述你的任务或 Prompt')
    return
  }
  const query: Record<string, string> = {
    rawPrompt: quickInput.value,
    type: 'optimize',
  }
  if (quickTool.value) query.targetTool = quickTool.value
  router.push({ path: '/generator', query })
}

function generatorIconPath(id: string): string {
  const map: Record<string, string> = {
    project: 'M3 7.5A1.5 1.5 0 0 1 4.5 6h4.879a1.5 1.5 0 0 1 1.06.44l1.122 1.12A1.5 1.5 0 0 0 12.62 8H19.5A1.5 1.5 0 0 1 21 9.5V18a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18V7.5Z',
    'cursor-rules': 'M5.636 4.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728ZM12 8.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V15a.75.75 0 0 1-1.5 0v-2.25H9a.75.75 0 0 1 0-1.5h2.25V9a.75.75 0 0 1 .75-.75Z',
    'claude-code': 'M12 2.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V9a.75.75 0 0 1-1.5 0V6.75H9a.75.75 0 0 1 0-1.5h2.25V3a.75.75 0 0 1 .75-.75ZM6.75 12a.75.75 0 0 1 .75-.75h9a.75.75 0 0 1 0 1.5h-9a.75.75 0 0 1-.75-.75Z',
    optimize: 'M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z',
  }
  return map[id] || map.project
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <!-- Hero: Quick Generation -->
      <section class="pt-16 pb-12 text-center">
        <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium mb-6 border" style="background:rgba(139,92,246,.1);border-color:rgba(139,92,246,.2);color:var(--accent-hover)">
          <span class="w-2 h-2 rounded-full bg-[var(--green)]" style="animation:pulse 2s infinite" />
          支持 Cursor / Claude Code / GPT / Gemini / DeepSeek / Qwen
        </div>

        <h1 class="text-[clamp(28px,4.5vw,52px)] font-black leading-[1.1] tracking-tight mb-4">
          为开发者生成可直接使用的<br>
          <span class="bg-gradient-to-r from-[var(--accent)] via-[var(--cyan)] to-[var(--accent-hover)] bg-clip-text text-transparent">AI 编程 Prompt</span>
        </h1>

        <p class="text-base" style="color:var(--text-secondary)">输入你的项目、任务或 Prompt，一键生成高质量编程 Prompt</p>

        <!-- Quick Input -->
        <div class="max-w-2xl mx-auto mt-8">
          <div class="glass p-2 rounded-2xl flex items-stretch gap-2" style="border-color:var(--border)">
            <input
              v-model="quickInput"
              type="text"
              class="flex-1 bg-transparent border-none outline-none px-4 py-3 text-[15px]"
              style="color:var(--text-primary)"
              placeholder="描述你的项目、任务或想优化的 Prompt..."
              @keyup.enter="goGeneratorWithInput"
            />
            <button class="btn btn-primary px-6 py-2.5 rounded-xl text-sm whitespace-nowrap" @click="goGeneratorWithInput">
              <svg class="inline-block mr-1.5" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
              立即生成
            </button>
          </div>
          <!-- Tool Quick Select -->
          <div class="flex flex-wrap items-center justify-center gap-2 mt-4">
            <span class="text-xs" style="color:var(--text-muted)">目标工具：</span>
            <span
              v-for="tool in aiTools" :key="tool.name"
              class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium cursor-pointer transition-all border"
              :class="quickTool === tool.name ? 'selected' : ''"
              :style="quickTool === tool.name
                ? { background: 'rgba(139,92,246,.12)', borderColor: 'var(--accent)', color: 'var(--accent-hover)' }
                : { background: 'transparent', borderColor: 'var(--border)', color: 'var(--text-muted)' }"
              @click="quickTool = quickTool === tool.name ? '' : tool.name"
            >
              <span class="w-1.5 h-1.5 rounded-full" :style="{ background: tool.color }"></span>
              {{ tool.name }}
            </span>
          </div>
        </div>

        <div class="flex gap-3 justify-center flex-wrap mt-6">
          <router-link to="/prompts" class="btn btn-ghost text-[14px] px-6 py-2.5 rounded-xl no-underline flex items-center gap-1.5">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>
            浏览模板库
          </router-link>
        </div>
      </section>

      <!-- Generator Entry Cards -->
      <div class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-5 my-8">
        <div
          v-for="(g, i) in generators" :key="g.id"
          class="glass p-7 relative overflow-hidden cursor-pointer group"
          @click="goGenerator(g.id)"
        >
          <div class="absolute top-0 left-0 right-0 h-0.5 bg-gradient-to-r opacity-0 group-hover:opacity-100 transition-opacity" :class="g.gradient" />
          <div
            class="w-12 h-12 rounded-xl flex items-center justify-center mb-4"
            :style="{ background: g.bg }"
          >
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" :style="{ color: i === 0 ? 'var(--accent)' : i === 1 ? 'var(--cyan)' : i === 2 ? 'var(--amber)' : 'var(--green)' }">
              <path :d="generatorIconPath(g.id)" />
            </svg>
          </div>
          <h3 class="text-[17px] font-bold mb-2">{{ g.title }}</h3>
          <p class="text-sm" style="color:var(--text-secondary)">{{ g.desc }}</p>
          <span
            class="absolute bottom-5 right-5 transition-all"
            style="color:var(--text-muted)"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:translate-x-0.5 transition-transform"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
          </span>
        </div>
      </div>

      <!-- Hot Templates -->
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-[22px] font-bold flex items-center gap-2">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8.5 14.5A2.5 2.5 0 0 0 11 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 1 1-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 0 0 2.5 2.5z"/></svg>
          热门模板
        </h2>
        <router-link to="/prompts" class="text-sm no-underline flex items-center gap-1" style="color:var(--accent)">
          查看全部
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
        </router-link>
      </div>
      <div class="grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-4 mb-14">
        <div v-if="prompt.loading" class="col-span-full text-center py-10" style="color:var(--text-muted)">加载中...</div>
        <router-link
          v-for="t in prompt.filteredTemplates.slice(0, 4)" :key="t.id"
          :to="`/prompts/${t.slug}`"
          class="no-underline"
        >
          <PromptCard :template="t" />
        </router-link>
      </div>

      <!-- Capabilities -->
      <section class="py-14 border-t" style="border-color:var(--border)">
        <h2 class="text-[22px] font-bold mb-8 flex items-center gap-2">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>
          核心能力
        </h2>
        <div class="grid grid-cols-[repeat(auto-fit,minmax(220px,1fr))] gap-6">
          <div v-for="cap in capabilities" :key="cap.title" class="glass text-center p-6">
            <div class="w-14 h-14 rounded-xl flex items-center justify-center mx-auto mb-4 border" style="background:rgba(139,92,246,.08);border-color:rgba(139,92,246,.12)">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="color:var(--accent)">
                <path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/>
              </svg>
            </div>
            <h4 class="text-[15px] font-semibold mb-1.5">{{ cap.title }}</h4>
            <p class="text-[13px]" style="color:var(--text-secondary)">{{ cap.desc }}</p>
          </div>
        </div>
      </section>

      <!-- CTA -->
      <section class="py-20 text-center">
        <div class="glass p-14 rounded-[20px] relative overflow-hidden" style="background:linear-gradient(135deg,rgba(139,92,246,.12),rgba(6,182,212,.08));border-color:rgba(139,92,246,.2)">
          <h2 class="text-[28px] font-extrabold mb-3">解锁无限生成次数</h2>
          <p class="text-[var(--text-secondary)] mb-7 text-base">升级 Pro 会员，每日 100 次生成额度，享受更多高级功能</p>
          <router-link to="/pricing" class="btn btn-primary text-[15px] px-7 py-3 rounded-xl no-underline">查看会员方案</router-link>
        </div>
      </section>
    </div>
  </div>
</template>
