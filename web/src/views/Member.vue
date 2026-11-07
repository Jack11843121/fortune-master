<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import PayModal from '../components/PayModal.vue'

const router = useRouter()
const me = ref(null)
const plans = ref(null)
const err = ref('')
const payOrder = ref(null)
const payTitle = ref('')
const needLogin = ref(false)

async function load() {
  try {
    me.value = await api.me()
    needLogin.value = !me.value.id
    plans.value = await api.plans()
  } catch (e) {
    err.value = e.message
  }
}
onMounted(load)

async function subscribe() {
  payTitle.value = '开通 VIP 会员'
  try {
    payOrder.value = await api.subscribe()
  } catch (e) {
    err.value = e.message
  }
}

async function buyQuota() {
  payTitle.value = 'AI 次数充值'
  try {
    payOrder.value = await api.buyQuota()
  } catch (e) {
    err.value = e.message
  }
}

function onPaid() {
  payOrder.value = null
  load()
}

function yuan(cents) {
  return ((cents || 0) / 100).toFixed(2)
}
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <button class="back" style="margin-bottom: 10px" @click="router.push('/me')">‹ 返回</button>

    <div v-if="needLogin" class="card" style="text-align: center">
      <div style="font-size: 40px">👑</div>
      <p style="margin: 8px 0">请先登录后再开通会员</p>
      <button class="btn btn-red" @click="router.push('/me')">去登录</button>
    </div>

    <div class="card" style="text-align: center; background: linear-gradient(150deg, #fdf6e9, #f7e8cd)">
      <div style="font-size: 40px">👑</div>
      <h2 style="margin: 6px 0 2px; color: var(--red-deep)">{{ me?.vip ? '尊享会员' : '开通会员' }}</h2>
      <p style="color: var(--muted); font-size: 14px">
        {{ me?.vip ? `会员有效期至 ${me.vip_until?.slice(0, 10) || ''}` : '月付会员，赠送 AI 算命师调用次数' }}
      </p>
    </div>

    <div class="card">
      <div class="section-head"><span class="ic">📊</span><h3>我的权益</h3></div>
      <div class="stat-row">
        <div class="stat">
          <div class="stat-num">{{ me?.ai_credits ?? 0 }}</div>
          <div class="stat-label">AI 剩余次数</div>
        </div>
        <div class="stat">
          <div class="stat-num">Lv.{{ me?.level ?? 1 }}</div>
          <div class="stat-label">成长等级</div>
        </div>
        <div class="stat">
          <div class="stat-num">{{ me?.growth ?? 0 }}</div>
          <div class="stat-label">成长值</div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="section-head"><span class="ic">💳</span><h3>会员套餐</h3></div>
      <div class="plan-box">
        <div class="plan-left">
          <div class="plan-name">VIP 连续包月</div>
          <div class="plan-desc">全站报告免费解锁 + 每月赠 {{ plans?.monthly_quota ?? 1000 }} 次 AI 算命师</div>
        </div>
        <div class="plan-right">
          <div class="plan-price">¥{{ yuan(plans?.monthly_price_cents) }}</div>
          <button class="btn btn-red btn-sm" @click="subscribe">{{ me?.vip ? '续费' : '立即开通' }}</button>
        </div>
      </div>
      <div class="plan-box" style="margin-top: 10px">
        <div class="plan-left">
          <div class="plan-name">AI 次数充值包</div>
          <div class="plan-desc">{{ plans?.quota_add ?? 1000 }} 次 / 包 · 不限制购买</div>
        </div>
        <div class="plan-right">
          <div class="plan-price">¥{{ yuan(plans?.quota_price_cents) }}</div>
          <button class="btn btn-ghost btn-sm" @click="buyQuota">去充值</button>
        </div>
      </div>
    </div>

    <p v-if="err" style="color: var(--red); font-size: 13px; text-align: center">{{ err }}</p>

    <div class="notice">✨ 会员按自然月计费，到期可续；AI 次数用完可随时充值。</div>

    <PayModal :show="!!payOrder" :order="payOrder" :title="payTitle" @close="payOrder = null" @paid="onPaid" />
  </div>
</template>