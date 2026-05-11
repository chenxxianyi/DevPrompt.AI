<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useAuthStore } from '@/store/auth'
import IosNavBar from '../components/IosNavBar.vue'
import IosGlassPanel from '../components/IosGlassPanel.vue'
import IosButton from '../components/IosButton.vue'
import IosSegmentedControl from '../components/IosSegmentedControl.vue'

const router = useRouter()
const ui = useUiStore()
const auth = useAuthStore()

const mode = ref<'login' | 'register'>('login')
const email = ref('')
const password = ref('')
const username = ref('')
const submitting = ref(false)

const modeOptions = [
  { value: 'login', label: '登录' },
  { value: 'register', label: '注册' },
]

async function handleSubmit() {
  if (submitting.value) return

  if (mode.value === 'login') {
    if (!email.value || !password.value) {
      ui.showToast('请填写邮箱和密码')
      return
    }
    submitting.value = true
    try {
      await auth.login({ email: email.value, password: password.value })
      ui.showToast('登录成功', 'success')
      router.push('/ios26')
    } catch (e: any) {
      ui.showToast(e?.message || '登录失败，请检查邮箱和密码')
    } finally {
      submitting.value = false
    }
  } else {
    if (!username.value || !email.value || !password.value) {
      ui.showToast('请填写所有字段')
      return
    }
    if (password.value.length < 6) {
      ui.showToast('密码长度不能少于 6 位')
      return
    }
    submitting.value = true
    try {
      await auth.register({ username: username.value, email: email.value, password: password.value })
      ui.showToast('注册成功', 'success')
      router.push('/ios26')
    } catch (e: any) {
      ui.showToast(e?.message || '注册失败，请稍后重试')
    } finally {
      submitting.value = false
    }
  }
}

function handleSocialLogin(provider: string) {
  ui.showToast(`${provider} 登录即将支持`)
}
</script>

<template>
  <div class="ios-login">
    <IosNavBar :title="mode === 'login' ? '登录' : '注册'" :show-back="true" back-to="/ios26" />

    <div class="ios-page">
      <div class="ios-login__hero">
        <h1 class="ios-text-large-title ios-login__title">
          {{ mode === 'login' ? '欢迎回来' : '创建账号' }}
        </h1>
        <p class="ios-text-callout ios-login__sub">
          {{ mode === 'login' ? '登录 DevPrompt.AI，开始高效编程' : '注册 DevPrompt.AI 账号' }}
        </p>
      </div>

      <IosSegmentedControl v-model="mode" :options="modeOptions" size="md" class="ios-login__mode" />

      <IosGlassPanel size="md" class="ios-login__form">
        <form @submit.prevent="handleSubmit" class="ios-login__form-inner">
          <div v-if="mode === 'register'" class="ios-field">
            <label class="ios-field__label">用户名</label>
            <input v-model="username" class="ios-input" type="text" placeholder="your_name" autocomplete="username">
          </div>
          <div class="ios-field">
            <label class="ios-field__label">邮箱</label>
            <input v-model="email" class="ios-input" type="email" placeholder="your@email.com" autocomplete="email">
          </div>
          <div class="ios-field">
            <label class="ios-field__label">密码</label>
            <input v-model="password" class="ios-input" type="password" placeholder="输入密码" :autocomplete="mode === 'login' ? 'current-password' : 'new-password'">
          </div>

          <IosButton type="submit" variant="filled" size="lg" block :loading="submitting" :disabled="submitting">
            {{ submitting ? '处理中...' : (mode === 'login' ? '登录' : '注册') }}
          </IosButton>
        </form>
      </IosGlassPanel>

      <div class="ios-login__divider">
        <span class="ios-login__divider-line" />
        <span class="ios-text-footnote ios-login__divider-text">或</span>
        <span class="ios-login__divider-line" />
      </div>

      <div class="ios-login__social">
        <IosButton variant="glass" size="md" @click="handleSocialLogin('GitHub')">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0 1 12 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/></svg>
          GitHub
        </IosButton>
        <IosButton variant="glass" size="md" @click="handleSocialLogin('Google')">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>
          Google
        </IosButton>
      </div>

      <p class="ios-text-footnote ios-login__terms">
        继续即表示同意 <a class="ios-login__link">服务条款</a> 与 <a class="ios-login__link">隐私政策</a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.ios-login {
  display: flex;
  flex-direction: column;
}

.ios-page {
  max-width: 480px;
  margin: 0 auto;
  padding: 24px var(--ios-page-padding) 48px;
  width: 100%;
}

.ios-login__hero {
  text-align: center;
  padding: 16px 0 20px;
}

.ios-login__title {
  margin: 0 0 8px;
  color: var(--ios-color-label-primary);
}

.ios-login__sub {
  margin: 0;
  color: var(--ios-color-label-secondary);
}

.ios-login__mode {
  margin-bottom: 16px;
}

.ios-login__form {
  margin-bottom: 16px;
}

.ios-login__form-inner {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ios-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ios-field__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--ios-color-label-secondary);
}

.ios-login__divider {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
}

.ios-login__divider-line {
  flex: 1;
  height: 1px;
  background: var(--ios-color-separator);
}

.ios-login__divider-text {
  color: var(--ios-color-label-tertiary);
}

.ios-login__social {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.ios-login__terms {
  text-align: center;
  margin: 16px 0 0;
  color: var(--ios-color-label-tertiary);
}

.ios-login__link {
  color: var(--ios-color-tint);
}
</style>
