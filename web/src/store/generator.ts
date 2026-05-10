import { defineStore } from 'pinia'
import { ref } from 'vue'
import { generateProject } from '@/api/generator'
import type { ProjectGeneratorParams } from '@/types'

export const useGeneratorStore = defineStore('generator', () => {
  const projectName = ref('')
  const projectType = ref('')
  const techStack = ref<string[]>([])
  const features = ref('')
  const targetAiTool = ref('')
  const generating = ref(false)
  const result = ref('')
  const error = ref('')

  async function generate() {
    generating.value = true
    error.value = ''
    try {
      const params: ProjectGeneratorParams = {
        projectName: projectName.value,
        projectType: projectType.value,
        techStack: techStack.value,
        features: features.value.split('\n').filter(f => f.trim()),
        targetAiTool: targetAiTool.value,
      }
      const res = await generateProject(params)
      result.value = res.data.data.output
    } catch (e: any) {
      error.value = e.message || '生成失败，请稍后重试'
      result.value = ''
    } finally {
      generating.value = false
    }
  }

  function reset() {
    projectName.value = ''
    projectType.value = ''
    techStack.value = []
    features.value = ''
    targetAiTool.value = ''
    generating.value = false
    result.value = ''
    error.value = ''
  }

  return { projectName, projectType, techStack, features, targetAiTool, generating, result, error, generate, reset }
})
