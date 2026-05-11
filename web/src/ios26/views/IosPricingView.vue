<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useUiStore } from '@/store/ui'
import { getMembershipPlans } from '@/api/membership'
import { createTrialRequest } from '@/api/trial'
import type { MembershipPlan } from '@/types'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosIcon from '../components/IosIcon.vue'

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
    /* silent */
  } finally {
    loading.value = false
  }
})

function getFeatures(plan: MembershipPlan): string[] {
  return plan.features || []
}

function handleAction(plan: MembershipPlan) {
  if (plan.code === 'free') {
    router.push('/ios26/login')
    return
  }
  if (!auth.isLoggedIn) {
    router.push('/ios26/login')
    return
  }
  createTrialRequest({ planCode: plan.code })
    .then(() => {
      ui.showToast(
        plan.code === 'enterprise'
          ? '已记录您的需求，销售团队将尽快联系'
          : '申请已提交，将尽快为您开通试用',
        'success',
      )
    })
    .catch((e: any) => {
      ui.showToast(e?.message || '提交失败，请稍后重试')
    })
}

function getButtonLabel(plan: MembershipPlan): string {
  if (plan.code === 'free') return '免费开始'
  if (plan.code === 'enterprise') return '联系购买'
  if (auth.user?.membershipLevel === plan.code) return '当前方案'
  return '申请试用'
}

function getButtonVariant(plan: MembershipPlan): 'filled' | 'tinted' | 'glass' {
  if (plan.code === 'pro') return 'filled'
  if (plan.code === 'enterprise') return 'glass'
  return 'tinted'
}

function isCurrent(plan: MembershipPlan): boolean {
  return !!auth.isLoggedIn && auth.user?.membershipLevel === plan.code
}
</script>

<template>
  <div class="ios-pricing">
    <IosNavBar large />

    <div class="ios-page">
      <div class="ios-pricing__hero">
        <h1 class="ios-text-large-title ios-pricing__title">选择适合你的方案</h1>
        <p class="ios-text-callout ios-pricing__sub">从免费开始，随时升级解锁更多能力</p>
      </div>

      <div v-if="loading" class="ios-pricing__loading">
        <span class="ios-spinner" /> 加载中...
      </div>

      <div v-else class="ios-pricing__grid">
        <IosGlassPanel
          v-for="plan in plans"
          :key="plan.id"
          size="md"
          class="ios-pricing__card"
          :class="{ 'is-featured': plan.code === 'pro' }"
        >
          <div v-if="plan.code === 'pro'" class="ios-pricing__badge">推荐</div>

          <div class="ios-text-caption-2 ios-pricing__plan-name">{{ plan.name }}</div>

          <div class="ios-pricing__price">
            <template v-if="plan.price === 0 && plan.code === 'free'">免费</template>
            <template v-else-if="plan.code === 'enterprise'">定制</template>
            <template v-else>
              <span class="ios-pricing__currency">¥</span>{{ plan.price }}<span class="ios-pricing__period">/月</span>
            </template>
          </div>

          <div class="ios-text-subheadline-emph ios-pricing__limit">
            {{ plan.dailyLimit >= 999999 || plan.dailyLimit < 0 ? '无限生成' : `每日 ${plan.dailyLimit} 次生成` }}
          </div>

          <ul class="ios-pricing__features">
            <li v-for="feature in getFeatures(plan)" :key="feature" class="ios-pricing__feature">
              <span class="ios-pricing__check">
                <IosIcon path="M20 6 9 17l-5-5" :size="14" :stroke="2.6" />
              </span>
              <span class="ios-text-subheadline">{{ feature }}</span>
            </li>
          </ul>

          <IosButton
            :variant="getButtonVariant(plan)"
            size="md"
            block
            :disabled="isCurrent(plan)"
            @click="handleAction(plan)"
          >
            {{ getButtonLabel(plan) }}
          </IosButton>
        </IosGlassPanel>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ios-pricing {
  display: flex;
  flex-direction: column;
}

.ios-pricing__hero {
  text-align: center;
  padding: 16px 0 24px;
}

.ios-pricing__title {
  margin: 0 0 8px;
  color: var(--ios-color-label-primary);
}

.ios-pricing__sub {
  margin: 0;
  color: var(--ios-color-label-secondary);
}

.ios-pricing__loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 0;
  color: var(--ios-color-label-secondary);
}

.ios-pricing__grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 16px;
}

.ios-pricing__card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ios-pricing__card.is-featured {
  border-color: var(--ios-color-tint);
  box-shadow: 0 18px 50px var(--ios-color-tint-soft-strong);
}

.ios-pricing__badge {
  position: absolute;
  top: 14px;
  right: 14px;
  background: var(--ios-color-tint);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: var(--ios-radius-pill);
  letter-spacing: 0.06px;
  z-index: 2;
}

.ios-pricing__plan-name {
  color: var(--ios-color-label-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.6px;
}

.ios-pricing__price {
  font-size: 40px;
  line-height: 1;
  font-weight: 800;
  letter-spacing: -1.2px;
  color: var(--ios-color-label-primary);
}

.ios-pricing__currency {
  font-size: 22px;
  vertical-align: top;
  margin-right: 2px;
  font-weight: 600;
}

.ios-pricing__period {
  font-size: 15px;
  color: var(--ios-color-label-tertiary);
  font-weight: 400;
  margin-left: 2px;
}

.ios-pricing__limit {
  color: var(--ios-color-systemTeal);
}

.ios-pricing__features {
  list-style: none;
  margin: 4px 0 12px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ios-pricing__feature {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--ios-color-label-primary);
}

.ios-pricing__check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--ios-color-systemGreen);
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
</style>
