export default defineNuxtPlugin({
  name: 'auth',
  setup() {
    const { restoreSession } = useAuth()
    // Строго после гидрации: сайт статический, в разметке пользователь всегда
    // не авторизован, и подстановка профиля раньше времени ломает совпадение.
    onNuxtReady(() => {
      void restoreSession()
    })
  },
})
