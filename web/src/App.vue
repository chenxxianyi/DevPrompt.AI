<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import NavBar from '@/components/NavBar.vue'
import AppFooter from '@/components/AppFooter.vue'
import ToastContainer from '@/components/ToastContainer.vue'

const auth = useAuthStore()
onMounted(() => {
  auth.initAuth()
  if (auth.isLoggedIn) {
    auth.fetchProfile()
  }
})
</script>

<template>
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
