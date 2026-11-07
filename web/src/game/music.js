// 五行音乐引擎（治愈连奏版）：Web Audio 实时合成五声音阶（宫商角徵羽），
// 每个五行对应一阕「连贯」的五音旋律（音符长音相互重叠、连奏不断点），
// 叠加低频和弦底垫 + 呼吸式音量起伏 + 回声，营造静心安神、涤荡心绪的治愈感。
// 五音对应五行、入五脏（《黄帝内经》五音疗疾理念）。

export const ELEMENTS = [
  { key: 'shui', wuxing: '水', wuyin: '羽', organ: '肾', mood: '柔润沉静', color: '#4a86d9', base: 220.00, tempo: 0.88, phrase: [4, 3, 1, 0, 1, 3, 4, 3, 1, 0, 1, 3, 4, 3, 1, 0], desc: '羽音属水、入肾，如流水潺潺，安神定志、滋养根本。' },
  { key: 'tu', wuxing: '土', wuyin: '宫', organ: '脾', mood: '安稳厚重', color: '#a8794f', base: 261.63, tempo: 0.82, phrase: [0, 1, 3, 4, 3, 1, 0, 1, 3, 1, 0, 3, 4, 3, 1, 0], desc: '宫音属土、入脾，如厚土载物，静思养神、心生安稳。' },
  { key: 'jin', wuxing: '金', wuyin: '商', organ: '肺', mood: '清亮收敛', color: '#c9a24b', base: 293.66, tempo: 0.80, phrase: [1, 3, 4, 5, 4, 3, 1, 3, 4, 3, 1, 4, 5, 4, 3, 1], desc: '商音属金、入肺，清冽澄澈，收敛心绪、吐故纳新。' },
  { key: 'mu', wuxing: '木', wuyin: '角', organ: '肝', mood: '舒发生机', color: '#4c9a5a', base: 329.63, tempo: 0.76, phrase: [2, 4, 5, 7, 5, 4, 2, 4, 5, 4, 2, 5, 7, 5, 4, 2], desc: '角音属木、入肝，如春木生发，舒展情志、疏解郁结。' },
  { key: 'huo', wuxing: '火', wuyin: '徵', organ: '心', mood: '明亮欢快', color: '#d94a3a', base: 392.00, tempo: 0.72, phrase: [3, 5, 7, 8, 7, 5, 3, 5, 7, 5, 3, 7, 8, 7, 5, 3], desc: '徵音属火、入心，明快热烈，振奋精神、愉悦心情。' },
]

let audio = null // { ctx, master, timer }

export function isPlaying() {
  return !!audio
}

export function play(el) {
  stop()
  const AC = window.AudioContext || window.webkitAudioContext
  if (!AC) return
  const ctx = new AC()
  if (ctx.state === 'suspended') ctx.resume()

  // 主输出
  const master = ctx.createGain()
  master.gain.value = 0.85
  master.connect(ctx.destination)

  // 「呼吸」：慢速 LFO 调制音量，整体起伏，助放松
  const breath = ctx.createGain()
  breath.gain.value = 0.75
  breath.connect(master)
  const lfo = ctx.createOscillator()
  lfo.frequency.value = 0.11
  const lfoGain = ctx.createGain()
  lfoGain.gain.value = 0.12
  lfo.connect(lfoGain)
  lfoGain.connect(breath.gain)
  lfo.start()

  // 空间回声
  const delay = ctx.createDelay(1.5)
  delay.delayTime.value = 0.55
  const fb = ctx.createGain()
  fb.gain.value = 0.30
  delay.connect(fb)
  fb.connect(delay)
  delay.connect(breath)

  // 低频和弦底垫：主音的下八度 + 五度 + 主音
  const pad = ctx.createGain()
  pad.gain.value = 0.15
  pad.connect(breath)
  for (const ratio of [0.5, 0.75, 1]) {
    const o = ctx.createOscillator()
    o.type = 'sine'
    o.frequency.value = el.base * ratio
    const g = ctx.createGain()
    g.gain.value = 0.7
    o.connect(g)
    g.connect(pad)
    o.start()
  }

  // 两个八度的五声音阶
  const semis = [0, 2, 4, 7, 9, 12, 14, 16, 19, 21]
  const scale = semis.map((s) => el.base * Math.pow(2, s / 12))
  const phrase = el.phrase
  const onset = 0.68 * (el.tempo || 0.8) // 每次起音间隔（秒）

  // 连奏单音：attack 柔和，音长 > 间隔，音符前后重叠 → 连贯不断
  function legato(freq, when, dur) {
    const env = ctx.createGain()
    env.gain.setValueAtTime(0.0001, when)
    env.gain.linearRampToValueAtTime(0.30, when + 0.14)
    env.gain.setValueAtTime(0.30, when + dur - 0.28)
    env.gain.linearRampToValueAtTime(0.0001, when + dur)

    const osc1 = ctx.createOscillator()
    osc1.type = 'sine'
    osc1.frequency.value = freq
    const g1 = ctx.createGain()
    g1.gain.value = 1
    osc1.connect(g1)
    g1.connect(env)

    const osc2 = ctx.createOscillator()
    osc2.type = 'triangle'
    osc2.frequency.value = freq * 2
    const g2 = ctx.createGain()
    g2.gain.value = 0.18
    osc2.connect(g2)
    g2.connect(env)

    env.connect(breath)
    env.connect(delay)

    osc1.start(when)
    osc1.stop(when + dur)
    osc2.start(when)
    osc2.stop(when + dur)
  }

  let idx = 0
  function step() {
    if (!audio) return
    const t = ctx.currentTime + 0.06
    const ix = phrase[idx % phrase.length]
    idx++
    legato(scale[ix], t, onset * 1.9) // 长音重叠，连贯连奏
    audio.timer = setTimeout(step, onset * 1000)
  }

  audio = { ctx, master, timer: null }
  step()
}

export function stop() {
  if (!audio) return
  const a = audio
  audio = null
  clearTimeout(a.timer)
  try {
    a.master.gain.setTargetAtTime(0, a.ctx.currentTime, 0.06)
  } catch (e) {}
  setTimeout(() => {
    try { a.ctx.close() } catch (e) {}
  }, 280)
}