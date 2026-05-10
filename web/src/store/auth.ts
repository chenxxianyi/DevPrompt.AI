import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginParams, RegisterParams } from '@/types'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function initAuth() {
    const saved = localStorage.getItem('user')
    if (saved && token.value) {
      try { user.value = JSON.parse(saved) } catch { clearAuth() }
    }
  }

  async function login(data: LoginParams) {
    loading.value = true
    try {
      const res = await apiLogin(data)
      setAuth(res.data.data.token, res.data.data.user)
      return res.data.data
    } finally {
      loading.value = false
    }
  }

  async function register(data: RegisterParams) {
    loading.value = true
    try {
      const res = await apiRegister(data)
      setAuth(res.data.data.token, res.data.data.user)
      return res.data.data
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await getProfile()
      user.value = res.data.data
      localStorage.setItem('user', JSON.stringify(res.data.data))
    } catch {
      clearAuth()
    }
  }

  return { token, user, loading, isLoggedIn, setAuth, clearAuth, initAuth, login, register, fetchProfile }
})
