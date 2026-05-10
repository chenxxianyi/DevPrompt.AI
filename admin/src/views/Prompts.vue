<template>
  <div>
    <h2 class="page-title">Prompt 模板管理</h2>
    <el-card>
      <div style="margin-bottom: 16px;">
        <el-button type="primary" @click="openDialog()">新增 Prompt</el-button>
      </div>
      <el-table :data="prompts" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">{{ row.category?.name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="200">
          <template #default="{ row }">
            <el-tag v-for="tag in row.tags" :key="tag" size="small" style="margin-right: 4px;">{{ tag }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="useCount" label="使用次数" width="90" />
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
            <el-popconfirm title="确定删除该 Prompt？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑 Prompt' : '新增 Prompt'" width="700px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="256" />
        </el-form-item>
        <el-form-item label="标识" prop="slug">
          <el-input v-model="form.slug" maxlength="256" />
        </el-form-item>
        <el-form-item label="分类" prop="categoryId">
          <el-select v-model="form.categoryId" filterable style="width: 100%">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-select v-model="form.tags" multiple filterable allow-create default-first-option style="width: 100%">
            <el-option v-for="t in form.tags" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" maxlength="1024" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="8" />
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
import { getPrompts, createPrompt, updatePrompt, deletePrompt } from '@/api/prompts'
import { getCategories } from '@/api/categories'
import type { PromptTemplate, PromptCategory } from '@/types'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const prompts = ref<PromptTemplate[]>([])
const categories = ref<PromptCategory[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()

const form = reactive({
  title: '',
  slug: '',
  categoryId: undefined as number | undefined,
  tags: [] as string[],
  description: '',
  content: '',
  status: 'active' as 'active' | 'disabled',
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
  categoryId: [{ required: true, message: '请选择分类', trigger: 'change' }],
}

async function fetchData() {
  loading.value = true
  try {
    const [promptsRes, catsRes] = await Promise.all([
      getPrompts({ page: page.value, pageSize: pageSize.value }),
      getCategories(),
    ])
    prompts.value = promptsRes.data.data.list
    total.value = promptsRes.data.data.total
    categories.value = catsRes.data.data.list
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function openDialog(row?: PromptTemplate) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.title = row.title
    form.slug = row.slug
    form.categoryId = row.categoryId
    form.tags = [...(row.tags || [])]
    form.description = row.description
    form.content = row.content
    form.status = row.status
  } else {
    isEdit.value = false
    editId.value = null
    form.title = ''
    form.slug = ''
    form.categoryId = undefined
    form.tags = []
    form.description = ''
    form.content = ''
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
      await updatePrompt(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createPrompt({ ...form })
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
    await deletePrompt(id)
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
.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
