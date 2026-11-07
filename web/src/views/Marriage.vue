<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import PersonForm from '../components/PersonForm.vue'
import { api } from '../api'

const router = useRouter()
const a = reactive({ name: '男方', gender: 1, year: 1992, month: 3, day: 5, hour: 12, minute: 0, is_lunar: false })
const b = reactive({ name: '女方', gender: 0, year: 1994, month: 7, day: 9, hour: 12, minute: 0, is_lunar: false })
const loading = ref(false)
const err = ref('')

async function submit() {
  loading.value = true
  err.value = ''
  try {
    const res = await api.createReading({ scenario: 'marriage', person_a: a, person_b: b })
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
    <h2 class="page-title">💞 合婚匹配</h2>
    <p class="page-sub">分别填写两人信息，看生肖、五行、喜用是否契合</p>

    <PersonForm v-model="a" title="甲方信息" />
    <PersonForm v-model="b" title="乙方信息" />

    <p v-if="err" style="color: var(--red); font-size: 14px">{{ err }}</p>

    <button class="btn btn-primary" :disabled="loading" @click="submit">
      {{ loading ? '匹配中…' : '开始合婚' }}
    </button>
  </div>
</template>