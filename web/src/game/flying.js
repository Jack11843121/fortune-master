// 飞行棋（灵签版）规则引擎，无 Vue 依赖。
// 2 人（红 1 / 蓝 2），各 4 架飞机；环形跑道 24 格；红起点 0，蓝起点 12。
// 飞机 pos：-2 已完成、-1 机场、0..23 跑道。起飞需掷 5 或 6。
// 落地特殊格：吉=再掷一次、签=抽灵签、凶=后退一步；落到对方棋子所在格=击落回机场。

export const TRACK = 24
export const N_PLANES = 4
export const BASE = { 1: 0, 2: 12 }
export const TAKE_OFF = [5, 6]
export const SPECIAL = { 3: 'ji', 6: 'qian', 9: 'ji', 15: 'ji', 18: 'xiong', 21: 'ji' }

export function newGame() {
  const planes = []
  for (const p of [1, 2]) {
    for (let i = 0; i < N_PLANES; i++) planes.push({ p, id: i, pos: -1, steps: 0 })
  }
  return { turn: 1, planes, finished: { 1: 0, 2: 0 }, winner: 0, lastRoll: 0, rollAgain: false, qian: false }
}

export function roll() {
  return 1 + Math.floor(Math.random() * 6)
}

// 当前回合方可操作的飞机（在跑道上的），以及能否起飞
export function options(g) {
  const me = g.planes.filter((x) => x.p === g.turn)
  const onTrack = me.filter((x) => x.pos >= 0)
  const canTakeoff = me.some((x) => x.pos === -1) && TAKE_OFF.includes(g.lastRoll)
  return { onTrack, canTakeoff }
}

// 移动（起飞或前进）；返回是否有效
export function move(g, plane, r) {
  if (!plane || plane.p !== g.turn) return false
  g.qian = false
  g.rollAgain = r === 6
  g.lastRoll = r

  if (plane.pos === -1) {
    // 起飞到起点
    plane.pos = BASE[g.turn]
    plane.steps = 0
  } else {
    plane.steps += r
    if (plane.steps >= TRACK) {
      // 走满一圈，到家
      plane.pos = -2
      plane.steps = TRACK
      g.finished[g.turn]++
      if (g.finished[g.turn] >= N_PLANES) {
        g.winner = g.turn
        g.rollAgain = false
        return true
      }
    } else {
      plane.pos = (BASE[g.turn] + plane.steps) % TRACK
      // 击落对方
      for (const o of g.planes) {
        if (o.p !== g.turn && o.pos === plane.pos) {
          o.pos = -1
          o.steps = 0
        }
      }
      // 特殊格
      const sp = SPECIAL[plane.pos]
      if (sp === 'ji') {
        g.rollAgain = true
      } else if (sp === 'xiong' && plane.steps > 1) {
        plane.steps -= 1
        plane.pos = (BASE[g.turn] + plane.steps) % TRACK
      } else if (sp === 'qian') {
        g.qian = true
      }
    }
  }

  if (!g.winner && !g.rollAgain) {
    g.turn = g.turn === 1 ? 2 : 1
  }
  return true
}