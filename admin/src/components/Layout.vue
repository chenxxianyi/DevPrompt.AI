<template>
  <div class="layout-container">
    <el-container style="height: 100vh">
      <el-aside :width="appStore.sidebarCollapsed ? '64px' : '220px'" class="layout-aside">
        <div class="logo-area">
          <span v-if="!appStore.sidebarCollapsed" class="logo-text">DevPrompt AI</span>
          <span v-else class="logo-text-short">DP</span>
        </div>
        <el-menu
          :default-active="route.path"
          :collapse="appStore.sidebarCollapsed"
          :router="true"
          background-color="#001529"
          text-color="#ffffffb3"
          active-text-color="#fff"
        >
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>控制台</span>
          </el-menu-item>
          <el-menu-item index="/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/categories">
            <el-icon><Collection /></el-icon>
            <span>分类管理</span>
          </el-menu-item>
          <el-menu-item index="/prompts">
            <el-icon><Document /></el-icon>
            <span>Prompt 管理</span>
          </el-menu-item>
          <el-menu-item index="/ai-models">
            <el-icon><Cpu /></el-icon>
            <span>AI 模型管理</span>
          </el-menu-item>
          <el-menu-item index="/ai-call-logs">
            <el-icon><List /></el-icon>
            <span>AI 调用日志</span>
          </el-menu-item>
          <el-menu-item index="/membership-plans">
            <el-icon><Coin /></el-icon>
            <span>会员套餐管理</span>
          </el-menu-item>
          <el-menu-item index="/orders">
            <el-icon><Tickets /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          <el-menu-item index="/trial-requests">
            <el-icon><MessageBox /></el-icon>
            <span>试用申请管理</span>
          </el-menu-item>
          <el-menu-item index="/project-types">
            <el-icon><Grid /></el-icon>
            <span>项目类型管理</span>
          </el-menu-item>
          <el-menu-item index="/prompt-recipes">
            <el-icon><MagicStick /></el-icon>
            <span>Prompt Recipe</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="layout-header">
          <div class="header-left">
            <el-icon class="collapse-btn" @click="appStore.toggleSidebar" style="cursor: pointer; font-size: 20px;">
              <Fold v-if="!appStore.sidebarCollapsed" />
              <Expand v-else />
            </el-icon>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-dropdown trigger="click">
              <span class="user-info">
                <el-avatar :size="32" icon="UserFilled" />
                <span class="username">{{ authStore.user?.username || '管理员' }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <el-main class="layout-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'
import {
  Odometer, User, Collection, Document, Cpu,
  List, Coin, Tickets, Fold, Expand, UserFilled, Grid, MessageBox, MagicStick,
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-aside {
  background-color: #001529;
  overflow-y: auto;
  overflow-x: hidden;
  transition: width 0.3s;
}

.logo-area {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-text-short {
  font-size: 14px;
}

.el-menu {
  border-right: none;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 20px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
  color: #333;
}

.layout-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}
</style>
