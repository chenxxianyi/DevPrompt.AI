<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const route = useRoute()
const auth = useAuthStore()

const navLinks = [
  { name: 'home', label: '首页', path: '/' },
  { name: 'prompts', label: '模板库', path: '/prompts' },
  { name: 'generator', label: '生成器', path: '/generator' },
  { name: 'pricing', label: '会员', path: '/pricing' },
]

const membershipBadge: Record<string, { label: string; color: string }> = {
  pro: { label: 'Pro', color: 'var(--accent)' },
  team: { label: 'Team', color: 'var(--cyan)' },
  enterprise: { label: '企业', color: 'var(--amber)' },
}
</script>

<template>
  <nav class="fixed top-0 left-0 right-0 z-50 border-b" style="background:rgba(6,8,15,.75);backdrop-filter:blur(20px) saturate(1.5);border-color:var(--border)">
    <div class="max-w-screen-xl mx-auto px-6 flex items-center justify-between h-16">
      <router-link to="/" class="flex items-center gap-2.5 no-underline">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center text-base font-black text-white" style="background:linear-gradient(135deg,var(--accent),var(--cyan));font-family:'JetBrains Mono',monospace">&gt;_</div>
        <span class="text-lg font-bold tracking-tight" style="color:var(--text-primary)">DevPrompt<span style="color:var(--accent)">.</span>AI</span>
      </router-link>

      <div class="hidden md:flex items-center gap-1">
        <router-link
          v-for="link in navLinks" :key="link.name"
          :to="link.path"
          class="relative px-4 py-2 rounded-lg text-sm font-medium no-underline transition-all duration-300"
          :class="route.name === link.name ? 'text-[var(--accent)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[rgba(139,92,246,.08)]'"
          :style="route.name === link.name ? 'background:rgba(139,92,246,.1)' : ''"
        >
          {{ link.label }}
          <span v-if="route.name === link.name" class="absolute bottom-0.5 left-1/2 -translate-x-1/2 w-4 h-0.5 rounded bg-[var(--accent)]" />
        </router-link>
      </div>

      <div class="flex items-center gap-3">
        <!-- iOS 风格切换入口 -->
        <router-link
          to="/ios26"
          class="theme-switch-btn"
          title="切换到 iOS 26 风格"
        >
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
            <line x1="12" y1="18" x2="12.01" y2="18"/>
          </svg>
          <span>新风版</span>
        </router-link>

        <template v-if="auth.isLoggedIn">
          <router-link
            to="/dashboard"
            class="hidden md:inline-flex items-center gap-1.5 text-sm no-underline px-3 py-1.5 rounded-lg transition-all"
            style="color:var(--text-secondary)"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            工作台
          </router-link>
          <template v-if="auth.user">
            <span
              v-if="auth.user.membershipLevel !== 'free'"
              class="text-[11px] font-bold px-2 py-0.5 rounded"
              :style="{
                background: (membershipBadge[auth.user.membershipLevel]?.color || 'var(--accent)') + '20',
                color: membershipBadge[auth.user.membershipLevel]?.color || 'var(--accent)',
              }"
            >
              {{ membershipBadge[auth.user.membershipLevel]?.label || auth.user.membershipLevel }}
            </span>
          </template>
          <span class="text-sm" style="color:var(--text-secondary)">{{ auth.user?.username }}</span>
          <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg" @click="auth.clearAuth()">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          </button>
        </template>
        <template v-else>
          <router-link to="/login" class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg no-underline">登录</router-link>
          <router-link to="/login" class="btn btn-primary text-xs px-3 py-1.5 rounded-lg no-underline">免费注册</router-link>
        </template>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.theme-switch-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px 5px 10px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
  text-decoration: none;
  letter-spacing: 0.01em;
  white-space: nowrap;
  color: #fff;
  background: linear-gradient(135deg, #7DD8FA 0%, #4FAAED 45%, #A78BFA 100%);
  box-shadow: 0 2px 10px rgba(79, 170, 237, 0.35), 0 1px 3px rgba(167, 139, 250, 0.2);
  transition: filter 0.2s ease, box-shadow 0.2s ease, transform 0.15s ease;
}

.theme-switch-btn:hover {
  filter: brightness(1.08) saturate(1.15);
  box-shadow: 0 4px 16px rgba(79, 170, 237, 0.5), 0 2px 6px rgba(167, 139, 250, 0.3);
  transform: translateY(-1px);
}

.theme-switch-btn:active {
  transform: translateY(0) scale(0.97);
}
</style>
