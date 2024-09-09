<script lang="ts" setup>
import type { User } from '~/types/user'
const { fetchUser } = useUser()

definePageMeta({
    middleware: "unauth"
});

let username = ref('')
let email = ref('')
let password = ref('')

let error = ref('')

const submitForm = async () => {
    const response = await useFetch<User>('/api/signup', {
        method: 'POST',
        body: {
            "username": username.value,
            "email": email.value,
            "password": password.value,
        }
    })

    if (response.error.value != null) {
        error.value = response.error.value.data.message
        setTimeout(() => error.value = "", 15000)
    } else {
        await fetchUser()
        await navigateTo('/home')
    }
}
</script>

<template>
    <div class="min-h-screen min-w-screen grid place-content-center bg-base">
        <div
            class="flex flex-col text-center bg-surface border border-muted/20 shadow-md px-10 py-8 rounded-2xl min-w-0 max-w-[313px]">
            <h2 class="font-semibold text-2xl mb-2">Signup</h2>
            <Input v-model="username" placeholder="Username..." />
            <Input v-model="email" placeholder="Email..." />
            <Input v-model="password" type="password" placeholder="Password..." />
            <p class="text-love">{{ error }}</p>
            <button @click="submitForm"
                class="py-2 px-4 my-2 bg-pine/10 text-pine rounded-md transition-colors hover:bg-pine/15 active:bg-pine/25">Login</button>
            <p>Or <NuxtLink to="/login" class="text-foam hover:underline">Log in</NuxtLink>
            </p>
        </div>
    </div>
</template>