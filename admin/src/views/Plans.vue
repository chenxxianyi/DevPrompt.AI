<template>
  <div>
    <h2 class="page-title">会员套餐管理</h2>
    <el-card>
      <div style="margin-bottom: 16px;">
        <el-button type="primary" @click="openDialog()">新增套餐</el-button>
      </div>
      <el-table :data="plans" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="套餐名称" width="130" />
        <el-table-column prop="code" label="标识" width="100" />
        <el-table-column prop="price" label="价格" width="90">
          <template #default="{ row }">¥{{ row.price }}</template>
        </el-table-column>
        <el-table-column prop="durationDays" label="时长(天)" width="90" />
        <el-table-column prop="dailyLimit" label="每日限制" width="90" />
        <el-table-column label="功能权益" min-width="250">
          <template #default="{ row }">
            <el-tag v-for="f in row.features" :key="f" size="small" style="margin: 2px 4px 2px 0;">{{ f }}</el-tag>
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
            <el-popconfirm title="确定删除该套餐？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑套餐' : '新增套餐'" width="550px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="套餐名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="标识" prop="code">
          <el-input v-model="form.code" />
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="时长(天)" prop="durationDays">
          <el-input-number v-model="form.durationDays" :min="0" />
        </el-form-item>
        <el-form-item label="每日限制" prop="dailyLimit">
          <el-input-number v-model="form.dailyLimit" :min="0" />
        </el-form-item>
        <el-form-item label="功能权益" prop="features">
          <el-select v-model="form.features" multiple filterable allow-create default-first-option style="width: 100%">
            <el-option v-for="f in form.features" :key="f" :label="f" :value="f" />
          </el-select>
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
import { getPlans, createPlan, updatePlan, deletePlan } from '@/api/plans'
import type { MembershipPlan } from '@/types'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const plans = ref<MembershipPlan[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()

const form = reactive({
  name: '',
  code: 'free' as 'free' | 'pro' | 'team' | 'enterprise',
  price: 0,
  durationDays: 30,
  dailyLimit: 5,
  features: [] as string[],
  status: 'active' as 'active' | 'disabled',
})

const rules = {
  name: [{ required: true, message: '请输入套餐名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入标识', trigger: 'blur' }],
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getPlans()
    plans.value = res.data.data.list
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function openDialog(row?: MembershipPlan) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.name = row.name
    form.code = row.code
    form.price = row.price
    form.durationDays = row.durationDays
    form.dailyLimit = row.dailyLimit
    form.features = [...(row.features || [])]
    form.status = row.status
  } else {
    isEdit.value = false
    editId.value = null
    form.name = ''
    form.code = 'free'
    form.price = 0
    form.durationDays = 30
    form.dailyLimit = 5
    form.features = []
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
      await updatePlan(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createPlan({ ...form })
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
    await deletePlan(id)
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
