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
    const res = await api.createReading({ scenario: 'wealth', person_a: form })
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
    <h2 class="page-title">💰 财运解析</h2>
    <p class="page-sub">财运走向 · 投资性格 · 股市板块倾向 · 求财开运</p>

    <PersonForm v-model="form" title="个人信息" />

    <p v-if="err" style="color: var(--red); font-size: 14px">{{ err }}</p>

    <button class="btn btn-primary" :disabled="loading" @click="submit">
      {{ loading ? '解析中…' : '开始解析财运' }}
    </button>

    <div class="notice">💡 财运与投资分析为性格参考，股市有风险，请量力而行、理性决策。</div>
  </div>
</template>