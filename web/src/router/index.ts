import { createRouter, createWebHistory } from 'vue-router'
import { ios26Routes } from '@/ios26/router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('@/views/HomeView.vue') },
    { path: '/prompts', name: 'prompts', component: () => import('@/views/PromptsView.vue') },
    { path: '/prompts/:slug', name: 'prompt-detail', component: () => import('@/views/PromptDetailView.vue') },
    { path: '/generator', name: 'generator', component: () => import('@/views/GeneratorView.vue') },
    { path: '/generator/history', name: 'history', component: () => import('@/views/HistoryView.vue') },
    { path: '/favorites', name: 'favorites', component: () => import('@/views/FavoritesView.vue') },
    { path: '/dashboard', name: 'dashboard', component: () => import('@/views/DashboardView.vue') },
    { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },
    { path: '/pricing', name: 'pricing', component: () => import('@/views/PricingView.vue') },
    // iOS 26 / Liquid Glass 新版前端入口
    ...ios26Routes,
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
