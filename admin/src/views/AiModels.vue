<template>
  <div>
    <h2 class="page-title">AI 模型管理</h2>
    <el-card>
      <div style="margin-bottom: 16px;">
        <el-button type="primary" @click="openDialog()">新增模型</el-button>
      </div>
      <el-table :data="models" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="displayName" label="显示名称" width="160" />
        <el-table-column prop="provider" label="Provider" width="100" />
        <el-table-column prop="modelName" label="模型标识" min-width="200" />
        <el-table-column prop="priority" label="优先级" width="80" />
        <el-table-column prop="timeoutSeconds" label="超时(s)" width="80" />
        <el-table-column label="默认" width="70">
          <template #default="{ row }">
            <el-tag v-if="row.isDefault" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除该模型？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑模型' : '新增模型'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="显示名称" prop="displayName">
          <el-input v-model="form.displayName" />
        </el-form-item>
        <el-form-item label="Provider" prop="provider">
          <el-input v-model="form.provider" />
        </el-form-item>
        <el-form-item label="模型标识" prop="modelName">
          <el-input v-model="form.modelName" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="form.priority" :min="0" />
        </el-form-item>
        <el-form-item label="超时(s)" prop="timeoutSeconds">
          <el-input-number v-model="form.timeoutSeconds" :min="1" :max="300" />
        </el-form-item>
        <el-form-item label="默认模型" prop="isDefault">
          <el-switch v-model="form.isDefault" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="form.status" active-value="active" inactive-value="disabled" />
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
import { getAiModels, createAiModel, updateAiModel, deleteAiModel } from '@/api/aiModels'
import type { AIModel } from '@/types'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const models = ref<AIModel[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()

const form = reactive({
  displayName: '',
  provider: '',
  modelName: '',
  priority: 0,
  timeoutSeconds: 60,
  isDefault: false,
  status: 'active' as 'active' | 'disabled',
})

const rules = {
  displayName: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  provider: [{ required: true, message: '请输入 Provider', trigger: 'blur' }],
  modelName: [{ required: true, message: '请输入模型标识', trigger: 'blur' }],
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getAiModels()
    models.value = res.data.data.list
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function openDialog(row?: AIModel) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.displayName = row.displayName
    form.provider = row.provider
    form.modelName = row.modelName
    form.priority = row.priority
    form.timeoutSeconds = row.timeoutSeconds
    form.isDefault = row.isDefault
    form.status = row.status
  } else {
    isEdit.value = false
    editId.value = null
    form.displayName = ''
    form.provider = ''
    form.modelName = ''
    form.priority = 0
    form.timeoutSeconds = 60
    form.isDefault = false
    form.status = 'active'
  }
  dialogVisible.value = true
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (isEdit.value && editId.value) {
      await updateAiModel(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createAiModel({ ...form })
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
    await deleteAiModel(id)
    ElMessage.success('删除成功')
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
