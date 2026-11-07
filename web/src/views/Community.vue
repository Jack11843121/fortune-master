<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const posts = ref([])
const total = ref(0)
const page = ref(1)
const cat = ref('')
const showCompose = ref(false)
const form = ref({ title: '', content: '', category: '命理趣谈' })
const err = ref('')
const cats = ['命理趣谈', '合婚故事', '求职问卦', '宝宝取名', '开运经验']

async function load() {
  try {
    const r = await api.listPosts({ page: page.value, category: cat.value })
    posts.value = r.posts
    total.value = r.total
  } catch (e) {
    err.value = e.message
  }
}
onMounted(load)

function pickCat(c) {
  cat.value = c
  page.value = 1
  load()
}

async function submit() {
  if (!api.isLoggedIn()) {
    router.push('/me')
    return
  }
  if (!form.value.title.trim() || !form.value.content.trim()) {
    err.value = '请填写标题与内容'
    return
  }
  try {
    await api.createPost(form.value)
    showCompose.value = false
    form.value = { title: '', content: '', category: '命理趣谈' }
    page.value = 1
    load()
  } catch (e) {
    err.value = e.message
  }
}

function more() {
  page.value++
  load()
}
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <button class="back" style="margin-bottom: 10px" @click="router.push('/me')">‹ 返回</button>

    <div class="card" style="display: flex; align-items: center; justify-content: space-between; padding: 14px 18px">
      <div>
        <h2 style="margin: 0">💬 交流社区</h2>
        <p style="color: var(--muted); font-size: 13px; margin: 4px 0 0">分享命理趣事 · 合婚故事 · 问卦交流</p>
      </div>
      <button class="btn btn-red btn-sm" style="width: auto" @click="showCompose = !showCompose">+ 发帖</button>
    </div>

    <div v-if="showCompose" class="card">
      <div class="field">
        <label>标题</label>
        <input class="input" v-model="form.title" placeholder="一句话标题" maxlength="120" />
      </div>
      <div class="field">
        <label>分类</label>
        <div class="seg">
          <button v-for="c in cats" :key="c" :class="{ on: form.category === c }" type="button" @click="form.category = c">{{ c }}</button>
        </div>
      </div>
      <div class="field">
        <label>内容</label>
        <textarea class="input" v-model="form.content" rows="4" placeholder="分享你的故事或问题…" maxlength="2000"></textarea>
      </div>
      <button class="btn btn-primary" @click="submit">发布</button>
    </div>

    <div style="display: flex; flex-wrap: wrap; gap: 6px; margin: 10px 0">
      <button class="btn btn-ghost btn-sm" style="width:auto" :class="{ on: cat === '' }" @click="pickCat('')">全部</button>
      <button v-for="c in cats" :key="c" class="btn btn-ghost btn-sm" style="width:auto" :class="{ on: cat === c }" @click="pickCat(c)">{{ c }}</button>
    </div>

    <p v-if="err" style="color: var(--red); font-size: 13px">{{ err }}</p>

    <div v-for="p in posts" :key="p.id" class="post-card" @click="router.push(`/community/${p.id}`)">
      <div style="display: flex; align-items: center; gap: 8px">
        <img v-if="p.avatar" :src="p.avatar" style="width: 28px; height: 28px; border-radius: 6px" />
        <span v-else style="width: 28px; height: 28px; border-radius: 6px; background: #f1e3c8; display: inline-flex; align-items: center; justify-content: center; font-size: 14px">🙂</span>
        <div style="display: flex; flex-direction: column">
          <span class="post-cat" style="align-self: flex-start">{{ p.category }}</span>
          <span style="font-size: 12px; color: var(--muted)">{{ p.nickname || '匿名' }}</span>
        </div>
      </div>
      <div class="post-title">{{ p.title }}</div>
      <div class="post-content">{{ p.content }}</div>
      <div class="post-meta">
        <span>{{ p.created_at?.slice(0, 10) }}</span>
        <span>👍 {{ p.likes }}</span>
      </div>
    </div>

    <div v-if="posts.length === 0" class="notice" style="text-align: center">还没有帖子，来抢个沙发吧 🎉</div>

    <button v-if="total > page * 20" class="btn btn-ghost" @click="more">加载更多</button>
  </div>
</template>