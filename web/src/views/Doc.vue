<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { docs } from '../data/legal'

const route = useRoute()
const router = useRouter()
const doc = computed(() => docs[route.params.key])
const others = computed(() => Object.keys(docs).filter((k) => k !== route.params.key))
</script>

<template>
  <div class="page" style="padding-top: 16px">
    <button class="back" style="margin-bottom: 10px" @click="router.back()">‹ 返回</button>

    <div v-if="doc" class="card">
      <h2 style="margin: 0 0 4px">{{ doc.title }}</h2>
      <div style="color: var(--muted); font-size: 12px; margin-bottom: 14px">更新日期：{{ doc.updated }}</div>
      <div class="legal-body">{{ doc.content }}</div>

      <div style="margin-top: 18px; border-top: 1px solid var(--line); padding-top: 12px">
        <div style="font-size: 13px; color: var(--muted); margin-bottom: 6px">相关文档</div>
        <router-link
          v-for="k in others"
          :key="k"
          :to="`/doc/${k}`"
          style="display: block; padding: 6px 0; font-size: 14px; color: var(--red-deep); text-decoration: none"
        >
          {{ docs[k].title }} ›
        </router-link>
      </div>
    </div>
    <div v-else class="card" style="text-align: center; color: var(--muted)">文档不存在</div>
  </div>
</template>