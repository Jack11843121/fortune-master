<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import PersonForm from '../components/PersonForm.vue'
import { api } from '../api'

const router = useRouter()
const candidate = reactive({ name: '', gender: 1, year: 1995, month: 8, day: 20, hour: 10, minute: 0, is_lunar: false })
const boss = reactive({ name: '', gender: 1, year: 1980, month: 1, day: 1, hour: 10, minute: 0, is_lunar: false })
const roles = ['技术研发', '销售市场', '管理领导', '财务审计', '创意设计', '行政人事', '运营执行']
const role = ref('技术研发')
const withBoss = ref(false)
const loading = ref(false)
const err = ref('')

async function submit() {
  loading.value = true
  err.value = ''
  try {
    const res = await api.createReading({
      scenario: 'hiring',
      person_a: candidate,
      person_b: withBoss.value ? boss : null,
      role: role.value,
    })
    router.push(`/result/${res.id}`)
  } catch (e) {
    err.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h2 class="page-title">💼 招聘看八字</h2>
    <p class="page-sub">评估候选人与岗位、团队的契合度，供招聘参考</p>

    <PersonForm v-model="candidate" title="候选人信息" show-name />

    <div class="field">
      <label>目标岗位</label>
      <select class="input" v-model="role">
        <option v-for="r in roles" :key="r" :value="r">{{ r }}</option>
      </select>
    </div>

    <div class="field">
      <label>是否填写老板/团队负责人八字（用于匹配）</label>
      <div class="seg">
        <button :class="{ on: !withBoss }" type="button" @click="withBoss = false">不填</button>
        <button :class="{ on: withBoss }" type="button" @click="withBoss = true">填写</button>
      </div>
    </div>

    <PersonForm v-if="withBoss" v-model="boss" title="老板/团队负责人信息" show-name />

    <p v-if="err" style="color: var(--red); font-size: 14px">{{ err }}</p>

    <button class="btn btn-primary" :disabled="loading" @click="submit">
      {{ loading ? '评估中…' : '开始评估' }}
    </button>
  </div>
</template>