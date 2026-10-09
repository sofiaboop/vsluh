// Звук барабана синтезируется на лету — параметры сняты с эталонного макета.

let ctx: AudioContext | null = null

function audio(): AudioContext | null {
  if (typeof window === 'undefined') return null

  if (!ctx) {
    const Ctor = window.AudioContext || (window as any).webkitAudioContext
    if (!Ctor) return null
    ctx = new Ctor()
  }
  // автоплей-политика: контекст просыпается только после жеста пользователя
  if (ctx.state === 'suspended') void ctx.resume()
  return ctx
}

/** Щелчок на каждую проехавшую строку. */
export function playTick(): void {
  const c = audio()
  if (!c) return

  const t = c.currentTime + 0.005
  const osc = c.createOscillator()
  const gain = c.createGain()

  osc.type = 'triangle'
  osc.frequency.setValueAtTime(1120, t)
  osc.frequency.exponentialRampToValueAtTime(760, t + 0.025)

  gain.gain.setValueAtTime(1e-4, t)
  gain.gain.exponentialRampToValueAtTime(0.028, t + 0.003)
  gain.gain.exponentialRampToValueAtTime(1e-4, t + 0.034)

  osc.connect(gain)
  gain.connect(c.destination)
  osc.start(t)
  osc.stop(t + 0.04)
}

const CHORD = [
  { frequency: 740, volume: 0.045, duration: 0.78 },
  { frequency: 1485, volume: 0.019, duration: 0.52 },
  { frequency: 2240, volume: 0.008, duration: 0.34 },
  { frequency: 2965, volume: 0.0035, duration: 0.22 },
]

/** Аккорд в момент остановки барабана. */
export function playStop(): void {
  const c = audio()
  if (!c) return

  const t = c.currentTime + 0.005
  for (const { frequency, volume, duration } of CHORD) {
    const osc = c.createOscillator()
    const gain = c.createGain()

    osc.type = 'sine'
    osc.frequency.setValueAtTime(frequency, t)

    gain.gain.setValueAtTime(1e-4, t)
    gain.gain.exponentialRampToValueAtTime(volume, t + 0.008)
    gain.gain.exponentialRampToValueAtTime(1e-4, t + duration)

    osc.connect(gain)
    gain.connect(c.destination)
    osc.start(t)
    osc.stop(t + duration + 0.03)
  }
}
