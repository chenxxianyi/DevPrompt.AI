<template>
  <div class="dashboard">
    <h2 class="page-title">控制台</h2>
    <el-row :gutter="20">
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-inner">
            <div class="stat-icon" style="background: #e6f7ff;">
              <el-icon :size="28" color="#1890ff"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.userCount }}</p>
              <p class="stat-label">用户总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-inner">
            <div class="stat-icon" style="background: #f6ffed;">
              <el-icon :size="28" color="#52c41a"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.promptCount }}</p>
              <p class="stat-label">Prompt 模板</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-inner">
            <div class="stat-icon" style="background: #fff7e6;">
              <el-icon :size="28" color="#fa8c16"><Cpu /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.callLogCount }}</p>
              <p class="stat-label">AI 调用次数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-inner">
            <div class="stat-icon" style="background: #f0f5ff;">
              <el-icon :size="28" color="#722ed1"><Coin /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.planCount }}</p>
              <p class="stat-label">会员套餐</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-inner">
            <div class="stat-icon" style="background: #fff0f6;">
              <el-icon :size="28" color="#eb2f96"><Grid /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.projectTypeCount }}</p>
              <p class="stat-label">项目类型</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboardStats } from '@/api/dashboard'
import type { DashboardStats } from '@/types'
import { UserFilled, Document, Cpu, Coin, Grid } from '@element-plus/icons-vue'

const stats = ref<DashboardStats>({
  userCount: 0,
  promptCount: 0,
  callLogCount: 0,
  planCount: 0,
  projectTypeCount: 0,
})

onMounted(async () => {
  try {
    const res = await getDashboardStats()
    stats.value = res.data.data
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
}
.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
  color: #1d2129;
}
.stat-card {
  border-radius: 8px;
}
.stat-inner {
  display: flex;
  align-items: center;
  gap: 16px;
}
.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
  line-height: 1.2;
}
.stat-label {
  font-size: 14px;
  color: #86909c;
  margin: 4px 0 0;
}
</style>
