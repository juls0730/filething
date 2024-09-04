<script setup>
if (useCookie("sessionToken").value) {
    await navigateTo('/')
}

let email = ref('')
let password = ref('')

const submitForm = async () => {
    await useFetch('/api/signup', {
        method: 'POST',
        body: {
            "email": email.value,
            "password": password.value,
        }
    })

    await navigateTo('/')
}
</script>

<template>
    <div class="min-h-screen min-w-screen grid place-content-center">
        <div class="flex flex-col max-w-72 w-full text-center">
            <input
                class="py-2 px-4 resize-none rounded-md my-2 border-gray-200 dark:border-dark-gray-800 focus:border-gray-300 dark:focus:border-dark-gray-800 border focus:outline-none placeholder:italic placeholder:text-slate-500 dark:placeholder:text-gray-300"
                v-model="email" type="email" placeholder="Email..." />
            <input
                class="py-2 px-4 resize-none rounded-md my-2 border-gray-200 dark:border-dark-gray-800 focus:border-gray-300 dark:focus:border-dark-gray-800 border focus:outline-none placeholder:italic placeholder:text-slate-500 dark:placeholder:text-gray-300"
                v-model="password" type="password" placeholder="Password..." />
            <button @click="submitForm" class="py-2 px-4 bg-blue-500 text-white rounded-md">Signup</button>
            <p>Or <NuxtLink to="/login" class="text-pink-600">Log in</NuxtLink>
            </p>
        </div>
    </div>
</template>