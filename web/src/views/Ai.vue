<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const msgs = ref([])
const input = ref('')
const session = ref('')
const busy = ref(false)
const credits = ref(0)
const err = ref('')
const box = ref(null)
const needLogin = ref(false)

onMounted(async () => {
  if (!api.isLoggedIn()) {
    needLogin.value = true
    return
  }
  try {
    const m = await api.me()
    credits.value = m.ai_credits || 0
    const hist = await api.aiHistory('')
    if (hist.messages?.length) {
      msgs.value = hist.messages.map((x) => ({ role: x.role, content: x.content, time: x.created_at }))
    }
  } catch (e) {
    err.value = e.message
  }
})

async function send() {
  const text = input.value.trim()
  if (!text || busy.value) return
  msgs.value.push({ role: 'user', content: text, time: now() })
  input.value = ''
  busy.value = true
  err.value = ''
  nextTick(scroll)
  try {
    const r = await api.aiChat({ session_id: session.value, message: text, reading_id: lastReadingId() })
    session.value = r.session_id
    msgs.value.push({ role: 'assistant', content: r.reply, time: now() })
    credits.value = r.ai_credits
  } catch (e) {
    err.value = e.message
  } finally {
    busy.value = false
    nextTick(scroll)
  }
}

function now() {
  const d = new Date()
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

function lastReadingId() {
  return localStorage.getItem('last_reading_id') || ''
}

function scroll() {
  if (box.value) box.value.scrollTop = box.value.scrollHeight
}
</script>

<template>
  <div v-if="needLogin" class="page" style="padding-top: 60px; text-align: center">
    <div style="font-size: 48px">🤖</div>
    <p style="color: var(--muted); margin: 10px 0">请先登录后使用 AI 算命师</p>
    <button class="btn btn-red" style="max-width: 240px; margin: 0 auto" @click="router.push('/me')">去登录</button>
  </div>

  <div v-else class="chat-page">
    <!-- 微信式聊天头部 -->
    <div class="chat-header">
      <button class="back" @click="router.push('/me')">‹</button>
      <div class="chat-avatar master">{{ credits > 0 ? '🧙' : '😴' }}</div>
      <div class="chat-head-info">
        <div class="chat-head-name">AI 算命大师</div>
        <div class="chat-head-status">{{ credits > 0 ? '在线 · 余 ' + credits + ' 次' : '次数已用完' }}</div>
      </div>
      <button class="chat-head-btn" @click="router.push('/member')">👑</button>
    </div>

    <!-- 对话区 -->
    <div ref="box" class="chat-body">
      <div v-if="msgs.length === 0" class="chat-empty">
        <div class="chat-avatar master" style="margin: 0 auto 10px; width: 64px; height: 64px; font-size: 36px">🧙</div>
        <div style="color: var(--muted); line-height: 2">
          您好，我是算命大师的 AI 命理师<br />事业、婚恋、财运、起名……尽管来问
        </div>
      </div>

      <div v-for="(m, i) in msgs" :key="i" class="chat-row" :class="m.role === 'user' ? 'user' : 'ai'">
        <div v-if="m.role !== 'user'" class="chat-avatar master">🧙</div>
        <div class="chat-col">
          <div class="chat-bubble">{{ m.content }}</div>
        </div>
        <div v-if="m.role === 'user'" class="chat-avatar me">🙂</div>
      </div>

      <div v-if="busy" class="chat-row ai">
        <div class="chat-avatar master">🧙</div>
        <div class="chat-bubble typing"><span></span><span></span><span></span></div>
      </div>

      <p v-if="err" style="color: var(--red); font-size: 12px; text-align: center; margin: 8px 0">{{ err }}</p>
    </div>

    <!-- 输入栏 -->
    <div class="chat-input-bar">
      <input class="chat-input" v-model="input" placeholder="问财运、问姻缘、问事业…" @keyup.enter="send" :disabled="credits <= 0" />
      <button class="chat-send" :disabled="busy || credits <= 0" @click="send">{{ busy ? '…' : '发送' }}</button>
    </div>

    <button v-if="credits <= 0" class="btn btn-primary" style="margin: 10px 12px" @click="router.push('/member')">次数用完了，去开通会员或充值 ›</button>
  </div>
</template>