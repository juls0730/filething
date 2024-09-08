<script setup lang="js">
const props = defineProps({
    path: {
        type: String,
        required: true
    }
});

const crumbs = computed(() => {
    const paths = props.path.split("/").filter(x => !!x);
    return paths.map((crumb, index) => {
        return {
            name: crumb,
            link: "/" + paths.slice(0, index + 1).join("/")
        };
    });
});
</script>

<template>
    <div class="flex flex-row">
        <span v-for="(crumb, index) in crumbs" class="flex items-center">
            <svg v-if="index != 0" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
                class="text-subtle mx-1">
                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="m9 6l6 6l-6 6" />
            </svg>
            <a class="hover:text-text" :class="index === crumbs.length - 1 ? 'text-foam' : 'text-subtle'"
                :href="crumb.link">{{
                    crumb.name }}</a>
        </span>
    </div>
</template>