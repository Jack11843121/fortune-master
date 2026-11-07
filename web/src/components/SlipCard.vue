<script setup>
import { LEVEL_COLOR } from '../data/slips'
defineProps({ slip: { type: Object, required: true } })
</script>

<template>
  <div class="slip">
    <div class="slip-ribbon" :style="{ background: LEVEL_COLOR[slip.level] }">{{ slip.level }}</div>

    <div class="slip-no">第 {{ slip.no }} 签 · {{ slip.title }}</div>

    <div class="slip-poem">
      <div v-for="(l, i) in slip.poem" :key="i" class="slip-line" :style="i % 2 === 1 ? 'margin-left:18px' : ''">
        {{ l }}
      </div>
    </div>

    <div class="slip-jie">
      <span class="slip-jie-tag">解曰</span>
      <span>{{ slip.jie }}</span>
    </div>

    <div class="slip-foot">
      <span class="slip-seal">算命大师</span>
      <span class="slip-note">愿君吉祥 · 心想事成</span>
    </div>
  </div>
</template>