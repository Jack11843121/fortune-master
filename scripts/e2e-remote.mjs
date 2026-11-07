// 发布后全量回归测试：node scripts/e2e-remote.mjs [BASE_URL]
// 覆盖：meta / 个人命理排盘 / 合婚 / 招聘 / 起名 / 下单-支付-解锁全流程 / 微信登录降级 / 前端静态资源
// 用法：默认测本地 http://localhost:8080；可传参数或设 BASE 环境变量指向线上，例如 node scripts/e2e-remote.mjs https://yourdomain.com
const BASE = process.argv[2] || process.env.BASE || 'http://localhost:8080'

let pass = 0
let fail = 0

function ok(msg) { pass++; console.log('  ✓ ' + msg) }
function bad(msg) { fail++; console.log('  ✗ ' + msg) }

async function req(method, path, body) {
  const res = await fetch(BASE + path, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: body != null ? JSON.stringify(body) : undefined,
  })
  const text = await res.text()
  let data
  try { data = JSON.parse(text) } catch { data = text }
  return { status: res.status, data }
}

function assert(cond, msg) { if (!cond) throw new Error(msg) }

async function test(name, fn) {
  try { await fn(); ok(name) } catch (e) { bad(`${name} —— ${e.message}`) }
}

async function main() {
  console.log('测试目标: ' + BASE + '\n')

  await test('meta 全局配置', async () => {
    const { status, data } = await req('GET', '/api/v1/meta')
    assert(status === 200, '状态 ' + status)
    assert(data.app_name === '算命大师', 'app_name=' + data.app_name)
    assert(data.pay_provider === 'mock', '支付通道 mock')
    assert(data.price_cents === 990, '价格 990')
  })

  await test('个人命理排盘', async () => {
    const r = await req('POST', '/api/v1/readings', {
      scenario: 'personal',
      person_a: { name: '张三', gender: 1, year: 1990, month: 6, day: 15, hour: 10, minute: 30, is_lunar: false },
    })
    assert(r.status === 201, '状态 ' + r.status)
    assert(r.data.bazi.pillars.length === 4, '四柱齐全')
    assert(r.data.bazi.day_master === '辛', '日主=辛 实际=' + r.data.bazi.day_master)
    assert(r.data.locked === true, '预览锁定')
    assert(r.data.sections.length > 0, '有预览内容')
  })

  await test('合婚匹配', async () => {
    const r = await req('POST', '/api/v1/readings', {
      scenario: 'marriage',
      person_a: { name: '男', gender: 1, year: 1992, month: 3, day: 5, hour: 12, minute: 0, is_lunar: false },
      person_b: { name: '女', gender: 0, year: 1994, month: 7, day: 9, hour: 12, minute: 0, is_lunar: false },
    })
    assert(r.status === 201, '状态 ' + r.status)
    const s = r.data.match.score.value
    assert(s >= 0 && s <= 100, '评分 0-100 实际=' + s)
  })

  await test('招聘看八字', async () => {
    const r = await req('POST', '/api/v1/readings', {
      scenario: 'hiring',
      person_a: { name: '候选人', gender: 1, year: 1995, month: 8, day: 20, hour: 10, minute: 0, is_lunar: false },
      person_b: { name: '老板', gender: 1, year: 1980, month: 1, day: 1, hour: 10, minute: 0, is_lunar: false },
      role: '技术研发',
    })
    assert(r.status === 201, '状态 ' + r.status)
    assert(r.data.bazi_b, '返回老板八字')
  })

  await test('宝宝起名', async () => {
    const r = await req('POST', '/api/v1/naming', {
      surname: '王', gender: 1, name_len: 2,
      birth: { year: 2025, month: 3, day: 8, hour: 9, minute: 0, is_lunar: false },
    })
    assert(r.status === 200, '状态 ' + r.status)
    assert(Array.isArray(r.data.suggestions) && r.data.suggestions.length > 0, '有起名建议')
  })

  await test('下单→支付→解锁 全流程', async () => {
    const r1 = await req('POST', '/api/v1/readings', {
      scenario: 'personal',
      person_a: { name: '流程', gender: 0, year: 1988, month: 1, day: 1, hour: 8, minute: 0, is_lunar: false },
    })
    const id = r1.data.id
    const o = await req('POST', '/api/v1/orders', { reading_id: id })
    assert(o.data.order_no, '返回订单号')
    assert(o.data.amount === 990, '金额 990')
    const pay = await req('POST', '/api/v1/pay/mock/confirm', { order_no: o.data.order_no })
    assert(pay.data.status === 'paid', '支付成功')
    const r2 = await req('GET', `/api/v1/readings/${id}`)
    assert(r2.data.locked === false, '已解锁')
    assert(r2.data.locked_titles.length === 0, '无锁定标题')
    assert(r2.data.sections.length >= 8, '全文完整段落 实际=' + r2.data.sections.length)
  })

  await test('微信登录(未配置应优雅降级)', async () => {
    const r = await req('POST', '/api/v1/wechat/login', { code: 'test' })
    assert(r.status === 400, '返回 400 而非 500 实际=' + r.status)
  })

  await test('前端首页 + 静态资源可加载', async () => {
    const res = await fetch(BASE + '/')
    const html = await res.text()
    assert(res.status === 200, '首页 200')
    assert(html.includes('assets/index-'), '引用构建资源')
    const m = html.match(/assets\/[^"']+\.(js|css)/g) || []
    assert(m.length > 0, '解析到资源引用')
    for (const a of m) {
      const rr = await fetch(BASE + '/' + a)
      assert(rr.status === 200, a + ' 可加载')
    }
  })

  // —— 会员 / AI / 社区 / 管理后台（新增能力回归）——

  async function authReq(token, method, path, body) {
    const res = await fetch(BASE + path, {
      method,
      headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
      body: body != null ? JSON.stringify(body) : undefined,
    })
    const text = await res.text()
    let data
    try { data = JSON.parse(text) } catch { data = text }
    return { status: res.status, data }
  }
  async function devLogin() {
    const r = await req('POST', '/api/v1/auth/dev', {})
    assert(r.data.token, '开发登录返回 token')
    return r.data.token
  }

  await test('账户登录（开发模式）', async () => {
    const token = await devLogin()
    const me = await authReq(token, 'GET', '/api/v1/auth/profile')
    assert(me.status === 200 && me.data.id, '登录后可获取资料')
    const noauth = await fetch(BASE + '/api/v1/auth/profile')
    assert(noauth.status === 401, '未登录访问资料被拒')
  })

  await test('财运 / 股市场景', async () => {
    const r = await req('POST', '/api/v1/readings', {
      scenario: 'wealth',
      person_a: { name: '张三', gender: 1, year: 1990, month: 6, day: 15, hour: 10, minute: 30, is_lunar: false },
    })
    assert(r.status === 201, '状态 ' + r.status)
    assert(r.data.scenario === 'wealth', '场景 wealth')
    assert(r.data.sections.length >= 2, '财运段落 实际=' + r.data.sections.length)
  })

  await test('会员订阅→支付→发货', async () => {
    const token = await devLogin()
    const o = await authReq(token, 'POST', '/api/v1/membership/subscribe', {})
    assert(o.data.order_no, '订阅订单号')
    const pay = await req('POST', '/api/v1/pay/mock/confirm', { order_no: o.data.order_no })
    assert(pay.data.status === 'paid', '订阅支付成功')
    const m = await authReq(token, 'GET', '/api/v1/me')
    assert(m.data.vip === true, '已是会员')
    assert(m.data.ai_credits === 1000, '发放 1000 次 实际=' + m.data.ai_credits)
  })

  await test('会员全站免解锁', async () => {
    const token = await devLogin()
    const o = await authReq(token, 'POST', '/api/v1/membership/subscribe', {})
    await req('POST', '/api/v1/pay/mock/confirm', { order_no: o.data.order_no })
    const ppl = { name: '会员', gender: 1, year: 1990, month: 6, day: 15, hour: 10, minute: 30, is_lunar: false }
    const r = await authReq(token, 'POST', '/api/v1/readings', { scenario: 'personal', person_a: ppl })
    assert(r.data.locked === false, '会员报告直接解锁 实际 locked=' + r.data.locked)
    assert(r.data.sections.length >= 8, '会员看到全文 实际=' + r.data.sections.length)
    const r2 = await req('POST', '/api/v1/readings', { scenario: 'personal', person_a: ppl })
    assert(r2.data.locked === true, '非会员报告锁定 实际 locked=' + r2.data.locked)
  })

  await test('AI 算命师（本地兜底）', async () => {
    const token = await devLogin()
    const o = await authReq(token, 'POST', '/api/v1/membership/subscribe', {})
    await req('POST', '/api/v1/pay/mock/confirm', { order_no: o.data.order_no })
    const a = await authReq(token, 'POST', '/api/v1/ai/chat', { message: '今年财运如何？' })
    assert(a.status === 200, '有回复 状态=' + a.status)
    assert(a.data.reply && a.data.reply.length > 5, '回复非空')
  })

  await test('交流社区', async () => {
    const token = await devLogin()
    const p = await authReq(token, 'POST', '/api/v1/community/posts', { title: '测试帖', content: '命理测试内容', category: '命理趣谈' })
    assert(p.status === 200 && p.data.id, '发帖成功')
    const list = await req('GET', '/api/v1/community/posts')
    assert(list.data.total >= 1, '帖子可见')
  })

  await test('管理后台登录 + 看板', async () => {
    const adminUser = process.env.ADMIN_USERNAME || 'admin'
    const adminPass = process.env.ADMIN_PASSWORD || 'admin123'
    const bad = await req('POST', '/api/v1/admin/login', { username: adminUser, password: 'wrong' })
    assert(bad.status === 401, '错误密码拒绝 实际=' + bad.status)
    const ok = await req('POST', '/api/v1/admin/login', { username: adminUser, password: adminPass })
    assert(ok.status === 200 && ok.data.token, '正确登录')
    const dash = await fetch(BASE + '/api/v1/admin/dashboard', { headers: { Authorization: 'Bearer ' + ok.data.token } })
    const dd = await dash.json()
    assert(dd.users >= 0, '看板数据返回')
    const noauth = await fetch(BASE + '/api/v1/admin/dashboard')
    assert(noauth.status === 401, '未登录拒绝')
  })

  console.log(`\n====================`)
  console.log(`结果: ${pass} 通过, ${fail} 失败`)
  if (fail > 0) process.exit(1)
}

main().catch(e => { console.error('脚本异常:', e); process.exit(1) })