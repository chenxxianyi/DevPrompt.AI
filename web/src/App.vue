<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import NavBar from '@/components/NavBar.vue'
import AppFooter from '@/components/AppFooter.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import BackToTop from '@/components/BackToTop.vue'

const auth = useAuthStore()
const route = useRoute()

// iOS26 路由组自带 NavBar / TabBar / Toast，需要隐藏旧版外壳
const isIos26 = computed(() => route.path.startsWith('/ios26'))

onMounted(() => {
  auth.initAuth()
  if (auth.isLoggedIn) {
    auth.fetchProfile()
  }
})
</script>

<template>
  <!-- iOS26: 由 IosAppShell 渲染独立外壳 -->
  <template v-if="isIos26">
    <router-view />
  </template>

  <!-- 旧版：保持原有结构和动效 -->
  <template v-else>
    <NavBar />
    <main class="pt-16">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <AppFooter />
    <ToastContainer />
    <BackToTop />
  </template>
</template>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.page-enter-from {
  opacity: 0;
  transform: translateY(6px);
}
.page-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
