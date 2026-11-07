<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from './api'
const route = useRoute()

// 自动演示登录：用户无感免注册，直接以演示身份使用全部功能
onMounted(async () => {
  if (!api.isLoggedIn()) {
    try {
      const r = await api.authDev({ nickname: '演示用户' })
      api.setToken(r.token)
    } catch (e) { /* 静默失败，不影响浏览 */ }
  }
})

const tabs = [
  { to: '/', ic: '🏠', label: '首页' },
  { to: '/fun', ic: '🎮', label: '算乐' },
  { to: '/me', ic: '👤', label: '我的' },
]

function isActive(path) {
  if (path === '/') {
    return ['/', '/personal', '/wealth', '/marriage', '/hiring', '/naming', '/result'].some((p) => route.path.startsWith(p))
  }
  if (path === '/me') {
    return ['/me', '/member', '/ai', '/community'].some((p) => route.path.startsWith(p))
  }
  return route.path.startsWith(path)
}
</script>

<template>
  <header class="header">
    <h1>算命大师</h1>
    <div class="sub">知命 · 明理 · 改运</div>
  </header>

  <router-view />

  <nav class="tabbar">
    <router-link v-for="t in tabs" :key="t.to" :to="t.to" :class="{ active: isActive(t.to) }">
      <span class="ic">{{ t.ic }}</span>{{ t.label }}
    </router-link>
  </nav>
</template>