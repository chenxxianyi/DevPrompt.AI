<template>
  <div>
    <h2 class="page-title">试用申请管理</h2>
    <el-card>
      <el-table :data="list" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="planCode" label="方案" width="110">
          <template #default="{ row }">
            <el-tag :type="planTagType(row.planCode)" size="small">{{ row.planCode.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="userId" label="用户ID" width="90" />
        <el-table-column prop="contact" label="联系方式" min-width="180" show-overflow-tooltip />
        <el-table-column prop="company" label="公司/团队" min-width="160" show-overflow-tooltip />
        <el-table-column prop="teamSize" label="团队规模" width="110" />
        <el-table-column prop="useCase" label="使用场景" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">{{ row.createdAt }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">处理</el-button>
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

    <el-dialog v-model="dialogVisible" title="处理试用申请" width="560px">
      <template v-if="current">
        <div class="detail-block">
          <p><strong>方案：</strong>{{ current.planCode.toUpperCase() }}</p>
          <p><strong>联系方式：</strong>{{ current.contact || '-' }}</p>
          <p><strong>公司/团队：</strong>{{ current.company || '-' }}</p>
          <p><strong>团队规模：</strong>{{ current.teamSize || '-' }}</p>
          <p><strong>使用场景：</strong>{{ current.useCase || '-' }}</p>
          <p><strong>用户留言：</strong>{{ current.message || '-' }}</p>
        </div>

        <el-form label-width="88px">
          <el-form-item label="处理状态">
            <el-select v-model="statusForm.status" style="width: 100%">
              <el-option label="已联系" value="contacted" />
              <el-option label="已通过" value="approved" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
          </el-form-item>
          <el-form-item label="后台备注">
            <el-input v-model="statusForm.adminNote" type="textarea" :rows="4" placeholder="记录沟通情况、审批说明或拒绝原因" />
          </el-form-item>
        </el-form>
      </template>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitStatus">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getTrialRequests, updateTrialRequestStatus } from '@/api/trialRequests'
import type { TrialRequest } from '@/types'

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const list = ref<TrialRequest[]>([])
const current = ref<TrialRequest | null>(null)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const statusForm = reactive({
  status: 'contacted' as TrialRequest['status'],
  adminNote: '',
})

function statusTagType(status: TrialRequest['status']) {
  const map: Record<TrialRequest['status'], string> = {
    pending: 'warning',
    contacted: 'primary',
    approved: 'success',
    rejected: 'danger',
  }
  return map[status] || 'info'
}

function statusLabel(status: TrialRequest['status']) {
  const map: Record<TrialRequest['status'], string> = {
    pending: '待处理',
    contacted: '已联系',
    approved: '已通过',
    rejected: '已拒绝',
  }
  return map[status] || status
}

function planTagType(code: TrialRequest['planCode']) {
  const map: Record<TrialRequest['planCode'], string> = {
    pro: 'primary',
    team: 'warning',
    enterprise: 'danger',
  }
  return map[code] || 'info'
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getTrialRequests({ page: page.value, pageSize: pageSize.value })
    const data = res.data.data
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function openDialog(row: TrialRequest) {
  current.value = row
  statusForm.status = row.status === 'pending' ? 'contacted' : row.status
  statusForm.adminNote = row.adminNote || ''
  dialogVisible.value = true
}

async function submitStatus() {
  if (!current.value) return
  saving.value = true
  try {
    await updateTrialRequestStatus(current.value.id, {
      status: statusForm.status,
      adminNote: statusForm.adminNote.trim(),
    })
    ElMessage.success(statusForm.status === 'approved' ? '处理成功，会员已开通' : '处理成功')
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
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

.detail-block {
  margin-bottom: 16px;
  padding: 12px 14px;
  border-radius: 8px;
  background: #f7f8fa;
  color: #4e5969;
  line-height: 1.7;
}

.detail-block p {
  margin: 0 0 4px;
}
</style>
