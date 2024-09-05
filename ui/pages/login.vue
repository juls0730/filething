<script setup>
console.log(useCookie("sessionToken").value)

if (useCookie("sessionToken").value) {
    await navigateTo('/')
}

let username_or_email = ref('')
let password = ref('')

let error = ref('')

const submitForm = async () => {
    let response = await useFetch('/api/login', {
        method: 'POST',
        body: {
            "username_or_email": username_or_email.value,
            "password": password.value,
        }
    })

    if (response.error.value != null) {
        console.log(response)
        error.value = response.error.value.data.message
        setTimeout(() => error.value = "", 15000)
    } else {
        await navigateTo('/')
    }
}
</script>

<template>
    <div class="min-h-screen min-w-screen grid place-content-center bg-base">
        <div
            class="flex flex-col text-center bg-surface border border-muted/20 shadow-md px-10 py-8 rounded-2xl min-w-0 max-w-[313px]">
            <h2 class="font-semibold text-2xl mb-2">Login</h2>
            <Input v-model="username_or_email" placeholder="Username or Email..." />
            <Input v-model="password" type="password" placeholder="Password..." />
            <p class="text-love">{{ error }}</p>
            <button @click="submitForm"
                class="py-2 px-4 my-2 bg-pine/10 text-pine rounded-md transition-colors hover:bg-pine/15 active:bg-pine/25">Login</button>
            <p>Or <NuxtLink to="/signup" class="text-foam hover:underline">Sign up</NuxtLink>
            </p>
        </div>
    </div>
</template>