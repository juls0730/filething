<script setup lang="ts">
import type { Plan, User } from '~/types/user';

definePageMeta({
    middleware: ["auth", "admin"],
    layout: "admin"
});

const route = useRoute();

let { data: user } = await useFetch<User>('/api/admin/users/' + route.params.id);

let username = ref(user.value?.username);
let email = ref(user.value?.email);
let password = ref('');
let plan_id = ref(user.value?.plan.id);
let is_admin = ref(user.value?.is_admin ? 'checked' : 'unchecked');

const updateUser = async () => {
    let body = {
        username: username.value,
        email: email.value,
        password: password.value as string || undefined,
        plan_id: plan_id.value,
        is_admin: is_admin.value === 'checked' ? true : false,
    }

    if (password.value === '') {
        delete body.password
    }

    await $fetch('/api/admin/users/edit/' + route.params.id, {
        method: "POST",
        body,
    })
}

let { data: plans } = await useFetch<Plan[]>('/api/admin/plans');
</script>

<template>
    <div class="w-full h-fit mb-4">
        <div class="overflow-hidden rounded-md border text-[15px]">
            <h4 class="bg-surface px-3.5 py-3 border-b">Edit User Account
            </h4>
            <div class="p-4">
                <label for="username" class="block max-w-64 text-sm">Username</label>
                <input v-model="username" id="username" placeholder="Username" class="w-full mb-2" />
                <label for="email" class="block max-w-64 text-sm">Email</label>
                <input v-model="email" id="email" placeholder="Email" class="w-full mb-2" />
                <div class="mb-2">
                    <label for="password" class="block max-w-64 text-sm">Password</label>
                    <input v-model="password" id="password" placeholder="Password" class="w-full" />
                    <p class="text-muted text-sm">Leave the password empty to keep it unchanged</p>
                </div>
                <label for="plan_id" class="block max-w-64 text-sm">Plan</label>
                <!-- select the one with the value of user.value.plan_id -->
                <select v-model="plan_id" id="plan_id" :selected="plan_id"
                    class="w-full max-w-64 px-4 py-2 rounded-md bg-overlay border hover:border-muted/40 focus:border-muted/60 cursor-pointer">
                    <option v-for="plan in plans" :key="plan.id" :value="plan.id">
                        {{ formatBytes(plan.max_storage) }}
                    </option>
                </select>
                <hr class="my-4" />
                <div class="flex items-center">
                    <Checkbox v-model="is_admin" id="is_admin" type="checkbox" class="mr-2" />
                    <label for="is_admin" class="text-sm">
                        Is Admin
                    </label>
                </div>
                <hr class="my-4" />
                <div>
                    <button
                        class="transition-bg bg-pine/10 text-pine px-3 py-2 rounded-md hover:bg-pine/15 active:bg-pine/25"
                        v-on:click="updateUser">
                        Update User
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>