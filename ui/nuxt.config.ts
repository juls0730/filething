// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    ssr: true,
    compatibilityDate: '2024-04-03',

    css: ['~/assets/css/main.css'],

    colorMode: {
        classSuffix: ''
    },

    nitro: {
        routeRules: {
            '/api/**': { proxy: 'http://localhost:1323/api/**' },
            '/test/**': { proxy: 'http://localhost:1323/api/**' },
        }
    },

    devtools: { enabled: true },

    modules: ['@nuxtjs/color-mode', '@nuxtjs/tailwindcss']
})