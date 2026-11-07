<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { randomSlip, huangliToday } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'
import { NAME, createBoard, legalMovesFor, inCheck, bestMove } from '../game/xiangqi'

const router = useRouter()

// —— 棋盘几何（自然尺寸，随后按屏幕缩放）——
const CELL = 44
const PAD = 30
const RIVER = 66
const PC = 38
const W = PAD * 2 + CELL * 8
const H = PAD * 2 + CELL * 8 + RIVER
const xp = (c) => PAD + c * CELL
const yp = (r) => PAD + r * CELL + (r > 4 ? RIVER - CELL : 0)

// —— 自适应缩放 ——
const scale = ref(1)
function fit() {
  const avail = (window.innerWidth || 375) - 32
  scale.value = Math.min(1, avail / W)
}
onMounted(() => {
  fit()
  window.addEventListener('resize', fit)
})
onBeforeUnmount(() => window.removeEventListener('resize', fit))

// —— 状态 ——
const board = ref(createBoard())
const turn = ref('r')
const selected = ref(null)
const moves = ref([])
const winner = ref(null)
const checking = ref(false)
const lastMove = ref(null)
const mode = ref('pvp') // pvp 双人 | ai 人机
const aiThinking = ref(false)

const slip = computed(() => (winner.value ? randomSlip() : null))
const hl = huangliToday()

function reset() {
  board.value = createBoard()
  turn.value = 'r'
  selected.value = null
  moves.value = []
  winner.value = null
  checking.value = false
  lastMove.value = null
}

function chosen(r, c) {
  const p = board.value[r][c]
  selected.value = null
  moves.value = []
  if (p && p.side === turn.value) {
    selected.value = { r, c }
    moves.value = legalMovesFor(board.value, turn.value)
      .filter((m) => m.from.r === r && m.from.c === c)
      .map((m) => ({ r: m.to.r, c: m.to.c, capture: !!board.value[m.to.r][m.to.c] }))
  }
}

function click(r, c) {
  if (winner.value) return
  if (mode.value === 'ai' && (turn.value !== 'r' || aiThinking.value)) return
  const p = board.value[r][c]
  if (selected.value && p && p.side === turn.value) {
    chosen(r, c)
    return
  }
  if (selected.value) {
    const mv = moves.value.find((m) => m.r === r && m.c === c)
    if (mv) {
      doMove(selected.value, { r, c })
      return
    }
  }
  chosen(r, c)
}

function doMove(from, to) {
  board.value[to.r][to.c] = board.value[from.r][from.c]
  board.value[from.r][from.c] = null
  selected.value = null
  moves.value = []
  lastMove.value = { from, to }
  const other = turn.value === 'r' ? 'b' : 'r'
  checking.value = inCheck(board.value, other)
  if (legalMovesFor(board.value, other).length === 0) {
    winner.value = turn.value
  } else {
    turn.value = other
    if (mode.value === 'ai' && other === 'b' && !winner.value) scheduleAI()
  }
}

function scheduleAI() {
  aiThinking.value = true
  setTimeout(() => {
    const m = bestMove(board.value, 'b', 2)
    aiThinking.value = false
    if (m) doMove(m.from, m.to)
  }, 300)
}

function switchMode(m) {
  mode.value = m
  reset()
}

// —— SVG 棋盘数据 ——
const vLinesFull = computed(() => Array.from({ length: 7 }, (_, i) => i + 1))
const vLinesHalf = [0, 8]
const palaces = [
  { x1: xp(3), y1: yp(0), x2: xp(5), y2: yp(2) }, { x1: xp(5), y1: yp(0), x2: xp(3), y2: yp(2) },
  { x1: xp(3), y1: yp(7), x2: xp(5), y2: yp(9) }, { x1: xp(5), y1: yp(7), x2: xp(3), y2: yp(9) },
]
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">象棋对弈</span>
      <button class="back" @click="reset">↺</button>
    </div>

    <div class="m3-status">
      <div class="seg" style="flex: 2">
        <button :class="{ on: mode === 'pvp' }" type="button" @click="switchMode('pvp')">双人对弈</button>
        <button :class="{ on: mode === 'ai' }" type="button" @click="switchMode('ai')">人机对战</button>
      </div>
    </div>

    <p class="m3-hint xq-status">
      宜<strong style="color: var(--green)">{{ hl.yi }}</strong> 忌<strong style="color: var(--red)">{{ hl.ji }}</strong>
      <span class="xq-cur" :class="{ check: checking }">
        {{ winner ? (winner === 'r' ? '红方' : '黑方') + '胜 🎉' : aiThinking ? '思考中…' : (turn === 'r' ? '红方' : '黑方') + '走棋' + (checking ? ' · 将军！' : '') }}
      </span>
    </p>

    <div class="xq-outer" :style="{ width: W * scale + 'px', height: H * scale + 'px' }">
      <div class="xq-board" :style="{ width: W + 'px', height: H + 'px', transform: 'scale(' + scale + ')' }">
        <svg :width="W" :height="H" :viewBox="'0 0 ' + W + ' ' + H">
          <rect :x="xp(0)" :y="yp(0)" :width="CELL * 8" :height="yp(9) - yp(0)" fill="#f0d9a0" stroke="#5b4326" stroke-width="2" />
          <line v-for="r in 10" :key="'h' + r" :x1="xp(0)" :y1="yp(r - 1)" :x2="xp(8)" :y2="yp(r - 1)" stroke="#5b4326" />
          <line v-for="c in vLinesFull" :key="'v' + c" :x1="xp(c)" :y1="yp(0)" :x2="xp(c)" :y2="yp(9)" stroke="#5b4326" />
          <line v-for="c in vLinesHalf" :key="'vhl' + c" :x1="xp(c)" :y1="yp(0)" :x2="xp(c)" :y2="yp(4)" stroke="#5b4326" />
          <line v-for="c in vLinesHalf" :key="'vhr' + c" :x1="xp(c)" :y1="yp(5)" :x2="xp(c)" :y2="yp(9)" stroke="#5b4326" />
          <line v-for="(pl, i) in palaces" :key="'p' + i" :x1="pl.x1" :y1="pl.y1" :x2="pl.x2" :y2="pl.y2" stroke="#5b4326" />
          <text :x="W * 0.28" :y="(yp(4) + yp(5)) / 2 + 5" font-size="16" fill="#5b4326" text-anchor="middle" style="font-family:serif">楚 河</text>
          <text :x="W * 0.72" :y="(yp(4) + yp(5)) / 2 + 5" font-size="16" fill="#5b4326" text-anchor="middle" style="font-family:serif">汉 界</text>
        </svg>

        <template v-for="(m, i) in moves" :key="'mv' + i">
          <div class="xq-dot" :class="{ cap: m.capture }" :style="{ left: xp(m.c) + 'px', top: yp(m.r) + 'px' }"></div>
        </template>

        <template v-for="(row, r) in board" :key="r">
          <template v-for="(p, c) in row" :key="r + '-' + c">
            <div
              v-if="p"
              class="xq-piece"
              :class="['side-' + p.side, selected && selected.r === r && selected.c === c ? 'sel' : '', lastMove && lastMove.to.r === r && lastMove.to.c === c ? 'last' : '']"
              :style="{ left: xp(c) - PC / 2 + 'px', top: yp(r) - PC / 2 + 'px', width: PC + 'px', height: PC + 'px' }"
              @click="click(r, c)"
            >
              {{ NAME[p.type][p.side] }}
            </div>
          </template>
        </template>

        <template v-for="(row, r) in board" :key="'hit' + r">
          <template v-for="(p, c) in row" :key="'hit' + r + '-' + c">
            <div v-if="!p" class="xq-hit" :style="{ left: xp(c) - 20 + 'px', top: yp(r) - 20 + 'px' }" @click="click(r, c)"></div>
          </template>
        </template>
      </div>
    </div>

    <p class="m3-hint" style="text-align:center">点击己方棋子，再点绿色落点走棋 · {{ mode === 'ai' ? '你执红先行，与电脑对弈' : '双人同屏轮流对弈' }}</p>

    <div class="modal-mask" v-if="winner">
      <div class="modal">
        <h3 style="margin-bottom: 12px">🎉 {{ winner === 'r' ? '红方' : '黑方' }}胜出，抽得灵签</h3>
        <SlipCard :slip="slip" />
        <button class="btn btn-red" style="margin-top: 16px" @click="reset">再来一局</button>
        <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="router.push('/fun')">返回算乐</button>
      </div>
    </div>
  </div>
</template>