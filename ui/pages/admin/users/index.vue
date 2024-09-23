<script setup lang="ts">
import { useUser } from "~/composables/useUser"
import type { User } from "~/types/user";
const { getUser } = useUser()

definePageMeta({
    middleware: ["auth", "admin"],
    layout: "admin"
});

let page = ref(0)

const { data: users } = await useFetch<User[]>('/api/admin/get-users/' + page.value);
const { data: usersCount } = await useFetch<{ total_users: number }>('/api/admin/get-total-users');

const fetchNextPage = async () => {
    page.value += 1;
    let moreUsers = await $fetch('/api/admin/get-users/' + page.value);
    console.log(moreUsers)
    users.value = users.value?.concat(moreUsers)
}
</script>

<template>
    <div class="w-full h-fit mb-4">
        <div class="overflow-hidden rounded-md border text-[15px]">
            <h4 class="bg-surface px-3.5 py-3 border-b">User Account Management (Total: {{ usersCount.total_users }})
            </h4>
            <div class="overflow-x-scroll max-w-full">
                <table class="min-w-full">
                    <thead>
                        <tr class="text-left">
                            <th class="py-2 px-4">ID</th>
                            <th class="py-2 px-4">Username</th>
                            <th class="py-2 px-4">Email Address</th>
                            <th class="py-2 px-4">Restricted</th>
                            <th class="py-2 px-4">Created</th>
                            <th class="py-2 px-4 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in users" class="border-t">
                            <td class="py-2 px-4 max-w-44" :title="user.id">{{ user.id }}</td>
                            <td class="py-2 px-4">
                                {{ user.username }}
                                <span v-if="user.is_admin"
                                    class="ml-2 text-xs bg-accent/10 text-accent py-1 px-2 rounded">Admin</span>
                            </td>
                            <td class="py-2 px-4">{{ user.email }} </td>
                            <td class="py-2 px-4">
                                <svg v-if="true" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                    viewBox="0 0 24 24">
                                    <path fill="none" stroke="currentColor" stroke-linecap="round"
                                        stroke-linejoin="round" stroke-width="2" d="M18 6L6 18M6 6l12 12" />
                                </svg>
                                <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                    viewBox="0 0 24 24">
                                    <path fill="none" stroke="currentColor" stroke-linecap="round"
                                        stroke-linejoin="round" stroke-width="2" d="m5 12l5 5L20 7" />
                                </svg>
                            </td>
                            <td class="py-2 px-4">{{ new Date(user.created_at).toLocaleDateString('en-US', {
                                year:
                                    'numeric', month: 'short', day: 'numeric'
                            }) }}</td>
                            <td class="py-2 px-4 h-full">
                                <div class="flex items-center justify-end">
                                    <NuxtLink :to="`/admin/users/${user.id}/edit`"></NuxtLink>
                                    <button
                                        class="my-auto hover:bg-muted/10 p-1 transition-bg active:bg-muted/20 rounded-md">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                            viewBox="0 0 24 24">
                                            <g class="stroke-blue-400/90" fill="none" stroke="currentColor"
                                                stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                                                <path d="M7 7H6a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2-2v-1" />
                                                <path d="M20.385 6.585a2.1 2.1 0 0 0-2.97-2.97L9 12v3h3zM16 5l3 3" />
                                            </g>
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="w-full h-full flex justify-center mt-4" v-if="users?.length != usersCount.total_users">
            <button class="bg-accent/10 text-accent px-2 py-1 rounded-md hover:" v-on:click="fetchNextPage()">Load
                More</button>
        </div>
    </div>
</template>