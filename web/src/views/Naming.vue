<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const form = reactive({ surname: '', gender: 1, nameLen: 2, year: 2025, month: 6, day: 1, hour: 12, minute: 0, is_lunar: false })
const result = ref(null)
const shown = ref(8)
const loading = ref(false)
const err = ref('')

const years = Array.from({ length: 2031 - 1980 + 1 }, (_, i) => 2031 - i)
const hours = Array.from({ length: 24 }, (_, i) => i)
const minutes = Array.from({ length: 60 }, (_, i) => i)

const wuxingColor = { 金: '#c9a24b', 木: '#4c9a5a', 水: '#4a86d9', 火: '#d94a3a', 土: '#a8794f' }

async function submit() {
  if (!form.surname.trim()) { err.value = '请先填写姓氏'; return }
  loading.value = true
  err.value = ''
  try {
    const res = await fetch('/api/v1/naming', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ surname: form.surname.trim(), gender: form.gender, name_len: form.nameLen, birth: { ...form, name: '' } }),
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok) throw new Error(data.error || '请求失败')
    result.value = data
    shown.value = 8
  } catch (e) {
    err.value = e.message
  } finally {
    loading.value = false
  }
}

function wx(chip) {
  return chip.split('+')
}
</script>

<template>
  <div class="page">
    <h2 class="page-title">📛 宝宝取名</h2>
    <p class="page-sub">按宝宝出生八字看五行喜用，推荐吉祥好名，寓意美好、五行补缺</p>

    <div class="card">
      <div class="field">
        <label>姓氏</label>
        <input class="input" v-model="form.surname" placeholder="如：王、李、欧阳" maxlength="2" />
      </div>
      <div class="field">
        <label>性别</label>
        <div class="seg">
          <button :class="{ on: form.gender === 1 }" type="button" @click="form.gender = 1">男孩</button>
          <button :class="{ on: form.gender === 0 }" type="button" @click="form.gender = 0">女孩</button>
        </div>
      </div>
      <div class="field">
        <label>名字字数</label>
        <div class="seg">
          <button :class="{ on: form.nameLen === 2 }" type="button" @click="form.nameLen = 2">双字名（如 王浩然）</button>
          <button :class="{ on: form.nameLen === 1 }" type="button" @click="form.nameLen = 1">单字名（如 王浩）</button>
        </div>
      </div>
      <div class="field">
        <label>宝宝出生日期</label>
        <div class="row">
          <select class="input" v-model="form.year"><option v-for="y in years" :key="y" :value="y">{{ y }}年</option></select>
          <select class="input" v-model="form.month"><option v-for="m in 12" :key="m" :value="m">{{ m }}月</option></select>
          <select class="input" v-model="form.day"><option v-for="d in 31" :key="d" :value="d">{{ d }}日</option></select>
        </div>
      </div>
      <div class="field">
        <label>出生时间（时辰）</label>
        <div class="row">
          <select class="input" v-model="form.hour"><option v-for="h in hours" :key="h" :value="h">{{ h }}时</option></select>
          <select class="input" v-model="form.minute"><option v-for="mi in minutes" :key="mi" :value="mi">{{ mi }}分</option></select>
        </div>
      </div>
      <p v-if="err" style="color: var(--red); font-size: 14px">{{ err }}</p>
      <button class="btn btn-primary" :disabled="loading" @click="submit">{{ loading ? '起名中…' : '开始起名' }}</button>
    </div>

    <div v-if="result">
      <div class="card" style="background: #fdf6e9">
        <div class="section-head"><span class="ic">🧧</span><h3>喜用神参考</h3></div>
        <p style="margin: 0; font-size: 14px; color: var(--muted)">
          宝宝日主为 <b style="color: var(--red-deep)">{{ result.day_master }}</b>，
          喜用神五行：<span class="tag" v-for="w in result.xi_yong" :key="w" :style="{ background: wuxingColor[w]+'22', color: wuxingColor[w] }">{{ w }}</span>
          ，取名宜多补这些五行。
        </p>
      </div>

      <div v-for="(s, i) in result.suggestions.slice(0, shown)" :key="i" class="card" style="padding: 16px">
        <div style="display: flex; align-items: flex-start; gap: 12px">
          <div style="flex: 1">
            <div style="font-size: 24px; font-weight: 800; color: var(--red-deep); letter-spacing: 3px">{{ s.name }}</div>
            <div style="font-size: 13px; color: var(--muted); margin-top: 2px">{{ s.pinyin }}</div>
          </div>
          <div style="text-align: center; min-width: 58px">
            <div style="font-size: 30px; font-weight: 800; color: var(--gold-deep)">{{ s.score }}</div>
            <div style="font-size: 11px; color: var(--muted)">推荐指数</div>
          </div>
        </div>
        <div style="margin: 8px 0 10px">
          <span class="tag" v-for="w in wx(s.wuxing)" :key="w" :style="{ background: wuxingColor[w] + '22', color: wuxingColor[w] }">{{ w }}</span>
        </div>
        <div style="font-size: 13.5px; line-height: 1.8; color: var(--ink); margin-bottom: 8px">
          <b style="color: var(--red-deep)">为什么推荐：</b>{{ s.reason }}
        </div>
        <div style="font-size: 13.5px; line-height: 1.8; color: var(--muted)">
          <b>字义详解：</b>{{ s.meaning }}
        </div>
      </div>

      <div v-if="shown < result.suggestions.length" style="text-align: center; margin: 8px 0">
        <button class="btn btn-ghost" @click="shown += 8">查看更多（已看 {{ shown }}/{{ result.suggestions.length }}）</button>
      </div>

      <p class="m3-hint" style="text-align: center">祝福宝宝健康快乐成长，名字承载父母的美好期许。</p>
    </div>
  </div>
</template>