import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('@/views/HomeView.vue') },
    { path: '/prompts', name: 'prompts', component: () => import('@/views/PromptsView.vue') },
    { path: '/prompts/:slug', name: 'prompt-detail', component: () => import('@/views/PromptDetailView.vue') },
    { path: '/generator', name: 'generator', component: () => import('@/views/GeneratorView.vue') },
    { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },
    { path: '/pricing', name: 'pricing', component: () => import('@/views/PricingView.vue') },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
