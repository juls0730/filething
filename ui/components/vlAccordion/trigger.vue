<script setup>
const item = inject('accordionItem');
const { toggleAccordion } = inject('accordion');
</script>

<template>
    <button
        class="vl-accordion-header focus-visible:outline-none focus-visible:ring focus-visible:ring-inset select-none cursor-pointer w-full"
        @click="toggleAccordion(item.index)">
        <div class="flex flex-1 justify-between items-center w-full" :id="`vueless-${item.index}`"
            :aria-controls="`vueless-${item.index}`" :aria-expanded="item.isOpen.value">
            <slot />
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="m6 9l6 6l6-6" />
            </svg>
        </div>
    </button>
</template>

<style>
.vl-accordion-item button div svg {
    transition: transform 300ms ease;
}

.vl-accordion-item[data-state="open"] button div svg {
    transform: rotate(180deg);
}
</style>