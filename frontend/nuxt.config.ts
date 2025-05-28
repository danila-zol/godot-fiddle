// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },
  css: ['~/assets/css/main.scss'],
  runtimeConfig: {
    public: {
      docsUrl: "/docs",
      forumUrl: "/forum",
      apiRoot: "https://d5df6jka59qn3n45eubv.yl4tuxdu.apigw.yandexcloud.net/game-hangar",
      apiDemosPrefix: "/v1/demos",
      enginePath: "/thirdparty/godot/"
    }
  },
  routeRules : {
    '/**' : {cors: true, headers: { 
      "Cross-Origin-Opener-Policy": "same-origin",
      "Cross-Origin-Embedder-Policy": "require-corp"
    }},
    '/new' : {ssr: false},
    '/games/*': {ssr: false}
  },
})
