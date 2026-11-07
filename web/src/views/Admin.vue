<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'

const token = ref(localStorage.getItem('admin_token') || '')
const username = ref('')
const password = ref('')
const loginErr = ref('')
const err = ref('')

const tab = ref('dash')
const dash = ref(null)
const users = ref([])
const orders = ref([])
const msgs = ref([])
const posts = ref([])
const logs = ref([])

async function login() {
  loginErr.value = ''
  try {
    const r = await api.adminLogin(username.value, password.value)
    token.value = r.token
    localStorage.setItem('admin_token', r.token)
    username.value = ''
    password.value = ''
    await load()
  } catch (e) {
    loginErr.value = e.message
  }
}

function logout() {
  localStorage.removeItem('admin_token')
  token.value = ''
}

async function load() {
  try {
    dash.value = await api.adminDashboard()
    users.value = await api.adminUsers()
    orders.value = await api.adminOrders()
    msgs.value = await api.adminConsultations()
    posts.value = await api.adminPosts()
    logs.value = await api.adminMemberships()
  } catch (e) {
    err.value = e.message
    if (err.value.includes('401') || err.value.includes('会话')) logout()
  }
}
onMounted(() => {
  if (token.value) load()
})

async function userAction(u, patch) {
  try {
    await api.adminUpdateUser(u.id, patch)
    await load()
  } catch (e) {
    err.value = e.message
  }
}

async function replyMsg(m) {
  const text = prompt('回复内容：', '')
  if (!text) return
  try {
    await api.adminReplyConsult(m.id, text)
    await load()
  } catch (e) {
    err.value = e.message
  }
}

async function moderate(p, status) {
  try {
    await api.adminModeratePost(p.id, status)
    await load()
  } catch (e) {
    err.value = e.message
  }
}

function yuan(c) {
  return ((c || 0) / 100).toFixed(2)
}
function roleLabel(r) {
  return { user: '用户', assistant: 'AI', admin: '管理员' }[r] || r
}
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <!-- 登录 -->
    <div v-if="!token" class="card" style="max-width: 360px; margin: 40px auto">
      <h3 style="text-align: center">🔐 管理后台</h3>
      <div class="field">
        <label>账号</label>
        <input class="input" v-model="username" placeholder="admin" />
      </div>
      <div class="field">
        <label>密码</label>
        <input class="input" type="password" v-model="password" placeholder="密码" @keyup.enter="login" />
      </div>
      <p v-if="loginErr" style="color: var(--red); font-size: 13px">{{ loginErr }}</p>
      <button class="btn btn-red" @click="login">登录</button>
    </div>

    <template v-else>
      <div class="card" style="display: flex; align-items: center; justify-content: space-between; padding: 12px 16px">
        <div><b>🧭 管理后台</b></div>
        <button class="btn btn-ghost btn-sm" style="width:auto" @click="logout">退出登录</button>
      </div>

      <div class="admin-nav">
        <button :class="{ on: tab === 'dash' }" @click="tab = 'dash'">数据看板</button>
        <button :class="{ on: tab === 'users' }" @click="tab = 'users'">会员管理</button>
        <button :class="{ on: tab === 'orders' }" @click="tab = 'orders'">支付订单</button>
        <button :class="{ on: tab === 'msgs' }" @click="tab = 'msgs'">日常咨询</button>
        <button :class="{ on: tab === 'posts' }" @click="tab = 'posts'">社区内容</button>
        <button :class="{ on: tab === 'logs' }" @click="tab = 'logs'">会员流水</button>
      </div>

      <p v-if="err" style="color: var(--red); font-size: 13px">{{ err }}</p>

      <!-- 看板 -->
      <div v-if="tab === 'dash'" class="card">
        <div class="stat-row">
          <div class="stat"><div class="stat-num">{{ dash?.users }}</div><div class="stat-label">总用户</div></div>
          <div class="stat"><div class="stat-num">{{ dash?.vip_users }}</div><div class="stat-label">会员</div></div>
          <div class="stat"><div class="stat-num">¥{{ yuan(dash?.revenue_cents) }}</div><div class="stat-label">累计收入</div></div>
        </div>
        <div class="stat-row" style="margin-top: 10px">
          <div class="stat"><div class="stat-num">{{ dash?.paid_orders }}</div><div class="stat-label">已支付订单</div></div>
          <div class="stat"><div class="stat-num">{{ dash?.ai_messages }}</div><div class="stat-label">AI 对话</div></div>
          <div class="stat"><div class="stat-num">{{ dash?.posts }}</div><div class="stat-label">社区帖子</div></div>
        </div>
      </div>

      <!-- 会员 -->
      <div v-else-if="tab === 'users'" class="card">
        <table class="tbl">
          <thead><tr><th>用户</th><th>会员</th><th>AI次数</th><th>等级</th><th>成长</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="u in users.users" :key="u.id">
              <td>{{ u.nickname || u.id.slice(0, 8) }}</td>
              <td><span class="badge" :class="u.vip_until && new Date(u.vip_until) > new Date() ? 'badge-green' : 'badge-gray'">{{ u.vip_until && new Date(u.vip_until) > new Date() ? '会员' : '非会员' }}</span></td>
              <td>{{ u.ai_credits }}</td>
              <td>Lv{{ u.level }}</td>
              <td>{{ u.growth }}</td>
              <td>
                <button class="btn btn-ghost btn-sm" style="width:auto" @click="userAction(u, { add_credits: 1000 })">+1000次</button>
                <button class="btn btn-ghost btn-sm" style="width:auto" @click="userAction(u, { add_vip_days: 30 })">+30天</button>
                <button class="btn btn-ghost btn-sm" style="width:auto" @click="userAction(u, { status: u.status === 'banned' ? 'active' : 'banned' })">{{ u.status === 'banned' ? '解封' : '封禁' }}</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 订单 -->
      <div v-else-if="tab === 'orders'" class="card">
        <table class="tbl">
          <thead><tr><th>订单号</th><th>类型</th><th>金额</th><th>状态</th><th>时间</th></tr></thead>
          <tbody>
            <tr v-for="o in orders.orders" :key="o.order_no">
              <td style="font-size: 12px">{{ o.order_no }}</td>
              <td>{{ o.type }}</td>
              <td>¥{{ yuan(o.amount) }}</td>
              <td><span class="badge" :class="o.status === 'paid' ? 'badge-green' : 'badge-gray'">{{ o.status }}</span></td>
              <td style="font-size: 12px">{{ o.created_at?.slice(0, 16) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 咨询 -->
      <div v-else-if="tab === 'msgs'" class="card">
        <table class="tbl">
          <thead><tr><th>角色</th><th>用户</th><th>内容</th><th>时间</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="m in msgs.messages" :key="m.id">
              <td><span class="badge" :class="m.role === 'user' ? 'badge-red' : m.role === 'assistant' ? 'badge-green' : 'badge-gray'">{{ roleLabel(m.role) }}</span></td>
              <td style="font-size: 12px">{{ m.user_id.slice(0, 8) }}</td>
              <td style="max-width: 260px; white-space: pre-wrap">{{ m.content }}</td>
              <td style="font-size: 12px">{{ m.created_at?.slice(0, 16) }}</td>
              <td><button v-if="m.role === 'user'" class="btn btn-ghost btn-sm" style="width:auto" @click="replyMsg(m)">回复</button></td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 社区 -->
      <div v-else-if="tab === 'posts'" class="card">
        <table class="tbl">
          <thead><tr><th>标题</th><th>分类</th><th>状态</th><th>赞</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="p in posts.posts" :key="p.id">
              <td>{{ p.title }}</td>
              <td>{{ p.category }}</td>
              <td><span class="badge" :class="p.status === 'published' ? 'badge-green' : 'badge-gray'">{{ p.status }}</span></td>
              <td>{{ p.likes }}</td>
              <td>
                <button v-if="p.status !== 'published'" class="btn btn-ghost btn-sm" style="width:auto" @click="moderate(p, 'published')">上架</button>
                <button v-if="p.status !== 'hidden'" class="btn btn-ghost btn-sm" style="width:auto" @click="moderate(p, 'hidden')">下架</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 流水 -->
      <div v-else-if="tab === 'logs'" class="card">
        <table class="tbl">
          <thead><tr><th>用户</th><th>类型</th><th>金额</th><th>次数</th><th>说明</th><th>时间</th></tr></thead>
          <tbody>
            <tr v-for="l in logs.logs" :key="l.id">
              <td style="font-size: 12px">{{ l.user_id.slice(0, 8) }}</td>
              <td>{{ l.kind }}</td>
              <td>¥{{ yuan(l.amount) }}</td>
              <td>{{ l.credits }}</td>
              <td>{{ l.note }}</td>
              <td style="font-size: 12px">{{ l.created_at?.slice(0, 16) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>