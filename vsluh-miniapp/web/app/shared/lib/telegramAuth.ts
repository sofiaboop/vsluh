// initData хранится вне Pinia: он нужен раньше, чем инициализируются сторы.
let initData = ''

export function setInitData(raw: string): void {
  initData = raw
}

export function getInitData(): string {
  return initData
}

/**
 * Поверхностная проверка формы строки. Подлинность определяет только бэкенд,
 * сверяя HMAC-подпись с токеном бота.
 */
export function looksLikeInitData(raw: string | undefined | null): raw is string {
  return Boolean(raw) && raw!.includes('user=') && raw!.includes('hash=')
}

/** API Gateway портит сырую строку в заголовке, поэтому кодируем в Base64. */
export function encodeInitData(raw: string): string {
  const bytes = new TextEncoder().encode(raw)
  let binary = ''
  bytes.forEach((b) => {
    binary += String.fromCharCode(b)
  })
  return btoa(binary)
}
