<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { SLIPS } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'

const router = useRouter()

const dice = ref([1, 1, 1])
const rolling = ref(false)
const result = ref(null) // { total, label, slip }

const DOTS = { 1: [4], 2: [0, 8], 3: [0, 4, 8], 4: [0, 2, 6, 8], 5: [0, 2, 4, 6, 8], 6: [0, 2, 3, 5, 6, 8] }
const pips = (v) => Array.from({ length: 9 }, (_, i) => DOTS[v].includes(i))
const rnd6 = () => Math.floor(Math.random() * 6) + 1
const sleep = (ms) => new Promise((r) => setTimeout(r, ms))

const GRADE = [
  { min: 16, star: 5, label: '大吉' },
  { min: 13, star: 4, label: '小吉' },
  { min: 10, star: 3, label: '中平' },
  { min: 7, star: 2, label: '小凶' },
  { min: 3, star: 1, label: '凶' },
]

function buildResult(total) {
  const g = GRADE.find((x) => total >= x.min)
  const pool = SLIPS.filter((s) => s.star === g.star)
  const slip = pool[Math.floor(Math.random() * pool.length)]
  return { total, label: g.label, slip }
}

async function roll() {
  if (rolling.value) return
  rolling.value = true
  result.value = null
  for (let i = 0; i < 9; i++) {
    dice.value = [rnd6(), rnd6(), rnd6()]
    await sleep(60 + i * 12)
  }
  const vals = [rnd6(), rnd6(), rnd6()]
  dice.value = vals
  const total = vals.reduce((a, b) => a + b, 0)
  rolling.value = false
  result.value = buildResult(total)
}

// —— 摇一摇（设备加速度）——
let motionReady = false
let lastAcc = null
let lastShake = 0

function onMotion(e) {
  const a = e.accelerationIncludingGravity
  if (!a) return
  const x = a.x || 0, y = a.y || 0, z = a.z || 0
  const now = Date.now()
  if (lastAcc) {
    const d = Math.abs(x - lastAcc.x) + Math.abs(y - lastAcc.y) + Math.abs(z - lastAcc.z)
    if (d > 32 && now - lastShake > 1400) {
      lastShake = now
      roll()
    }
  }
  lastAcc = { x, y, z }
}

function setupMotion() {
  if (motionReady || typeof DeviceMotionEvent === 'undefined') return
  if (typeof DeviceMotionEvent.requestPermission === 'function') {
    DeviceMotionEvent.requestPermission()
      .then(() => {
        window.addEventListener('devicemotion', onMotion)
        motionReady = true
      })
      .catch(() => {})
  } else {
    window.addEventListener('devicemotion', onMotion)
    motionReady = true
  }
}

function onTap() {
  setupMotion() // iOS 需在用户手势中申请权限
  roll()
}

onMounted(() => setupMotion())
onBeforeUnmount(() => window.removeEventListener('devicemotion', onMotion))
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">摇色子</span>
      <button class="back" @click="roll">↺</button>
    </div>

    <div class="dice-stage" @click="onTap">
      <div class="dice-row">
        <div v-for="(v, i) in dice" :key="i" class="dice" :class="{ shaking: rolling }">
          <span v-for="j in 9" :key="j" class="pip-slot">
            <span v-if="pips(v)[j - 1]" class="pip"></span>
          </span>
        </div>
      </div>
      <p class="dice-hint">{{ rolling ? '摇动中……' : '点击屏幕 或 摇晃手机' }}</p>
    </div>

    <div v-if="result" class="card" style="text-align: center">
      <div class="dice-total">点数 {{ result.total }} · {{ result.label }}</div>
      <SlipCard :slip="result.slip" />
    </div>

    <p class="m3-hint" style="text-align: center">摇出好兆头 · 愿君心想事成</p>
  </div>
</template>