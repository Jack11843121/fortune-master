<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import PersonForm from '../components/PersonForm.vue'
import { api } from '../api'

const router = useRouter()
const form = reactive({ name: '', gender: 1, year: 1990, month: 6, day: 15, hour: 12, minute: 0, is_lunar: false })
const loading = ref(false)
const err = ref('')

async function submit() {
  loading.value = true
  err.value = ''
  try {
    const res = await api.createReading({ scenario: 'personal', person_a: form })
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
    <h2 class="page-title">☯️ 个人命理</h2>
    <p class="page-sub">填写出生信息，为您排出八字命盘并详细解读</p>

    <PersonForm v-model="form" title="个人信息" />

    <p v-if="err" style="color: var(--red); font-size: 14px">{{ err }}</p>

    <button class="btn btn-primary" :disabled="loading" @click="submit">
      {{ loading ? '排盘中…' : '开始算命' }}
    </button>
  </div>
</template>