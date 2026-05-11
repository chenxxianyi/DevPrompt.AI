import type { RouteRecordRaw } from 'vue-router'

/**
 * iOS26 路由组：以 /ios26 为根，使用 IosAppShell 作为外壳。
 * 旧版路由不受影响。
 */
export const ios26Routes: RouteRecordRaw[] = [
  {
    path: '/ios26',
    component: () => import('./components/IosAppShell.vue'),
    children: [
      {
        path: '',
        name: 'ios26-home',
        component: () => import('./views/IosHomeView.vue'),
      },
      {
        path: 'prompts',
        name: 'ios26-prompts',
        component: () => import('./views/IosPromptsView.vue'),
      },
      {
        path: 'prompts/:slug',
        name: 'ios26-prompt-detail',
        component: () => import('./views/IosPromptDetailView.vue'),
      },
      {
        path: 'generator',
        name: 'ios26-generator',
        component: () => import('./views/IosGeneratorView.vue'),
      },
      {
        path: 'generator/history',
        name: 'ios26-history',
        component: () => import('./views/IosHistoryView.vue'),
      },
      {
        path: 'dashboard',
        name: 'ios26-dashboard',
        component: () => import('./views/IosDashboardView.vue'),
      },
      {
        path: 'login',
        name: 'ios26-login',
        component: () => import('./views/IosLoginView.vue'),
      },
      {
        path: 'pricing',
        name: 'ios26-pricing',
        component: () => import('./views/IosPricingView.vue'),
      },
    ],
  },
]
