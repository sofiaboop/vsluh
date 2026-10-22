/** Запись ответа с микрофона. Формат выбираем по тому, что умеет браузер. */
export interface Recording {
  blob: Blob
  /** Контейнер записи: его ждёт бэкенд, из MIME-типа он не выводится. */
  format: string
}

/** Кандидаты в порядке предпочтения: opus компактнее, mp4 нужен Safari. */
const CANDIDATES: Array<{ mime: string; format: string }> = [
  { mime: 'audio/webm;codecs=opus', format: 'webm' },
  { mime: 'audio/webm', format: 'webm' },
  { mime: 'audio/mp4', format: 'mp4' },
  { mime: 'audio/ogg;codecs=opus', format: 'ogg' },
]

export function isRecordingSupported(): boolean {
  return (
    typeof window !== 'undefined' &&
    typeof MediaRecorder !== 'undefined' &&
    Boolean(navigator.mediaDevices?.getUserMedia)
  )
}

function pickFormat(): { mime: string; format: string } | null {
  for (const candidate of CANDIDATES) {
    if (MediaRecorder.isTypeSupported(candidate.mime)) return candidate
  }
  return null
}

/** Открывает микрофон. Разрешение спрашиваем заранее, а не в момент старта речи. */
export async function openMicrophone(): Promise<MediaStream> {
  return navigator.mediaDevices.getUserMedia({
    audio: { echoCancellation: true, noiseSuppression: true },
  })
}

/**
 * Пишет дорожку из уже открытого потока. Поток не закрываем: им распоряжается
 * вызывающая сторона, иначе повторная практика снова спросит разрешение.
 */
export function recordStream(stream: MediaStream): {
  stop: () => Promise<Recording>
} {
  const picked = pickFormat()
  if (!picked) throw new Error('recording is not supported')

  const chunks: BlobPart[] = []
  const recorder = new MediaRecorder(stream, { mimeType: picked.mime, audioBitsPerSecond: 32000 })
  recorder.ondataavailable = (event) => {
    if (event.data.size > 0) chunks.push(event.data)
  }
  recorder.start(1000)

  return {
    stop: () =>
      new Promise<Recording>((resolve, reject) => {
        if (recorder.state === 'inactive') {
          resolve({ blob: new Blob(chunks, { type: picked.mime }), format: picked.format })
          return
        }
        recorder.onerror = () => reject(new Error('recording failed'))
        recorder.onstop = () =>
          resolve({ blob: new Blob(chunks, { type: picked.mime }), format: picked.format })
        recorder.stop()
      }),
  }
}

/** Переводит запись в base64 без префикса data-URI. */
export function toBase64(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onerror = () => reject(new Error('cannot read recording'))
    reader.onload = () => {
      const result = String(reader.result)
      resolve(result.slice(result.indexOf('base64,') + 'base64,'.length))
    }
    reader.readAsDataURL(blob)
  })
}
