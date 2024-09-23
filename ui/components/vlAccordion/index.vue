<script setup>
const props = defineProps({
    type: {
        type: String,
        default: "multiple"
    },
    defaultValue: {
        type: String
    }
})

const accordion = ref(null);

const accordionItems = ref([])

function toggleAccordion(index) {
    const item = accordionItems.value[index];
    if (props.type === "single") {
        // close everything but the one we just opened
        accordionItems.value.forEach((item, i) => {
            if (i === index) return;
            item.isOpen = false;
        })
    }
    if (item.hidden) {
        item.hidden = false;
    }
    item.isOpen = !item.isOpen;
}

const registerAccordionItem = (value) => {
    const item = { isOpen: ref(false), hidden: ref(true), value }
    accordionItems.value.push(item);

    return { index: accordionItems.value.indexOf(item), isOpen: item.isOpen, hidden: item.hidden, value };
};

const unregisterAccordionItem = (index) => {
    accordionItems.value.splice(index, 1);
};

provide('accordion', { registerAccordionItem, unregisterAccordionItem, toggleAccordion })

function keydown(event) {
    const headers = Array.from(accordion.value.querySelectorAll(".vl-accordion-header"));

    if (event.key === "ArrowUp") {
        event.preventDefault();
        const focusedElement = document.activeElement;
        const currentIndex = headers.indexOf(focusedElement);
        const nextIndex = currentIndex > 0 ? currentIndex - 1 : headers.length - 1;
        console.log(nextIndex, headers)
        const nextButton = headers[nextIndex];
        nextButton.focus();
        return;
    }

    if (event.key === "ArrowDown") {
        event.preventDefault();
        const focusedElement = document.activeElement;
        const currentIndex = headers.indexOf(focusedElement);
        const nextIndex = currentIndex < headers.length - 1 ? currentIndex + 1 : 0;
        const nextButton = headers[nextIndex];
        console.log(nextIndex, headers)
        nextButton.focus();
        return;
    }

    if (event.key === "End") {
        event.preventDefault();
        return headers[headers.length - 1].focus();
    }

    if (event.key === "Home") {
        event.preventDefault();
        return headers[0].focus();
    }
}

onMounted(() => {
    if (!!props.defaultValue) {
        const item = accordionItems.value.filter(item => item.value === props.defaultValue)[0];
        item.isOpen = true;
        item.hidden = false;
    }

})

watch(props, () => {
    if (!!props.defaultValue) {
        const item = accordionItems.value.filter(item => item.value === props.defaultValue)[0];
        item.isOpen = true;
        item.hidden = false;
    }
})
</script>

<template>
    <div class="vl-accordion" ref="accordion" @keydown="keydown($event)">
        <slot />
    </div>
</template>

<style>
.vl-accordion-content {
    overflow: hidden;
    transform-origin: top center;
    height: 0;
}


.vl-accordion-content[data-state="closed"] {
    animation: 300ms cubic-bezier(0.25, 1, 0.5, 1) 0s 1 normal forwards running closeAccordion;
}

.vl-accordion-content[data-state="open"] {
    animation: 300ms cubic-bezier(0.25, 1, 0.5, 1) 0s 1 normal forwards running openAccordion;
}

@keyframes closeAccordion {
    0% {
        height: var(--vueless-accordion-content-height);
    }

    100% {
        height: 0px;
    }
}

@keyframes openAccordion {
    0% {
        height: 0px;
    }

    100% {
        height: var(--vueless-accordion-content-height);
    }
}
</style>