<script setup>
defineProps({ shapes: { type: Array, required: true } })

// 用 stroke-dashoffset 做「笔迹生长」：每个形状像用笔一笔一划勾出来
function style(i) {
  return {
    strokeDasharray: 600,
    strokeDashoffset: 600,
    animation: `sketchDraw 0.6s ease ${i * 0.38}s forwards`,
  }
}

const pts = (p) => p.map(([x, y]) => `${x},${y}`).join(' ')
</script>

<template>
  <svg viewBox="0 0 100 100" class="sketch-svg">
    <template v-for="(s, i) in shapes" :key="i">
      <circle v-if="s.t === 'c'" :cx="s.cx" :cy="s.cy" :r="s.r" :style="style(i)" />
      <line v-else-if="s.t === 'l'" :x1="s.x1" :y1="s.y1" :x2="s.x2" :y2="s.y2" :style="style(i)" />
      <rect v-else-if="s.t === 'r'" :x="s.x" :y="s.y" :width="s.w" :height="s.h" :style="style(i)" />
      <ellipse v-else-if="s.t === 'e'" :cx="s.cx" :cy="s.cy" :rx="s.rx" :ry="s.ry" :style="style(i)" />
      <polygon v-else-if="s.t === 'p'" :points="pts(s.pts)" :style="style(i)" />
      <path v-else-if="s.t === 'path'" :d="s.d" :style="style(i)" />
    </template>
  </svg>
</template>