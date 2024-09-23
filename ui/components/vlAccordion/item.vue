<script setup>
const props = defineProps({
    value: {
        type: String,
        default: "",
    }
})
const accordion = inject('accordion')

const item = accordion.registerAccordionItem(props.value);

provide('accordionItem', item);
</script>

<template>
    <div class="vl-accordion-item" :data-state="(item.isOpen.value) ? 'open' : 'closed'">
        <slot />
    </div>
</template>