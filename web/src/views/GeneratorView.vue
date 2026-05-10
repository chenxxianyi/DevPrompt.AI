<script setup lang="ts">
import { useUiStore } from '@/store/ui'
import { useGeneratorStore } from '@/store/generator'
import { techStacks, aiTools } from '@/mock/data'

const ui = useUiStore()
const gen = useGeneratorStore()

function toggleTech(name: string) {
  const idx = gen.techStack.indexOf(name)
  if (idx > -1) gen.techStack.splice(idx, 1)
  else gen.techStack.push(name)
}

function selectAi(name: string) {
  gen.targetAiTool = name
}

async function generate() {
  if (!gen.projectName || !gen.projectType || !gen.techStack.length || !gen.targetAiTool || !gen.features) {
    ui.showToast('请填写所有必填项')
    return
  }

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
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="py-8 pb-16">
        <h1 class="text-[28px] font-extrabold mb-2">🏗️ 项目 Prompt 生成器</h1>
        <p class="mb-8 text-[15px]" style="color:var(--text-secondary)">填写项目信息，AI 将为你生成专业的开发 Prompt</p>

        <!-- Form -->
        <div class="glass p-7 mb-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-8">
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">项目名称</label>
              <input v-model="gen.projectName" class="form-input" type="text" placeholder="例如：DevPrompt AI">
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">项目类型</label>
              <select v-model="gen.projectType" class="form-input cursor-pointer">
                <option value="">请选择项目类型</option>
                <option value="SaaS">SaaS</option>
                <option value="电商">电商平台</option>
                <option value="社交">社交应用</option>
                <option value="工具">工具类</option>
                <option value="企业内部">企业内部系统</option>
              </select>
            </div>
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">技术栈（可多选）</label>
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
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">功能需求</label>
              <textarea v-model="gen.features" class="form-input min-h-[100px] resize-y leading-relaxed" placeholder="描述项目的主要功能需求，每行一个..."></textarea>
            </div>
            <div class="flex flex-col gap-1.5 md:col-span-2">
              <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">目标 AI 工具</label>
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
          <button class="btn btn-primary text-[15px] px-7 py-3 rounded-xl" :disabled="gen.generating" @click="generate">
            <span v-if="gen.generating" class="spinner"></span>
            {{ gen.generating ? '生成中...' : '✦ 生成 Prompt' }}
          </button>
        </div>

        <!-- Result -->
        <div v-if="gen.result" class="mt-8">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold">✨ 生成结果</h3>
            <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg" @click="copyResult">📋 复制结果</button>
          </div>
          <div class="glass p-6 min-h-[200px] border relative" style="background:var(--bg-elevated)">
            <pre class="text-[13.5px] leading-relaxed whitespace-pre-wrap break-words" style="font-family:'JetBrains Mono',monospace;color:var(--text-primary)">{{ gen.result }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
