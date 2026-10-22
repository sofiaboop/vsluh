import {
  isRecordingSupported,
  openMicrophone,
  recordStream,
  type Recording,
} from '~/shared/lib/audioRecorder'

/** idle — микрофон не открывали, denied — доступ не дали, error — не смогли записать. */
export type RecorderState = 'idle' | 'ready' | 'recording' | 'denied' | 'unsupported' | 'error'

// Поток и запись живут вне компонента: экран ответа размонтируется раньше,
// чем запись успевает закрыться и уехать на разбор.
let stream: MediaStream | null = null
let session: { stop: () => Promise<Recording> } | null = null
let recording: Recording | null = null
// Остановка записи асинхронная: экран разбора открывается раньше, чем
// MediaRecorder отдаёт последний кусок, и должен дождаться именно его.
let stopping: Promise<void> | null = null

export function useRecorder() {
  const state = useState<RecorderState>('recorder-state', () => 'idle')

  /** Спрашивает доступ к микрофону заранее — до начала отсчёта времени ответа. */
  async function prepare(): Promise<boolean> {
    if (!isRecordingSupported()) {
      state.value = 'unsupported'
      return false
    }
    if (stream && stream.getAudioTracks().some((track) => track.readyState === 'live')) {
      state.value = 'ready'
      return true
    }

    try {
      stream = await openMicrophone()
      state.value = 'ready'
      return true
    } catch {
      // Отказ в доступе — обычный сценарий: практика продолжается без разбора.
      state.value = 'denied'
      return false
    }
  }

  function start() {
    if (!stream || state.value === 'denied' || state.value === 'unsupported') return
    try {
      recording = null
      stopping = null
      session = recordStream(stream)
      state.value = 'recording'
    } catch {
      state.value = 'error'
    }
  }

  /** Останавливает запись и запоминает её до отправки на разбор. */
  async function stop(): Promise<Recording | null> {
    if (!session) return null
    const current = session
    session = null
    stopping = current
      .stop()
      .then((result) => {
        recording = result
        state.value = 'ready'
      })
      .catch(() => {
        recording = null
        state.value = 'error'
      })
    await stopping
    return recording
  }

  /** Забирает запись: повторная отправка той же дорожки не нужна. */
  async function take(): Promise<Recording | null> {
    if (stopping) await stopping
    const result = recording
    recording = null
    return result
  }

  /** Закрывает микрофон — иначе индикатор записи висит в браузере после практики. */
  function release() {
    session = null
    stopping = null
    recording = null
    stream?.getTracks().forEach((track) => track.stop())
    stream = null
    if (state.value !== 'denied' && state.value !== 'unsupported') state.value = 'idle'
  }

  return { state, prepare, start, stop, take, release }
}
