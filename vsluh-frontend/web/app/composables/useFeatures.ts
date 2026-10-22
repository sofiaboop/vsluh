import { FEATURES } from '~/shared/constants/features'

/**
 * Фичефлаги с учётом контура: часть кусков VER 3 обкатывается на dev
 * раньше, чем включается на проде (`NUXT_PUBLIC_FEATURE_AUTH`).
 */
export function useFeatures() {
  const { featureAuth } = useRuntimeConfig().public
  return computed(() => ({ ...FEATURES, auth: FEATURES.auth || featureAuth }))
}
