<template>
  <div>
    <h2 class="page-title">用户管理</h2>
    <el-card>
      <el-table :data="users" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户名" width="140" />
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="role" label="角色" width="90">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
              {{ row.role === 'admin' ? '管理员' : '用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="membershipLevel" label="会员等级" width="110">
          <template #default="{ row }">
            <el-tag :type="membershipTagType(row.membershipLevel)" size="small">
              {{ membershipLabel(row.membershipLevel) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="dailyGenerateCount" label="今日生成" width="90" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="180">
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
import { getUsers } from '@/api/users'
import type { User } from '@/types'

const loading = ref(false)
const users = ref<User[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

function membershipTagType(level: string) {
  const map: Record<string, string> = { free: 'info', pro: 'primary', team: 'warning', enterprise: 'danger' }
  return map[level] || 'info'
}

function membershipLabel(level: string) {
  const map: Record<string, string> = { free: '免费', pro: 'Pro', team: 'Team', enterprise: 'Enterprise' }
  return map[level] || level
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getUsers({ page: page.value, pageSize: pageSize.value })
    const data = res.data.data
    users.value = data.list
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
