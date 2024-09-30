<script setup lang="ts">
import type { NuxtError } from '#app';
import type { User } from '~/types/user';

definePageMeta({
    middleware: ["auth", "admin"],
    layout: "admin"
});

let username = ref('')
let email = ref('')
let password = ref('')

let error = ref('')

let timeout: NodeJS.Timeout;
const submitForm = async () => {
    let { data, error: fetchError } = await useAsyncData<User, NuxtError<{ message: string }>>(
        () => $fetch('/api/admin/users/new', {
            method: 'POST',
            body: {
                "username": username.value,
                "email": email.value,
                "password": password.value,
            }
        })
    )

    if (fetchError.value != null && fetchError.value.data !== undefined) {
        error.value = fetchError.value.data.message
        timeout = setTimeout(() => error.value = "", 15000)
    } else if (data.value !== null) {
        await navigateTo('/admin/users')
    }
}

onUnmounted(() => {
    clearTimeout(timeout)
})
</script>

<template>
    <div class="w-full h-fit mb-4">
        <div class="overflow-hidden rounded-md border text-[15px]">
            <h4 class="bg-surface px-3.5 py-3 border-b">Create User Account
            </h4>
            <div class="p-4">
                <label for="username" class="block max-w-64 text-sm">Username</label>
                <Input v-model="username" :value="username" id="username" placeholder="Username" autocomplete="off" class="w-full mb-2" />
                <label for="email" class="block max-w-64 text-sm">Email</label>
                <Input v-model="email" :value="email" id="email" placeholder="Email" autocomplete="off" class="w-full mb-2" />
                <label for="password" class="block max-w-64 text-sm">Password</label>
                <Input v-model="password" :value="password" id="password" type="password" placeholder="Password" autocomplete="off" class="w-full mb-2" />
                <p class="text-love mb-2">{{ error }}</p>
                <div>
                    <button
                        class="transition-bg bg-pine/10 text-pine px-3 py-2 rounded-md hover:bg-pine/15 active:bg-pine/25"
                        v-on:click="submitForm">
                        Create User Account
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>