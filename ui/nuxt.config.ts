// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    ssr: true,
    compatibilityDate: '2024-04-03',

    css: ['~/assets/css/main.css'],

    colorMode: {
        classSuffix: ''
    },

    devtools: { enabled: true },

    experimental: {
        buildCache: true,
    },

    modules: [
        '@nuxtjs/color-mode',
    ],

    postcss: {
        plugins: {
            tailwindcss: {},
            autoprefixer: {},
        },
    },
})
