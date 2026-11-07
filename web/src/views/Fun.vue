<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { dailySlip } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'

const router = useRouter()
const revealed = ref(false)
const shaking = ref(false)
const slip = ref(dailySlip())

// —— 每日一签：同一天只能抽一次，状态记住在本地 ——
function todayKey() {
  const d = new Date()
  return `${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()}`
}

function checkDrawn() {
  if (localStorage.getItem('daily_slip_date') === todayKey()) {
    slip.value = dailySlip()
    revealed.value = true
  }
}
checkDrawn()

const games = [
  { to: '/fun/dice', e: '🎲', t: '摇色子', d: '点击屏幕或摇晃手机，掷色子问吉凶' },
  { to: '/fun/guess', e: '🖼️', t: '我画你猜', d: '系统出图你来猜，看图猜词、连对有加成' },
  { to: '/fun/flying', e: '✈️', t: '飞行棋', d: '掷骰起飞、击落对方、抽签归家' },
  { to: '/fun/match3', e: '💎', t: '五行消消乐', d: '三消五行，连锁连消，达标抽灵签' },
  { to: '/fun/gomoku', e: '⚫', t: '五子棋', d: '双人对弈或人机，先五连者胜' },
  { to: '/fun/xiangqi', e: '🐴', t: '象棋对弈', d: '楚河汉界，落子对弈，胜负抽签' },
  { to: '/fun/weiqi', e: '⚫⚪', t: '围棋', d: '黑白阴阳，落子提子，终局看阴阳之势' },
  { to: '/fun/music', e: '🎵', t: '五行音乐', d: '五音对应五行，静心聆听、调衡养五脏' },
]

function draw() {
  if (shaking.value || revealed.value) return
  shaking.value = true
  setTimeout(() => {
    shaking.value = false
    slip.value = dailySlip()
    revealed.value = true
    localStorage.setItem('daily_slip_date', todayKey())
  }, 700)
}
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <div class="card" style="display: flex; align-items: center; gap: 14px; padding: 14px 18px; background: linear-gradient(150deg, #fdf6e9, #f7e8cd)">
      <div style="font-size: 34px; line-height: 1">🎐</div>
      <div style="text-align: left">
        <h2 style="margin: 0; color: var(--red-deep); letter-spacing: 2px; font-size: 18px; line-height: 1.3">算乐 · 边玩边算</h2>
        <p style="color: var(--muted); font-size: 13px; margin: 2px 0 0; line-height: 1.4">游戏里藏天机，通关对弈，抽支灵签乐一乐</p>
      </div>
    </div>

    <!-- 每日一签 -->
    <div class="card" style="text-align: center">
      <div class="section-head" style="justify-content: center"><span class="ic">🔮</span><h3>每日一签</h3></div>
      <div v-if="!revealed" style="padding: 10px 0">
        <div class="qiantong" :class="{ shaking }" style="font-size: 52px">🪔</div>
        <p style="color: var(--muted); font-size: 14px">今日运势，抽一支看看？</p>
        <button class="btn btn-red btn-sm" style="max-width: 220px; margin: 8px auto" @click="draw">抽取今日灵签</button>
      </div>
      <div v-else style="padding: 6px 0">
        <SlipCard :slip="slip" />
        <p style="font-size: 12px; color: var(--muted); margin-top: 6px">明日再来 · 好运常伴</p>
      </div>
    </div>

    <!-- 游戏列表 -->
    <div v-for="g in games" :key="g.to" class="scene-card" @click="router.push(g.to)">
      <div class="sc-icon">
        <svg v-if="g.to === '/fun/weiqi'" class="go-icon" viewBox="0 0 36 36" aria-hidden="true">
          <rect x="1.5" y="1.5" width="33" height="33" rx="6" fill="#e6b877" stroke="#b98f4e" stroke-width="2" />
          <g stroke="#7c5a2a" stroke-width="1.2">
            <line x1="9" y1="9" x2="9" y2="27" />
            <line x1="18" y1="9" x2="18" y2="27" />
            <line x1="27" y1="9" x2="27" y2="27" />
            <line x1="9" y1="9" x2="27" y2="9" />
            <line x1="9" y1="18" x2="27" y2="18" />
            <line x1="9" y1="27" x2="27" y2="27" />
          </g>
          <circle cx="18" cy="18" r="4" fill="#222" />
          <circle cx="9" cy="9" r="4" fill="#fdfdf6" stroke="#9a9a9a" stroke-width="0.8" />
        </svg>
        <span v-else class="emoji">{{ g.e }}</span>
      </div>
      <div>
        <h3>{{ g.t }}</h3>
        <p>{{ g.d }}</p>
      </div>
      <div class="arrow">›</div>
    </div>

    <div class="notice">🎮 游戏适中为宜，久坐伤身，劳逸结合更安康。</div>
  </div>
</template>