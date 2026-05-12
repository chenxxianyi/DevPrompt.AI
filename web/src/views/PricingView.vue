<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
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
const dialogVisible = ref(false)
const submitting = ref(false)
const selectedPlan = ref<MembershipPlan | null>(null)

const trialForm = reactive({
  contact: '',
  company: '',
  teamSize: '',
  useCase: '',
  message: '',
})

onMounted(async () => {
  try {
    const res = await getMembershipPlans()
    plans.value = res.data.data
  } catch {
    ui.showToast('加载会员方案失败，请稍后重试')
  } finally {
    loading.value = false
  }
})

function getFeatures(plan: MembershipPlan) {
  return plan.features || []
}

function getButtonLabel(plan: MembershipPlan) {
  if (plan.code === 'free') return '免费开始'
  if (plan.code === 'enterprise') return '申请咨询'
  if (auth.user?.membershipLevel === plan.code) return '当前方案'
  return '申请试用'
}

function getButtonClass(plan: MembershipPlan) {
  if (plan.code === 'pro') return 'btn-primary'
  if (plan.code === 'team') return 'btn-cyan'
  return 'btn-ghost'
}

function openTrialDialog(plan: MembershipPlan) {
  selectedPlan.value = plan
  trialForm.contact = auth.user?.email || ''
  trialForm.company = ''
  trialForm.teamSize = ''
  trialForm.useCase = ''
  trialForm.message = ''
  dialogVisible.value = true
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
  if (auth.user?.membershipLevel === plan.code) {
    return
  }
  openTrialDialog(plan)
}

async function submitTrial() {
  if (!selectedPlan.value) return
  if (!trialForm.contact.trim()) {
    ui.showToast('请填写联系方式')
    return
  }

  submitting.value = true
  try {
    await createTrialRequest({
      planCode: selectedPlan.value.code,
      contact: trialForm.contact.trim(),
      company: trialForm.company.trim(),
      teamSize: trialForm.teamSize.trim(),
      useCase: trialForm.useCase.trim(),
      message: trialForm.message.trim(),
    })
    dialogVisible.value = false
    ui.showToast(
      selectedPlan.value.code === 'enterprise'
        ? '已记录您的需求，销售团队将尽快联系您'
        : '申请已提交，我们将尽快为您开通试用',
      'success',
    )
  } catch (e: any) {
    ui.showToast(e.message || '提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="relative z-[1] min-h-screen">
    <div class="max-w-screen-xl mx-auto px-6">
      <div class="text-center py-14 pb-20">
        <h1 class="text-[36px] font-black mb-2.5 tracking-tight">选择适合你的方案</h1>
        <p class="text-base mb-12" style="color:var(--text-secondary)">从免费开始，随时升级解锁更多生成能力和团队效率</p>

        <div v-if="loading" class="text-center py-10" style="color:var(--text-muted)">加载中...</div>

        <div v-else class="grid grid-cols-[repeat(auto-fit,minmax(260px,1fr))] gap-5 max-w-[1100px] mx-auto">
          <div
            v-for="plan in plans"
            :key="plan.id"
            class="glass p-8 text-left relative overflow-hidden"
            :style="plan.code === 'pro' ? { borderColor: 'var(--accent)', boxShadow: '0 0 40px var(--accent-glow)' } : {}"
          >
            <div
              v-if="plan.code === 'pro'"
              class="absolute top-4 -right-7 bg-[var(--accent)] text-white text-xs font-bold py-1 px-9 rotate-45"
            >
              推荐
            </div>

            <div class="text-xs font-semibold uppercase tracking-wider mb-2" style="color:var(--text-muted)">{{ plan.name }}</div>

            <div class="text-[36px] font-black mb-1">
              <template v-if="plan.price === 0">免费</template>
              <template v-else-if="plan.code === 'enterprise'">定制</template>
              <template v-else>
                <span class="text-lg font-semibold">¥</span>{{ plan.price }}
                <span class="text-sm font-normal" style="color:var(--text-muted)">/月</span>
              </template>
            </div>

            <div class="text-sm font-medium mb-5" style="color:var(--cyan)">
              {{ plan.dailyLimit >= 999999 ? '无限生成' : `每日 ${plan.dailyLimit} 次生成` }}
            </div>

            <ul class="list-none mb-6">
              <li
                v-for="feature in getFeatures(plan)"
                :key="feature"
                class="flex items-center gap-2 py-1.5 text-sm"
                style="color:var(--text-secondary)"
              >
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

    <div v-if="dialogVisible" class="fixed inset-0 z-50 flex items-center justify-center px-4" style="background:rgba(3,7,18,.72)">
      <div class="glass w-full max-w-[560px] p-6">
        <div class="flex items-start justify-between gap-4 mb-5">
          <div>
            <h2 class="text-xl font-bold">申请 {{ selectedPlan?.name }} 试用</h2>
            <p class="text-sm mt-1" style="color:var(--text-secondary)">补充一些信息，方便我们尽快联系并为你开通。</p>
          </div>
          <button class="btn btn-ghost text-xs px-3 py-1.5 rounded-lg" @click="dialogVisible = false">关闭</button>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1.5 md:col-span-2">
            <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">联系方式</label>
            <input v-model="trialForm.contact" class="form-input" type="text" placeholder="邮箱 / 微信 / 手机">
          </div>
          <div class="flex flex-col gap-1.5">
            <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">公司 / 团队</label>
            <input v-model="trialForm.company" class="form-input" type="text" placeholder="例如：某某工作室">
          </div>
          <div class="flex flex-col gap-1.5">
            <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">团队规模</label>
            <select v-model="trialForm.teamSize" class="form-input cursor-pointer">
              <option value="">请选择</option>
              <option value="1人">1人</option>
              <option value="2-5人">2-5人</option>
              <option value="6-20人">6-20人</option>
              <option value="20+人">20+人</option>
            </select>
          </div>
          <div class="flex flex-col gap-1.5 md:col-span-2">
            <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">使用场景</label>
            <input v-model="trialForm.useCase" class="form-input" type="text" placeholder="例如：研发提效、Prompt 规范、团队协作">
          </div>
          <div class="flex flex-col gap-1.5 md:col-span-2">
            <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">补充说明</label>
            <textarea v-model="trialForm.message" class="form-input min-h-[100px] resize-y" placeholder="描述你的业务场景、日常使用工具、期望解决的问题..."></textarea>
          </div>
        </div>

        <div class="flex justify-end gap-3 mt-6">
          <button class="btn btn-ghost text-sm px-4 py-2.5 rounded-xl" @click="dialogVisible = false">取消</button>
          <button class="btn btn-primary text-sm px-5 py-2.5 rounded-xl" :disabled="submitting" @click="submitTrial">
            {{ submitting ? '提交中...' : '提交申请' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
