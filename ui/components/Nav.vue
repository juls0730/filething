<script setup lang="ts">
let colorMode = useColorMode();

const changeTheme = () => {
    if (colorMode.preference === "dark") {
        // from dark => light
        colorMode.preference = "light"
    } else if (colorMode.preference === "light") {
        // from light => system
        colorMode.preference = "system";
    } else {
        // from system => dark
        colorMode.preference = "dark";
    }

    return;
}
</script>

<template>
    <header class="flex h-[var(--nav-height)] px-4 justify-center sticky top-0 z-50 border-b bg-base">
        <div class="flex w-full items-center justify-between space-x-2.5">
            <p
                class="-ml-2.5 flex shrink-0 items-center px-2.5 py-1.5 focus:outline-none focus:ring rounded-m font-semiboldd">
                filething
            </p>
        </div>
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
                    <button
                        class="flex items-center px-3 h-8 text-[15px] font-semibold transition-bg duration-300 hover:bg-muted/10 focus:outline-none focus:ring focus:ring-inset rounded-md"
                        @click="changeTheme">
                        <span class="inline-block">
                            <svg v-if="$colorMode.preference === 'dark'" xmlns="http://www.w3.org/2000/svg" width="22"
                                height="22" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2" d="M12 3h.393a7.5 7.5 0 0 0 7.92 12.446A9 9 0 1 1 12 2.992z" />
                            </svg>
                            <svg v-else-if="$colorMode.preference === 'light'" xmlns="http://www.w3.org/2000/svg"
                                width="22" height="22" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M14.828 14.828a4 4 0 1 0-5.656-5.656a4 4 0 0 0 5.656 5.656m-8.485 2.829l-1.414 1.414M6.343 6.343L4.929 4.929m12.728 1.414l1.414-1.414m-1.414 12.728l1.414 1.414M4 12H2m10-8V2m8 10h2m-10 8v2" />
                            </svg>
                            <svg v-else xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 256 256">
                                <path fill="currentColor"
                                    d="M208 36H48a28 28 0 0 0-28 28v112a28 28 0 0 0 28 28h160a28 28 0 0 0 28-28V64a28 28 0 0 0-28-28Zm4 140a4 4 0 0 1-4 4H48a4 4 0 0 1-4-4V64a4 4 0 0 1 4-4h160a4 4 0 0 1 4 4Zm-40 52a12 12 0 0 1-12 12H96a12 12 0 0 1 0-24h64a12 12 0 0 1 12 12Z" />
                            </svg>
                        </span>
                    </button>
                </li>
            </ul>
        </nav>
    </header>
</template>