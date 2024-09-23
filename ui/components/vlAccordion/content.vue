<script setup>
const item = inject('accordionItem');
const contentHeight = ref(0);
const content = ref(null);

let timeout;
watch(item.hidden, () => {
    if (!item.hidden.value) {
        timeout = setTimeout(() => {
            let styles = window.getComputedStyle(content.value);
            let margin = parseFloat(styles['marginTop']) +
                parseFloat(styles['marginBottom']);

            contentHeight.value = content.value.offsetHeight + margin;
        })
    }
})

const attrs = useAttrs();

defineOptions({
    inheritAttrs: false
});

onUnmounted(() => {
    clearTimeout(timeout)
})
</script>

<template>
    <div :id="`vueless-${item.index}`" role="region" :aria-labelledby="`vueless-${item.index}`"
        class="vl-accordion-content" :style="`--vueless-accordion-content-height: ${contentHeight}px`"
        :data-state="(item.isOpen.value) ? 'open' : 'closed'"
        @animationend="(!item.isOpen.value) ? item.hidden.value = true : ''" :hidden="item.hidden.value">
        <div ref="content" v-bind="attrs">
            <slot />
        </div>
    </div>
</template>