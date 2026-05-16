<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useGeneratorStore, type GeneratorTab } from '@/store/generator'
import { techStacks, aiTools } from '@/mock/data'
import { getProjectTypes } from '@/api/generator'
import { getPromptBySlug } from '@/api/prompts'
import type { ProjectType } from '@/types'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosSegmentedControl from '../components/IosSegmentedControl.vue'
import IosIcon from '../components/IosIcon.vue'
import IosSelect from '../components/IosSelect.vue'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()
const gen = useGeneratorStore()
const projectTypes = ref<ProjectType[]>([])

const tabs = [
  { value: 'project' as GeneratorTab, label: '项目' },
  { value: 'cursor-rules' as GeneratorTab, label: 'Cursor' },
  { value: 'claude-code' as GeneratorTab, label: 'Claude' },
  { value: 'optimize' as GeneratorTab, label: '优化' },
]

const optimizeLevels = [
  { value: 'basic', label: '基础' },
  { value: 'professional', label: '专业' },
  { value: 'expert', label: '专家' },
]

const codeStyles = ['简洁实用', '严谨规范', '函数式', '面向对象', 'TypeScript 严格模式']

const qualityModes = [
  { value: 'concise', label: '简洁' },
  { value: 'standard', label: '标准' },
  { value: 'expert', label: '专家' },
]

const outputFormats = [
  { value: 'markdown', label: 'Markdown' },
  { value: 'checklist', label: '清单' },
  { value: 'json', label: 'JSON' },
  { value: 'plain', label: '纯 Prompt' },
]

onMounted(async () => {
  if (route.query.type && tabs.some(t => t.value === route.query.type)) {
    gen.setTab(route.query.type as GeneratorTab)
  } else if (route.query.type === 'history') {
    router.push('/ios26/generator/history')
    return
  }
  if (route.query.rawPrompt) {
    gen.rawPrompt = String(route.query.rawPrompt)
  }
  if (route.query.targetTool) {
    const tool = String(route.query.targetTool)
    if (gen.activeTab === 'optimize') gen.optimizeTargetTool = tool
    else if (gen.activeTab === 'project') gen.targetAiTool = tool
  }
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
    /* silent */
  }
})

watch(() => route.query.type, (newType: unknown) => {
  if (newType && tabs.some(t => t.value === newType)) {
    gen.setTab(newType as GeneratorTab)
  }
})

function setTab(value: string) {
  gen.setTab(value as GeneratorTab)
}

function toggleTech(name: string) {
  const idx = gen.techStack.indexOf(name)
  if (idx > -1) gen.techStack.splice(idx, 1)
  else gen.techStack.push(name)
}

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

async function handleGenerate() {
  await gen.generate()
  if (gen.result) ui.showToast('生成成功', 'success')
  else if (gen.error) ui.showToast(gen.error)
}

function copyResult() {
  navigator.clipboard.writeText(gen.result)
  ui.showToast('已复制', 'success')
}

function exportMarkdown() {
  const blob = new Blob([gen.result], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `prompt-${Date.now()}.md`
  a.click()
  URL.revokeObjectURL(url)
  ui.showToast('已导出', 'success')
}
</script>

<template>
  <div class="ios-gen">
    <IosNavBar large>
      <template #trailing>
        <IosButton variant="tinted" size="sm" to="/ios26/generator/history">
          <IosIcon path="M3 12a9 9 0 1 0 18 0 9 9 0 0 0-18 0Zm9-5v5l3 2" :size="14" :stroke="2" />
          历史记录
        </IosButton>
      </template>
    </IosNavBar>

    <div class="ios-page">
      <!-- Tabs -->
      <IosSegmentedControl
        :model-value="gen.activeTab"
        :options="tabs"
        size="md"
        class="ios-gen__tabs"
        @update:model-value="setTab"
      />

      <!-- Two-column workspace: form (left) + result (right) -->
      <div class="ios-gen__workspace">
        <div class="ios-gen__left">
      <!-- ============ Project ============ -->
      <IosGlassPanel v-if="gen.activeTab === 'project'" size="md" class="ios-gen__panel">
        <h3 class="ios-text-headline ios-gen__panel-title">项目 Prompt</h3>

        <div class="ios-field">
          <label class="ios-field__label">项目名称 <span class="ios-field__req">*</span></label>
          <input v-model="gen.projectName" class="ios-input" type="text" placeholder="例如：DevPrompt AI">
        </div>

        <div class="ios-field">
          <label class="ios-field__label">项目类型 <span class="ios-field__req">*</span></label>
          <IosSelect
            v-model="gen.projectType"
            :options="projectTypes.map(pt => ({ value: pt.value, label: pt.name }))"
            placeholder="请选择项目类型"
          />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">技术栈 <span class="ios-field__req">*</span></label>
          <div class="ios-chips">
            <button
              v-for="t in techStacks"
              :key="t"
              class="ios-chip ios-press"
              :class="{ 'is-active': gen.techStack.includes(t) }"
              @click="toggleTech(t)"
              type="button"
            >{{ t }}</button>
          </div>
        </div>

        <div class="ios-field">
          <label class="ios-field__label">功能需求 <span class="ios-field__req">*</span></label>
          <textarea v-model="gen.features" class="ios-input" rows="5" placeholder="描述项目的主要功能需求，每行一个..." />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">目标 AI 工具 <span class="ios-field__req">*</span></label>
          <div class="ios-chips">
            <button
              v-for="tool in aiTools"
              :key="tool.name"
              class="ios-chip ios-press"
              :class="{ 'is-active': gen.targetAiTool === tool.name }"
              @click="gen.targetAiTool = tool.name"
              type="button"
            >
              <span class="ios-chip__dot" :style="{ background: tool.color }" />
              {{ tool.name }}
            </button>
          </div>
        </div>
      </IosGlassPanel>

      <!-- ============ Cursor Rules ============ -->
      <IosGlassPanel v-else-if="gen.activeTab === 'cursor-rules'" size="md" class="ios-gen__panel">
        <h3 class="ios-text-headline ios-gen__panel-title">Cursor Rules</h3>

        <div class="ios-field">
          <label class="ios-field__label">编程语言 <span class="ios-field__req">*</span></label>
          <input v-model="gen.cursorLanguage" class="ios-input" type="text" placeholder="例如：TypeScript">
        </div>

        <div class="ios-field">
          <label class="ios-field__label">框架</label>
          <input v-model="gen.cursorFramework" class="ios-input" type="text" placeholder="例如：React、Vue3">
        </div>

        <div class="ios-field">
          <label class="ios-field__label">代码风格 <span class="ios-field__req">*</span></label>
          <IosSelect
            v-model="gen.cursorCodeStyle"
            :options="codeStyles.map(s => ({ value: s, label: s }))"
            placeholder="请选择代码风格"
          />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">规则偏好（每行一条）</label>
          <textarea v-model="gen.cursorRules" class="ios-input" rows="4" placeholder="例如：&#10;禁止使用 any 类型&#10;组件使用函数式声明" />
        </div>
      </IosGlassPanel>

      <!-- ============ Claude Code ============ -->
      <IosGlassPanel v-else-if="gen.activeTab === 'claude-code'" size="md" class="ios-gen__panel">
        <h3 class="ios-text-headline ios-gen__panel-title">Claude Code</h3>

        <div class="ios-field">
          <label class="ios-field__label">任务目标 <span class="ios-field__req">*</span></label>
          <input v-model="gen.claudeTask" class="ios-input" type="text" placeholder="例如：开发一个用户认证模块">
        </div>

        <div class="ios-field">
          <label class="ios-field__label">项目上下文</label>
          <textarea v-model="gen.claudeContext" class="ios-input" rows="4" placeholder="描述当前项目的技术栈、架构、目录结构等..." />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">具体要求（每行一条）</label>
          <textarea v-model="gen.claudeRequirements" class="ios-input" rows="4" placeholder="例如：&#10;使用 Golang + Gin&#10;遵循 RESTful API 设计" />
        </div>
      </IosGlassPanel>

      <!-- ============ Optimize ============ -->
      <IosGlassPanel v-else-if="gen.activeTab === 'optimize'" size="md" class="ios-gen__panel">
        <h3 class="ios-text-headline ios-gen__panel-title">Prompt 优化</h3>

        <div class="ios-field">
          <label class="ios-field__label">原始 Prompt <span class="ios-field__req">*</span></label>
          <textarea v-model="gen.rawPrompt" class="ios-input" rows="6" placeholder="粘贴你想要优化的原始 Prompt..." />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">优化级别</label>
          <IosSegmentedControl
            :model-value="gen.optimizeLevel"
            :options="optimizeLevels"
            size="md"
            @update:model-value="(v: string) => gen.optimizeLevel = v as 'basic' | 'professional' | 'expert'"
          />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">目标 AI 工具</label>
          <div class="ios-chips">
            <button
              v-for="tool in aiTools"
              :key="tool.name"
              class="ios-chip ios-press"
              :class="{ 'is-active': gen.optimizeTargetTool === tool.name }"
              @click="gen.optimizeTargetTool = tool.name === gen.optimizeTargetTool ? '' : tool.name"
              type="button"
            >
              <span class="ios-chip__dot" :style="{ background: tool.color }" />
              {{ tool.name }}
            </button>
          </div>
        </div>
      </IosGlassPanel>

      <!-- ============ Quality Options ============ -->
      <IosGlassPanel size="md" class="ios-gen__panel">
        <h3 class="ios-text-headline ios-gen__panel-title">生成质量</h3>

        <div class="ios-field">
          <label class="ios-field__label">详细程度</label>
          <IosSegmentedControl
            :model-value="gen.qualityMode"
            :options="qualityModes"
            size="sm"
            @update:model-value="(v: string) => gen.qualityMode = v as 'concise' | 'standard' | 'expert'"
          />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">输出格式</label>
          <IosSegmentedControl
            :model-value="gen.outputFormat"
            :options="outputFormats"
            size="sm"
            @update:model-value="(v: string) => gen.outputFormat = v as 'markdown' | 'checklist' | 'json' | 'plain'"
          />
        </div>

        <div class="ios-field">
          <label class="ios-field__label">高级选项</label>
          <div class="ios-toggle-group">
            <label class="ios-toggle">
              <input type="checkbox" v-model="gen.includeAcceptanceCriteria">
              <span>验收标准</span>
            </label>
            <label class="ios-toggle">
              <input type="checkbox" v-model="gen.includeRiskCheck">
              <span>风险检查</span>
            </label>
            <label class="ios-toggle">
              <input type="checkbox" v-model="gen.includeTestPlan">
              <span>测试建议</span>
            </label>
            <label class="ios-toggle">
              <input type="checkbox" v-model="gen.includeDeploymentNotes">
              <span>部署建议</span>
            </label>
          </div>
        </div>
      </IosGlassPanel>

      <!-- Generate CTA (inside left column) -->
      <div class="ios-gen__cta">
        <IosButton
          variant="filled"
          size="lg"
          block
          :loading="gen.generating"
          :disabled="!formValid || gen.generating"
          @click="handleGenerate"
        >
          <IosIcon v-if="!gen.generating" path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z" :size="18" :stroke="2" />
          {{ gen.generating ? '生成中...' : '生成 Prompt' }}
        </IosButton>
        <IosButton
          v-if="gen.result"
          variant="plain"
          size="md"
          @click="gen.resetAll()"
        >清空结果</IosButton>
      </div>
        </div><!-- /.ios-gen__left -->

        <!-- Right column: result panel -->
        <div class="ios-gen__right">
          <!-- Empty state placeholder -->
          <div v-if="!gen.result && !gen.error && !gen.generating" class="ios-gen__right-empty">
            <IosIcon path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z" :size="36" :stroke="1.4" />
            <p>填写左侧表单，点击生成</p>
          </div>

          <!-- Generating spinner -->
          <div v-if="gen.generating" class="ios-gen__right-empty">
            <IosIcon path="M21 12a9 9 0 1 1-6.219-8.56" :size="36" :stroke="2" class="ios-spin" />
            <p>正在生成...</p>
          </div>

          <!-- ====== Result ====== -->
          <section v-if="gen.result" class="ios-gen__result">
            <header class="ios-gen__result-head">
              <h3 class="ios-text-title-3">生成结果</h3>
              <div class="ios-gen__result-actions">
                <IosButton variant="tinted" size="sm" @click="copyResult">
                  <IosIcon path="M9 9h13v13H9z M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" :size="14" :stroke="2" />
                  复制
                </IosButton>
                <IosButton
                  v-if="gen.activeTab === 'project' || gen.activeTab === 'claude-code'"
                  variant="tinted"
                  size="sm"
                  @click="exportMarkdown"
                >
                  <IosIcon path="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4 M7 10l5 5 5-5 M12 15V3" :size="14" :stroke="2" />
                  导出 MD
                </IosButton>
                <IosButton variant="tinted" size="sm" @click="handleGenerate">
                  <IosIcon path="M23 4v6h-6 M20.49 15a9 9 0 1 1-2.12-9.36L23 10" :size="14" :stroke="2" />
                  重新生成
                </IosButton>
              </div>
            </header>

            <!-- Optimize comparison -->
            <div v-if="gen.activeTab === 'optimize' && gen.originalPromptForCompare" class="ios-gen__compare">
              <IosGlassPanel size="md" class="ios-gen__compare-card">
                <div class="ios-gen__compare-head">
                  <span class="ios-gen__compare-dot is-before" />
                  <span class="ios-text-subheadline-emph">优化前</span>
                </div>
                <pre class="ios-mono ios-gen__compare-pre">{{ gen.originalPromptForCompare }}</pre>
              </IosGlassPanel>
              <IosGlassPanel size="md" class="ios-gen__compare-card is-after">
                <div class="ios-gen__compare-head">
                  <span class="ios-gen__compare-dot is-after" />
                  <span class="ios-text-subheadline-emph" style="color: var(--ios-color-tint)">优化后</span>
                </div>
                <pre class="ios-mono ios-gen__compare-pre">{{ gen.result }}</pre>
              </IosGlassPanel>
            </div>

            <!-- Single result -->
            <IosGlassPanel v-else size="md" class="ios-gen__result-panel">
              <div v-if="gen.activeTab === 'cursor-rules'" class="ios-gen__hint">
                <IosIcon path="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2v6h6" :size="13" :stroke="2" />
                建议保存为 <code class="ios-gen__code">.cursorrules</code> 放在项目根目录
              </div>
              <div v-if="gen.activeTab === 'claude-code'" class="ios-gen__hint">
                <IosIcon path="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2v6h6" :size="13" :stroke="2" />
                可直接用于 Claude Code CLI，或保存为 <code class="ios-gen__code">CLAUDE.md</code>
              </div>
              <pre class="ios-mono ios-gen__result-pre">{{ gen.result }}</pre>
            </IosGlassPanel>
          </section>

          <!-- Error -->
          <div v-if="gen.error && !gen.result" class="ios-gen__error">
            <IosIcon path="M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20Zm-3-13l6 6m0-6-6 6" :size="18" :stroke="2" />
            {{ gen.error }}
          </div>
        </div><!-- /.ios-gen__right -->
      </div><!-- /.ios-gen__workspace -->
    </div>
  </div>
</template>

<style scoped>
.ios-gen {
  display: flex;
  flex-direction: column;
}

.ios-gen__tabs {
  margin-bottom: 20px;
}

/* ===== Two-column workspace (desktop) ===== */
.ios-gen__workspace {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 20px;
  align-items: stretch;
}

.ios-gen__left {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.ios-gen__right {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 0;
}

/* Empty / loading state in right column */
.ios-gen__right-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  border-radius: var(--ios-radius-xl);
  border: 1.5px dashed var(--ios-color-separator);
  color: var(--ios-color-label-tertiary);
  font-size: 14px;
}

.ios-spin {
  animation: ios-spin 1s linear infinite;
}

@keyframes ios-spin {
  to { transform: rotate(360deg); }
}

/* On mobile / narrow: single column */
@media (max-width: 900px) {
  .ios-gen__workspace {
    grid-template-columns: 1fr;
    align-items: start;
  }

  .ios-gen__right-empty {
    display: none;
  }
}

.ios-gen__panel {
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ios-gen__panel-title {
  margin: 0 0 4px;
  color: var(--ios-color-label-primary);
}

.ios-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ios-field__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--ios-color-label-secondary);
  letter-spacing: -0.08px;
}

.ios-field__req {
  color: var(--ios-color-systemRed);
}

.ios-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.ios-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 12px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-color-fill-quaternary);
  border: 1px solid transparent;
  color: var(--ios-color-label-primary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
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

.ios-gen__cta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
  margin-bottom: 24px;
}

.ios-gen__cta > * {
  width: 100%;
  max-width: 420px;
}

.ios-gen__result {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  margin-top: 8px;
}

.ios-gen__result-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.ios-gen__result-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.ios-gen__result-head h3 {
  margin: 0;
  color: var(--ios-color-label-primary);
}

.ios-gen__result-actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.ios-gen__compare {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 12px;
}

.ios-gen__compare-card.is-after {
  border-color: var(--ios-color-tint);
  box-shadow: 0 18px 40px var(--ios-color-tint-soft-strong);
}

.ios-gen__compare-head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
}

.ios-gen__compare-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--ios-color-label-tertiary);
}

.ios-gen__compare-dot.is-after {
  background: var(--ios-color-tint);
  animation: ios-pulse 2.4s ease-in-out infinite;
}

.ios-gen__compare-pre {
  margin: 0;
  font-size: 12.5px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 320px;
  overflow-y: auto;
  color: var(--ios-color-label-primary);
}

.ios-gen__hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 12px;
}

.ios-gen__code {
  background: var(--ios-color-fill-quaternary);
  padding: 1px 6px;
  border-radius: 4px;
  font-family: 'SF Mono', monospace;
  font-size: 11px;
  color: var(--ios-color-tint);
}

.ios-gen__result-pre {
  flex: 1;
  margin: 0;
  font-size: 13px;
  line-height: 1.65;
  white-space: pre-wrap;
  word-break: break-word;
  color: var(--ios-color-label-primary);
  overflow-y: auto;
  min-height: 0;
}

.ios-gen__error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
  border-radius: var(--ios-radius-md);
  background: rgba(255, 59, 48, 0.1);
  border: 1px solid rgba(255, 59, 48, 0.25);
  color: var(--ios-color-systemRed);
  font-size: 14px;
  font-weight: 500;
  margin-top: 12px;
}

.ios-toggle-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}

.ios-toggle {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--ios-color-label-primary);
  cursor: pointer;
  user-select: none;
}

.ios-toggle input[type="checkbox"] {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  accent-color: var(--ios-color-tint);
  cursor: pointer;
}
</style>
