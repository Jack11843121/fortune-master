<script setup>
import { ref, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { SKETCHES } from '../data/sketches'
import { randomSlip } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'
import SketchDraw from '../components/SketchDraw.vue'

const router = useRouter()
const DURATION = 600

const pool = ref([])
const idx = ref(0)
const seq = ref(0) // 每换一题 +1，用于让 SketchDraw 重新作画
const score = ref(0)
const combo = ref(0)
const correct = ref(0)
const timeLeft = ref(DURATION)
const guess = ref('')
const feedback = ref('')
const started = ref(false)
const finished = ref(false)
const hintShown = ref(false)
const drawing = ref(false)
let timer = null
let drawTimer = null

const cur = computed(() => pool.value[idx.value] || null) // { word, cat, shapes }
const timeText = computed(() => `${Math.floor(timeLeft.value / 60)}:${String(timeLeft.value % 60).padStart(2, '0')}`)
const wordLen = computed(() => (cur.value ? cur.value.word.length : 0))
const slip = computed(() => (finished.value ? randomSlip() : null))
const grade = computed(() => {
  const c = correct.value
  if (c >= 15) return '火眼金睛'
  if (c >= 10) return '慧眼如炬'
  if (c >= 6) return '眼力不错'
  if (c >= 3) return '初窥门径'
  return '再接再厉'
})

function start() {
  pool.value = [...SKETCHES].sort(() => Math.random() - 0.5)
  idx.value = 0
  seq.value++
  score.value = 0
  combo.value = 0
  correct.value = 0
  timeLeft.value = DURATION
  guess.value = ''
  feedback.value = ''
  hintShown.value = false
  started.value = true
  finished.value = false
  clearInterval(timer)
  timer = setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) end()
  }, 1000)
  markDrawing()
}

function end() {
  clearInterval(timer)
  finished.value = true
}

function next() {
  guess.value = ''
  hintShown.value = false
  feedback.value = ''
  idx.value = (idx.value + 1) % pool.value.length
  if (idx.value === 0) pool.value = [...SKETCHES].sort(() => Math.random() - 0.5)
  seq.value++
  markDrawing()
}

function markDrawing() {
  const n = cur.value ? cur.value.shapes.length : 1
  drawing.value = true
  clearTimeout(drawTimer)
  drawTimer = setTimeout(() => (drawing.value = false), Math.round(n * 380 + 650))
}

function submit() {
  if (finished.value || !cur.value) return
  const g = guess.value.trim()
  if (!g) return
  if (g === cur.value.word) {
    combo.value++
    score.value += 10 + combo.value * 2
    correct.value++
    feedback.value = combo.value >= 2 ? `✅ 答对了！连对 ×${combo.value}` : '✅ 答对了！'
    next()
  } else {
    feedback.value = '❌ 不对，再想想～'
    guess.value = ''
  }
}

function trySubmit(e) {
  if (e.isComposing || e.keyCode === 229) return
  submit()
}

function hint() {
  if (finished.value || !cur.value) return
  if (!hintShown.value) {
    hintShown.value = true
    score.value = Math.max(0, score.value - 2)
  }
}

function skip() {
  if (finished.value) return
  combo.value = 0
  next()
}

function redraw() {
  seq.value++
  markDrawing()
}

onBeforeUnmount(() => {
  clearInterval(timer)
  clearTimeout(drawTimer)
})
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">我画你猜</span>
      <button class="back" @click="redraw">↺</button>
    </div>

    <!-- 未开始 -->
    <div v-if="!started" class="card" style="text-align: center">
      <div style="font-size: 46px">🖌️</div>
      <h3 style="margin: 8px 0; color: var(--red-deep)">系统一笔一划作画，你来猜</h3>
      <p style="color: var(--muted); font-size: 14px; line-height: 1.8; text-align: left">
        · 系统用笔迹一笔笔勾出图案，你看轮廓猜它是什么<br />
        · 输入完整词语，答对得分、连对有加成<br />
        · 「提示」露出类别（扣 2 分），「换一题」跳过<br />
        · 限时 10 分钟，看你能猜中多少
      </p>
      <button class="btn btn-red" style="margin-top: 14px" @click="start">开始挑战</button>
    </div>

    <!-- 游戏中 -->
    <template v-else-if="!finished">
      <div class="m3-status">
        <div class="m3-cell">⏱ {{ timeText }}</div>
        <div class="m3-cell">⭐ {{ score }}</div>
        <div class="m3-cell">✅ {{ correct }} 题</div>
      </div>

      <div class="card" style="text-align: center; padding: 20px 16px">
        <div class="draw-stage">
          <SketchDraw v-if="cur" :key="seq" :shapes="cur.shapes" />
        </div>
        <div style="font-size: 14px; color: var(--muted); margin-top: 8px">
          {{ drawing ? '🎨 系统作画中…' : '' }} {{ hintShown ? cur[2] : '？' }} · 共 {{ wordLen }} 字
        </div>
        <div style="min-height: 20px; font-size: 14px; color: var(--red-deep); margin: 6px 0">{{ feedback }}</div>
        <div style="display: flex; gap: 8px; margin-top: 8px">
          <input class="input" v-model="guess" placeholder="请输入完整词语…" @keydown.enter="trySubmit" />
          <button class="btn btn-primary btn-sm" style="width: auto; padding: 12px 18px" @click="submit">提交</button>
        </div>
        <p style="font-size: 12px; color: var(--muted); margin: 8px 0 0">请输入完整的词语（如「太阳」），不是只填一个空</p>
        <div style="display: flex; gap: 8px; margin-top: 12px">
          <button class="btn btn-ghost btn-sm" @click="hint">💡 提示（-2）</button>
          <button class="btn btn-ghost btn-sm" @click="skip">↷ 换一题</button>
        </div>
      </div>
    </template>

    <!-- 结束 -->
    <div v-else class="card" style="text-align: center">
      <div style="font-size: 40px">🏆</div>
      <h3 style="margin: 6px 0">{{ grade }}</h3>
      <p style="color: var(--muted); font-size: 15px">猜对 <b style="color: var(--red-deep)">{{ correct }}</b> 题 · 得分 {{ score }}</p>
      <SlipCard :slip="slip" />
      <button class="btn btn-red" style="margin-top: 14px" @click="start">再来一局</button>
      <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="router.push('/fun')">返回算乐</button>
    </div>
  </div>
</template>