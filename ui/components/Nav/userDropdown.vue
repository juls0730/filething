<script setup lang="ts">
const logout = async () => {
    await $fetch('/api/logout', {
        method: "POST"
    })
    useCookie("sessionToken").value = null
    useUser().resetUser()

    navigateTo("/login")
}

defineProps({
    user: {
        type: Object,
        required: true
    },
    changeTheme: {
        type: Function,
        required: true,
    }
})
</script>

<template>
    <div
        class="invisible z-10 w-fit h-fit absolute -right-[4px] top-full opacity-0 group-hover:visible group-focus-within:visible group-focus-within:scale-100 group-focus-within:opacity-100 group-hover:scale-100 group-hover:opacity-100 transition">
        <div class="mt-1 w-64 origin-top-right scale-[.97] rounded-xl bg-surface shadow-lg">
            <div class="max-w-64 overflow-hidden text-ellipsis p-2">
                <p class="text-lg font-semibold">{{ user.username }}</p>
                <p class="text-subtle text-xs">{{ user.email }}</p>
                <p class="text-subtle text-xs">
                    you have {{ formatBytes(user.plan.max_storage) }} of storage
                </p>
            </div>
            <hr />
            <ul class="py-2 flex flex-col gap-x-1">
                <li class="select-none mx-2">
                    <button v-on:click="changeTheme"
                        class="flex items-center hover:bg-muted/10 active:bg-muted/20 transition-bg w-full px-2 py-1 rounded-md focus-visible:outline-none focus-visible:ring focus-visible:ring-inset">
                        <span class="mr-2">
                            <svg v-if="$colorMode.preference === 'dark'" xmlns="http://www.w3.org/2000/svg" width="18"
                                height="18" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2" d="M12 3h.393a7.5 7.5 0 0 0 7.92 12.446A9 9 0 1 1 12 2.992z" />
                            </svg>
                            <svg v-else-if="$colorMode.preference === 'light'" xmlns="http://www.w3.org/2000/svg"
                                width="18" height="18" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M14.828 14.828a4 4 0 1 0-5.656-5.656a4 4 0 0 0 5.656 5.656m-8.485 2.829l-1.414 1.414M6.343 6.343L4.929 4.929m12.728 1.414l1.414-1.414m-1.414 12.728l1.414 1.414M4 12H2m10-8V2m8 10h2m-10 8v2" />
                            </svg>
                            <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 256 256">
                                <path fill="currentColor"
                                    d="M208 36H48a28 28 0 0 0-28 28v112a28 28 0 0 0 28 28h160a28 28 0 0 0 28-28V64a28 28 0 0 0-28-28Zm4 140a4 4 0 0 1-4 4H48a4 4 0 0 1-4-4V64a4 4 0 0 1 4-4h160a4 4 0 0 1 4 4Zm-40 52a12 12 0 0 1-12 12H96a12 12 0 0 1 0-24h64a12 12 0 0 1 12 12Z" />
                            </svg>
                        </span>
                        Change Theme
                    </button>
                </li>
                <hr class="my-2" />
                <li v-if="user.is_admin" class="select-none mx-2">
                    <NuxtLink
                        class="flex items-center hover:bg-muted/10 active:bg-muted/20 transition-bg w-full px-2 py-1 rounded-md focus-visible:outline-none focus-visible:ring focus-visible:ring-inset"
                        to="/admin">
                        <span class="mr-2">
                            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M3 7a3 3 0 0 1 3-3h12a3 3 0 0 1 3 3v2a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3zm0 8a3 3 0 0 1 3-3h12a3 3 0 0 1 3 3v2a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3zm4-7v.01M7 16v.01M11 8h6m-6 8h6" />
                            </svg>
                        </span>
                        Site Administration
                    </NuxtLink>
                </li>
                <hr v-if="user.is_admin" class="my-2" />
                <li class="select-none mx-2">
                    <button
                        class="flex items-center hover:bg-muted/10 active:bg-muted/20 transition-bg w-full px-2 py-1 rounded-md focus-visible:outline-none focus-visible:ring focus-visible:ring-inset"
                        v-on:click="logout">
                        <span class="mr-2">
                            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
                                <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2">
                                    <path
                                        d="M14 8V6a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2-2v-2" />
                                    <path d="M9 12h12l-3-3m0 6l3-3" />
                                </g>
                            </svg>
                        </span>
                        Logout
                    </button>
                </li>
            </ul>
        </div>
    </div>
</template>