import { useUser } from '~/composables/useUser'

// We have server side things that does effectively this, but that wont stop SPA navigation
export default defineNuxtRouteMiddleware(async (to, from) => {
    if (import.meta.server) {
        return
    }

    const { getUser } = useUser()
    const user = await getUser()

    if (user.id) {
        return navigateTo('/home')
    }
})
