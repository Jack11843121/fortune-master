// 统一 API 层：登录后自动携带 Bearer 令牌；管理后台单独带 admin token。

let token = localStorage.getItem('auth_token') || ''

function baseHeaders(extra = {}) {
  const h = { 'Content-Type': 'application/json' }
  if (token) h['Authorization'] = 'Bearer ' + token
  return { ...h, ...extra }
}

async function request(path, opts = {}) {
  let res
  try {
    res = await fetch(path, {
      method: opts.method || 'GET',
      headers: baseHeaders(opts.headers),
      body: opts.body != null ? JSON.stringify(opts.body) : undefined,
    })
  } catch (e) {
    throw new Error('网络请求失败，请检查网络后重试')
  }
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data.error || `请求失败（${res.status}）`)
  }
  return data
}

function adminRequest(path, opts = {}) {
  const t = localStorage.getItem('admin_token') || ''
  return request(path, { ...opts, headers: { Authorization: t ? `Bearer ${t}` : '' } })
}

export const api = {
  // —— 会话 ——
  setToken(t) { token = t; localStorage.setItem('auth_token', t) },
  clearToken() { token = ''; localStorage.removeItem('auth_token') },
  isLoggedIn() { return !!token },

  // —— 登录 ——
  authDev: (body) => request('/api/v1/auth/dev', { method: 'POST', body }),
  authWechat: (code) => request('/api/v1/auth/wechat', { method: 'POST', body: { code } }),
  authWechatOAuth: (code) => request('/api/v1/auth/wechat/oauth', { method: 'POST', body: { code } }),
  authWechatUrl: (redirect) => request('/api/v1/auth/wechat/url?redirect=' + encodeURIComponent(redirect)),
  authLogout: () => request('/api/v1/auth/logout', { method: 'POST' }),
  authDelete: () => request('/api/v1/auth/delete', { method: 'POST' }),
  authProfile: () => request('/api/v1/auth/profile'),

  // —— 基础 ——
  meta: () => request('/api/v1/meta'),
  createReading: (body) => request('/api/v1/readings', { method: 'POST', body }),
  getReading: (id) => request(`/api/v1/readings/${id}`),
  createOrder: (readingId) => request('/api/v1/orders', { method: 'POST', body: { reading_id: readingId } }),
  mockPay: (orderNo) => request('/api/v1/pay/mock/confirm', { method: 'POST', body: { order_no: orderNo } }),
  getOrder: (orderNo) => request(`/api/v1/orders/${orderNo}`),

  // —— 会员 ——
  me: () => request('/api/v1/me'),
  plans: () => request('/api/v1/membership/plans'),
  subscribe: () => request('/api/v1/membership/subscribe', { method: 'POST', body: {} }),
  buyQuota: () => request('/api/v1/membership/quota', { method: 'POST', body: {} }),

  // —— AI 算命师 ——
  aiChat: (body) => request('/api/v1/ai/chat', { method: 'POST', body }),
  aiHistory: (session) => request(`/api/v1/ai/history?session_id=${session || ''}`),

  // —— 交流社区 ——
  listPosts: (params = {}) => request('/api/v1/community/posts?' + new URLSearchParams(params).toString()),
  createPost: (body) => request('/api/v1/community/posts', { method: 'POST', body }),
  getPost: (id) => request(`/api/v1/community/posts/${id}`),
  likePost: (id) => request(`/api/v1/community/posts/${id}/like`, { method: 'POST' }),
  createComment: (id, content) => request(`/api/v1/community/posts/${id}/comments`, { method: 'POST', body: { content } }),

  // —— 内容分享 ——
  share: (readingId, platform = 'wechat') => request('/api/v1/share', { method: 'POST', body: { reading_id: readingId, platform } }),

  // —— 管理后台 ——
  adminLogin: (username, password) => request('/api/v1/admin/login', { method: 'POST', body: { username, password } }),
  adminMe: () => adminRequest('/api/v1/admin/me'),
  adminDashboard: () => adminRequest('/api/v1/admin/dashboard'),
  adminUsers: (params = {}) => adminRequest('/api/v1/admin/users?' + new URLSearchParams(params).toString()),
  adminUpdateUser: (id, body) => adminRequest(`/api/v1/admin/users/${id}`, { method: 'PATCH', body }),
  adminOrders: (params = {}) => adminRequest('/api/v1/admin/orders?' + new URLSearchParams(params).toString()),
  adminMemberships: (params = {}) => adminRequest('/api/v1/admin/memberships?' + new URLSearchParams(params).toString()),
  adminConsultations: (params = {}) => adminRequest('/api/v1/admin/consultations?' + new URLSearchParams(params).toString()),
  adminReplyConsult: (id, content) => adminRequest(`/api/v1/admin/consultations/${id}/reply`, { method: 'POST', body: { content } }),
  adminPosts: (params = {}) => adminRequest('/api/v1/admin/posts?' + new URLSearchParams(params).toString()),
  adminModeratePost: (id, status) => adminRequest(`/api/v1/admin/posts/${id}`, { method: 'PATCH', body: { status } }),
}