<script setup lang="ts">
import type { User } from "~/types/user";

definePageMeta({
    middleware: ["auth", "admin"],
    layout: "admin"
});

let page = ref(0)

const { data } = await useFetch<{ users: User[], total_users: number }>('/api/admin/users?page=' + page.value);

if (data.value === null) {
    throw new Error("Failed to fetch users");
}

// let { users, total_users } = data.value;
let users = ref(data.value.users);
let total_users = ref(data.value.total_users);
const fetchNextPage = async () => {
    page.value += 1;
    let { users: moreUsers } = await $fetch<{ users: User[], total_users: number }>('/api/admin/users?page=' + page.value);
    console.log(moreUsers)
    users.value = users.value?.concat(moreUsers)
}
</script>

<template>
    <div class="w-full h-fit mb-4">
        <div class="overflow-hidden rounded-md border text-[15px]">
            <div class="flex bg-surface border-b items-center justify-between px-3.5 ">
                <h4 class="py-3 w-fit">User Account Management (Total: {{ total_users }})
                </h4>
                <NuxtLink to="/admin/users/new">
                    <button
                        class="transition-bg bg-pine/10 text-pine px-2 py-1.5 rounded-md hover:bg-pine/15 active:bg-pine/25 h-fit text-xs"
                        v-on:click="updateUser">
                        Create User Account
                    </button>
                </NuxtLink>
            </div>
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
                            <td class="py-2 px-4 max-w-44 whitespace-nowrap overflow-hidden text-ellipsis"
                                :title="user.id">
                                {{ user.id }}
                            </td>
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
                                    <NuxtLink :to="`/admin/users/${user.id}/edit`">
                                        <button
                                            class="my-auto hover:bg-muted/10 p-1 transition-bg active:bg-muted/20 rounded-md">
                                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                                viewBox="0 0 24 24">
                                                <g class="stroke-blue-400/90" fill="none" stroke="currentColor"
                                                    stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                                                    <path d="M7 7H6a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2-2v-1" />
                                                    <path
                                                        d="M20.385 6.585a2.1 2.1 0 0 0-2.97-2.97L9 12v3h3zM16 5l3 3" />
                                                </g>
                                            </svg>
                                        </button>
                                    </NuxtLink>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="w-full h-full flex justify-center mt-4" v-if="users?.length != total_users">
            <button class="transition-bg bg-pine/10 text-pine px-2 py-1 rounded-md hover:bg-pine/15 active:bg-pine/25"
                v-on:click="fetchNextPage()">Load
                More</button>
        </div>
    </div>
</template>