<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { randomSlip } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'
import { createBoard, tryPlace, aiMove, countScore, boardString } from '../game/weiqi'

const router = useRouter()
const STARS = [[2, 2], [2, 6], [6, 2], [6, 6], [4, 4]]
const KOMI = 6.5

const board = ref(createBoard())
const turn = ref(1) // 1 黑 2 白
const mode = ref('pvp') // pvp | ai
const captured = ref({ 1: 0, 2: 0 })
const passCount = ref(0)
const lastMove = ref(null)
const finished = ref(false)
const winner = ref(null)
const scoreResult = ref(null)
const aiThinking = ref(false)
let history = []

reset()
function reset() {
  board.value = createBoard()
  turn.value = 1
  captured.value = { 1: 0, 2: 0 }
  passCount.value = 0
  lastMove.value = null
  finished.value = false
  winner.value = null
  scoreResult.value = null
  aiThinking.value = false
  history = [boardString(board.value)]
}

const isStar = (r, c) => STARS.some(([sr, sc]) => sr === r && sc === c)

function placeStone(r, c, color) {
  if (finished.value || winner.value) return false
  const koStr = history.length >= 2 ? history[history.length - 2] : null
  const res = tryPlace(board.value, color, r, c, koStr)
  if (!res) return false
  board.value = res.board
  captured.value = { ...captured.value, [color]: captured.value[color] + res.captured }
  lastMove.value = { r, c }
  history.push(boardString(board.value))
  passCount.value = 0
  turn.value = color === 1 ? 2 : 1
  return true
}

function tap(r, c) {
  if (mode.value === 'ai' && (turn.value !== 1 || aiThinking.value)) return
  if (placeStone(r, c, turn.value) && mode.value === 'ai' && !finished.value) scheduleAI()
}

function scheduleAI() {
  aiThinking.value = true
  setTimeout(() => {
    const koStr = history.length >= 2 ? history[history.length - 2] : null
    const m = aiMove(board.value, 2, koStr)
    aiThinking.value = false
    if (finished.value || winner.value) return
    if (m) placeStone(m.r, m.c, 2)
    else {
      history.push(boardString(board.value))
      passCount.value++
      if (passCount.value >= 2) finish()
      else turn.value = 1
    }
  }, 280)
}

function pass() {
  if (finished.value || winner.value) return
  if (mode.value === 'ai' && (turn.value !== 1 || aiThinking.value)) return
  history.push(boardString(board.value))
  passCount.value++
  if (passCount.value >= 2) finish()
  else {
    turn.value = turn.value === 1 ? 2 : 1
    if (mode.value === 'ai' && turn.value === 2) scheduleAI()
  }
}

function resign() {
  if (finished.value) return
  winner.value = turn.value === 1 ? 2 : 1
  finished.value = true
}

function finish() {
  finished.value = true
  const s = countScore(board.value)
  const white = s.white + KOMI
  scoreResult.value = { black: s.black, white, blackRaw: s.black, whiteRaw: s.white }
  winner.value = s.black > white ? 1 : 2
}

function switchMode(m) {
  mode.value = m
  reset()
}

const slip = computed(() => (winner.value ? randomSlip() : null))

const balance = computed(() => {
  if (!scoreResult.value) return ''
  const diff = scoreResult.value.black - scoreResult.value.white
  if (Math.abs(diff) < 5) return '阴阳调和，势均力敌，和风细雨。'
  return diff > 0 ? '黑（阴）势大，聚敛沉稳，厚积薄发。' : '白（阳）势广，舒展光明，先声夺人。'
})
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">围棋</span>
      <button class="back" @click="reset">↺</button>
    </div>

    <div class="m3-status">
      <div class="seg" style="flex: 2">
        <button :class="{ on: mode === 'pvp' }" type="button" @click="switchMode('pvp')">双人对弈</button>
        <button :class="{ on: mode === 'ai' }" type="button" @click="switchMode('ai')">人机对战</button>
      </div>
    </div>
    <p class="m3-hint" style="text-align: center; margin: 8px 0 2px">
      黑提 <b>{{ captured[1] }}</b> 子 · 白提 <b>{{ captured[2] }}</b> 子
      ｜ {{ winner ? (winner === 1 ? '黑方' : '白方') + '胜 🎉' : aiThinking ? '思考中…' : (turn === 1 ? '黑方' : '白方') + '走棋' }}
    </p>

    <div class="wq-board">
      <template v-for="(row, r) in board" :key="r">
        <button
          v-for="(v, c) in row"
          :key="r + '-' + c"
          type="button"
          class="wq-cell"
          :class="{ last: lastMove && lastMove.r === r && lastMove.c === c }"
          @click="tap(r, c)"
        >
          <span v-if="!v && isStar(r, c)" class="wq-star"></span>
          <span v-if="v" class="stone" :class="v === 1 ? 'black' : 'white'"></span>
        </button>
      </template>
    </div>

    <div style="display: flex; gap: 8px; margin: 6px 0">
      <button class="btn btn-ghost btn-sm" @click="pass">停一手</button>
      <button class="btn btn-ghost btn-sm" style="color: var(--red)" @click="resign">认输</button>
      <button class="btn btn-primary btn-sm" @click="finish">终局算目</button>
    </div>

    <div class="modal-mask" v-if="finished">
      <div class="modal">
        <h3 style="margin-bottom: 6px">🎉 {{ winner === 1 ? '黑方' : '白方' }}胜出</h3>
        <p v-if="scoreResult" style="color: var(--muted); font-size: 14px; margin: 0 0 4px">
          黑 {{ scoreResult.blackRaw }} 子 · 白 {{ scoreResult.whiteRaw }} 子（贴 {{ KOMI }} 目）
        </p>
        <p style="font-size: 14px; color: var(--red-deep); margin: 4px 0 10px">{{ balance }}</p>
        <SlipCard :slip="slip" />
        <button class="btn btn-red" style="margin-top: 16px" @click="reset">再来一局</button>
        <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="router.push('/fun')">返回算乐</button>
      </div>
    </div>
  </div>
</template>