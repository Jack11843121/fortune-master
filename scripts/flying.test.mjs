// 飞行棋引擎单测：node scripts/flying.test.mjs
import { newGame, move, options } from '../web/src/game/flying.js'

let failed = 0
function ok(cond, msg) {
  if (cond) console.log('  ✓ ' + msg)
  else { console.error('  ✗ ' + msg); failed++ }
}

// 1. 起飞
{
  const g = newGame()
  g.lastRoll = 6
  ok(options(g).canTakeoff === true, '掷 6 可从机场起飞')
  const p = g.planes.find((x) => x.p === 1 && x.pos === -1)
  move(g, p, 6)
  ok(p.pos === 0 && p.steps === 0, '红方飞机起飞到起点 0')
}

// 2. 前进 + 击落
{
  const g = newGame()
  const red = g.planes.find((x) => x.p === 1 && x.pos === -1)
  g.lastRoll = 6
  move(g, red, 6) // 红到 0
  const blue = g.planes.find((x) => x.p === 2 && x.pos === -1)
  blue.pos = 3
  blue.steps = 3 // 蓝在 3 号位
  g.turn = 1
  g.lastRoll = 3
  move(g, red, 3) // 红 0 → 3
  ok(red.pos === 3, '红前进到 3 号位')
  ok(blue.pos === -1, '红落地击落蓝方飞机（回机场）')
}

// 3. 走满一圈到家
{
  const g = newGame()
  const red = g.planes.find((x) => x.p === 1 && x.pos === -1)
  g.lastRoll = 6
  move(g, red, 6) // 起飞 steps=0
  g.turn = 1
  move(g, red, 24) // 直接跳满一圈（测试到家分支）
  ok(red.pos === -2, '走满一圈到家（pos=-2）')
  ok(g.finished[1] === 1, '红方已完成 1 架')
}

// 4. 掷 6 再掷一次
{
  const g = newGame()
  const red = g.planes.find((x) => x.p === 1 && x.pos === -1)
  g.lastRoll = 6
  move(g, red, 6)
  ok(g.rollAgain === true, '掷 6 后可再掷一次')
}

console.log(failed === 0 ? '\n全部通过 ✅' : `\n${failed} 项失败 ❌`)
process.exit(failed === 0 ? 0 : 1)