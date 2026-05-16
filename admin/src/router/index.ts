import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/components/Layout.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        { path: 'dashboard', name: 'Dashboard', meta: { title: '控制台' }, component: () => import('@/views/Dashboard.vue') },
        { path: 'users', name: 'Users', meta: { title: '用户管理' }, component: () => import('@/views/Users.vue') },
        { path: 'categories', name: 'Categories', meta: { title: '分类管理' }, component: () => import('@/views/Categories.vue') },
        { path: 'prompts', name: 'Prompts', meta: { title: 'Prompt 管理' }, component: () => import('@/views/Prompts.vue') },
        { path: 'ai-models', name: 'AiModels', meta: { title: 'AI 模型管理' }, component: () => import('@/views/AiModels.vue') },
        { path: 'ai-call-logs', name: 'CallLogs', meta: { title: 'AI 调用日志' }, component: () => import('@/views/CallLogs.vue') },
        { path: 'membership-plans', name: 'Plans', meta: { title: '会员套餐管理' }, component: () => import('@/views/Plans.vue') },
        { path: 'orders', name: 'Orders', meta: { title: '订单管理' }, component: () => import('@/views/Orders.vue') },
        { path: 'trial-requests', name: 'TrialRequests', meta: { title: '试用申请管理' }, component: () => import('@/views/TrialRequests.vue') },
        { path: 'project-types', name: 'ProjectTypes', meta: { title: '项目类型管理' }, component: () => import('@/views/ProjectTypes.vue') },
        { path: 'prompt-recipes', name: 'PromptRecipes', meta: { title: 'Prompt Recipe 管理' }, component: () => import('@/views/PromptRecipes.vue') },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth !== false && !auth.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && auth.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
