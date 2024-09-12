<script setup lang="ts">
let user = await useUser().getUser()

defineEmits(["update:filenav"])
defineProps(["filenav"])
</script>

<template>
    <header class="flex h-[var(--nav-height)] px-4 justify-center sticky top-0 z-10 border-b bg-base">
        <div class="flex w-full items-center justify-between space-x-2.5">
            <p
                class="-ml-2.5 flex shrink-0 items-center px-2.5 py-1.5 focus:outline-none focus:ring rounded-m font-semiboldd">
                filething
            </p>
        </div>
        <nav class="flex md:hidden">
            <ul class="flex items-center gap-3" role="list">
                <li>
                    <span class="group relative flex items-center">
                        <button
                            class="flex items-center px-3 h-8 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0-8 0M6 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2" />
                            </svg>
                            <span class="group-focus-within:rotate-180 group-hover:rotate-180 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                                    <path fill="none" stroke="currentColor" stroke-linecap="round"
                                        stroke-linejoin="round" stroke-width="2" d="m6 9l6 6l6-6" />
                                </svg>
                            </span>
                        </button>
                        <NavUserDropdown :user="user" />
                    </span>
                </li>
                <li class="h-6 border-r"></li>
                <li>
                    <button v-on:click="$emit('update:filenav', !filenav)"
                        class="flex items-center px-3 h-8 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                stroke-width="2" d="M4 6h16M7 12h13m-10 6h10" />
                        </svg>
                    </button>
                </li>
            </ul>
        </nav>
        <nav class="hidden md:flex" aria-label="Main">
            <ul class="flex items-center gap-3" role="list">
                <li>
                    <a href="#"
                        class="px-2.5 py-1.5 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md">Link</a>
                </li>
                <li>
                    <a href="#"
                        class="px-2.5 py-1.5 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md">Link</a>
                </li>
                <li>
                    <a href="#"
                        class="px-2.5 py-1.5 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md">Link</a>
                </li>
                <li class="h-6 border-r"></li>
                <li>
                    <span class="group relative flex items-center">
                        <button
                            class="flex items-center px-3 h-8 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0-8 0M6 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2" />
                            </svg>
                            <span class="group-focus-within:rotate-180 group-hover:rotate-180 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                                    <path fill="none" stroke="currentColor" stroke-linecap="round"
                                        stroke-linejoin="round" stroke-width="2" d="m6 9l6 6l6-6" />
                                </svg>
                            </span>
                        </button>
                        <NavUserDropdown :user="user" />
                    </span>
                </li>
            </ul>
        </nav>
    </header>
    <div></div>
</template>