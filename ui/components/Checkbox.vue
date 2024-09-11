<template>
    <div v-on:click="toggle()" class="w-5 h-5 border rounded cursor-pointer flex items-center justify-center"
        :class="state === 'unchecked' ? 'hover:bg-muted/5 active:bg-muted/15' : 'bg-accent/10 hover:bg-accent/15 active:bg-accent/25 text-accent'">
        <div v-if="state === 'some'" class="w-8/12 h-0.5 bg-current rounded-full"></div>
        <span v-else-if="state === 'checked'">
            <svg class="w-full h-full" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="m5 12l5 5L20 7" />
            </svg>
        </span>
    </div>
</template>

<script lang="ts" setup>
const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])
let state = computed(() => props.modelValue)

const toggle = () => {
    const newState = state.value === 'unchecked' ? 'checked' : 'unchecked'
    emit('update:modelValue', newState);
}
</script>