// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    // in SSG mode, we dont want to prerender pages since that will do actual API calls
    // and cause pages like /admin to literally be a redirect to /login
    ssr: process.env.RENDERING_MODE !== 'static',
    compatibilityDate: '2024-04-03',

    css: ['~/assets/css/main.css'],

    colorMode: {
        classSuffix: ''
    },

    nitro: {
        routeRules: {
            '/api/**': { proxy: 'http://localhost:1323/api/**' },
        },
        // these routes never change so we can statically prerender them
        prerender: {
            routes: ['/login', '/signup', '/']
        }
    },

    devtools: { enabled: true },

    modules: ['@nuxtjs/color-mode', '@nuxtjs/tailwindcss']
})