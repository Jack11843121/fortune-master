<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { randomSlip } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'

const router = useRouter()
const N = 15
const board = ref(new Array(N * N).fill(0)) // 0 空 1 黑 2 白
const turn = ref(1) // 1 黑 2 白
const mode = ref('pvp') // pvp | ai
const winner = ref(0)
const winLine = ref([])
const lastMove = ref(-1)
const status = ref('')

const curName = computed(() => (turn.value === 1 ? '黑方' : '白方'))

init()
function init() {
  board.value = new Array(N * N).fill(0)
  turn.value = 1
  winner.value = 0
  winLine.value = []
  lastMove.value = -1
  status.value = ''
}

const idx = (r, c) => r * N + c
const inBoard = (r, c) => r >= 0 && r < N && c >= 0 && c < N

function play(i) {
  if (winner.value || board.value[i] !== 0) return
  board.value[i] = turn.value
  lastMove.value = i
  const line = winAt(i)
  if (line) {
    winner.value = turn.value
    winLine.value = line
  } else if (board.value.every((v) => v !== 0)) {
    winner.value = 3 // 和棋
  } else {
    turn.value = turn.value === 1 ? 2 : 1
    if (mode.value === 'ai' && turn.value === 2 && !winner.value) {
      setTimeout(aiMove, 260)
    }
  }
}

function winAt(i) {
  const r = Math.floor(i / N)
  const c = i % N
  const color = board.value[i]
  const dirs = [[1, 0], [0, 1], [1, 1], [1, -1]]
  for (const [dr, dc] of dirs) {
    const line = [i]
    let rr = r + dr, cc = c + dc
    while (inBoard(rr, cc) && board.value[idx(rr, cc)] === color) line.push(idx(rr, cc)), (rr += dr), (cc += dc)
    rr = r - dr, cc = c - dc
    while (inBoard(rr, cc) && board.value[idx(rr, cc)] === color) line.push(idx(rr, cc)), (rr -= dr), (cc -= dc)
    if (line.length >= 5) return line
  }
  return null
}

// —— 简单 AI（启发式评分，不深度搜索）——
function lineScore(i, dr, dc, color) {
  const r = Math.floor(i / N)
  const c = i % N
  let cnt = 1
  let block = 0
  let rr = r + dr, cc = c + dc
  while (inBoard(rr, cc) && board.value[idx(rr, cc)] === color) cnt++, (rr += dr), (cc += dc)
  if (!inBoard(rr, cc) || board.value[idx(rr, cc)] !== 0) block++
  rr = r - dr, cc = c - dc
  while (inBoard(rr, cc) && board.value[idx(rr, cc)] === color) cnt++, (rr -= dr), (cc -= dc)
  if (!inBoard(rr, cc) || board.value[idx(rr, cc)] !== 0) block++
  const open = 2 - block
  if (cnt >= 5) return 100000
  if (cnt === 4) return open === 2 ? 10000 : open === 1 ? 1000 : 0
  if (cnt === 3) return open === 2 ? 1500 : open === 1 ? 300 : 0
  if (cnt === 2) return open === 2 ? 100 : open === 1 ? 30 : 0
  return open === 2 ? 10 : 2
}

function evalCell(i, color) {
  const dirs = [[1, 0], [0, 1], [1, 1], [1, -1]]
  let total = 0
  for (const [dr, dc] of dirs) total += lineScore(i, dr, dc, color)
  return total
}

function aiMove() {
  if (winner.value || turn.value !== 2) return
  let best = -1
  let bestIdx = -1
  for (let i = 0; i < N * N; i++) {
    if (board.value[i] !== 0) continue
    const sc = evalCell(i, 2) + evalCell(i, 1) * 1.1
    if (sc > best) {
      best = sc
      bestIdx = i
    }
  }
  if (bestIdx >= 0) play(bestIdx)
}

function switchMode(m) {
  mode.value = m
  init()
}

const slip = computed(() => (winner.value === 1 || winner.value === 2 ? randomSlip() : null))
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">五子棋</span>
      <button class="back" @click="init">↺</button>
    </div>

    <div class="m3-status">
      <div class="seg" style="flex: 2">
        <button :class="{ on: mode === 'pvp' }" type="button" @click="switchMode('pvp')">双人对弈</button>
        <button :class="{ on: mode === 'ai' }" type="button" @click="switchMode('ai')">人机对战</button>
      </div>
    </div>
    <p class="m3-hint" style="margin-top: 10px">
      <template v-if="winner === 0">{{ curName }}落子（{{ mode === 'pvp' ? '轮流在同屏落子' : '你执黑先行' }}）</template>
      <template v-else-if="winner === 3">双方和棋</template>
      <template v-else>{{ winner === 1 ? '黑方' : '白方' }}胜出 🎉</template>
    </p>

    <div class="gomoku-board">
      <button
        v-for="i in N * N"
        :key="i"
        class="g-cell"
        :class="{ win: winLine.includes(i - 1), last: lastMove === i - 1 }"
        @click="play(i - 1)"
      >
        <span v-if="board[i - 1]" class="stone" :class="board[i - 1] === 1 ? 'black' : 'white'"></span>
      </button>
    </div>

    <div class="modal-mask" v-if="winner && winner !== 0">
      <div class="modal">
        <template v-if="winner === 3">
          <h3>🤝 和棋</h3>
          <p style="color: var(--muted)">棋逢对手，难分伯仲，再来一局。</p>
        </template>
        <template v-else>
          <h3>🎉 {{ winner === 1 ? '黑方' : '白方' }}胜出</h3>
          <p style="color: var(--muted); margin-bottom: 12px">胜者抽得灵签——</p>
          <SlipCard :slip="slip" />
        </template>
        <button class="btn btn-red" style="margin-top: 12px" @click="init">再来一局</button>
        <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="router.push('/fun')">返回算乐</button>
      </div>
    </div>
  </div>
</template>