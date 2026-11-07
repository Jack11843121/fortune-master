<script setup>
import { ref, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { newGame, move, options, roll, TRACK, N_PLANES, SPECIAL } from '../game/flying'
import { randomSlip } from '../data/slips'
import SlipCard from '../components/SlipCard.vue'

const router = useRouter()
const DICE = ['⚀', '⚁', '⚂', '⚃', '⚄', '⚅']
const SPMARK = { ji: '吉', qian: '签', xiong: '凶' }
const COLORS = { 1: '#f0a23a', 2: '#4a86d9' }

const game = ref(newGame())
const mode = ref('pvp')
const rolled = ref(false)
const selectable = ref({ onTrack: [], canTakeoff: false })
const aiBusy = ref(false)
const slip = ref(null)
let slipTimer = null

const cells = Array.from({ length: TRACK }, (_, i) => i)
const dice = computed(() => (game.value.lastRoll ? DICE[game.value.lastRoll - 1] : '🎲'))
const curName = computed(() => (game.value.turn === 1 ? '红方' : '蓝方'))
const winnerName = computed(() => (game.value.winner === 1 ? '红方' : '蓝方'))

// 动态提示：直白告诉玩家这一步该做什么
const hint = computed(() => {
  if (game.value.winner) return '🎉 ' + winnerName.value + '获胜！'
  if (mode.value === 'ai' && game.value.turn === 2) return '对方思考中，请稍候…'
  if (!rolled.value) return '点「掷骰子」摇骰子（掷到 5 或 6 才能起飞）'
  if (selectable.value.canTakeoff) return '已掷到 ' + dice.value + '，点「' + curName.value + '机场」里的 ✈️ 起飞'
  if (selectable.value.onTrack.length) return '点跑道上高亮的 ✈ 前进'
  return '这步没有棋子可动，自动轮到对方'
})

function restart() {
  game.value = newGame()
  rolled.value = false
  selectable.value = { onTrack: [], canTakeoff: false }
  slip.value = null
}

function planesAt(pos) {
  return game.value.planes.filter((x) => x.pos === pos)
}
function parked(p) {
  return game.value.planes.filter((x) => x.p === p && x.pos === -1)
}
function isSelectable(pl) {
  if (!pl) return false // 机场飞机起飞后，格子可能为空，防御访问 undefined
  if (!rolled.value || pl.p !== game.value.turn) return false
  if (pl.pos === -1) return selectable.value.canTakeoff
  return selectable.value.onTrack.includes(pl)
}

function humanRoll() {
  if (game.value.winner || aiBusy.value) return
  if (mode.value === 'ai' && game.value.turn === 2) return
  const r = roll()
  game.value.lastRoll = r
  const o = options(game.value)
  selectable.value = o
  rolled.value = true
  if (o.onTrack.length === 0 && !o.canTakeoff) {
    setTimeout(() => {
      rolled.value = false
      passTurn()
    }, 900)
  }
}

function passTurn() {
  game.value.turn = game.value.turn === 1 ? 2 : 1
  if (mode.value === 'ai' && game.value.turn === 2 && !game.value.winner) doAIMove()
}

function choose(pl) {
  if (!isSelectable(pl) || !game.value.lastRoll) return
  move(game.value, pl, game.value.lastRoll)
  rolled.value = false
  selectable.value = { onTrack: [], canTakeoff: false }
  afterMove()
}

function afterMove() {
  if (game.value.winner) return
  if (game.value.qian) showSlip()
  if (mode.value === 'ai' && game.value.turn === 2 && !game.value.rollAgain && !game.value.winner) doAIMove()
}

function showSlip() {
  slip.value = randomSlip()
  clearTimeout(slipTimer)
  slipTimer = setTimeout(() => (slip.value = null), 2600)
}

function doAIMove() {
  if (game.value.winner || mode.value !== 'ai' || game.value.turn !== 2) return
  aiBusy.value = true
  setTimeout(() => {
    const r = roll()
    game.value.lastRoll = r
    const o = options(game.value)
    let target = null
    if (o.canTakeoff) {
      target = game.value.planes.find((x) => x.p === 2 && x.pos === -1)
    } else if (o.onTrack.length) {
      target = o.onTrack[Math.floor(Math.random() * o.onTrack.length)]
    }
    if (target) {
      move(game.value, target, r)
      afterMove()
    } else {
      passTurn()
    }
    aiBusy.value = false
    if (game.value.turn === 2 && game.value.rollAgain && !game.value.winner) doAIMove()
  }, 650)
}

function switchMode(m) {
  mode.value = m
  restart()
  // 人机：若红（人）先手在 1，无需操作；开局即红先
}

const winSlip = computed(() => (game.value.winner ? randomSlip() : null))

onBeforeUnmount(() => clearTimeout(slipTimer))
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">飞行棋</span>
      <button class="back" @click="restart">↺</button>
    </div>

    <div class="m3-status" style="margin-bottom: 10px">
      <div class="seg" style="flex: 2">
        <button :class="{ on: mode === 'pvp' }" type="button" @click="switchMode('pvp')">双人同屏</button>
        <button :class="{ on: mode === 'ai' }" type="button" @click="switchMode('ai')">人机对战</button>
      </div>
    </div>

    <!-- 比分 + 骰子 -->
    <div class="card" style="padding: 12px 16px">
      <div style="display: flex; align-items: center; justify-content: space-between">
        <div>
          <div>🔴 红方：完成 <b>{{ game.finished[1] }}</b>/{{ N_PLANES }}</div>
          <div>🔵 蓝方：完成 <b>{{ game.finished[2] }}</b>/{{ N_PLANES }}</div>
        </div>
        <div style="text-align: center">
          <div class="fq-dice">{{ dice }}</div>
          <div style="font-size: 12px; color: var(--muted)">{{ game.winner ? winnerName + '胜！' : '轮到 ' + curName }}</div>
        </div>
      </div>
    </div>

    <!-- 机场 -->
    <div class="card" style="padding: 12px 16px">
      <div class="fly-side-row">
        <div class="fly-airport">
          <span style="font-weight: 600; color: #d94a3a">红机场</span>
          <span class="slot" v-for="i in N_PLANES" :key="'r' + i" :class="{ sel: isSelectable(parked(1)[i - 1]) }" @click="choose(parked(1)[i - 1])">
            <template v-if="parked(1)[i - 1]">✈️</template>
          </span>
        </div>
        <div class="fly-airport">
          <span style="font-weight: 600; color: #4a86d9">蓝机场</span>
          <span class="slot" v-for="i in N_PLANES" :key="'b' + i" :class="{ sel: isSelectable(parked(2)[i - 1]) }" @click="choose(parked(2)[i - 1])">
            <template v-if="parked(2)[i - 1]">✈️</template>
          </span>
        </div>
      </div>
    </div>

    <!-- 跑道 -->
    <div class="card" style="padding: 12px">
      <div class="fly-board">
        <div v-for="i in cells" :key="i" class="fly-cell" :class="{ 'base-red': i === 0, 'base-blue': i === 12 }">
          <span style="font-size: 9px">{{ i === 0 ? '红' : i === 12 ? '蓝' : i }}</span>
          <span v-if="SPECIAL[i]" class="fly-sp" :class="SPECIAL[i]">{{ SPMARK[SPECIAL[i]] }}</span>
          <div class="fly-planes">
            <span v-for="pl in planesAt(i)" :key="pl.p + '-' + pl.id" class="fly-plane" :class="pl.p === 1 ? 'red' : 'blue'" :style="isSelectable(pl) ? 'outline:2px solid var(--gold-deep)' : ''" @click="choose(pl)">✈</span>
          </div>
        </div>
      </div>
      <p style="font-size: 11px; color: var(--muted); text-align: center; margin: 8px 0 0">吉=再掷 · 签=抽灵签 · 凶=后退一步</p>
    </div>

    <!-- 操作 -->
    <button class="btn btn-red" :disabled="rolled || aiBusy || game.winner" @click="humanRoll">
      掷骰子　{{ game.lastRoll ? dice : '' }}
    </button>
    <p style="font-size: 13px; color: var(--red-deep); text-align: center; margin: 8px 0; font-weight: 600">
      {{ hint }}
    </p>

    <!-- 玩法说明 -->
    <div class="card" style="margin-top: 4px">
      <div class="section-head"><span class="ic">📖</span><h3>怎么玩</h3></div>
      <ol style="margin: 0; padding-left: 20px; color: var(--muted); font-size: 13px; line-height: 1.9">
        <li>点「掷骰子」摇骰子</li>
        <li>掷到 <b>5 或 6</b>，点自己机场（红/蓝）里的 ✈️ 起飞</li>
        <li>飞机落到跑道起点，之后掷几点就走几格</li>
        <li>落在对方棋子那格，能把它打回机场</li>
        <li>格子上「吉」=再掷一次，「签」=抽支灵签，「凶」=后退一步</li>
        <li>4 架飞机全部先绕完一圈到家 → 获胜</li>
      </ol>
    </div>

    <!-- 抽签提示 -->
    <div class="modal-mask" v-if="slip">
      <div class="modal">
        <h3 style="margin-bottom: 12px">🎁 落在签格，抽得灵签</h3>
        <SlipCard :slip="slip" />
        <button class="btn btn-red" style="margin-top: 12px" @click="slip = null">继续</button>
      </div>
    </div>

    <!-- 胜负 -->
    <div class="modal-mask" v-if="game.winner">
      <div class="modal">
        <h3 style="margin-bottom: 12px">🎉 {{ winnerName }}胜出，抽得灵签</h3>
        <SlipCard :slip="winSlip" />
        <button class="btn btn-red" style="margin-top: 12px" @click="restart">再来一局</button>
        <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="router.push('/fun')">返回算乐</button>
      </div>
    </div>
  </div>
</template>