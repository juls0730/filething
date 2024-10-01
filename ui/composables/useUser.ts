import type { Plan, User } from '~/types/user'
import { useFetch } from '#app'

const uninitializedUser = {
    id: "",
    username: "",
    email: "",
    plan: <Plan>{
        id: 0,
        max_storage: 0
    },
    usage: 0,
    created_at: "",
    is_admin: false
}

export const useUser = () => {
    // Global state for storing the user
    const user = useState('user', () => { return { fetched: false, user: uninitializedUser } })

    // Fetch the user only if it's uninitialized (i.e., null)
    const getUser = async () => {
        if (!user.value.fetched && useCookie('sessionToken').value) {
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
        } catch (e: any) {
            console.error(e.message)
            user.value.user = uninitializedUser
        }
    }

    // Manually set the user (e.g., after login/signup)
    const setUser = (userData: User) => {
        user.value.user = userData
    }

    // Clear the user data (e.g., on logout)
    const resetUser = () => {
        user.value.user = uninitializedUser
    }

    return {
        getUser,
        setUser,
        resetUser,
        fetchUser,
    }
}