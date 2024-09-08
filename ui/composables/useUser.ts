import type { User } from '~/types/user'
import { useFetch } from '#app'

export const useUser = () => {
    // Global state for storing the user
    const user = useState('user', () => { return { fetched: false, user: <User>{} } })

    // Fetch the user only if it's uninitialized (i.e., null)
    const getUser = async () => {
        if (!user.value.fetched && import.meta.client) {
            await fetchUser()
        }

        return user.value.user
    }

    const fetchUser = async () => {
        try {
            const { data, error } = await useFetch<User, { message: string }>('/api/user')
            user.value.fetched = true

            if (error.value || !data.value) {
                throw new Error('Failed to fetch user')
            }

            user.value.user = data.value
        } catch (e) {
            console.error(e.message)
            user.value.user = {}
        }
    }

    // Manually set the user (e.g., after login/signup)
    const setUser = (userData: User) => {
        user.value.user = userData
    }

    // Clear the user data (e.g., on logout)
    const resetUser = () => {
        user.value.user = {}
    }

    return {
        getUser,
        setUser,
        resetUser,
        fetchUser,
    }
}