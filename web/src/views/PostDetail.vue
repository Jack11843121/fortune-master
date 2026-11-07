<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'

const route = useRoute()
const router = useRouter()
const post = ref(null)
const comments = ref([])
const content = ref('')
const err = ref('')

async function load() {
  try {
    const r = await api.getPost(route.params.id)
    post.value = r.post
    comments.value = r.comments
  } catch (e) {
    err.value = e.message
  }
}
onMounted(load)

async function like() {
  await api.likePost(route.params.id)
  post.value.likes++
}

async function comment() {
  if (!api.isLoggedIn()) {
    router.push('/me')
    return
  }
  const c = content.value.trim()
  if (!c) return
  try {
    await api.createComment(route.params.id, c)
    content.value = ''
    load()
  } catch (e) {
    err.value = e.message
  }
}
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <button class="back" @click="router.push('/community')">‹ 返回社区</button>

    <p v-if="err" style="color: var(--red); font-size: 13px">{{ err }}</p>

    <template v-if="post">
      <div class="card">
        <div><span class="post-cat">{{ post.category }}</span></div>
        <h2 style="margin: 10px 0">{{ post.title }}</h2>
        <div style="color: var(--ink); line-height: 1.8; white-space: pre-wrap">{{ post.content }}</div>
        <div class="post-meta">
          <span>{{ post.created_at?.slice(0, 10) }}</span>
          <button class="btn btn-ghost btn-sm" style="width:auto" @click="like">👍 {{ post.likes }}</button>
        </div>
      </div>

      <div class="card">
        <div class="section-head"><span class="ic">💬</span><h3>评论 {{ comments.length }}</h3></div>
        <div v-for="c in comments" :key="c.id" style="border-bottom: 1px solid var(--line); padding: 8px 0">
          <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 4px">
            <img v-if="c.avatar" :src="c.avatar" style="width: 26px; height: 26px; border-radius: 6px" />
            <span v-else style="width: 26px; height: 26px; border-radius: 6px; background: #f1e3c8; display: inline-flex; align-items: center; justify-content: center; font-size: 13px">🙂</span>
            <span style="font-size: 13px; font-weight: 600">{{ c.nickname || '匿名' }}</span>
          </div>
          <div style="color: var(--ink); font-size: 14px">{{ c.content }}</div>
          <div style="color: var(--muted); font-size: 12px; margin-top: 4px">{{ c.created_at?.slice(0, 16) }}</div>
        </div>
        <div v-if="comments.length === 0" style="color: var(--muted); font-size: 13px; text-align: center; padding: 10px 0">暂无评论</div>
        <div style="display: flex; gap: 8px; margin-top: 10px">
          <input class="input" v-model="content" placeholder="写下你的看法…" @keyup.enter="comment" />
          <button class="btn btn-red btn-sm" style="width:auto" @click="comment">评论</button>
        </div>
      </div>
    </template>
  </div>
</template>