export default defineNuxtPlugin({
  name: 'auth',
  // не блокируем гидрацию: профиль подтягивается фоном, интерфейс уже отрисован
  parallel: true,
  async setup() {
    await useAuth().restoreSession()
  },
})
