<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()

const scenes = [
  { to: '/personal', e: '☯️', t: '个人命理', d: '看八字 · 五行 · 性格事业 · 感情健康 · 大运流年' },
  { to: '/wealth', e: '💰', t: '财运解析', d: '财运走向 · 投资性格 · 股市板块 · 求财开运' },
  { to: '/marriage', e: '💞', t: '合婚匹配', d: '两人八字合不合 · 缘分深浅 · 相处建议' },
  { to: '/hiring', e: '💼', t: '招聘看八字', d: '员工性格特质 · 岗位契合度 · 团队匹配' },
  { to: '/naming', e: '📛', t: '宝宝取名', d: '按八字喜用神起名 · 寓意吉祥 · 五行补缺' },
]
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <div class="card" style="text-align: center; background: linear-gradient(150deg, #fdf6e9, #f7e8cd)">
      <div style="font-size: 40px">🔮</div>
      <h2 style="margin: 6px 0 4px; color: var(--red-deep)">大师 · 八字命理</h2>
      <p style="color: var(--muted); font-size: 13px; margin: 0">传统排盘 · 通俗解读 · 先看预览，满意再解锁</p>
    </div>

    <div v-for="s in scenes" :key="s.to" class="scene-card" @click="router.push(s.to)">
      <div class="sc-icon"><span class="emoji">{{ s.e }}</span></div>
      <div>
        <h3>{{ s.t }}</h3>
        <p>{{ s.d }}</p>
      </div>
      <div class="arrow">›</div>
    </div>

    <div class="notice">🔮 传承千年命理智慧，助您知命明理、趋吉避凶、顺遂如意。</div>
  </div>
</template>