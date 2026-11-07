<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaziTable from '../components/BaziTable.vue'
import PayModal from '../components/PayModal.vue'
import { api } from '../api'

const route = useRoute()
const router = useRouter()
const data = ref(null)
const copied = ref(false)
const loading = ref(true)
const err = ref('')
const payOrder = ref(null)
const showPay = ref(false)

const icons = {
  compass: '🧭', gem: '💎', heart: '❤️', paw: '🐾', briefcase: '💼',
  coins: '💰', users: '👥', shield: '🛡️', book: '📖', sparkles: '✨',
  trend: '📈', calendar: '📅', star: '⭐',
}

async function load() {
  loading.value = true
  err.value = ''
  try {
    data.value = await api.getReading(route.params.id)
    localStorage.setItem('last_reading_id', route.params.id)
  } catch (e) {
    err.value = e.message
  } finally {
    loading.value = false
  }
}
onMounted(load)

const priceText = computed(() => ((data.value?.price_cents || 0) / 100).toFixed(2))

async function unlock() {
  err.value = ''
  try {
    const order = await api.createOrder(data.value.id)
    if (order.already_paid) {
      await load()
      return
    }
    payOrder.value = order
    showPay.value = true
  } catch (e) {
    err.value = e.message
  }
}

function onPaid() {
  showPay.value = false
  load()
}

async function shareReading() {
  const link = `${location.origin}${location.pathname}#/result/${route.params.id}`
  try {
    await navigator.clipboard?.writeText(link)
    copied.value = true
    setTimeout(() => (copied.value = false), 1500)
  } catch (e) {}
  api.share(route.params.id).catch(() => {})
  if (navigator.share) {
    try {
      await navigator.share({ title: data.value?.title || '算命大师', url: link })
    } catch (e) {}
  }
}

function askAI() {
  router.push('/ai')
}
</script>

<template>
  <div class="page">
    <div v-if="loading" class="skeleton">
      <div class="spin"></div>
      <p>正在排盘推演…</p>
    </div>

    <div v-else-if="err" class="card">
      <p style="color: var(--red)">{{ err }}</p>
      <button class="btn btn-ghost" @click="$router.push('/')">返回首页</button>
    </div>

    <template v-else>
      <h2 class="page-title" style="margin-top: 16px">{{ data.title }}</h2>
      <div style="display: flex; gap: 8px; margin-bottom: 12px">
        <button class="btn btn-ghost btn-sm" style="width:auto" @click="shareReading">{{ copied ? '已复制链接' : '分享' }}</button>
        <button class="btn btn-ghost btn-sm" style="width:auto" @click="askAI">🤖 问问 AI 大师</button>
      </div>

      <!-- 匹配评分 -->
      <div v-if="data.match" class="card">
        <div class="score-hero">
          <div class="score-num">{{ data.match.score.value }}</div>
          <div style="color: var(--muted); font-size: 14px">综合匹配度</div>
          <span class="score-grade">{{ data.match.score.grade }}</span>
        </div>
      </div>

      <!-- 单/双人排盘 -->
      <BaziTable :bazi="data.bazi" />
      <BaziTable v-if="data.bazi_b" :bazi="data.bazi_b" />

      <!-- 匹配各维度 -->
      <div v-if="data.match" class="card">
        <div class="section-head"><span class="ic">🔍</span><h3>匹配维度解析</h3></div>
        <div v-for="a in data.match.aspects" :key="a.key" class="aspect">
          <div class="a-head">
            <span>{{ a.good ? '✅' : '⚠️' }}</span>{{ a.title }}
            <span class="tag" :style="{ background: a.good ? '#e4f2e7' : '#fbe9e4', color: a.good ? 'var(--green)' : 'var(--red)' }">
              {{ a.good ? '吉' : '需注意' }}
            </span>
          </div>
          <div class="a-body">{{ a.body }}</div>
        </div>

        <div v-if="data.match.conclusion" class="card" style="box-shadow:none; border:1px dashed var(--gold); margin-top:8px">
          <div class="section-head"><span class="ic">💡</span><h3>综合结论</h3></div>
          <p style="margin:0">{{ data.match.conclusion }}</p>
          <div class="section-head" style="margin-top:14px"><span class="ic">📌</span><h3>给你的建议</h3></div>
          <p style="margin:0">{{ data.match.suggestion }}</p>
        </div>
      </div>

      <!-- 解读段落（个人命理） -->
      <div v-for="s in data.sections" :key="s.key" class="card section-card">
        <div class="section-head">
          <span class="ic">{{ icons[s.icon] || '📖' }}</span>
          <h3>{{ s.title }}</h3>
          <span v-if="s.score" class="score-pill">{{ s.score.value }}分 {{ s.score.grade }}</span>
        </div>
        <div class="section-body">{{ s.body }}</div>
      </div>

      <!-- 解锁入口 -->
      <div v-if="data.locked" class="card" style="border: 2px solid var(--gold-deep); background: #fdf6e9">
        <div class="section-head"><span class="ic">🔒</span><h3>还有 {{ data.locked_count }} 项深度解析待解锁</h3></div>

        <div v-if="data.teaser" class="teaser">
          <div class="teaser-tag">深度解析 · 预览</div>
          <div class="teaser-body">{{ data.teaser }}</div>
        </div>

        <div class="locked-list">
          <div v-for="t in data.locked_titles || []" :key="t" class="locked-item">🔒 {{ t }}</div>
        </div>

        <button class="btn btn-red" @click="unlock">解锁全部深度解析　¥{{ priceText }}</button>
        <p style="font-size: 12px; color: var(--muted); text-align: center; margin: 10px 0 0">
          含 {{ data.locked_count }} 大板块 · 解锁后永久可看本次结果
        </p>
      </div>

      <!-- 已解锁提示 -->
      <div v-else class="notice">✅ 已解锁全文，感谢支持。愿君知命改运、顺遂如意。</div>

      <p v-if="err" style="color: var(--red); font-size: 14px">{{ err }}</p>
    </template>

    <PayModal :show="showPay" :order="payOrder" @close="showPay = false" @paid="onPaid" />
  </div>
</template>