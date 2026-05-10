<template>
  <div>
    <h2 class="page-title">项目类型管理</h2>
    <el-card>
      <div style="margin-bottom: 16px;">
        <el-button type="primary" @click="openDialog()">新增类型</el-button>
      </div>
      <el-table :data="types" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="显示名称" width="150" />
        <el-table-column prop="value" label="值" width="120" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="70" />
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
            <el-popconfirm title="确定删除该项目类型？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑类型' : '新增类型'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="显示名称" prop="name">
          <el-input v-model="form.name" maxlength="64" />
        </el-form-item>
        <el-form-item label="值" prop="value">
          <el-input v-model="form.value" maxlength="64" placeholder="用作选项 value" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" maxlength="256" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
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
import { getProjectTypes, createProjectType, updateProjectType, deleteProjectType } from '@/api/projectTypes'
import type { ProjectType } from '@/types'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const types = ref<ProjectType[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()

const form = reactive({
  name: '',
  value: '',
  description: '',
  sort: 0,
  status: 'active' as 'active' | 'disabled',
})

const rules = {
  name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  value: [{ required: true, message: '请输入值', trigger: 'blur' }],
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getProjectTypes()
    types.value = res.data.data.list
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function openDialog(row?: ProjectType) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.name = row.name
    form.value = row.value
    form.description = row.description
    form.sort = row.sort
    form.status = row.status
  } else {
    isEdit.value = false
    editId.value = null
    form.name = ''
    form.value = ''
    form.description = ''
    form.sort = 0
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
      await updateProjectType(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createProjectType({ ...form })
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
    await deleteProjectType(id)
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
