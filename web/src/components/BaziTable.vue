<script setup>
import { computed } from 'vue'

const props = defineProps({ bazi: { type: Object, required: true } })

const wuxingColor = { 金: '#c9a24b', 木: '#4c9a5a', 水: '#4a86d9', 火: '#d94a3a', 土: '#a8794f' }
const order = ['金', '木', '水', '火', '土']

const bars = computed(() => {
  const c = props.bazi.wuxing_count || {}
  const max = Math.max(1, ...order.map((w) => c[w] || 0))
  return order.map((w) => ({ w, n: c[w] || 0, pct: Math.round(((c[w] || 0) / max) * 100), color: wuxingColor[w] }))
})
</script>

<template>
  <div class="card">
    <div class="section-head"><span class="ic">📜</span><h3>八字命盘</h3></div>

    <p style="margin: 0 0 4px; color: var(--muted); font-size: 14px">
      公历：{{ bazi.solar_text }}　农历：{{ bazi.lunar_text }}
    </p>
    <p style="margin: 0; color: var(--muted); font-size: 14px">
      生肖：{{ bazi.animal }}　日主：<b style="color: var(--red-deep)">{{ bazi.day_master }}（{{ bazi.day_master_wuxing }}）</b>　{{ bazi.body_strength === '强' ? '身强' : bazi.body_strength === '弱' ? '身弱' : '中和平和' }}
    </p>

    <div class="bazi-wrap" style="margin-top: 12px">
      <table class="bazi-table">
        <thead>
          <tr>
            <th></th>
            <th v-for="p in bazi.pillars" :key="p.name">{{ p.name }}</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="small">十神</td>
            <td v-for="p in bazi.pillars" :key="p.name" class="small">{{ p.shi_shen }}</td>
          </tr>
          <tr>
            <td class="small">天干</td>
            <td v-for="p in bazi.pillars" :key="p.name" class="gz">{{ p.gan }}</td>
          </tr>
          <tr>
            <td class="small">地支</td>
            <td v-for="p in bazi.pillars" :key="p.name" class="gz">{{ p.zhi }}</td>
          </tr>
          <tr>
            <td class="small">藏干</td>
            <td v-for="p in bazi.pillars" :key="p.name" class="small">{{ (p.hidden_gan || []).join(' ') }}</td>
          </tr>
          <tr>
            <td class="small">纳音</td>
            <td v-for="p in bazi.pillars" :key="p.name" class="small">{{ p.na_yin }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div style="margin-top: 14px">
      <div style="font-size: 13px; color: var(--muted); margin-bottom: 8px">
        五行分布　喜用：<span class="tag" v-for="w in bazi.xi_yong" :key="w" :style="{ background: wuxingColor[w] + '22', color: wuxingColor[w] }">{{ w }}</span>
      </div>
      <div class="wuxing-row">
        <div class="wuxing-item" v-for="b in bars" :key="b.w">
          <div class="wuxing-bar"><div class="wuxing-fill" :style="{ width: b.pct + '%', background: b.color }"></div></div>
          <div style="font-size: 12px; margin-top: 4px">{{ b.w }} {{ b.n }}</div>
        </div>
      </div>
    </div>
  </div>
</template>