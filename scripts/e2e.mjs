// 端到端自测脚本：node scripts/e2e.mjs
// 需要服务已在 http://localhost:8080 运行。
const BASE = 'http://localhost:8080'

async function req(method, path, body) {
  const res = await fetch(BASE + path, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: body != null ? JSON.stringify(body) : undefined,
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(`${method} ${path} -> ${res.status} ${JSON.stringify(data)}`)
  return data
}

function assert(cond, msg) {
  if (!cond) throw new Error('断言失败: ' + msg)
  console.log('  ✓ ' + msg)
}

async function main() {
  console.log('== 个人命理 ==')
  const r1 = await req('POST', '/api/v1/readings', {
    scenario: 'personal',
    person_a: { name: '张三', gender: 1, year: 1990, month: 6, day: 15, hour: 10, minute: 30, is_lunar: false },
  })
  assert(r1.id, '创建成功返回 id')
  assert(r1.locked === true, '预览状态为锁定')
  assert(r1.title.includes('张三'), '标题包含姓名中文 ' + r1.title)
  assert(r1.bazi.day_master === '辛', '日主正确（辛）')
  assert(r1.bazi.pillars.length === 4, '四柱齐全')
  assert(r1.bazi.pillars[0].gan_wuxing === '金', '年干五行正确（金，非合并值）')
  assert(r1.bazi.pillars[0].zhi_wuxing === '火', '年支五行正确（火）')
  assert(r1.sections.length < 8, '预览段数少于全文')
  assert(r1.locked_titles.length > 0, '存在锁定标题列表')
  console.log('  预览段数=%d 锁定项=%d', r1.sections.length, r1.locked_count)

  console.log('== 下单（模拟） ==')
  const o = await req('POST', '/api/v1/orders', { reading_id: r1.id })
  assert(o.order_no, '返回订单号')
  assert(o.provider === 'mock', '支付通道为 mock')
  assert(o.amount === 990, '金额正确 990 分')

  console.log('== 模拟支付 ==')
  const pay = await req('POST', '/api/v1/pay/mock/confirm', { order_no: o.order_no })
  assert(pay.status === 'paid', '支付状态为 paid')

  console.log('== 再次读取（应已解锁） ==')
  const r2 = await req('GET', `/api/v1/readings/${r1.id}`)
  assert(r2.locked === false, '已解锁')
  assert(r2.sections.length === 8, '全文 8 段全部返回')
  assert(r2.locked_titles.length === 0, '无锁定标题')

  console.log('== 合婚 ==')
  const m = await req('POST', '/api/v1/readings', {
    scenario: 'marriage',
    person_a: { name: '男方', gender: 1, year: 1992, month: 3, day: 5, hour: 12, minute: 0, is_lunar: false },
    person_b: { name: '女方', gender: 0, year: 1994, month: 7, day: 9, hour: 12, minute: 0, is_lunar: false },
  })
  assert(m.match, '存在匹配结果')
  assert(m.match.score.value >= 0 && m.match.score.value <= 100, '合婚评分在 0-100')
  assert(m.match.aspects.length === 2, '免费看前 2 个维度')
  assert(!m.match.conclusion, '结论在付费前隐藏')
  console.log('  合婚评分=%d(%s)', m.match.score.value, m.match.score.grade)

  console.log('== 招聘 ==')
  const h = await req('POST', '/api/v1/readings', {
    scenario: 'hiring',
    person_a: { name: '候选人', gender: 1, year: 1995, month: 8, day: 20, hour: 10, minute: 0, is_lunar: false },
    person_b: { name: '老板', gender: 1, year: 1980, month: 1, day: 1, hour: 10, minute: 0, is_lunar: false },
    role: '技术研发',
  })
  assert(h.match, '存在招聘评估结果')
  assert(h.bazi_b, '返回老板八字')
  console.log('  招聘评分=%d(%s) 结论(付费前应为空)=%s', h.match.score.value, h.match.score.grade, h.match.conclusion || '(空)')

  console.log('\n全部通过 ✅')
}

main().catch((e) => {
  console.error('\n失败 ❌', e.message)
  process.exit(1)
})