<script setup lang="ts">
import { useUser } from '~/composables/useUser'
const { getUser } = useUser()

const user = await getUser();
const route = useRoute();

const accordionMapping: Record<string, string> = {
    '/admin/users': '',
    '/admin/config/settings': 'item-2',
    '/admin': 'item-1',
};

const getActiveAccordion = () => {
    const path = Object.keys(accordionMapping).find(key => route.path.startsWith(key));
    return path ? accordionMapping[path] : null;
};

const isActiveLink = (path: string) => route.path === path;
</script>

<template>
    <div class="w-full">
        <Nav :user="user" />
        <div class="py-4 px-4 lg:px-8 flex flex-col md:flex-row gap-4">
            <div class="w-auto md:w-60 flex-shrink-0">
                <aside class="rounded-md border overflow-hidden w-full h-fit">
                    <div class="px-4 py-3.5 bg-surface border-b">Admin Settings</div>
                    <VlAccordion type="single" :defaultValue="getActiveAccordion() || undefined">
                        <VlAccordionItem value="item-1" class="text-sm">
                            <VlAccordionTrigger class="transition-bg hover:bg-muted/10 px-4 py-3.5">
                                Maintenance
                            </VlAccordionTrigger>
                            <VlAccordionContent class="mt-1 mb-2">
                                <div class="text-xs">
                                    <NuxtLink to="/admin"
                                        class="w-full indent-4 px-4 py-1.5 hover:text-text text-subtle block"
                                        :class="isActiveLink('/admin') ? 'text-text' : 'text-subtle'">
                                        Dashboard
                                    </NuxtLink>
                                </div>
                            </VlAccordionContent>
                        </VlAccordionItem>
                        <VlAccordionItem value="item-2" class="text-sm">
                            <VlAccordionTrigger class="transition-bg hover:bg-muted/10 px-4 py-3.5">
                                Configuration
                            </VlAccordionTrigger>
                            <VlAccordionContent class="mt-1 mb-2">
                                <div class="text-xs">
                                    <NuxtLink to="/admin/config/settings"
                                        class="w-full indent-4 px-4 py-1.5 hover:text-text block"
                                        :class="isActiveLink('/admin/config/settings') ? 'text-text' : 'text-subtle'">
                                        Settings
                                    </NuxtLink>
                                </div>
                            </VlAccordionContent>
                        </VlAccordionItem>
                        <NuxtLink
                            class="vl-accordion-header focus-visible:outline-none focus-visible:ring focus-visible:ring-inset flex flex-1 justify-between items-center w-full transition-bg px-4 py-3.5 text-sm"
                            :class="route.path.startsWith('/admin/users') ? 'bg-muted/15' : 'hover:bg-muted/10'"
                            to="/admin/users">
                            Users
                        </NuxtLink>
                    </VlAccordion>
                </aside>
            </div>
            <slot />
        </div>
    </div>
</template>

<style>
.vl-accordion .vl-accordion-item:not(:last-child) {
    border-bottom-width: 1px;
}
</style>