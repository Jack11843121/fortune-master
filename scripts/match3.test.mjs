// 消消乐核心逻辑单测：node scripts/match3.test.mjs
import { makeBoard, findMatches, hasAnyMove, findHint, dropAndRefill, ROWS, COLS } from '../web/src/game/match3.js'

let failed = 0
function ok(cond, msg) {
  if (cond) console.log('  ✓ ' + msg)
  else { console.error('  ✗ ' + msg); failed++ }
}

// 1. 生成 100 个新棋盘，每个都无初始三连，且至少有一处可交换
for (let i = 0; i < 100; i++) {
  const b = makeBoard()
  if (findMatches(b).size !== 0) { ok(false, '棋盘有初始三连'); break }
}
ok(true, '多次生成：棋盘均无初始三连')
for (let i = 0; i < 100; i++) {
  const b = makeBoard()
  if (!hasAnyMove(b)) { ok(false, '存在无解棋盘'); break }
}
ok(true, '多次生成：棋盘均有至少一步可消')

// 2. 三连消除检测
const g = [
  [0, 1, 2, 3, 4, 0],
  [0, 1, 2, 3, 4, 1],
  [0, 1, 2, 3, 4, 2],
  [1, 2, 3, 4, 0, 1],
  [2, 3, 4, 0, 1, 2],
  [3, 4, 0, 1, 2, 3],
]
// 第一列前三格均为 0 → 竖向三连
const m = findMatches(g)
ok(m.has('0,0') && m.has('1,0') && m.has('2,0'), '竖向三连被识别')

// 3. findHint 给出的交换确实能消除
for (let i = 0; i < 20; i++) {
  const b = makeBoard()
  const h = findHint(b)
  if (!h) { ok(false, 'findHint 返回空'); break }
  const g = b.map((r) => [...r])
  ;[g[h.a.r][h.a.c], g[h.b.r][h.b.c]] = [g[h.b.r][h.b.c], g[h.a.r][h.a.c]]
  if (findMatches(g).size === 0) { ok(false, 'findHint 给出的交换不能消除'); break }
}
ok(true, 'findHint 每次都能给出可消除的一步')

// 4. 下落补位后所有格子均 >=0
const g2 = [
  [-1, 1, 2, 3, 4, 0],
  [0, -1, 2, 3, 4, 1],
  [0, 1, -1, 3, 4, 2],
  [1, 2, 3, 4, 0, 1],
  [2, 3, 4, 0, 1, 2],
  [3, 4, 0, 1, 2, 3],
]
dropAndRefill(g2)
const allFilled = g2.every((row) => row.every((v) => v >= 0 && v < 5))
ok(allFilled && g2.length === ROWS && g2[0].length === COLS, '下落补位后棋盘完整且取值合法')

console.log(failed === 0 ? '\n全部通过 ✅' : `\n${failed} 项失败 ❌`)
process.exit(failed === 0 ? 0 : 1)