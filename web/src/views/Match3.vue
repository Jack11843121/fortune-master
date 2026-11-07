<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { randomSlip } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'
import { makeBoard, findMatches, hasAnyMove, findHint, resolveAsync, key } from '../game/match3'

const router = useRouter()

const ELEMENTS = [
  { name: '金', sym: '🪙', color: '#c9a24b' },
  { name: '木', sym: '🌳', color: '#4c9a5a' },
  { name: '水', sym: '💧', color: '#4a86d9' },
  { name: '火', sym: '🔥', color: '#d94a3a' },
  { name: '土', sym: '⛰️', color: '#a8794f' },
]
const TARGET = 300
const MOVES = 22

const grid = ref(makeBoard())
const selected = ref(null)
const score = ref(0)
const moves = ref(MOVES)
const matched = ref(new Set())
const hintCells = ref(new Set())
const invalidCells = ref(new Set())
const busy = ref(false)
const comboMsg = ref('')
const shuffleToast = ref(false)
const result = ref(null)

function init() {
  grid.value = makeBoard()
  selected.value = null
  score.value = 0
  moves.value = MOVES
  matched.value = new Set()
  hintCells.value = new Set()
  invalidCells.value = new Set()
  busy.value = false
  comboMsg.value = ''
  shuffleToast.value = false
  result.value = null
}

function isSel(r, c) {
  return selected.value && selected.value.r === r && selected.value.c === c
}
function boom(r, c) {
  return matched.value.has(key(r, c))
}
function hinted(r, c) {
  return hintCells.value.has(key(r, c))
}
function invalid(r, c) {
  return invalidCells.value.has(key(r, c))
}

function doHint() {
  if (busy.value || result.value) return
  const h = findHint(grid.value)
  if (h) {
    hintCells.value = new Set([key(h.a.r, h.a.c), key(h.b.r, h.b.c)])
    setTimeout(() => (hintCells.value = new Set()), 2600)
  }
}

function flashInvalid(a, b) {
  invalidCells.value = new Set([key(a.r, a.c), key(b.r, b.c)])
  setTimeout(() => (invalidCells.value = new Set()), 420)
}

async function tap(r, c) {
  if (busy.value || result.value || moves.value <= 0) return
  const s = selected.value
  if (!s) { selected.value = { r, c }; return }
  if (s.r === r && s.c === c) { selected.value = null; return }
  if (Math.abs(s.r - r) + Math.abs(s.c - c) !== 1) { selected.value = { r, c }; return }

  const g = grid.value
  ;[g[s.r][s.c], g[r][c]] = [g[r][c], g[s.r][s.c]]
  if (findMatches(g).size === 0) {
    ;[g[s.r][s.c], g[r][c]] = [g[r][c], g[s.r][s.c]]
    selected.value = null
    flashInvalid(s, { r, c })
    return
  }

  selected.value = null
  hintCells.value = new Set()
  busy.value = true
  moves.value--
  let gained = 0
  let chains = 1
  try {
    const res = await resolveAsync(g, (cleared) => {
      matched.value = new Set(cleared.map(([rr, cc]) => key(rr, cc)))
    })
    gained = res.gained
    chains = res.chains
    matched.value = new Set()
  } finally {
    busy.value = false
  }
  score.value += gained
  comboMsg.value = chains >= 2 ? `连锁 ×${chains}！+${gained}` : `消除 +${gained}`

  if (!hasAnyMove(g)) {
    grid.value = makeBoard()
    shuffleToast.value = true
    setTimeout(() => (shuffleToast.value = false), 1600)
  }

  if (score.value >= TARGET) result.value = { win: true, slip: randomSlip() }
  else if (moves.value <= 0) result.value = { win: false }
}
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">五行消消乐</span>
      <button class="back" @click="init">↺</button>
    </div>

    <div class="m3-status">
      <div class="m3-cell">🎯 目标 {{ TARGET }}</div>
      <div class="m3-cell">⭐ {{ score }}</div>
      <div class="m3-cell">🫳 {{ moves }} 步</div>
    </div>
    <div class="m3-combo" :class="{ show: comboMsg }">{{ comboMsg }}</div>
    <div class="m3-combo" :class="{ show: shuffleToast }" style="color: var(--muted)">暂无可消，已自动洗牌</div>

    <div class="m3-board">
      <template v-for="(row, r) in grid" :key="r">
        <button
          v-for="(v, c) in row"
          :key="r + '-' + c"
          type="button"
          class="m3-tile"
          :class="{ sel: isSel(r, c), boom: boom(r, c), hint: hinted(r, c), invalid: invalid(r, c) }"
          :style="{ background: 'linear-gradient(160deg,' + ELEMENTS[v].color + '33,' + ELEMENTS[v].color + '18)', color: ELEMENTS[v].color, borderColor: ELEMENTS[v].color + '55' }"
          @click="tap(r, c)"
        >
          <span class="m3-sym">{{ ELEMENTS[v].sym }}</span>
          <span class="m3-name">{{ ELEMENTS[v].name }}</span>
        </button>
      </template>
    </div>

    <div style="display: flex; gap: 10px; align-items: center; justify-content: center; margin: 4px 0">
      <p class="m3-hint" style="margin: 0">点击相邻两块交换，凑成三个及以上同五行即可消除</p>
      <button class="btn btn-ghost btn-sm" style="width: auto; padding: 8px 14px" @click="doHint">💡 提示</button>
    </div>

    <div class="modal-mask" v-if="result">
      <div class="modal">
        <template v-if="result.win">
          <h3 style="margin-bottom: 14px">🎉 达成目标，抽得灵签</h3>
          <SlipCard :slip="result.slip" />
        </template>
        <template v-else>
          <h3>步数用完啦</h3>
          <p style="color: var(--muted)">本次得分 {{ score }}，差一点到 {{ TARGET }}，再来一局！</p>
        </template>
        <button class="btn btn-red" style="margin-top: 16px" @click="init">再来一局</button>
        <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="router.push('/fun')">返回算乐</button>
      </div>
    </div>
  </div>
</template>