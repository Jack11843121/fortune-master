// 中国象棋纯规则引擎（无 Vue 依赖），供组件与单测共用。
// 棋盘为 10x9 二维数组：b[r][c] = null 或 { side: 'r'|'b', type: 'K'|'A'|'E'|'H'|'R'|'C'|'P' }
// 行 0 为黑方（上方），行 9 为红方（下方），红先。

export const NAME = {
  K: { r: '帅', b: '将' }, A: { r: '仕', b: '士' }, E: { r: '相', b: '象' },
  H: { r: '马', b: '马' }, R: { r: '车', b: '车' }, C: { r: '炮', b: '砲' }, P: { r: '兵', b: '卒' },
}

export const inBoard = (r, c) => r >= 0 && r < 10 && c >= 0 && c < 9
export const inPalace = (side, r, c) => c >= 3 && c <= 5 && (side === 'b' ? r >= 0 && r <= 2 : r >= 7 && r <= 9)
export const ownHalf = (side, r) => (side === 'b' ? r <= 4 : r >= 5)

export function createBoard() {
  const b = Array.from({ length: 10 }, () => new Array(9).fill(null))
  const back = ['R', 'H', 'E', 'A', 'K', 'A', 'E', 'H', 'R']
  for (let c = 0; c < 9; c++) {
    b[0][c] = { side: 'b', type: back[c] }
    b[9][c] = { side: 'r', type: back[c] }
  }
  b[2][1] = { side: 'b', type: 'C' }
  b[2][7] = { side: 'b', type: 'C' }
  b[7][1] = { side: 'r', type: 'C' }
  b[7][7] = { side: 'r', type: 'C' }
  for (let c = 0; c < 9; c += 2) {
    b[3][c] = { side: 'b', type: 'P' }
    b[6][c] = { side: 'r', type: 'P' }
  }
  return b
}

export const clone = (b) => b.map((row) => row.map((p) => (p ? { ...p } : null)))

// 返回某棋子当前可到达（移动或吃子）的坐标对数组 [[r,c], ...]
export function pseudoMoves(b, r, c) {
  const p = b[r][c]
  if (!p) return []
  const side = p.side
  const out = []
  const add = (nr, nc) => {
    if (inBoard(nr, nc) && !(b[nr][nc] && b[nr][nc].side === side)) out.push([nr, nc])
  }
  const dir4 = [[1, 0], [-1, 0], [0, 1], [0, -1]]
  switch (p.type) {
    case 'K':
      for (const [dr, dc] of dir4) if (inPalace(side, r + dr, c + dc)) add(r + dr, c + dc)
      break
    case 'A':
      for (const [dr, dc] of [[1, 1], [1, -1], [-1, 1], [-1, -1]])
        if (inPalace(side, r + dr, c + dc)) add(r + dr, c + dc)
      break
    case 'E':
      for (const [dr, dc] of [[2, 2], [2, -2], [-2, 2], [-2, -2]]) {
        const nr = r + dr, nc = c + dc
        if (inBoard(nr, nc) && ownHalf(side, nr) && !b[r + dr / 2][c + dc / 2]) add(nr, nc)
      }
      break
    case 'H':
      for (const [dr, dc] of [[-2, -1], [-2, 1], [2, -1], [2, 1], [-1, -2], [-1, 2], [1, -2], [1, 2]]) {
        const lr = Math.abs(dr) === 2 ? r + dr / 2 : r
        const lc = Math.abs(dc) === 2 ? c + dc / 2 : c
        if (inBoard(lr, lc) && !b[lr][lc]) add(r + dr, c + dc)
      }
      break
    case 'R':
      for (const [dr, dc] of dir4) {
        let nr = r + dr, nc = c + dc
        while (inBoard(nr, nc)) {
          if (!b[nr][nc]) { out.push([nr, nc]); nr += dr; nc += dc; continue }
          if (b[nr][nc].side !== side) out.push([nr, nc])
          break
        }
      }
      break
    case 'C':
      for (const [dr, dc] of dir4) {
        let nr = r + dr, nc = c + dc, screen = false
        while (inBoard(nr, nc)) {
          const cell = b[nr][nc]
          if (!screen) {
            if (!cell) out.push([nr, nc])
            else screen = true
          } else if (cell) {
            if (cell.side !== side) out.push([nr, nc])
            break
          }
          nr += dr; nc += dc
        }
      }
      break
    case 'P': {
      const fwd = side === 'r' ? -1 : 1
      const crossed = side === 'r' ? r <= 4 : r >= 5
      add(r + fwd, c)
      if (crossed) { add(r, c - 1); add(r, c + 1) }
      break
    }
  }
  return out
}

export function findKing(b, side) {
  for (let r = 0; r < 10; r++)
    for (let c = 0; c < 9; c++)
      if (b[r][c] && b[r][c].side === side && b[r][c].type === 'K') return { r, c }
  return null
}

export function isSquareAttacked(b, r, c, bySide) {
  for (let rr = 0; rr < 10; rr++)
    for (let cc = 0; cc < 9; cc++) {
      const p = b[rr][cc]
      if (p && p.side === bySide && pseudoMoves(b, rr, cc).some((m) => m[0] === r && m[1] === c)) return true
    }
  return false
}

export function generalsFacing(b) {
  const rk = findKing(b, 'r')
  const bk = findKing(b, 'b')
  if (!rk || !bk || rk.c !== bk.c) return false
  for (let r = Math.min(rk.r, bk.r) + 1; r < Math.max(rk.r, bk.r); r++) if (b[r][rk.c]) return false
  return true
}

export function inCheck(b, side) {
  const k = findKing(b, side)
  if (!k) return false
  return isSquareAttacked(b, k.r, k.c, side === 'r' ? 'b' : 'r')
}

export function isLegalAfter(b, side) {
  return !generalsFacing(b) && !inCheck(b, side)
}

// 返回一方所有合法走法 [{ from:{r,c}, to:{r,c} }, ...]
export function legalMovesFor(b, side) {
  const res = []
  for (let r = 0; r < 10; r++)
    for (let c = 0; c < 9; c++) {
      const p = b[r][c]
      if (!p || p.side !== side) continue
      for (const [nr, nc] of pseudoMoves(b, r, c)) {
        const nb = clone(b)
        nb[nr][nc] = nb[r][c]
        nb[r][c] = null
        if (isLegalAfter(nb, side)) res.push({ from: { r, c }, to: { r: nr, c: nc } })
      }
    }
  return res
}

// —— 简易人机 AI：极小极大 + α-β 剪枝，默认搜索 2 层 ——
const VALUE = { K: 10000, R: 900, C: 450, H: 400, E: 200, A: 200, P: 100 }

// 静态评估：从 side 视角（正 = 对 side 有利）
export function evaluate(b, side) {
  let score = 0
  for (let r = 0; r < 10; r++) {
    for (let c = 0; c < 9; c++) {
      const p = b[r][c]
      if (!p) continue
      let v = VALUE[p.type]
      if (p.type === 'P') {
        const crossed = p.side === 'r' ? r <= 4 : r >= 5
        if (crossed) v += 60
      }
      score += p.side === side ? v : -v
    }
  }
  return score
}

function applyMove(b, from, to) {
  const nb = clone(b)
  nb[to.r][to.c] = nb[from.r][from.c]
  nb[from.r][from.c] = null
  return nb
}

function minimax(b, depth, alpha, beta, maximizing, aiSide) {
  const side = maximizing ? aiSide : (aiSide === 'r' ? 'b' : 'r')
  const moves = legalMovesFor(b, side)
  if (moves.length === 0) return maximizing ? -99999 : 99999
  if (depth === 0) return evaluate(b, aiSide)
  if (maximizing) {
    let best = -Infinity
    for (const m of moves) {
      const nb = applyMove(b, m.from, m.to)
      best = Math.max(best, minimax(nb, depth - 1, alpha, beta, false, aiSide))
      alpha = Math.max(alpha, best)
      if (beta <= alpha) break
    }
    return best
  }
  let best = Infinity
  for (const m of moves) {
    const nb = applyMove(b, m.from, m.to)
    best = Math.min(best, minimax(nb, depth - 1, alpha, beta, true, aiSide))
    beta = Math.min(beta, best)
    if (beta <= alpha) break
  }
  return best
}

// 返回 aiSide 的最优走法 {from,to}；无棋可走返回 null
export function bestMove(b, aiSide, depth = 2) {
  const moves = legalMovesFor(b, aiSide)
  if (moves.length === 0) return null
  let best = null
  let bestScore = -Infinity
  for (const m of moves) {
    const nb = applyMove(b, m.from, m.to)
    const sc = minimax(nb, depth - 1, -Infinity, Infinity, false, aiSide)
    if (sc > bestScore) {
      bestScore = sc
      best = m
    }
  }
  return best
}