<script setup>
const props = defineProps({
    placeholder: String,
    type: String | undefined,
})
const emit = defineEmits(['update:modelValue'])

function updateValue(value) {
    emit('update:modelValue', value)
}
</script>

<template>
    <input
        class="py-2 px-4 resize-none bg-overlay rounded-md border hover:border-muted/40 focus:border-muted/60 placeholder:italic placeholder:text-subtle transition-[border-color] max-w-64"
        :placeholder="placeholder" :type="type" v-on:input="updateValue($event.target.value)" />
</template>