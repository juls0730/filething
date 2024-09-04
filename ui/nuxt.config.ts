// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',

    css: ['~/assets/css/main.css'],

    ssr: true,

    modules: [
        '@nuxtjs/color-mode',
    ],

    colorMode: {
        classSuffix: ''
    },

    devtools: { enabled: true },

    postcss: {
        plugins: {
            tailwindcss: {},
            autoprefixer: {},
        },
    },
})
