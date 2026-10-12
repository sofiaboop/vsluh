/**
 * Ужимает шрифт элемента, пока текст не поместится в коробку.
 * Потолок шрифта — `font-size` из CSS, поэтому макетные clamp() остаются источником истины.
 *
 * Высота коробки — `max-height` из CSS. С `fillScreen` она дополнительно ограничивается
 * местом, которое реально осталось на экране после шапки, соседей и отступов.
 */
export function useFitText(
  target: Ref<HTMLElement | null>,
  source: () => unknown,
  options: { min?: number; fillScreen?: boolean } = {},
) {
  const minSize = options.min ?? 14
  let range: Range | null = null

  /** Высота строк текста: scrollHeight врёт при плотном line-height. */
  function textHeight(el: HTMLElement) {
    range ??= document.createRange()
    range.selectNodeContents(el)
    return range.getBoundingClientRect().height
  }

  function outerHeight(el: Element) {
    const style = getComputedStyle(el)
    if (style.position === 'absolute' || style.position === 'fixed' || style.display === 'none') {
      return 0
    }
    return (
      el.getBoundingClientRect().height +
      Number.parseFloat(style.marginTop) +
      Number.parseFloat(style.marginBottom)
    )
  }

  /** Сколько вертикального места осталось элементу внутри экрана практики. */
  function freeHeight(el: HTMLElement) {
    const section = el.parentElement
    const screen = section?.parentElement
    const shell = screen?.parentElement
    if (!section || !screen || !shell) return Number.POSITIVE_INFINITY

    const screenStyle = getComputedStyle(screen)
    let free =
      window.innerHeight -
      Number.parseFloat(screenStyle.paddingTop) -
      Number.parseFloat(screenStyle.paddingBottom)

    for (const sibling of shell.children) {
      if (sibling !== screen) free -= outerHeight(sibling)
    }
    for (const sibling of section.children) {
      if (sibling !== el) free -= outerHeight(sibling)
    }

    const style = getComputedStyle(el)
    return free - Number.parseFloat(style.marginTop) - Number.parseFloat(style.marginBottom)
  }

  function fit() {
    const el = target.value
    if (!el) return

    el.style.removeProperty('font-size')
    const style = getComputedStyle(el)
    const max = Math.round(Number.parseFloat(style.fontSize))
    if (!max) return

    const declared = Number.parseFloat(style.maxHeight)
    const limits = [Number.isFinite(declared) ? declared : Number.POSITIVE_INFINITY]
    if (options.fillScreen) limits.push(freeHeight(el))

    // на совсем низком экране коробка уходит в минус — там честнее дать прокрутку
    const box = Math.max(Math.min(...limits), minSize * 2)
    if (!Number.isFinite(box) || textHeight(el) <= box) return

    let low = minSize
    let high = max
    let best = minSize
    while (low <= high) {
      const mid = Math.floor((low + high) / 2)
      el.style.fontSize = `${mid}px`
      if (textHeight(el) <= box) {
        best = mid
        low = mid + 1
      } else {
        high = mid - 1
      }
    }
    el.style.fontSize = `${best}px`
  }

  onMounted(() => {
    fit()
    // Manrope подгружается асинхронно и меняет метрики уже после первой отрисовки
    document.fonts?.ready.then(fit)
  })

  watch(source, () => nextTick(fit), { flush: 'post' })
  useEventListener(window, 'resize', fit)
  useEventListener(window, 'orientationchange', fit)
}
