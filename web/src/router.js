import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'home', component: () => import('./views/Home.vue') },
  { path: '/master', name: 'master', component: () => import('./views/Master.vue') },
  { path: '/personal', name: 'personal', component: () => import('./views/Personal.vue') },
  { path: '/marriage', name: 'marriage', component: () => import('./views/Marriage.vue') },
  { path: '/hiring', name: 'hiring', component: () => import('./views/Hiring.vue') },
  { path: '/naming', name: 'naming', component: () => import('./views/Naming.vue') },
  { path: '/wealth', name: 'wealth', component: () => import('./views/Wealth.vue') },
  { path: '/result/:id', name: 'result', component: () => import('./views/Result.vue') },
  { path: '/me', name: 'me', component: () => import('./views/Me.vue') },
  { path: '/doc/:key', name: 'doc', component: () => import('./views/Doc.vue') },
  { path: '/member', name: 'member', component: () => import('./views/Member.vue') },
  { path: '/ai', name: 'ai', component: () => import('./views/Ai.vue') },
  { path: '/community', name: 'community', component: () => import('./views/Community.vue') },
  { path: '/community/:id', name: 'post', component: () => import('./views/PostDetail.vue') },
  { path: '/admin', name: 'admin', component: () => import('./views/Admin.vue') },
  { path: '/fun', name: 'fun', component: () => import('./views/Fun.vue') },
  { path: '/fun/dice', name: 'dice', component: () => import('./views/Dice.vue') },
  { path: '/fun/music', name: 'music', component: () => import('./views/Music.vue') },
  { path: '/fun/guess', name: 'guess', component: () => import('./views/Guess.vue') },
  { path: '/fun/flying', name: 'flying', component: () => import('./views/Flying.vue') },
  { path: '/fun/weiqi', name: 'weiqi', component: () => import('./views/Weiqi.vue') },
  { path: '/fun/match3', name: 'match3', component: () => import('./views/Match3.vue') },
  { path: '/fun/gomoku', name: 'gomoku', component: () => import('./views/Gomoku.vue') },
  { path: '/fun/xiangqi', name: 'xiangqi', component: () => import('./views/Xiangqi.vue') },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

// 旧缓存兜底：新版发布后，若浏览器里的旧包用 hash 路由跳转去加载已被删除的旧 chunk，
// 会报「Failed to fetch dynamically imported module」。此时自动整页刷新一次，载入新版。
router.onError((error) => {
  const msg = error?.message || String(error)
  if (msg.includes('dynamically imported module') || msg.includes('Failed to fetch') || msg.includes('Loading chunk')) {
    if (!sessionStorage.getItem('stale_reload')) {
      sessionStorage.setItem('stale_reload', '1')
      location.reload()
    }
  }
})

export default router