<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ELEMENTS, play, stop } from '../game/music'

const router = useRouter()
const current = ref(null)
const playing = ref(false)

function toggle(el) {
  if (current.value === el.key && playing.value) {
    stop()
    current.value = null
    playing.value = false
  } else {
    play(el)
    current.value = el.key
    playing.value = true
  }
}

function stopAll() {
  stop()
  current.value = null
  playing.value = false
}

onBeforeUnmount(() => stop())
</script>

<template>
  <div class="page" style="padding-top: 0">
    <div class="game-top">
      <button class="back" @click="router.push('/fun')">‹</button>
      <span class="gt-title">五行音乐</span>
      <button class="back" @click="stopAll">■</button>
    </div>

    <div class="card" style="text-align: center; background: linear-gradient(150deg, #fdf6e9, #f7e8cd)">
      <div style="font-size: 38px">{{ playing ? '🎼' : '🎵' }}</div>
      <h3 style="margin: 4px 0; color: var(--red-deep)">五音疗疾 · 静心安神</h3>
      <p style="color: var(--muted); font-size: 13px; margin: 0">
        宫商角徵羽，对应土金木火水、入脾肺肝心肾。静心安坐，听一阕连贯的五音曲，涤荡心绪、滋养五脏。
        可先到「个人命理」看喜用神，选对应五行，更有裨益。
      </p>
    </div>

    <div
      v-for="el in ELEMENTS"
      :key="el.key"
      class="mu-card"
      :class="{ playing: current === el.key && playing }"
      @click="toggle(el)"
    >
      <div class="mu-dot" :style="{ background: el.color }">{{ el.wuyin }}</div>
      <div style="flex: 1">
        <div class="mu-title">
          <b>{{ el.wuxing }}</b> · {{ el.mood }}
          <span style="color: var(--muted); font-weight: 400">｜入{{ el.organ }}</span>
        </div>
        <div class="mu-desc">{{ el.desc }}</div>
      </div>
      <div class="mu-state">{{ current === el.key && playing ? '♪ 播放中' : '▶' }}</div>
    </div>

    <p class="m3-hint" style="text-align: center">点击五行卡片即奏一阕连贯的五音曲，再点一次停止。</p>
  </div>
</template>