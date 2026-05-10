<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useUiStore } from '@/store/ui'
import { getMembershipPlans } from '@/api/membership'
import { createTrialRequest } from '@/api/trial'
import type { MembershipPlan } from '@/types'

const router = useRouter()
const auth = useAuthStore()
const ui = useUiStore()
const plans = ref<MembershipPlan[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await getMembershipPlans()
    plans.value = res.data.data
  } catch {
    // 接口不可用时静默处理
  } finally {
    loading.value = false
  }
})

function getFeatures(plan: MembershipPlan): string[] {
  return plan.features || []
}

function handleAction(plan: MembershipPlan) {
  if (plan.code === 'free') {
    router.push('/login')
    return
  }
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  // 提交试用申请
  createTrialRequest({ planCode: plan.code })
    .then(() => {
      ui.showToast(plan.code === 'enterprise' ? '已记录您的需求，销售团队将尽快联系您' : '申请已提交，我们将尽快为您开通试用', 'success')
    })
    .catch((e: any) => {
      ui.showToast(e.message || '提交失败，请稍后重试')
    })
}

function getButtonLabel(plan: MembershipPlan): string {
  if (plan.code === 'free') return '免费开始'
  if (plan.code === 'enterprise') return '联系购买'
  if (auth.user?.membershipLevel === plan.code) return '当前方案'
  return '申请试用'
}

function getButtonClass(plan: MembershipPlan): string {
  if (plan.code === 'pro') return 'btn-primary'
  if (plan.code === 'team') return 'btn-cyan'
  if (plan.code === 'enterprise') return 'btn-ghost'
  return 'btn-ghost'
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="text-center py-14 pb-20">
        <h1 class="text-[36px] font-black mb-2.5 tracking-tight">选择适合你的方案</h1>
        <p class="text-base mb-12" style="color:var(--text-secondary)">从免费开始，随时升级解锁更多能力</p>

        <div v-if="loading" class="text-center py-10" style="color:var(--text-muted)">加载中...</div>

        <div v-else class="grid grid-cols-[repeat(auto-fit,minmax(260px,1fr))] gap-5 max-w-[1100px] mx-auto">
          <div
            v-for="plan in plans" :key="plan.id"
            class="glass p-8 text-left relative overflow-hidden"
            :class="plan.code === 'pro' ? 'featured' : ''"
            :style="plan.code === 'pro'
              ? { borderColor: 'var(--accent)', boxShadow: '0 0 40px var(--accent-glow)' }
              : {}"
          >
            <!-- Featured badge -->
            <div v-if="plan.code === 'pro'"
              class="absolute top-4 -right-7 bg-[var(--accent)] text-white text-xs font-bold py-1 px-9 rotate-45"
            >推荐</div>

            <div class="text-xs font-semibold uppercase tracking-wider mb-2" style="color:var(--text-muted)">{{ plan.name }}</div>

            <div class="text-[36px] font-black mb-1">
              <template v-if="plan.price === 0">免费</template>
              <template v-else-if="plan.code === 'enterprise'">定制</template>
              <template v-else>
                <span class="text-lg font-semibold">¥</span>{{ plan.price }}<span class="text-sm font-normal" style="color:var(--text-muted)">/月</span>
              </template>
            </div>

            <div class="text-sm font-medium mb-5" style="color:var(--cyan)">
              {{ plan.dailyLimit >= 999999 ? '无限生成' : `每日 ${plan.dailyLimit} 次生成` }}
            </div>

            <ul class="list-none mb-6">
              <li v-for="feature in getFeatures(plan)" :key="feature" class="flex items-center gap-2 py-1.5 text-sm" style="color:var(--text-secondary)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="var(--green)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                {{ feature }}
              </li>
            </ul>

            <button
              class="w-full flex justify-center text-sm py-2.5 rounded-xl"
              :class="getButtonClass(plan)"
              :disabled="auth.isLoggedIn && auth.user?.membershipLevel === plan.code"
              @click="handleAction(plan)"
            >
              {{ getButtonLabel(plan) }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
