<template>
  <div>
    <h2 class="page-title">Prompt Recipe 管理</h2>
    <el-card>
      <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap;">
        <el-button type="primary" @click="openDialog()">新建 Recipe</el-button>
        <el-select v-model="filterType" placeholder="按类型筛选" clearable style="width: 150px" @change="fetchData">
          <el-option label="项目 Prompt" value="project" />
          <el-option label="Cursor Rules" value="cursor-rules" />
          <el-option label="Claude Code" value="claude-code" />
          <el-option label="Prompt 优化" value="optimize" />
        </el-select>
        <el-select v-model="filterStatus" placeholder="按状态筛选" clearable style="width: 130px" @change="fetchData">
          <el-option label="草稿" value="draft" />
          <el-option label="启用" value="active" />
          <el-option label="禁用" value="disabled" />
        </el-select>
      </div>
      <el-table :data="recipes" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ typeLabels[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="targetTool" label="目标工具" width="100" />
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column label="默认" width="70">
          <template #default="{ row }">
            <el-tag v-if="row.isDefault" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="statusTagType[row.status]" size="small">
              {{ statusLabels[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
            <el-button v-if="!row.isDefault && row.status === 'active'" size="small" type="warning" @click="handleSetDefault(row.id)">设为默认</el-button>
            <el-popconfirm title="确定删除该 Recipe？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑 Recipe' : '新建 Recipe'" width="700px" top="5vh">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="Recipe 名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="选择生成类型" style="width: 100%">
            <el-option label="项目 Prompt" value="project" />
            <el-option label="Cursor Rules" value="cursor-rules" />
            <el-option label="Claude Code" value="claude-code" />
            <el-option label="Prompt 优化" value="optimize" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标工具" prop="targetTool">
          <el-select v-model="form.targetTool" placeholder="选择目标工具" style="width: 100%">
            <el-option label="通用" value="common" />
            <el-option label="Cursor" value="cursor" />
            <el-option label="Claude Code" value="claude-code" />
            <el-option label="GPT" value="gpt" />
            <el-option label="Gemini" value="gemini" />
            <el-option label="DeepSeek" value="deepseek" />
            <el-option label="Qwen" value="qwen" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本" prop="version">
          <el-input v-model="form.version" placeholder="例如：v1、v2" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="Recipe 描述" />
        </el-form-item>
        <el-form-item label="System Prompt" prop="systemPrompt">
          <el-input v-model="form.systemPrompt" type="textarea" :rows="8" placeholder="系统 Prompt，发送给 AI 的核心指令" />
        </el-form-item>
        <el-form-item label="用户模板" prop="userTemplate">
          <el-input v-model="form.userTemplate" type="textarea" :rows="4" placeholder="用户 Prompt 模板（留空则使用默认输入序列化）" />
        </el-form-item>
        <el-form-item label="输出结构" prop="outputSchema">
          <el-input v-model="form.outputSchema" type="textarea" :rows="3" placeholder="输出 JSON Schema（可选）" />
        </el-form-item>
        <el-form-item label="质量评估" prop="qualityRubric">
          <el-input v-model="form.qualityRubric" type="textarea" :rows="3" placeholder="质量评估标准（可选）" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="草稿" value="draft" />
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="设为默认" prop="isDefault">
          <el-switch v-model="form.isDefault" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getPromptRecipes, createPromptRecipe, updatePromptRecipe, deletePromptRecipe, setDefaultPromptRecipe } from '@/api/promptRecipes'
import type { PromptRecipe } from '@/types'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const recipes = ref<PromptRecipe[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()
const filterType = ref('')
const filterStatus = ref('')

const typeLabels: Record<string, string> = {
  'project': '项目 Prompt',
  'cursor-rules': 'Cursor Rules',
  'claude-code': 'Claude Code',
  'optimize': 'Prompt 优化',
}

const statusLabels: Record<string, string> = {
  'draft': '草稿',
  'active': '启用',
  'disabled': '禁用',
}

const statusTagType: Record<string, string> = {
  'draft': 'info',
  'active': 'success',
  'disabled': 'danger',
}

const form = reactive({
  name: '',
  type: 'project' as PromptRecipe['type'],
  targetTool: 'common',
  version: 'v1',
  description: '',
  systemPrompt: '',
  userTemplate: '',
  outputSchema: '',
  qualityRubric: '',
  status: 'draft' as PromptRecipe['status'],
  isDefault: false,
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  systemPrompt: [{ required: true, message: '请输入 System Prompt', trigger: 'blur' }],
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getPromptRecipes({
      type: filterType.value || undefined,
      status: filterStatus.value || undefined,
    })
    recipes.value = res.data.data.list
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function openDialog(row?: PromptRecipe) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    Object.assign(form, {
      name: row.name,
      type: row.type,
      targetTool: row.targetTool,
      version: row.version,
      description: row.description,
      systemPrompt: row.systemPrompt,
      userTemplate: row.userTemplate,
      outputSchema: row.outputSchema,
      qualityRubric: row.qualityRubric,
      status: row.status,
      isDefault: row.isDefault,
    })
  } else {
    isEdit.value = false
    editId.value = null
    Object.assign(form, {
      name: '',
      type: 'project',
      targetTool: 'common',
      version: 'v1',
      description: '',
      systemPrompt: '',
      userTemplate: '',
      outputSchema: '',
      qualityRubric: '',
      status: 'draft',
      isDefault: false,
    })
  }
  dialogVisible.value = true
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (isEdit.value && editId.value) {
      await updatePromptRecipe(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createPromptRecipe({ ...form })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } catch {
    // ignore
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deletePromptRecipe(id)
    ElMessage.success('删除成功')
    await fetchData()
  } catch {
    // ignore
  }
}

async function handleSetDefault(id: number) {
  try {
    await setDefaultPromptRecipe(id)
    ElMessage.success('已设为默认')
    await fetchData()
  } catch {
    // ignore
  }
}

onMounted(fetchData)
</script>

<style scoped>
.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #1d2129;
}
</style>
