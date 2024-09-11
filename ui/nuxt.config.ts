// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    ssr: process.env.NODE_ENV === 'production' ? true : false,
    compatibilityDate: '2024-04-03',

    css: ['~/assets/css/main.css'],

    colorMode: {
        classSuffix: ''
    },

    devtools: { enabled: true },

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
