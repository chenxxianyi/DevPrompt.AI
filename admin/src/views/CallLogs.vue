<template>
  <div>
    <h2 class="page-title">AI 调用日志</h2>
    <el-card>
      <el-table :data="logs" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="userId" label="用户ID" width="80" />
        <el-table-column prop="provider" label="Provider" width="100" />
        <el-table-column prop="model" label="模型" width="140" />
        <el-table-column prop="requestType" label="请求类型" width="140" />
        <el-table-column prop="promptTokens" label="输入 Token" width="100" />
        <el-table-column prop="completionTokens" label="输出 Token" width="100" />
        <el-table-column prop="totalTokens" label="总 Token" width="90" />
        <el-table-column prop="latencyMs" label="耗时(ms)" width="90" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="errorMessage" label="错误信息" min-width="160" show-overflow-tooltip />
        <el-table-column label="调用时间" width="180">
          <template #default="{ row }">{{ row.createdAt }}</template>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCallLogs } from '@/api/callLogs'
import type { AICallLog } from '@/types'

const loading = ref(false)
const logs = ref<AICallLog[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function fetchData() {
  loading.value = true
  try {
    const res = await getCallLogs({ page: page.value, pageSize: pageSize.value })
    const data = res.data.data
    logs.value = data.list
    total.value = data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
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
