import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  generateProject,
  generateCursorRules,
  generateClaudeCode,
  optimizePrompt,
} from '@/api/generator'
import type {
  ProjectGeneratorParams,
  CursorRulesParams,
  ClaudeCodeParams,
  OptimizeParams,
} from '@/types'

export type GeneratorTab = 'project' | 'cursor-rules' | 'claude-code' | 'optimize'

export const useGeneratorStore = defineStore('generator', () => {
  // Shared state
  const activeTab = ref<GeneratorTab>('project')
  const generating = ref(false)
  const result = ref('')
  const error = ref('')
  const lastType = ref<GeneratorTab>('project')

  // --- Project Generator ---
  const projectName = ref('')
  const projectType = ref('')
  const techStack = ref<string[]>([])
  const features = ref('')
  const targetAiTool = ref('')

  // --- Cursor Rules ---
  const cursorLanguage = ref('')
  const cursorFramework = ref('')
  const cursorCodeStyle = ref('')
  const cursorRules = ref('')

  // --- Claude Code ---
  const claudeTask = ref('')
  const claudeContext = ref('')
  const claudeRequirements = ref('')

  // --- Optimize ---
  const rawPrompt = ref('')
  const optimizeLevel = ref<'basic' | 'professional' | 'expert'>('professional')
  const optimizeTargetTool = ref('')
  // Store the original prompt for comparison display
  const originalPromptForCompare = ref('')
  // Flag: features content came from a template (skip line splitting)
  const templateLoaded = ref(false)

  function setTab(tab: GeneratorTab) {
    activeTab.value = tab
    result.value = ''
    error.value = ''
    originalPromptForCompare.value = ''
  }

  async function generate() {
    generating.value = true
    error.value = ''
    result.value = ''
    originalPromptForCompare.value = ''

    try {
      const tab = activeTab.value
      lastType.value = tab

      switch (tab) {
        case 'project': {
          if (!projectName.value || !projectType.value || !techStack.value.length || !targetAiTool.value || !features.value) {
            throw new Error('请填写所有必填项')
          }
          const params: ProjectGeneratorParams = {
            projectName: projectName.value,
            projectType: projectType.value,
            techStack: techStack.value,
            features: templateLoaded.value
              ? [features.value]
              : features.value.split('\n').filter((f: string) => f.trim()),
            targetAiTool: targetAiTool.value,
          }
          templateLoaded.value = false
          const projectRes = await generateProject(params)
          result.value = projectRes.data.data.output
          break
        }
        case 'cursor-rules': {
          if (!cursorLanguage.value || !cursorCodeStyle.value) {
            throw new Error('请填写语言和代码风格')
          }
          const params: CursorRulesParams = {
            language: cursorLanguage.value,
            framework: cursorFramework.value,
            codeStyle: cursorCodeStyle.value,
            rules: cursorRules.value.split('\n').filter((r: string) => r.trim()),
          }
          const cursorRes = await generateCursorRules(params)
          result.value = cursorRes.data.data.output
          break
        }
        case 'claude-code': {
          if (!claudeTask.value) {
            throw new Error('请填写任务目标')
          }
          const params: ClaudeCodeParams = {
            task: claudeTask.value,
            context: claudeContext.value,
            requirements: claudeRequirements.value.split('\n').filter((r: string) => r.trim()),
          }
          const claudeRes = await generateClaudeCode(params)
          result.value = claudeRes.data.data.output
          break
        }
        case 'optimize': {
          if (!rawPrompt.value) {
            throw new Error('请填写需要优化的 Prompt')
          }
          const params: OptimizeParams = {
            rawPrompt: rawPrompt.value,
            targetTool: optimizeTargetTool.value,
            optimizeLevel: optimizeLevel.value,
          }
          originalPromptForCompare.value = rawPrompt.value
          const optimizeRes = await optimizePrompt(params)
          result.value = optimizeRes.data.data.output
          break
        }
      }
    } catch (e: any) {
      error.value = e?.response?.data?.message || e?.message || '生成失败，请稍后重试'
      result.value = ''
    } finally {
      generating.value = false
    }
  }

  function resetAll() {
    // Project
    projectName.value = ''
    projectType.value = ''
    techStack.value = []
    features.value = ''
    targetAiTool.value = ''
    // Cursor
    cursorLanguage.value = ''
    cursorFramework.value = ''
    cursorCodeStyle.value = ''
    cursorRules.value = ''
    // Claude
    claudeTask.value = ''
    claudeContext.value = ''
    claudeRequirements.value = ''
    // Optimize
    rawPrompt.value = ''
    optimizeLevel.value = 'professional'
    optimizeTargetTool.value = ''
    originalPromptForCompare.value = ''
    templateLoaded.value = false
    // Shared
    result.value = ''
    error.value = ''
    generating.value = false
  }

  function parseHistoryInput(input: string | Record<string, any>): Record<string, any> {
    if (typeof input !== 'string') return input
    try {
      return JSON.parse(input)
    } catch {
      return {}
    }
  }

  function loadFromHistory(item: { type: GeneratorTab; input: string | Record<string, any> }) {
    setTab(item.type)
    const input = parseHistoryInput(item.input)
    switch (item.type) {
      case 'project':
        projectName.value = input.projectName || ''
        projectType.value = input.projectType || ''
        techStack.value = input.techStack || []
        features.value = Array.isArray(input.features) ? input.features.join('\n') : (input.features || '')
        targetAiTool.value = input.targetAiTool || ''
        break
      case 'cursor-rules':
        cursorLanguage.value = input.language || ''
        cursorFramework.value = input.framework || ''
        cursorCodeStyle.value = input.codeStyle || ''
        cursorRules.value = Array.isArray(input.rules) ? input.rules.join('\n') : (input.rules || '')
        break
      case 'claude-code':
        claudeTask.value = input.task || ''
        claudeContext.value = input.context || ''
        claudeRequirements.value = Array.isArray(input.requirements) ? input.requirements.join('\n') : (input.requirements || '')
        break
      case 'optimize':
        rawPrompt.value = input.rawPrompt || ''
        optimizeLevel.value = input.optimizeLevel || 'professional'
        optimizeTargetTool.value = input.targetTool || ''
        break
    }
  }

  function loadFromTemplate(template: { title: string; content: string }) {
    // 当从模板带入时，将 template.content 作为完整正文发送，不按行拆分
    setTab('project')
    projectName.value = template.title
    features.value = template.content
    templateLoaded.value = true
  }

  function loadRawPrompt(text: string) {
    setTab('optimize')
    rawPrompt.value = text
  }

  return {
    activeTab, generating, result, error, lastType, originalPromptForCompare,
    // Project
    projectName, projectType, techStack, features, targetAiTool,
    // Cursor
    cursorLanguage, cursorFramework, cursorCodeStyle, cursorRules,
    // Claude
    claudeTask, claudeContext, claudeRequirements,
    // Optimize
    rawPrompt, optimizeLevel, optimizeTargetTool,
    // Methods
    setTab, generate, resetAll, loadFromHistory, loadFromTemplate, loadRawPrompt,
  }
})
