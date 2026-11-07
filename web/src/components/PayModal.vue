<script setup>
import { ref } from 'vue'
import { api } from '../api'

const props = defineProps({
  show: Boolean,
  order: Object, // { order_no, amount, provider, pay_params, mock_pay_url, pay_url }
  title: { type: String, default: '解锁全文' },
  desc: { type: String, default: '解锁后即可查看完整内容，可反复查看。' },
})
const emit = defineEmits(['close', 'paid'])

const paying = ref(false)
const err = ref('')
const copied = ref(false)

async function pay() {
  if (!props.order) return
  err.value = ''
  paying.value = true
  try {
    const o = props.order
    if (o.provider === 'mock') {
      await api.mockPay(o.order_no)
      emit('paid')
    } else if (o.provider === 'wechat') {
      const pp = o.pay_params || {}
      if (pp.h5_url) {
        window.location.href = pp.h5_url
      } else if (pp.paySign) {
        payJSAPI(pp)
      } else if (pp.code_url) {
        err.value = '请在微信中扫码支付（Native）'
      } else {
        err.value = '未识别的微信支付参数'
      }
    } else if (o.provider === 'alipay') {
      if (o.pay_url) {
        window.location.href = o.pay_url
      } else {
        err.value = '未识别的支付宝参数'
      }
    }
  } catch (e) {
    err.value = e.message
  } finally {
    paying.value = false
  }
}

function payJSAPI(pp) {
  function invoke() {
    WeixinJSBridge.invoke(
      'getBrandWCPayRequest',
      {
        appId: pp.appId,
        timeStamp: pp.timeStamp,
        nonceStr: pp.nonceStr,
        package: pp.package,
        signType: pp.signType || 'RSA',
        paySign: pp.paySign,
      },
      (res) => {
        if (res && res.err_msg === 'get_brand_wcpay_request:ok') {
          emit('paid')
        } else {
          err.value = '支付未完成或已取消'
        }
      },
    )
  }
  if (typeof WeixinJSBridge === 'undefined') {
    document.addEventListener('WeixinJSBridgeReady', invoke, false)
  } else {
    invoke()
  }
}

function copyCode() {
  if (!props.order?.pay_params?.code_url) return
  navigator.clipboard?.writeText(props.order.pay_params.code_url).then(() => {
    copied.value = true
    setTimeout(() => (copied.value = false), 1500)
  })
}
</script>

<template>
  <div class="modal-mask" v-if="show" @click.self="emit('close')">
    <div class="modal">
      <h3>🔒 {{ title }}</h3>
      <div class="price">¥{{ ((order?.amount || 0) / 100).toFixed(2) }}</div>
      <div class="desc">{{ desc }}</div>

      <p v-if="err" style="color: var(--red); font-size: 13px">{{ err }}</p>

      <div v-if="order?.provider === 'wechat' && order?.pay_params?.code_url" style="margin-bottom: 12px">
        <div style="font-size: 12px; color: var(--muted); margin-bottom: 6px">请复制以下链接到微信中打开，或扫码支付</div>
        <div style="display: flex; gap: 8px">
          <input class="input" style="font-size: 12px" :value="order.pay_params.code_url" readonly />
          <button class="btn btn-ghost btn-sm" style="width: auto; white-space: nowrap" @click="copyCode">{{ copied ? '已复制' : '复制链接' }}</button>
        </div>
      </div>

      <button class="btn btn-red" :disabled="paying" @click="pay">
        {{ paying ? '处理中…' : order?.provider === 'mock' ? '模拟支付（测试用）' : '立即支付' }}
      </button>
      <button class="btn btn-ghost btn-sm" style="margin-top: 10px" @click="emit('close')">取消</button>
      <p style="font-size: 12px; color: var(--muted); margin: 12px 0 0">
        {{ order?.provider === 'mock' ? '当前为本地模拟支付，用于跑通流程。' : order?.provider === 'alipay' ? '即将跳转支付宝完成支付。' : '支付由微信支付提供。' }}
      </p>
    </div>
  </div>
</template>