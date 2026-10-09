import { MODES, MOODS, TOPICS, WHEEL, type SpeakingMode, type TopicId } from '~/shared/constants/practice'
import { QUESTIONS } from '~/shared/data/questions'
import { playStop } from '~/shared/lib/sound'

export type Stage = 'topics' | 'wheel' | 'prepare' | 'speak' | 'finished'

export interface WheelRow {
  topic: string
  question: string
}

function pick<T>(list: T[]): T {
  return list[Math.floor(Math.random() * list.length)]!
}

// Таймеры живут на уровне модуля: экран-инициатор размонтируется раньше, чем они отработают
let startDelay: ReturnType<typeof setTimeout> | undefined
let spinFallback: ReturnType<typeof setTimeout> | undefined
let tick: ReturnType<typeof setInterval> | undefined

export function usePractice() {
  const stage = useState<Stage>('practice-stage', () => 'topics')
  const mode = useState<SpeakingMode>('practice-mode', () => 'spontaneous')
  const selected = useState<TopicId[]>('practice-topics', () => [])
  const moodEnabled = useState('practice-mood-on', () => false)
  const mood = useState('practice-mood', () => '')

  const rows = useState<WheelRow[]>('practice-rows', () => [])
  const spinning = useState('practice-spinning', () => false)
  const settled = useState('practice-settled', () => false)
  // какая строка подсвечена: пока барабан крутится, она догоняет центр
  const selectedIndex = useState('practice-selected-index', () => WHEEL.startIndex)

  const remaining = useState('practice-remaining', () => 0)
  const paused = useState('practice-paused', () => false)

  const modeOption = computed(() => MODES.find((m) => m.id === mode.value)!)
  const allSelected = computed(() => selected.value.length === TOPICS.length)
  const canContinue = computed(() => selected.value.length > 0)
  const chosen = computed<WheelRow | null>(() => rows.value[WHEEL.selectedIndex] ?? null)
  const totalSec = computed(() =>
    stage.value === 'speak' ? modeOption.value.speakSec : modeOption.value.prepareSec,
  )
  /** Доля оставшегося времени в градусах — ею маскируется подсветка таймера. */
  const timerAngle = computed(() =>
    totalSec.value ? (360 * remaining.value) / totalSec.value : 0,
  )

  function toggleTopic(id: TopicId) {
    selected.value = selected.value.includes(id)
      ? selected.value.filter((t) => t !== id)
      : [...selected.value, id]
  }

  function toggleAll() {
    selected.value = allSelected.value ? [] : TOPICS.map((t) => t.id)
  }

  function buildRows(): WheelRow[] {
    const pool = selected.value.flatMap((id) => {
      const label = TOPICS.find((t) => t.id === id)!.label
      return QUESTIONS[id].map((question) => ({ topic: label, question }))
    })
    return Array.from({ length: WHEEL.rows }, () => pick(pool))
  }

  function spin() {
    if (spinning.value) return
    rows.value = buildRows()
    mood.value = moodEnabled.value ? pick(MOODS) : ''
    selectedIndex.value = WHEEL.startIndex
    settled.value = false
    spinning.value = true

    // в фоновой вкладке анимация не доигрывает и animationend не приходит
    clearTimeout(spinFallback)
    spinFallback = setTimeout(finishSpin, WHEEL.spinMs + 400)
  }

  /** Барабан доехал: фиксируем выбор и даём финальный аккорд. */
  function finishSpin() {
    if (!spinning.value) return
    clearTimeout(spinFallback)
    selectedIndex.value = WHEEL.selectedIndex
    settled.value = true
    spinning.value = false
    playStop()
  }

  /** Подсветка строки, проехавшей центр. */
  function highlight(index: number) {
    selectedIndex.value = Math.min(index, WHEEL.selectedIndex)
  }

  function startWheel() {
    if (!canContinue.value) return
    stage.value = 'wheel'
    rows.value = buildRows()
    selectedIndex.value = WHEEL.startIndex
    settled.value = false
    // короткая пауза, чтобы экран успел появиться до старта барабана
    clearTimeout(startDelay)
    startDelay = setTimeout(spin, 220)
  }

  function stopTimer() {
    clearInterval(tick)
    tick = undefined
  }

  function runTimer(onDone: () => void) {
    stopTimer()
    paused.value = false
    tick = setInterval(() => {
      if (paused.value) return
      remaining.value -= 1
      if (remaining.value <= 0) {
        stopTimer()
        onDone()
      }
    }, 1000)
  }

  function startPrepare() {
    stage.value = 'prepare'
    remaining.value = modeOption.value.prepareSec
    runTimer(startSpeak)
  }

  function startSpeak() {
    stage.value = 'speak'
    remaining.value = modeOption.value.speakSec
    runTimer(finish)
  }

  function finish() {
    stopTimer()
    stage.value = 'finished'
  }

  function togglePause() {
    paused.value = !paused.value
  }

  function backToTopics() {
    stopTimer()
    clearTimeout(startDelay)
    clearTimeout(spinFallback)
    spinning.value = false
    settled.value = false
    stage.value = 'topics'
  }

  /** Новый вопрос по тем же темам. */
  function again() {
    stopTimer()
    startWheel()
  }

  /** Тот же вопрос ещё раз. */
  function repeatTopic() {
    startPrepare()
  }

  return {
    stage,
    mode,
    selected,
    moodEnabled,
    mood,
    rows,
    spinning,
    settled,
    selectedIndex,
    remaining,
    paused,
    modeOption,
    allSelected,
    canContinue,
    chosen,
    timerAngle,
    toggleTopic,
    toggleAll,
    startWheel,
    spin,
    finishSpin,
    highlight,
    startPrepare,
    startSpeak,
    finish,
    togglePause,
    backToTopics,
    again,
    repeatTopic,
  }
}
