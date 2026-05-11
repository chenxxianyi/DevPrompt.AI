<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useUiStore } from '@/store/ui'
import { useGeneratorStore, type GeneratorTab } from '@/store/generator'
import { getHistory } from '@/api/generator'
import { getGenerateStats } from '@/api/user'
import type { GeneratedPrompt } from '@/types'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosListRow from '../components/IosListRow.vue'
import IosIcon from '../components/IosIcon.vue'

const router = useRouter()
const auth = useAuthStore()
const ui = useUiStore()
const gen = useGeneratorStore()

const loading = ref(false)
const loadFailed = ref(false)
const stats = ref({ totalUsed: 0, dailyLimit: 5 })
const recentList = ref<GeneratedPrompt[]>([])

const typeLabelMap: Record<string, string> = {
  project: '项目',
  'cursor-rules': 'Cursor',
  'claude-code': 'Claude',
  optimize: '优化',
}

const remainingQuota = computed(() => {
  const limit = stats.value.dailyLimit
  if (limit >= 999999) return Infinity
  return Math.max(0, limit - stats.value.totalUsed)
})

const quotaPercent = computed(() => {
  const limit = stats.value.dailyLimit
  if (limit >= 999999) return 0
  return Math.min(100, Math.round((stats.value.totalUsed / limit) * 100))
})

const showUpgradeTip = computed(() => {
  if (!auth.isLoggedIn) return false
  if (auth.user?.membershipLevel !== 'free') return false
  return remainingQuota.value !== Infinity && remainingQuota.value <= 5
})

onMounted(async () => {
  if (!auth.isLoggedIn) {
    router.push('/ios26/login')
    return
  }
  loading.value = true
  try {
    const statsRes = await getGenerateStats()
    stats.value.totalUsed = statsRes.data.data.dailyUsed
    stats.value.dailyLimit = statsRes.data.data.dailyLimit
    const res = await getHistory({ page: 1, pageSize: 5 })
    recentList.value = res.data.data.list
  } catch {
    loadFailed.value = true
    stats.value.dailyLimit = 0
  } finally {
    loading.value = false
  }
})

function formatTime(t: string): string {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function regenerate(item: GeneratedPrompt) {
  gen.loadFromHistory({
    type: item.type as GeneratorTab,
    input: item.input,
  })
  router.push({ path: '/ios26/generator', query: { type: item.type } })
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  ui.showToast('已复制', 'success')
}

function logout() {
  auth.clearAuth()
  ui.showToast('已退出登录', 'success')
  router.push('/ios26')
}

const memberLabel = computed(() => {
  const lvl = auth.user?.membershipLevel
  if (!lvl || lvl === 'free') return '免费用户'
  if (lvl === 'pro') return 'Pro 会员'
  if (lvl === 'team') return 'Team 会员'
  if (lvl === 'enterprise') return '企业会员'
  return lvl
})
</script>

<template>
  <div class="ios-dash">
    <IosNavBar large />

    <div class="ios-page">
      <!-- Profile card -->
      <IosGlassPanel size="md" class="ios-dash__profile">
        <div class="ios-dash__avatar">
          {{ (auth.user?.username || 'U').slice(0, 1).toUpperCase() }}
        </div>
        <div class="ios-dash__profile-info">
          <div class="ios-text-title-3">{{ auth.user?.username || '访客' }}</div>
          <div class="ios-text-subheadline ios-dash__profile-email">{{ auth.user?.email || '-' }}</div>
          <span class="ios-tag ios-dash__member-tag">{{ memberLabel }}</span>
        </div>
      </IosGlassPanel>

      <!-- Loading -->
      <div v-if="loading" class="ios-dash__state"><span class="ios-spinner" /> 加载中...</div>

      <template v-else>
        <!-- Quota -->
        <IosGlassPanel size="md" class="ios-dash__quota">
          <div class="ios-dash__quota-head">
            <span class="ios-text-subheadline-emph">今日生成额度</span>
            <span v-if="!loadFailed && remainingQuota !== Infinity" class="ios-text-caption-1 ios-dash__quota-meta">
              剩余 {{ remainingQuota }} / {{ stats.dailyLimit }}
            </span>
          </div>
          <div v-if="loadFailed" class="ios-text-callout ios-dash__quota-fail">加载失败</div>
          <template v-else>
            <div class="ios-dash__quota-num">
              <span v-if="remainingQuota === Infinity">∞</span>
              <template v-else>
                <span class="ios-dash__quota-num-main">{{ remainingQuota }}</span>
                <span class="ios-dash__quota-num-sub">/ {{ stats.dailyLimit }}</span>
              </template>
            </div>
            <div v-if="remainingQuota !== Infinity" class="ios-dash__bar">
              <div
                class="ios-dash__bar-fill"
                :style="{ width: quotaPercent + '%', background: quotaPercent > 80 ? 'var(--ios-color-systemRed)' : 'var(--ios-color-tint)' }"
              />
            </div>
          </template>
        </IosGlassPanel>

        <!-- Upgrade tip -->
        <IosGlassPanel v-if="showUpgradeTip" size="md" class="ios-dash__upgrade">
          <div class="ios-dash__upgrade-row">
            <span class="ios-dash__upgrade-icon">
              <IosIcon path="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z M12 9v4 M12 17h.01" :size="20" :stroke="2" />
            </span>
            <div class="ios-dash__upgrade-text">
              <div class="ios-text-subheadline-emph">免费额度即将用尽</div>
              <p class="ios-text-footnote ios-dash__upgrade-desc">升级 Pro 获得更多生成次数和高级功能</p>
            </div>
            <IosButton variant="filled" size="sm" to="/ios26/pricing">查看方案</IosButton>
          </div>
        </IosGlassPanel>

        <!-- Quick actions -->
        <section class="ios-section">
          <header class="ios-section__header">
            <h3 class="ios-text-title-3">快捷生成</h3>
          </header>
          <div class="ios-list">
            <IosListRow disclosure to="/ios26/generator?type=project">
              <template #leading>
                <span class="ios-dash__row-icon" style="background:var(--ios-color-systemIndigo)">
                  <IosIcon path="M3 7.5A1.5 1.5 0 0 1 4.5 6h4.879a1.5 1.5 0 0 1 1.06.44l1.122 1.12A1.5 1.5 0 0 0 12.62 8H19.5A1.5 1.5 0 0 1 21 9.5V18a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18V7.5Z" :size="18" :stroke="1.8" />
                </span>
              </template>
              项目 Prompt
            </IosListRow>
            <IosListRow disclosure to="/ios26/generator?type=cursor-rules">
              <template #leading>
                <span class="ios-dash__row-icon" style="background:var(--ios-color-systemTeal)">
                  <IosIcon path="M5.636 4.364a9 9 0 0 1 12.728 0 9 9 0 0 1 0 12.728 9 9 0 0 1-12.728 0 9 9 0 0 1 0-12.728Z" :size="18" :stroke="1.8" />
                </span>
              </template>
              Cursor Rules
            </IosListRow>
            <IosListRow disclosure to="/ios26/generator?type=optimize">
              <template #leading>
                <span class="ios-dash__row-icon" style="background:var(--ios-color-systemGreen)">
                  <IosIcon path="M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z" :size="18" :stroke="1.8" />
                </span>
              </template>
              Prompt 优化
            </IosListRow>
            <IosListRow disclosure to="/ios26/prompts">
              <template #leading>
                <span class="ios-dash__row-icon" style="background:var(--ios-color-systemOrange)">
                  <IosIcon path="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20" :size="18" :stroke="1.8" />
                </span>
              </template>
              浏览模板库
            </IosListRow>
          </div>
        </section>

        <!-- Recent history -->
        <section class="ios-section">
          <header class="ios-section__header">
            <h3 class="ios-text-title-3">最近生成</h3>
            <router-link to="/ios26/generator/history" class="ios-section__more">
              查看全部
              <IosIcon path="M9 18l6-6-6-6" :size="14" :stroke="2.4" />
            </router-link>
          </header>
          <div v-if="recentList.length === 0" class="ios-dash__empty">
            <p class="ios-text-subheadline">暂无生成记录</p>
            <IosButton variant="tinted" size="sm" to="/ios26/generator">去生成</IosButton>
          </div>
          <div v-else class="ios-dash__recent">
            <IosGlassPanel
              v-for="item in recentList"
              :key="item.id"
              size="md"
              interactive
              @click="regenerate(item)"
            >
              <div class="ios-dash__recent-row">
                <div class="ios-dash__recent-info">
                  <div class="ios-dash__recent-head">
                    <span class="ios-tag">{{ typeLabelMap[item.type] || item.type }}</span>
                    <span v-if="item.title" class="ios-text-subheadline-emph ios-dash__recent-title">{{ item.title }}</span>
                  </div>
                  <p class="ios-text-footnote ios-dash__recent-excerpt">{{ item.output?.slice(0, 100) || '（无内容）' }}</p>
                  <span class="ios-text-caption-1 ios-dash__recent-time">{{ formatTime(item.createdAt) }}</span>
                </div>
                <div class="ios-dash__recent-actions">
                  <IosButton variant="tinted" size="sm" @click.stop="copy(item.output)">复制</IosButton>
                </div>
              </div>
            </IosGlassPanel>
          </div>
        </section>

        <!-- Settings list -->
        <section class="ios-section">
          <div class="ios-list">
            <IosListRow disclosure to="/ios26/pricing">
              <template #leading>
                <span class="ios-dash__row-icon" style="background:var(--ios-color-systemPink)">
                  <IosIcon path="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" :size="18" :stroke="1.8" />
                </span>
              </template>
              会员方案
            </IosListRow>
            <IosListRow destructive @click="logout">
              退出登录
            </IosListRow>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<style scoped>
.ios-dash {
  display: flex;
  flex-direction: column;
}

.ios-dash__profile {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 16px;
}

.ios-dash__avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--ios-color-tint), var(--ios-color-systemBlue));
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.5px;
  flex-shrink: 0;
}

.ios-dash__profile-info {
  flex: 1;
  min-width: 0;
}

.ios-dash__profile-email {
  margin: 2px 0 6px;
  color: var(--ios-color-label-secondary);
}

.ios-dash__member-tag {
  display: inline-flex;
}

.ios-dash__state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 0;
  color: var(--ios-color-label-secondary);
}

.ios-dash__quota {
  margin-bottom: 16px;
}

.ios-dash__quota-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.ios-dash__quota-meta {
  color: var(--ios-color-label-tertiary);
}

.ios-dash__quota-fail {
  color: var(--ios-color-systemRed);
}

.ios-dash__quota-num {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 12px;
}

.ios-dash__quota-num-main {
  font-size: 36px;
  font-weight: 800;
  letter-spacing: -1px;
  color: var(--ios-color-tint);
}

.ios-dash__quota-num-sub {
  font-size: 18px;
  color: var(--ios-color-label-tertiary);
}

.ios-dash__bar {
  height: 6px;
  background: var(--ios-color-fill-tertiary);
  border-radius: 9999px;
  overflow: hidden;
}

.ios-dash__bar-fill {
  height: 100%;
  border-radius: inherit;
  transition: width 0.5s var(--ios-ease-gentle);
}

.ios-dash__upgrade {
  margin-bottom: 16px;
  background: rgba(255, 149, 0, 0.1);
  border-color: rgba(255, 149, 0, 0.3);
}

.ios-dash__upgrade-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ios-dash__upgrade-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 149, 0, 0.2);
  color: var(--ios-color-systemOrange);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ios-dash__upgrade-text {
  flex: 1;
  min-width: 0;
}

.ios-dash__upgrade-desc {
  margin: 2px 0 0;
  color: var(--ios-color-label-secondary);
}

.ios-section {
  margin-top: 20px;
}

.ios-section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px 12px;
}

.ios-section__header h3 {
  margin: 0;
  color: var(--ios-color-label-primary);
}

.ios-section__more {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 15px;
  color: var(--ios-color-tint);
  font-weight: 500;
}

.ios-dash__row-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.ios-dash__empty {
  text-align: center;
  padding: 30px 0;
}

.ios-dash__empty p {
  margin: 0 0 12px;
  color: var(--ios-color-label-secondary);
}

.ios-dash__recent {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.ios-dash__recent-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ios-dash__recent-info {
  flex: 1;
  min-width: 0;
}

.ios-dash__recent-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.ios-dash__recent-title {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--ios-color-label-primary);
  flex: 1;
  min-width: 0;
}

.ios-dash__recent-excerpt {
  margin: 0 0 4px;
  color: var(--ios-color-label-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ios-dash__recent-time {
  color: var(--ios-color-label-tertiary);
}

.ios-dash__recent-actions {
  flex-shrink: 0;
}

.ios-tag {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: var(--ios-radius-pill);
  background: var(--ios-color-tint-soft);
  color: var(--ios-color-tint);
  letter-spacing: 0.06px;
  white-space: nowrap;
}
</style>
