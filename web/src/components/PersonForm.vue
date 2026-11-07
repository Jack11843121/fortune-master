<script setup>
defineProps({
  title: { type: String, default: '出生信息' },
  showName: { type: Boolean, default: true },
})
const model = defineModel({ type: Object, required: true })

const years = Array.from({ length: 2031 - 1930 + 1 }, (_, i) => 2031 - i) // 2031 → 1930（倒序）
const hours = Array.from({ length: 24 }, (_, i) => i)
const minutes = Array.from({ length: 60 }, (_, i) => i)
</script>

<template>
  <div class="card">
    <div class="section-head">
      <span class="ic">👤</span>
      <h3>{{ title }}</h3>
    </div>

    <div class="field" v-if="showName">
      <label>称呼（选填）</label>
      <input class="input" v-model="model.name" placeholder="如：张先生" />
    </div>

    <div class="field">
      <label>性别</label>
      <div class="seg">
        <button :class="{ on: model.gender === 1 }" type="button" @click="model.gender = 1">男</button>
        <button :class="{ on: model.gender === 0 }" type="button" @click="model.gender = 0">女</button>
      </div>
    </div>

    <div class="field">
      <label>历法</label>
      <div class="seg">
        <button :class="{ on: !model.is_lunar }" type="button" @click="model.is_lunar = false">公历（阳历）</button>
        <button :class="{ on: model.is_lunar }" type="button" @click="model.is_lunar = true">农历（阴历）</button>
      </div>
    </div>

    <div class="field">
      <label>出生日期</label>
      <div class="row">
        <select class="input" v-model="model.year">
          <option v-for="y in years" :key="y" :value="y">{{ y }}年</option>
        </select>
        <select class="input" v-model="model.month">
          <option v-for="m in 12" :key="m" :value="m">{{ m }}月</option>
        </select>
        <select class="input" v-model="model.day">
          <option v-for="d in 31" :key="d" :value="d">{{ d }}日</option>
        </select>
      </div>
    </div>

    <div class="field">
      <label>出生时间（时辰）</label>
      <div class="row">
        <select class="input" v-model="model.hour">
          <option v-for="h in hours" :key="h" :value="h">{{ h }}时</option>
        </select>
        <select class="input" v-model="model.minute">
          <option v-for="mi in minutes" :key="mi" :value="mi">{{ mi }}分</option>
        </select>
      </div>
      <div style="font-size: 12px; color: var(--muted); margin-top: 5px">
        若记不清具体时间，可选中午 12 时
      </div>
    </div>
  </div>
</template>