<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const ui = useUiStore()
const auth = useAuthStore()

const isLogin = ref(true)
const email = ref('')
const password = ref('')
const username = ref('')
const remember = ref(false)
const submitting = ref(false)

function handleSocialLogin(provider: string) {
  ui.showToast(`${provider} 登录即将支持`)
}

async function handleSubmit() {
  if (submitting.value) return

  if (isLogin.value) {
    if (!email.value || !password.value) {
      ui.showToast('请填写邮箱和密码')
      return
    }
    submitting.value = true
    try {
      await auth.login({ email: email.value, password: password.value })
      ui.showToast('登录成功', 'success')
      router.push('/')
    } catch (e: any) {
      ui.showToast(e.message || '登录失败，请检查邮箱和密码')
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
      router.push('/')
    } catch (e: any) {
      ui.showToast(e.message || '注册失败，请稍后重试')
    } finally {
      submitting.value = false
    }
  }
}
</script>

<template>
  <div class="relative z-[1] min-h-screen flex items-center justify-center px-6" style="min-height:calc(100vh - 64px)">
    <div class="glass w-full max-w-[420px] p-10">
      <h2 class="text-2xl font-extrabold mb-1.5 text-center">{{ isLogin ? '欢迎回来' : '创建账号' }}</h2>
      <p class="text-center text-sm mb-7" style="color:var(--text-secondary)">
        {{ isLogin ? '登录 DevPrompt AI，开始高效编程' : '注册 DevPrompt AI 账号' }}
      </p>

      <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
        <div v-if="!isLogin" class="flex flex-col gap-2">
          <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">用户名</label>
          <input v-model="username" class="form-input py-3" type="text" placeholder="your_name">
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">邮箱</label>
          <input v-model="email" class="form-input py-3" type="email" placeholder="your@email.com">
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">密码</label>
          <input v-model="password" class="form-input py-3" type="password" placeholder="输入密码">
        </div>

        <div v-if="isLogin" class="flex items-center justify-between text-[13px]">
          <label class="flex items-center gap-1.5 cursor-pointer" style="color:var(--text-secondary)">
            <input v-model="remember" type="checkbox"> 记住我
          </label>
          <a class="no-underline cursor-pointer" style="color:var(--accent)">忘记密码？</a>
        </div>

        <button type="submit" class="btn btn-primary w-full py-3 text-[15px] justify-center rounded-xl mt-1" :disabled="submitting">
          <span v-if="submitting" class="spinner"></span>
          {{ submitting ? '处理中...' : (isLogin ? '登 录' : '注 册') }}
        </button>
      </form>

      <div class="flex items-center gap-3 my-5 text-[13px]" style="color:var(--text-muted)">
        <span class="flex-1 h-px" style="background:var(--border)"></span>
        <span>或</span>
        <span class="flex-1 h-px" style="background:var(--border)"></span>
      </div>

      <div class="flex gap-2.5">
        <button class="flex-1 py-2.5 rounded-xl border text-sm cursor-pointer transition-all flex items-center justify-center gap-2 opacity-70 hover:opacity-100 relative" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);font-family:inherit" @click="handleSocialLogin('GitHub')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0 1 12 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/></svg>
        </button>
        <button class="flex-1 py-2.5 rounded-xl border text-sm cursor-pointer transition-all flex items-center justify-center gap-2 opacity-70 hover:opacity-100 relative" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);font-family:inherit" @click="handleSocialLogin('Google')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>
        </button>
        <button class="flex-1 py-2.5 rounded-xl border text-sm cursor-pointer transition-all flex items-center justify-center gap-2 opacity-70 hover:opacity-100 relative" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);font-family:inherit" @click="handleSocialLogin('微信')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.045c.134 0 .24-.11.24-.245 0-.06-.024-.12-.04-.178l-.325-1.233a.492.492 0 0 1 .177-.553C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-7.062-6.122zM17.23 12.3c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm-4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982z"/></svg>
        </button>
      </div>

      <p class="text-center mt-5 text-sm" style="color:var(--text-secondary)">
        {{ isLogin ? '还没有账号？' : '已有账号？' }}
        <a class="font-semibold cursor-pointer no-underline" style="color:var(--accent)" @click="isLogin = !isLogin">
          {{ isLogin ? '立即注册' : '立即登录' }}
        </a>
      </p>
    </div>
  </div>
</template>
