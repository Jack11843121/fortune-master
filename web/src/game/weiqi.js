// 围棋（9 路）纯规则引擎，无 Vue 依赖，供组件与单测共用。
// 棋盘 N×N 二维数组：0 空、1 黑、2 白。

export const N = 9
const DIRS = [[1, 0], [-1, 0], [0, 1], [0, -1]]
const inBoard = (r, c) => r >= 0 && r < N && c >= 0 && c < N

export function createBoard() {
  return Array.from({ length: N }, () => new Array(N).fill(0))
}

export const clone = (b) => b.map((r) => [...r])

// 连通同色棋串（块）
export function groupAt(b, r, c) {
  const color = b[r][c]
  if (!color) return []
  const seen = new Set([r + ',' + c])
  const stack = [[r, c]]
  const cells = []
  while (stack.length) {
    const [rr, cc] = stack.pop()
    cells.push([rr, cc])
    for (const [dr, dc] of DIRS) {
      const nr = rr + dr, nc = cc + dc
      const k = nr + ',' + nc
      if (inBoard(nr, nc) && !seen.has(k) && b[nr][nc] === color) {
        seen.add(k)
        stack.push([nr, nc])
      }
    }
  }
  return cells
}

// 一块棋的气（空邻点集合）
export function liberties(b, group) {
  const libs = new Set()
  for (const [r, c] of group) {
    for (const [dr, dc] of DIRS) {
      const nr = r + dr, nc = c + dc
      if (inBoard(nr, nc) && b[nr][nc] === 0) libs.add(nr + ',' + nc)
    }
  }
  return libs
}

export function boardString(b) {
  return b.map((r) => r.join('')).join('/')
}

// 落子：成功返回 { board, captured }，非法（占位/自杀/打劫）返回 null。
// koStr：若落子后局面等于该串则判为打劫违规（通常传“上一手之前的局面串”）。
export function tryPlace(b, color, r, c, koStr) {
  if (!inBoard(r, c) || b[r][c] !== 0) return null
  const nb = clone(b)
  nb[r][c] = color
  let captured = 0
  const enemy = color === 1 ? 2 : 1
  for (const [dr, dc] of DIRS) {
    const nr = r + dr, nc = c + dc
    if (inBoard(nr, nc) && nb[nr][nc] === enemy) {
      const g = groupAt(nb, nr, nc)
      if (liberties(nb, g).size === 0) {
        for (const [gr, gc] of g) nb[gr][gc] = 0
        captured += g.length
      }
    }
  }
  const own = groupAt(nb, r, c)
  if (liberties(nb, own).size === 0) return null // 自杀
  if (koStr && boardString(nb) === koStr) return null // 打劫
  return { board: nb, captured }
}

// 所有合法落点（供 AI）
export function legalMoves(b, color, koStr) {
  const res = []
  for (let r = 0; r < N; r++) {
    for (let c = 0; c < N; c++) {
      if (b[r][c] !== 0) continue
      const t = tryPlace(b, color, r, c, koStr)
      if (t) res.push({ r, c, captured: t.captured })
    }
  }
  return res
}

// 简易 AI：优先提子，其次靠己方/中心，带随机扰动
export function aiMove(b, color, koStr) {
  const moves = legalMoves(b, color, koStr)
  if (moves.length === 0) return null
  let best = null
  let bestScore = -Infinity
  const enemy = color === 1 ? 2 : 1
  for (const m of moves) {
    let score
    if (m.captured > 0) {
      score = 1000 + m.captured * 20
    } else {
      let own = 0, foe = 0
      for (const [dr, dc] of DIRS) {
        const nr = m.r + dr, nc = m.c + dc
        if (!inBoard(nr, nc)) continue
        if (b[nr][nc] === color) own++
        else if (b[nr][nc] === enemy) foe++
      }
      const dist = Math.abs(m.r - 4) + Math.abs(m.c - 4)
      score = own * 9 + foe * 6 + (10 - dist) + Math.random() * 4
    }
    if (score > bestScore) {
      bestScore = score
      best = m
    }
  }
  return best
}

// 面积计分（中国规则）：棋盘上的子 + 单方围起的空地，返回 { black, white }
export function countScore(b) {
  const visited = new Set()
  let black = 0, white = 0
  for (let r = 0; r < N; r++) {
    for (let c = 0; c < N; c++) {
      if (b[r][c] === 1) { black++; continue }
      if (b[r][c] === 2) { white++; continue }
      const k = r + ',' + c
      if (visited.has(k)) continue
      const stack = [[r, c]]
      const region = []
      const touched = { 1: false, 2: false }
      visited.add(k)
      while (stack.length) {
        const [rr, cc] = stack.pop()
        region.push([rr, cc])
        for (const [dr, dc] of DIRS) {
          const nr = rr + dr, nc = cc + dc
          if (!inBoard(nr, nc)) continue
          if (b[nr][nc] === 0) {
            const kk = nr + ',' + nc
            if (!visited.has(kk)) { visited.add(kk); stack.push([nr, nc]) }
          } else if (b[nr][nc] === 1) touched[1] = true
          else touched[2] = true
        }
      }
      if (touched[1] && !touched[2]) black += region.length
      else if (touched[2] && !touched[1]) white += region.length
      // 双色皆触或皆无：中性地，不计
    }
  }
  return { black, white }
}