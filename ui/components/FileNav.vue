<script lang="ts" setup>
import { useUser } from '~/composables/useUser'
const { getUser } = useUser()

const props = defineProps({
    usageBytes: {
        type: Number,
        required: true,
    }
})

const user = await getUser()

const route = useRoute();

let capacityBytes = ref(user.plan.max_storage);

const radius = 13;
const circumference = 2 * Math.PI * radius;

const percentage = computed(() => {
    return (props.usageBytes / capacityBytes.value);
});

const offset = computed(() => {
    return circumference - percentage.value * circumference;
});
const usage = computed(() => {
    return formatBytes(props.usageBytes)
});
const capacity = computed(() => {
    return formatBytes(capacityBytes.value)
});

const isAllFilesActive = computed(() => route.path === '/home');

const isInFolder = computed(() => route.path.startsWith('/home/') && route.path !== '/home');
</script>

<template>
    <aside class="h-screen flex flex-col w-56 pt-3 bg-surface border-r z-50 md:z-20">
        <a href="#main"
            class="absolute w-fit -translate-x-full top-0 px-2 py-4 bg-surface border  opacity-0 focus-within:translate-x-0 focus-within:opacity-100">
            Skip to content
        </a>
        <div class="pl-9 h-14 flex items-center">
            <h2>Home</h2>
        </div>
        <div class="p-4 flex-grow">
            <ul class="flex flex-col gap-y-2">
                <li>
                    <NuxtLink to="/home"
                        class="flex py-1.5 px-4 rounded-lg transition-bg duration-300 hover:bg-muted/10 focus-visible:outline-none focus-visible:ring focus-visible:ring-inset"
                        :class="{ 'bg-muted/10': isAllFilesActive }">
                        <div class="flex relative">
                            <svg class="m-0.5 mr-2" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                                viewBox="0 0 256 256">
                                <g fill="currentColor">
                                    <path d="M208 72v112a8 8 0 0 1-8 8h-24v-88l-40-40H80V40a8 8 0 0 1 8-8h80Z"
                                        opacity=".2" />
                                    <path
                                        d="m213.66 66.34l-40-40A8 8 0 0 0 168 24H88a16 16 0 0 0-16 16v16H56a16 16 0 0 0-16 16v144a16 16 0 0 0 16 16h112a16 16 0 0 0 16-16v-16h16a16 16 0 0 0 16-16V72a8 8 0 0 0-2.34-5.66ZM168 216H56V72h76.69L168 107.31V216Zm32-32h-16v-80a8 8 0 0 0-2.34-5.66l-40-40A8 8 0 0 0 136 56H88V40h76.69L200 75.31Zm-56-32a8 8 0 0 1-8 8H88a8 8 0 0 1 0-16h48a8 8 0 0 1 8 8Zm0 32a8 8 0 0 1-8 8H88a8 8 0 0 1 0-16h48a8 8 0 0 1 8 8Z" />
                                </g>
                            </svg>
                            All files
                            <div class="absolute -left-1.5 top-px bottom-px bg-accent w-[2px]"
                                :class="{ 'hidden': !isAllFilesActive }"></div>
                        </div>
                    </NuxtLink>
                </li>
                <!-- <li class="flex flex-col">
                    <NuxtLink to="/home/name"
                        class="flex py-1.5 px-4 rounded-lg transition-bg duration-300 hover:bg-muted/10"
                        :class="{ 'bg-muted/10': isInFolder }">
                        <div class="flex relative">
                            <svg v-if="isInFolder" class="m-0.5 mr-2" xmlns="http://www.w3.org/2000/svg" width="20"
                                height="20" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="m5 19l2.757-7.351A1 1 0 0 1 8.693 11H21a1 1 0 0 1 .986 1.164l-.996 5.211A2 2 0 0 1 19.026 19za2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h4l3 3h7a2 2 0 0 1 2 2v2" />
                            </svg>
                            <svg v-else class="m-0.5 mr-2" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                                viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M5 4h4l3 3h7a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2" />
                            </svg>
                            Folders
                            <div class="absolute -left-1.5 top-px bottom-px bg-accent w-[2px]"
                                :class="{ 'hidden': !isInFolder }"></div>
                        </div>
                    </NuxtLink>
                </li> -->
            </ul>
        </div>
        <div class="m-2 w-[calc(100%-16px)]">
            <div class="p-3 bg-overlay border rounded-lg flex items-end">
                <svg width="32" height="32" class="-rotate-90 mr-2" xmlns="http://www.w3.org/2000/svg">
                    <!-- Background Track -->
                    <circle class="stroke-accent/20" cx="16" cy="16" :r="radius" fill="none" stroke-width="3" />
                    <!-- Progress Track -->
                    <circle class="stroke-accent" cx="16" cy="16" :r="radius" fill="none" stroke-width="3"
                        :stroke-dasharray="circumference" :stroke-dashoffset="offset" stroke-linecap="round" />
                </svg>
                <p class="text-sm h-min"> {{ usage }} of {{ capacity }}</p>
            </div>
        </div>
    </aside>
</template>