// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      docsUrl: "/docs",
      forumUrl: "/forum",
      apiRoot: "https://d5df6jka59qn3n45eubv.yl4tuxdu.apigw.yandexcloud.net/game-hangar",
      apiDemosPrefix: "/v1/demos",
    }
  }
})
