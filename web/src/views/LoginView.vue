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
        <button class="flex-1 py-2.5 rounded-xl border text-sm cursor-pointer transition-all flex items-center justify-center gap-2" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);font-family:inherit">
          🐙 GitHub
        </button>
        <button class="flex-1 py-2.5 rounded-xl border text-sm cursor-pointer transition-all flex items-center justify-center gap-2" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);font-family:inherit">
          🔵 Google
        </button>
        <button class="flex-1 py-2.5 rounded-xl border text-sm cursor-pointer transition-all flex items-center justify-center gap-2" style="background:var(--bg-elevated);border-color:var(--border);color:var(--text-secondary);font-family:inherit">
          💬 微信
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
