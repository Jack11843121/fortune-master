// 围棋引擎单测：node scripts/weiqi.test.mjs
import { createBoard, tryPlace, aiMove, countScore, boardString } from '../web/src/game/weiqi.js'

let failed = 0
function ok(cond, msg) {
  if (cond) console.log('  ✓ ' + msg)
  else { console.error('  ✗ ' + msg); failed++ }
}

// 1. 提子：白子封住黑子最后一口气后提掉
{
  const b = createBoard()
  b[1][1] = 1 // 黑
  b[0][1] = 2; b[1][0] = 2; b[1][2] = 2 // 白封三面
  const res = tryPlace(b, 2, 2, 1, null) // 白落最后一气
  ok(res && res.captured === 1, '白子提掉被围黑子')
  ok(res.board[1][1] === 0, '被提黑子从棋盘消失')
}

// 2. 自杀：无气且提不到子 → 非法
{
  const b = createBoard()
  b[4][3] = 2; b[4][5] = 2; b[3][4] = 2; b[5][4] = 2
  const res = tryPlace(b, 1, 4, 4, null)
  ok(res === null, '无气落子为自杀，判非法')
}

// 3. 打劫机制：落子后局面等于 koStr 时判非法
{
  const b = createBoard()
  b[0][0] = 1
  const normal = tryPlace(b, 2, 0, 1, null)
  ok(normal != null, 'koStr 不相同时正常落子')
  const ko = tryPlace(b, 2, 0, 1, boardString(normal.board)) // koStr 命中落子后局面
  ok(ko === null, 'koStr 命中时落子被判打劫非法')
}

// 4. 地盘：单色围起的空格计入该方（白 8 子围成盒，中心 1 目 + 外围黑子使外部中性地）
{
  const b = createBoard()
  for (const [r, c] of [[3,3],[3,4],[3,5],[4,3],[4,5],[5,3],[5,4],[5,5]]) b[r][c] = 2
  b[0][0] = 1 // 黑子在外部，使外部空区同时触黑白 → 中性地
  const s = countScore(b)
  ok(s.white === 9 && s.black === 1, '白 8 子 + 1 目围空 = 9，黑 1 子 = 1')
}

// 5. AI 返回合法落点
{
  const m = aiMove(createBoard(), 1, null)
  ok(m && m.r >= 0 && m.r < 9 && m.c >= 0 && m.c < 9, 'AI 返回合法落点')
}

console.log(failed === 0 ? '\n全部通过 ✅' : `\n${failed} 项失败 ❌`)
process.exit(failed === 0 ? 0 : 1)