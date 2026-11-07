// 中国象棋规则引擎单测：node scripts/xiangqi.test.mjs
import { createBoard, legalMovesFor, inCheck, pseudoMoves, generalsFacing, evaluate, bestMove } from '../web/src/game/xiangqi.js'

let failed = 0
function ok(cond, msg) {
  if (cond) console.log('  ✓ ' + msg)
  else { console.error('  ✗ ' + msg); failed++ }
}
const has = (list, r, c) => list.some((m) => m[0] === r && m[1] === c)

// 1. 初始局面双方都有合法走法
const b = createBoard()
ok(legalMovesFor(b, 'r').length > 0, '红方初始有合法走法')
ok(legalMovesFor(b, 'b').length > 0, '黑方初始有合法走法')

// 2. 车直线走，被己方兵挡住
const rookMoves = pseudoMoves(b, 9, 0)
ok(has(rookMoves, 8, 0), '红车可上走 (8,0)')
ok(has(rookMoves, 7, 0), '红车可走到 (7,0)')
ok(!has(rookMoves, 6, 0), '红车被己方兵挡住，不能到 (6,0)')

// 3. 马腿别马脚
const horseMoves = pseudoMoves(b, 9, 1)
ok(has(horseMoves, 7, 0) && has(horseMoves, 7, 2), '红马可跳 (7,0)/(7,2)')
ok(!has(horseMoves, 8, 3) && !has(horseMoves, 8, -1), '被象/边界别马腿，不能跳 (8,3)')

// 4. 炮无架不能隔打
const cannonMoves = pseudoMoves(b, 2, 1)
ok(has(cannonMoves, 6, 1), '黑炮可走空格 (6,1)')
ok(!has(cannonMoves, 7, 1), '黑炮前无炮架，不能隔山打到 (7,1)')

// 5. 帅/将 限九宫
const kingMoves = pseudoMoves(b, 9, 4)
ok(has(kingMoves, 8, 4), '红帅可上 (8,4)')
ok(!has(kingMoves, 7, 4), '红帅一步之内不能直接到 (7,4)')

// 6. 将军检测：红车直指黑将（清空中间遮挡）
const bb = createBoard()
for (let r = 1; r <= 4; r++) bb[r][4] = null
bb[4][4] = { side: 'r', type: 'R' }
ok(inCheck(bb, 'b') === true, '红车直指黑将，黑方被将军')
ok(legalMovesFor(bb, 'b').length > 0, '黑将可躲开，未被将死')

// 7. 将帅对脸（飞将）检测：清空两将之间
const bc = createBoard()
for (let r = 1; r <= 8; r++) bc[r][4] = null
ok(generalsFacing(bc) === true, '将帅同列无遮挡 = 对脸（非法局面）')

// 8. 静态评估：初始对称局面为 0
ok(evaluate(createBoard(), 'r') === 0, '初始局面红黑均势，评估为 0')

// 9. 人机 AI 返回一步合法走法
const bm = bestMove(createBoard(), 'r', 2)
const reds = legalMovesFor(createBoard(), 'r')
const legal = bm && reds.some((m) => m.from.r === bm.from.r && m.from.c === bm.from.c && m.to.r === bm.to.r && m.to.c === bm.to.c)
ok(!!legal, 'AI 返回的走法是红方合法走法')

console.log(failed === 0 ? '\n全部通过 ✅' : `\n${failed} 项失败 ❌`)
process.exit(failed === 0 ? 0 : 1)