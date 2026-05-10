<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePromptStore } from '@/store/prompt'
import PromptCard from '@/components/PromptCard.vue'

const router = useRouter()
const prompt = usePromptStore()

onMounted(async () => {
  await prompt.fetchTemplates()
})

const generators = [
  { icon: '🏗️', title: '项目 Prompt 生成', desc: '输入项目信息，自动生成完整的项目开发 Prompt', page: '/generator' },
  { icon: '📐', title: 'Cursor Rules 生成', desc: '根据语言、框架和代码风格，一键生成 .cursorrules 配置', page: '/generator' },
  { icon: '🤖', title: 'Claude Code Prompt', desc: '针对 Claude Code 场景，生成结构化的任务 Prompt', page: '/generator' },
  { icon: '⚡', title: 'Prompt 优化器', desc: '将模糊的 Prompt 优化为专业、精确的版本', page: '/generator' },
]

const capabilities = [
  { icon: '🎯', title: '多工具适配', desc: '支持 Cursor、Claude Code、GPT、Gemini、DeepSeek、Qwen' },
  { icon: '📚', title: '模板库', desc: '精选高质量 Prompt 模板，覆盖前端、后端、DevOps 等多个领域' },
  { icon: '🔄', title: '一键优化', desc: '将粗糙的 Prompt 优化为专业版本，提升 AI 输出质量' },
  { icon: '📋', title: '一键复制', desc: '生成结果支持一键复制和 Markdown 渲染，即拿即用' },
]
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">

      <!-- Hero -->
      <section class="pt-20 pb-14 text-center">
        <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium mb-6 border" style="background:rgba(139,92,246,.1);border-color:rgba(139,92,246,.2);color:var(--accent-hover)">
          <span class="w-2 h-2 rounded-full bg-[var(--green)]" style="animation:pulse 2s infinite" />
          支持 Cursor / Claude Code / GPT / Gemini / DeepSeek / Qwen
        </div>
        <h1 class="text-[clamp(36px,5.5vw,64px)] font-black leading-[1.1] tracking-tight mb-5">
          为开发者打造的<br>
          <span class="bg-gradient-to-r from-[var(--accent)] via-[var(--cyan)] to-[var(--accent-hover)] bg-clip-text text-transparent">AI Prompt 生成平台</span>
        </h1>
        <p class="text-lg text-[var(--text-secondary)] max-w-lg mx-auto mb-9 leading-relaxed">
          一键生成高质量编程 Prompt，支持 Cursor Rules、Claude Code、项目 Prompt 和 Prompt 优化，让 AI 编程效率提升 10 倍
        </p>
        <div class="flex gap-3 justify-center flex-wrap">
          <router-link to="/generator" class="btn btn-primary text-[15px] px-7 py-3 rounded-xl no-underline">✦ 开始生成</router-link>
          <router-link to="/prompts" class="btn btn-ghost text-[15px] px-7 py-3 rounded-xl no-underline">浏览模板库 →</router-link>
        </div>
      </section>

      <!-- Generator Entry Cards -->
      <div class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-5 my-14">
        <div
          v-for="(g, i) in generators" :key="i"
          class="glass p-7 relative overflow-hidden cursor-pointer group"
          @click="router.push(g.page)"
        >
          <div class="absolute top-0 left-0 right-0 h-0.5 bg-gradient-to-r from-[var(--accent)] to-[var(--cyan)] opacity-0 group-hover:opacity-100 transition-opacity" />
          <div
            class="w-12 h-12 rounded-xl flex items-center justify-center text-[22px] mb-4"
            :class="[
              i === 0 ? 'bg-[rgba(139,92,246,.15)]' : '',
              i === 1 ? 'bg-[rgba(6,182,212,.15)]' : '',
              i === 2 ? 'bg-[rgba(245,158,11,.15)]' : '',
              i === 3 ? 'bg-[rgba(16,185,129,.15)]' : '',
            ]"
          >{{ g.icon }}</div>
          <h3 class="text-[17px] font-bold mb-2">{{ g.title }}</h3>
          <p class="text-sm text-[var(--text-secondary)] leading-relaxed">{{ g.desc }}</p>
          <span class="absolute bottom-5 right-5 text-lg text-[var(--text-muted)] group-hover:text-[var(--accent)] group-hover:translate-x-1 transition-all">→</span>
        </div>
      </div>

      <!-- Hot Templates -->
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-[22px] font-bold">🔥 热门模板</h2>
        <router-link to="/prompts" class="text-sm text-[var(--accent)] no-underline flex items-center gap-1 hover:text-[var(--accent-hover)]">查看全部 →</router-link>
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
        <h2 class="text-[22px] font-bold mb-8">核心能力</h2>
        <div class="grid grid-cols-[repeat(auto-fit,minmax(220px,1fr))] gap-6">
          <div v-for="cap in capabilities" :key="cap.title" class="glass text-center p-6">
            <div class="w-14 h-14 rounded-xl flex items-center justify-center text-2xl mx-auto mb-4 border" style="background:rgba(139,92,246,.08);border-color:rgba(139,92,246,.12)">{{ cap.icon }}</div>
            <h4 class="text-[15px] font-semibold mb-1.5">{{ cap.title }}</h4>
            <p class="text-[13px] text-[var(--text-secondary)]">{{ cap.desc }}</p>
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
