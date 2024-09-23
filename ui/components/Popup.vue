<script lang="ts" setup>
defineProps(['modelValue', 'header'])
const emit = defineEmits(['update:modelValue'])
</script>

<template>
    <div class="grid place-content-center absolute top-0 left-0 bottom-0 right-0 z-40"
        :class="{ 'hidden': !modelValue }">
        <div v-on:click=" $emit('update:modelValue', !modelValue)"
            class="absolute top-0 left-0 bottom-0 right-0 bg-base/40">
        </div>

        <transition name="scale-fade">
            <div v-if="modelValue"
                class="bg-surface rounded-xl border shadow-md p-6 transition-[transform,opacity] duration-[250ms] origin-center z-50 w-screen h-screen sm:w-[600px] sm:h-auto">
                <div class="flex justify-between mb-2 items-center">
                    <h3 class="text-xl font-semibold">{{ header }}</h3>
                    <button v-on:click=" $emit('update:modelValue', !modelValue)"
                        class="p-1 border h-fit rounded-md hover:bg-muted/10 active:bg-muted/20 transition-bg focus-visible:outline-none focus-visible:ring focus-visible:ring-inset">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                stroke-width="2" d="M18 6L6 18M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                <slot />
            </div>
        </transition>
    </div>
</template>

<style scoped>
.scale-fade-enter-from,
.scale-fade-leave-to {
    opacity: 0;
    transform: scale(0.90);
}

.scale-fade-enter-to {
    opacity: 1;
    transform: scale(1);
}
</style>