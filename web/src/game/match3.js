// 五行消消乐核心逻辑（无 Vue 依赖），供组件与单测共用。
// 棋盘为 ROWS×COLS 二维数组，取值 0..TYPES-1（对应五行），-1 表示待下落。

export const ROWS = 6
export const COLS = 6
export const TYPES = 5

export const rand = () => Math.floor(Math.random() * TYPES)
export const key = (r, c) => r + ',' + c

// 生成一个「无初始三连」的随机棋盘
export function fillNoMatch() {
  const g = Array.from({ length: ROWS }, () => new Array(COLS).fill(0))
  for (let r = 0; r < ROWS; r++) {
    for (let c = 0; c < COLS; c++) {
      let v
      do {
        v = rand()
      } while (
        (c >= 2 && g[r][c - 1] === v && g[r][c - 2] === v) ||
        (r >= 2 && g[r - 1][c] === v && g[r - 2][c] === v)
      )
      g[r][c] = v
    }
  }
  return g
}

// 找出所有可消除的三连（返回 key 集合）
export function findMatches(g) {
  const s = new Set()
  for (let r = 0; r < ROWS; r++) {
    for (let c = 0; c < COLS; c++) {
      const v = g[r][c]
      if (v < 0) continue
      if (c <= COLS - 3 && g[r][c + 1] === v && g[r][c + 2] === v) {
        let k = c
        while (k < COLS && g[r][k] === v) s.add(key(r, k++))
      }
      if (r <= ROWS - 3 && g[r + 1][c] === v && g[r + 2][c] === v) {
        let k = r
        while (k < ROWS && g[k][c] === v) s.add(key(k++, c))
      }
    }
  }
  return s
}

// 下落并补新（原地修改）
export function dropAndRefill(g) {
  for (let c = 0; c < COLS; c++) {
    const arr = []
    for (let r = ROWS - 1; r >= 0; r--) if (g[r][c] >= 0) arr.push(g[r][c])
    for (let r = ROWS - 1; r >= 0; r--) {
      const i = ROWS - 1 - r
      g[r][c] = i < arr.length ? arr[i] : rand()
    }
  }
}

// 是否至少存在一步可消除的交换
export function hasAnyMove(g) {
  for (let r = 0; r < ROWS; r++) {
    for (let c = 0; c < COLS; c++) {
      if (c + 1 < COLS) {
        ;[g[r][c], g[r][c + 1]] = [g[r][c + 1], g[r][c]]
        const ok = findMatches(g).size > 0
        ;[g[r][c], g[r][c + 1]] = [g[r][c + 1], g[r][c]]
        if (ok) return true
      }
      if (r + 1 < ROWS) {
        ;[g[r][c], g[r + 1][c]] = [g[r + 1][c], g[r][c]]
        const ok = findMatches(g).size > 0
        ;[g[r][c], g[r + 1][c]] = [g[r + 1][c], g[r][c]]
        if (ok) return true
      }
    }
  }
  return false
}

// 找一步可消除的交换（供“提示”高亮），返回 { a:{r,c}, b:{r,c} } 或 null
export function findHint(g) {
  for (let r = 0; r < ROWS; r++) {
    for (let c = 0; c < COLS; c++) {
      if (c + 1 < COLS) {
        ;[g[r][c], g[r][c + 1]] = [g[r][c + 1], g[r][c]]
        const ok = findMatches(g).size > 0
        ;[g[r][c], g[r][c + 1]] = [g[r][c + 1], g[r][c]]
        if (ok) return { a: { r, c }, b: { r, c: c + 1 } }
      }
      if (r + 1 < ROWS) {
        ;[g[r][c], g[r + 1][c]] = [g[r + 1][c], g[r][c]]
        const ok = findMatches(g).size > 0
        ;[g[r][c], g[r + 1][c]] = [g[r + 1][c], g[r][c]]
        if (ok) return { a: { r, c }, b: { r: r + 1, c } }
      }
    }
  }
  return null
}

// 生成一个一定有解的棋盘
export function makeBoard() {
  let g
  let n = 0
  do {
    g = fillNoMatch()
    n++
  } while (!hasAnyMove(g) && n < 200)
  return g
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms))

// 连锁消除（带回调，供组件做动画）；返回 { gained, chains }
export async function resolveAsync(g, onClear) {
  let gained = 0
  let chains = 0
  for (;;) {
    const m = findMatches(g)
    if (m.size === 0) break
    chains++
    gained += m.size * 10 * chains
    const cleared = [...m].map((k) => k.split(',').map(Number))
    if (onClear) onClear(cleared)
    await sleep(240)
    for (const [r, c] of cleared) g[r][c] = -1
    dropAndRefill(g)
    await sleep(80)
  }
  return { gained, chains }
}