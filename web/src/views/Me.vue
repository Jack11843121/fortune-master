<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const me = ref(null)
const loggedIn = ref(false)
const err = ref('')
const isDev = ref(false)

onMounted(async () => {
  // 判断是否本地/演示环境（mock 支付），用于控制「测试登录」是否展示
  try {
    const meta = await api.meta()
    isDev.value = meta.pay_provider === 'mock'
  } catch (e) {}
  // OAuth 回调：微信/支付宝把 code 带回
  const code = new URLSearchParams(location.search).get('code')
  const pending = localStorage.getItem('oauth_pending')
  if (code && pending === 'wechat') {
    try {
      const r = await api.authWechatOAuth(code)
      api.setToken(r.token)
      me.value = r.user
      loggedIn.value = true
    } catch (e) {
      err.value = e.message
    } finally {
      localStorage.removeItem('oauth_pending')
      history.replaceState({}, '', location.pathname + location.hash)
    }
    return
  }
  await loadProfile()
})

async function loadProfile() {
  if (!api.isLoggedIn()) {
    loggedIn.value = false
    me.value = null
    return
  }
  try {
    const p = await api.authProfile()
    me.value = p
    loggedIn.value = true
  } catch (e) {
    api.clearToken()
    loggedIn.value = false
  }
}

async function loginWechat() {
  try {
    localStorage.setItem('oauth_pending', 'wechat')
    const r = await api.authWechatUrl(location.origin + location.pathname)
    location.href = r.url
  } catch (e) {
    err.value = e.message
  }
}

async function loginDev() {
  err.value = ''
  try {
    const r = await api.authDev({})
    api.setToken(r.token)
    me.value = r.user
    loggedIn.value = true
  } catch (e) {
    err.value = e.message
  }
}

const menus = [
  { to: '/member', e: '👑', t: '会员中心', d: '开通会员 · AI 次数 · 成长等级' },
  { to: '/ai', e: '🤖', t: 'AI 算命师', d: '随问随答，资深命理师在线开解' },
  { to: '/community', e: '💬', t: '交流社区', d: '分享命理趣事、合婚故事、问卦交流' },
  { to: '/doc/agreement', e: '📄', t: '用户协议', d: '服务协议 · 隐私政策 · 免责声明' },
]
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <!-- 未登录 -->
    <div v-if="!loggedIn" class="card" style="text-align: center">
      <div style="font-size: 46px">🙂</div>
      <h3 style="margin: 8px 0 4px">登录后体验完整服务</h3>
      <p style="color: var(--muted); font-size: 13px">会员、AI 算命师、社区发帖都需登录账户</p>
      <button class="btn btn-green" style="max-width: 260px; margin: 12px auto 8px" @click="loginWechat">微信登录</button>
      <button v-if="isDev" class="btn btn-ghost btn-sm" style="max-width: 260px; margin: 0 auto" @click="loginDev">测试登录（本地演示）</button>
      <p style="color: var(--muted); font-size: 12px; margin: 12px 0 0; line-height: 1.6">
        登录后会员权益、AI 次数、社区内容将绑定到你的微信账号，换设备、换浏览器数据都不会丢。
      </p>
      <p v-if="err" style="color: var(--red); font-size: 12px; margin: 10px 0 0">{{ err }}</p>
    </div>

    <!-- 已登录 -->
    <template v-else>
      <div class="card" style="display: flex; align-items: center; gap: 14px; background: linear-gradient(150deg, #fdf6e9, #f7e8cd)">
        <img v-if="me?.avatar" :src="me.avatar" class="chat-avatar" style="width: 56px; height: 56px; border-radius: 8px" />
        <div v-else class="chat-avatar me" style="width: 56px; height: 56px; font-size: 30px">🙂</div>
        <div style="flex: 1">
          <div style="display: flex; align-items: center; gap: 8px">
            <b style="font-size: 18px">{{ me?.nickname || '命主' }}</b>
            <span v-if="me?.vip" class="badge badge-red">VIP</span>
          </div>
          <div style="color: var(--muted); font-size: 13px; margin-top: 3px">
            AI 余 {{ me?.ai_credits ?? 0 }} 次 · Lv.{{ me?.level ?? 1 }} · 成长 {{ me?.growth ?? 0 }}
          </div>
        </div>
        <button v-if="!me?.vip" class="btn btn-red btn-sm" style="width: auto" @click="router.push('/member')">开通会员</button>
      </div>
    </template>

    <div class="card" style="padding: 4px 16px">
      <div v-for="m in menus" :key="m.to" class="scene-card" @click="router.push(m.to)">
        <div class="sc-icon"><span class="emoji">{{ m.e }}</span></div>
        <div style="flex: 1">
          <h3>{{ m.t }}</h3>
          <p>{{ m.d }}</p>
        </div>
        <div class="arrow">›</div>
      </div>
    </div>

    <div style="text-align: center; font-size: 12px; color: var(--muted); padding: 8px 0 20px; line-height: 1.8">
      © 2026 算命大师 · 传统文化娱乐与心理参考
    </div>
  </div>
</template>