<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()

const scenes = [
  { to: '/personal', e: '☯️', t: '个人命理', d: '八字五行 · 性格事业 · 感情健康 · 大运流年', tag: '热门' },
  { to: '/wealth', e: '💰', t: '财运解析', d: '财运走向 · 投资性格 · 股市板块 · 求财开运', tag: '赚钱' },
  { to: '/marriage', e: '💞', t: '合婚匹配', d: '两人合不合 · 缘分深浅 · 相处建议', tag: '姻缘' },
  { to: '/hiring', e: '💼', t: '招聘看八字', d: '员工特质 · 岗位契合 · 团队匹配', tag: '老板' },
  { to: '/naming', e: '📛', t: '宝宝取名', d: '按八字喜用神起名 · 寓意吉祥 · 五行补缺', tag: '起名' },
]
</script>

<template>
  <div class="page" style="padding-top: 12px">
    <!-- 顶部大标题 -->
    <div class="card" style="text-align: center; background: linear-gradient(150deg, #fdf6e9, #f4e3c0); padding: 20px 18px">
      <div style="font-size: 42px">🀄</div>
      <h2 style="margin: 8px 0 4px; color: var(--red-deep); letter-spacing: 2px">知命 · 明理 · 从容生活</h2>
      <p style="color: var(--muted); font-size: 13px; margin: 0">传统八字排盘 · 通俗易懂 · 全功能免费体验</p>
    </div>

    <!-- 算命功能入口 -->
    <div v-for="s in scenes" :key="s.to" class="scene-card" @click="router.push(s.to)">
      <div class="sc-icon"><span class="emoji">{{ s.e }}</span></div>
      <div style="flex: 1">
        <h3 style="display: flex; align-items: center; gap: 6px">
          {{ s.t }}
          <span class="tag" style="margin: 0">{{ s.tag }}</span>
        </h3>
        <p>{{ s.d }}</p>
      </div>
      <div class="arrow">›</div>
    </div>

    <!-- 为什么选择我们 -->
    <div class="card">
      <div class="section-head"><span class="ic">✨</span><h3>为什么选择我们</h3></div>
      <ul style="margin: 0; padding-left: 18px; color: var(--muted); font-size: 14px; line-height: 2">
        <li>专业排盘：节气、干支、藏干、纳音、大运流年一应俱全</li>
        <li>通俗解读：不玄乎，用人人能懂的话讲清你的命理</li>
        <li>全功能免费：命盘、解读、大运流年全部免费开放</li>
        <li>场景丰富：个人、财运、合婚、招聘、起名……一个入口全搞定</li>
      </ul>
    </div>

    <div class="notice">🔮 传承千年命理智慧，助您知命明理、趋吉避凶、顺遂如意。</div>

    <!-- 合作广告 -->
    <div class="card" style="background: linear-gradient(160deg, #232a3a, #37445f); color: #fff; border: none">
      <div style="font-size: 21px; font-weight: 700; letter-spacing: 1px">🤝 定制软件 · 欢迎合作</div>
      <p style="color: #e8c98a; font-size: 15px; margin: 6px 0 12px; letter-spacing: 1.5px">技术赋能 · 项目落实 · 企业数字化改革 · 增效提质</p>
      <ul style="margin: 0; padding-left: 18px; color: #dfe6f1; font-size: 14px; line-height: 1.9">
        <li>为老年人定制：大字版、易上手的娱乐 / 健康 / 陪伴应用</li>
        <li>为公司定制：小程序 / H5 / App / 管理后台 / AI 系统</li>
        <li>为小朋友定制：寓教于乐的互动学习应用</li>
        <li>AI Agent 落地：办公自动化、产品智能助手、企业 AI 应用</li>
      </ul>
      <p style="color: #b8c4d8; font-size: 13px; margin: 12px 0 0; line-height: 1.7">
        本站即交付案例——八字排盘、会员体系、AI 对话、社区、管理后台一应俱全，全功能免费体验，可深度合作。
      </p>
      <div style="margin-top: 14px; padding: 12px 14px; background: rgba(255,255,255,0.08); border-radius: 10px; text-align: center">
        <div style="font-size: 13px; color: #b8c4d8">有需求？欢迎邮箱洽谈，高质量交付</div>
        <div style="font-size: 16px; color: #ffd98a; font-weight: 700; margin-top: 4px; letter-spacing: 1px">📮 11843121@qq.com</div>
      </div>
    </div>
  </div>
</template>