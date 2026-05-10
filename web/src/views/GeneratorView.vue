<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useGeneratorStore } from '@/store/generator'
import type { GeneratorTab } from '@/store/generator'
import { techStacks, aiTools } from '@/mock/data'
import { getProjectTypes } from '@/api/generator'
import { getPromptBySlug } from '@/api/prompts'
import type { ProjectType } from '@/types'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()
const gen = useGeneratorStore()
const projectTypes = ref<ProjectType[]>([])

const tabs = [
  { id: 'project' as GeneratorTab, label: '项目 Prompt', icon: 'folder' },
  { id: 'cursor-rules' as GeneratorTab, label: 'Cursor Rules', icon: 'cursor' },
  { id: 'claude-code' as GeneratorTab, label: 'Claude Code', icon: 'claude' },
  { id: 'optimize' as GeneratorTab, label: 'Prompt 优化', icon: 'wand' },
]

onMounted(async () => {
  // 检查路由参数中是否有 type，自动切换对应 tab
  if (route.query.type && (tabs.some(t => t.id === route.query.type) || route.query.type === 'history')) {
    if (route.query.type === 'history') {
      router.push('/generator/history')
      return
    }
    gen.setTab(route.query.type as GeneratorTab)
  }
  // 如果路由参数中有 rawPrompt，填入优化器
  if (route.query.rawPrompt) {
    gen.rawPrompt = String(route.query.rawPrompt)
  }
  // 如果路由参数中有 targetTool，填充对应字段
  if (route.query.targetTool) {
    const tool = String(route.query.targetTool)
    if (gen.activeTab === 'optimize') {
      gen.optimizeTargetTool = tool
    } else if (gen.activeTab === 'project') {
      gen.targetAiTool = tool
    }
  }
  // 如果路由参数中有 templateSlug，拉取完整模板内容
  if (route.query.templateSlug) {
    try {
      const res = await getPromptBySlug(String(route.query.templateSlug))
      const tpl = res.data.data
      if (tpl) {
        gen.projectName = tpl.title
        gen.features = tpl.content
      }
    } catch {
      ui.showToast('模板加载失败')
    }
  }
  try {
    const res = await getProjectTypes()
    projectTypes.value = res.data.data
  } catch {
    // 静默处理
  }
})

// 当 route query 变化时同步 tab
watch(() => route.query.type, (newType) => {
  if (newType && tabs.some(t => t.id === newType)) {
    gen.setTab(newType as GeneratorTab)
  }
})

function toggleTech(name: string) {
  const idx = gen.techStack.indexOf(name)
  if (idx > -1) gen.techStack.splice(idx, 1)
  else gen.techStack.push(name)
}

function selectAi(name: string) {
  gen.targetAiTool = name
}

async function handleGenerate() {
  await gen.generate()
  if (gen.result) {
    ui.showToast('生成成功', 'success')
  } else if (gen.error) {
    ui.showToast(gen.error)
  }
}

function copyResult() {
  navigator.clipboard.writeText(gen.result)
  ui.showToast('已复制到剪贴板', 'success')
}

function exportMarkdown() {
  const blob = new Blob([gen.result], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `prompt-${Date.now()}.md`
  a.click()
  URL.revokeObjectURL(url)
  ui.showToast('已导出为 Markdown 文件', 'success')
}

function regenerate() {
  handleGenerate()
}

// 获取表单是否有效（用于按钮状态）
const formValid = computed(() => {
  switch (gen.activeTab) {
    case 'project':
      return !!(gen.projectName && gen.projectType && gen.techStack.length && gen.targetAiTool && gen.features)
    case 'cursor-rules':
      return !!(gen.cursorLanguage && gen.cursorCodeStyle)
    case 'claude-code':
      return !!gen.claudeTask
    case 'optimize':
      return !!gen.rawPrompt
    default:
      return false
  }
})

// 获取 tab 图标 SVG
function tabIcon(id: string): string {
  const icons: Record<string, string> = {
    folder: 'M3 7.5A1.5 1.5 0 0 1 4.5 6h4.879a1.5 1.5 0 0 1 1.06.44l1.122 1.12A1.5 1.5 0 0 0 12.62 8H19.5A1.5 1.5 0 0 1 21 9.5V18a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18V7.5Z',
    cursor: 'M5.636 4.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728ZM12 8.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V15a.75.75 0 0 1-1.5 0v-2.25H9a.75.75 0 0 1 0-1.5h2.25V9a.75.75 0 0 1 .75-.75Z',
    claude: 'M12 2.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V9a.75.75 0 0 1-1.5 0V6.75H9a.75.75 0 0 1 0-1.5h2.25V3a.75.75 0 0 1 .75-.75ZM6.75 12a.75.75 0 0 1 .75-.75h9a.75.75 0 0 1 0 1.5h-9a.75.75 0 0 1-.75-.75ZM5.636 16.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728ZM12 15.75a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V22.5a.75.75 0 0 1-1.5 0v-2.25H9a.75.75 0 0 1 0-1.5h2.25V16.5a.75.75 0 0 1 .75-.75Z',
    wand: 'M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z',
  }
  return icons[id] || icons.folder
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="py-8 pb-16">
        <!-- 页面标题 -->
        <div class="flex items-center justify-between mb-6">
          <div>
            <h1 class="text-[28px] font-extrabold tracking-tight">生成器工作台</h1>
            <p class="mt-1 text-[15px]" style="color:var(--text-secondary)">
              选择生成类型，AI 将为你生成专业 Prompt
            </p>
          </div>
          <router-link
            to="/generator/history"
            class="hidden md:inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-[13px] font-medium no-underline transition-all"
            style="background:rgba(139,92,246,.1);color:var(--accent);border:1px solid transparent"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            历史记录
          </router-link>
        </div>

        <!-- Tab Bar -->
        <div class="flex flex-wrap gap-1 mb-8 p-1 rounded-xl" style="background:var(--bg-elevated);border:1px solid var(--border)">
          <button
            v-for="tab in tabs" :key="tab.id"
            class="flex items-center gap-2 px-5 py-2.5 rounded-lg text-[14px] font-medium transition-all cursor-pointer border-none"
            :class="gen.activeTab === tab.id ? 'active' : ''"
            :style="gen.activeTab === tab.id
              ? { background: 'var(--bg-base)', color: 'var(--accent-hover)', boxShadow: '0 1px 3px rgba(0,0,0,.15)' }
              : { background: 'transparent', color: 'var(--text-secondary)' }"
            @click="gen.setTab(tab.id)"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path :d="tabIcon(tab.icon)" />
            </svg>
            {{ tab.label }}
          </button>
          <router-link
            to="/generator/history"
            class="md:hidden flex items-center gap-2 px-5 py-2.5 rounded-lg text-[14px] font-medium no-underline transition-all"
            style="background:transparent;color:var(--text-secondary)"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            历史记录
          </router-link>
        </div>

        <!-- ========== Project Prompt Form ========== -->
        <div v-if="gen.activeTab === 'project'" class="glass p-7 mb-6">
          <h3 class="text-[17px] font-bold mb-6 flex items-center gap-2">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7.5A1.5 1.5 0 0 1 4.5 6h4.879a1.5 1.5 0 0 1 1.06.44l1.122 1.12A1.5 1.5 0 0 0 12.62 8H19.5A1.5 1.5 0 0 1 21 9.5V18a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18V7.5Z"/></svg>
            项目 Prompt 生成
          </h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-8">
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">项目名称 <span style="color:var(--rose)">*</span></label>
              <input v-model="gen.projectName" class="form-input" type="text" placeholder="例如：DevPrompt AI">
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">项目类型 <span style="color:var(--rose)">*</span></label>
              <select v-model="gen.projectType" class="form-input cursor-pointer">
                <option value="">请选择项目类型</option>
                <option v-for="pt in projectTypes" :key="pt.value" :value="pt.value">{{ pt.name }}</option>
              </select>
            </div>
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">技术栈（可多选） <span style="color:var(--rose)">*</span></label>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="t in techStacks" :key="t"
                  class="px-3.5 py-1.5 rounded-lg text-[13px] font-medium cursor-pointer transition-all border"
                  :class="gen.techStack.includes(t) ? 'selected' : ''"
                  :style="gen.techStack.includes(t)
                    ? { background: 'rgba(139,92,246,.12)', borderColor: 'var(--accent)', color: 'var(--accent-hover)' }
                    : { background: 'var(--bg-elevated)', borderColor: 'var(--border)', color: 'var(--text-secondary)' }"
                  @click="toggleTech(t)"
                >{{ t }}</span>
              </div>
            </div>
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">功能需求 <span style="color:var(--rose)">*</span></label>
              <textarea v-model="gen.features" class="form-input min-h-[100px] resize-y leading-relaxed" placeholder="描述项目的主要功能需求，每行一个..."></textarea>
            </div>
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">目标 AI 工具 <span style="color:var(--rose)">*</span></label>
              <div class="flex flex-wrap gap-2.5">
                <span
                  v-for="tool in aiTools" :key="tool.name"
                  class="px-4 py-2.5 rounded-xl text-sm font-medium cursor-pointer transition-all border flex items-center gap-2"
                  :class="gen.targetAiTool === tool.name ? 'selected' : ''"
                  :style="gen.targetAiTool === tool.name
                    ? { background: 'rgba(139,92,246,.12)', borderColor: 'var(--accent)', color: 'var(--accent-hover)', boxShadow: '0 0 16px var(--accent-glow)' }
                    : { background: 'var(--bg-elevated)', borderColor: 'var(--border)', color: 'var(--text-secondary)' }"
                  @click="selectAi(tool.name)"
                >
                  <span class="w-2 h-2 rounded-full" :style="{ background: tool.color }"></span>
                  {{ tool.name }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- ========== Cursor Rules Form ========== -->
        <div v-if="gen.activeTab === 'cursor-rules'" class="glass p-7 mb-6">
          <h3 class="text-[17px] font-bold mb-6 flex items-center gap-2">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M5.636 4.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728ZM12 8.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V15a.75.75 0 0 1-1.5 0v-2.25H9a.75.75 0 0 1 0-1.5h2.25V9a.75.75 0 0 1 .75-.75Z"/></svg>
            Cursor Rules 生成
          </h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-8">
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">编程语言 <span style="color:var(--rose)">*</span></label>
              <input v-model="gen.cursorLanguage" class="form-input" type="text" placeholder="例如：TypeScript、Python">
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">框架</label>
              <input v-model="gen.cursorFramework" class="form-input" type="text" placeholder="例如：React、Vue3、Gin">
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">代码风格 <span style="color:var(--rose)">*</span></label>
              <select v-model="gen.cursorCodeStyle" class="form-input cursor-pointer">
                <option value="">请选择代码风格</option>
                <option value="简洁实用">简洁实用</option>
                <option value="严谨规范">严谨规范</option>
                <option value="函数式">函数式</option>
                <option value="面向对象">面向对象</option>
                <option value="TypeScript 严格模式">TypeScript 严格模式</option>
              </select>
            </div>
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">规则偏好（每行一条）</label>
              <textarea v-model="gen.cursorRules" class="form-input min-h-[80px] resize-y leading-relaxed" placeholder="例如：&#10;禁止使用 any 类型&#10;组件使用函数式声明&#10;CSS 使用 Tailwind 原子类"></textarea>
            </div>
          </div>
        </div>

        <!-- ========== Claude Code Form ========== -->
        <div v-if="gen.activeTab === 'claude-code'" class="glass p-7 mb-6">
          <h3 class="text-[17px] font-bold mb-6 flex items-center gap-2">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2.25a.75.75 0 0 1 .75.75v2.25H15a.75.75 0 0 1 0 1.5h-2.25V9a.75.75 0 0 1-1.5 0V6.75H9a.75.75 0 0 1 0-1.5h2.25V3a.75.75 0 0 1 .75-.75ZM6.75 12a.75.75 0 0 1 .75-.75h9a.75.75 0 0 1 0 1.5h-9a.75.75 0 0 1-.75-.75Z"/></svg>
            Claude Code 任务 Prompt
          </h3>
          <div class="grid grid-cols-1 gap-5 mb-8">
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">任务目标 <span style="color:var(--rose)">*</span></label>
              <input v-model="gen.claudeTask" class="form-input" type="text" placeholder="例如：开发一个用户认证模块">
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">项目上下文</label>
              <textarea v-model="gen.claudeContext" class="form-input min-h-[80px] resize-y leading-relaxed" placeholder="描述当前项目的技术栈、架构、目录结构等上下文信息..."></textarea>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">具体要求（每行一条）</label>
              <textarea v-model="gen.claudeRequirements" class="form-input min-h-[80px] resize-y leading-relaxed" placeholder="例如：&#10;使用 Golang + Gin 框架&#10;遵循 RESTful API 设计&#10;包含单元测试"></textarea>
            </div>
          </div>
        </div>

        <!-- ========== Optimize Form ========== -->
        <div v-if="gen.activeTab === 'optimize'" class="glass p-7 mb-6">
          <h3 class="text-[17px] font-bold mb-6 flex items-center gap-2">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
            Prompt 优化器
          </h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-8">
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">原始 Prompt <span style="color:var(--rose)">*</span></label>
              <textarea v-model="gen.rawPrompt" class="form-input min-h-[120px] resize-y leading-relaxed" placeholder="粘贴你想要优化的原始 Prompt..."></textarea>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">优化级别</label>
              <select v-model="gen.optimizeLevel" class="form-input cursor-pointer">
                <option value="basic">基础优化</option>
                <option value="professional">专业级优化</option>
                <option value="expert">专家级优化</option>
              </select>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">目标 AI 工具</label>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="tool in aiTools" :key="tool.name"
                  class="px-3.5 py-1.5 rounded-lg text-[13px] font-medium cursor-pointer transition-all border"
                  :class="gen.optimizeTargetTool === tool.name ? 'selected' : ''"
                  :style="gen.optimizeTargetTool === tool.name
                    ? { background: 'rgba(139,92,246,.12)', borderColor: 'var(--accent)', color: 'var(--accent-hover)' }
                    : { background: 'var(--bg-elevated)', borderColor: 'var(--border)', color: 'var(--text-secondary)' }"
                  @click="gen.optimizeTargetTool = tool.name"
                >
                  <span class="w-1.5 h-1.5 rounded-full inline-block mr-1.5" :style="{ background: tool.color }"></span>
                  {{ tool.name }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Generate Button -->
        <div class="flex items-center gap-3">
          <button
            class="btn btn-primary text-[15px] px-7 py-3 rounded-xl"
            :disabled="gen.generating || !formValid"
            @click="handleGenerate"
          >
            <span v-if="gen.generating" class="spinner"></span>
            <template v-else>
              <svg class="inline-block mr-1.5" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z"/></svg>
            </template>
            {{ gen.generating ? '生成中...' : '生成 Prompt' }}
          </button>
          <button
            v-if="gen.result"
            class="btn btn-ghost text-sm px-5 py-3 rounded-xl"
            @click="gen.resetAll()"
          >
            清空结果
          </button>
        </div>

        <!-- ========== Result Section ========== -->
        <div v-if="gen.result" class="mt-8">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold flex items-center gap-2">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/></svg>
              生成结果
            </h3>
            <div class="flex items-center gap-2">
              <!-- 复制 -->
              <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg flex items-center gap-1.5" @click="copyResult">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                复制
              </button>
              <!-- 导出 Markdown（项目/Claude Code 类型） -->
              <button v-if="gen.activeTab === 'project' || gen.activeTab === 'claude-code'" class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg flex items-center gap-1.5" @click="exportMarkdown">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                导出 Markdown
              </button>
              <!-- 重新生成 -->
              <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg flex items-center gap-1.5" @click="regenerate">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
                重新生成
              </button>
            </div>
          </div>

          <!-- === Optimize: 对比展示 === -->
          <div v-if="gen.activeTab === 'optimize' && gen.originalPromptForCompare" class="mb-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="glass p-5" style="border-color:var(--border)">
                <div class="flex items-center gap-2 mb-3">
                  <span class="w-2 h-2 rounded-full" style="background:var(--text-muted)"></span>
                  <span class="text-sm font-semibold">优化前</span>
                </div>
                <pre class="text-[13px] leading-relaxed whitespace-pre-wrap break-words max-h-[300px] overflow-y-auto" style="font-family:'JetBrains Mono',monospace;color:var(--text-secondary)">{{ gen.originalPromptForCompare }}</pre>
              </div>
              <div class="glass p-5" style="border-color:var(--accent);border-width:1.5px">
                <div class="flex items-center gap-2 mb-3">
                  <span class="w-2 h-2 rounded-full" style="background:var(--accent);animation:pulse 2s infinite"></span>
                  <span class="text-sm font-semibold" style="color:var(--accent-hover)">优化后</span>
                </div>
                <pre class="text-[13px] leading-relaxed whitespace-pre-wrap break-words max-h-[300px] overflow-y-auto" style="font-family:'JetBrains Mono',monospace;color:var(--text-primary)">{{ gen.result }}</pre>
              </div>
            </div>
          </div>

          <!-- === 其他类型：普通展示 === -->
          <div v-else class="glass p-6 min-h-[200px] border relative" style="background:var(--bg-elevated);border-color:var(--border)">
            <!-- Cursor Rules 特殊标注 -->
            <div v-if="gen.activeTab === 'cursor-rules'" class="mb-3 pb-3 border-b text-xs font-medium" style="border-color:var(--border);color:var(--text-secondary)">
              <span class="flex items-center gap-1.5">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                建议保存为 <code style="background:rgba(139,92,246,.15);padding:0 6px;border-radius:4px;color:var(--accent)">.cursorrules</code> 文件放在项目根目录
              </span>
            </div>
            <!-- Claude Code 特殊标注 -->
            <div v-if="gen.activeTab === 'claude-code'" class="mb-3 pb-3 border-b text-xs font-medium" style="border-color:var(--border);color:var(--text-secondary)">
              <span class="flex items-center gap-1.5">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                可复制后直接用于 Claude Code CLI，或保存为 <code style="background:rgba(139,92,246,.15);padding:0 6px;border-radius:4px;color:var(--accent)">CLAUDE.md</code>
              </span>
            </div>
            <pre class="text-[13.5px] leading-relaxed whitespace-pre-wrap break-words" style="font-family:'JetBrains Mono',monospace;color:var(--text-primary)">{{ gen.result }}</pre>
          </div>
        </div>

        <!-- Error -->
        <div v-if="gen.error && !gen.result" class="mt-6 p-4 rounded-xl" style="background:rgba(244,63,94,.1);border:1px solid rgba(244,63,94,.2);color:var(--rose)">
          <div class="flex items-center gap-2 text-sm font-medium">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            {{ gen.error }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
