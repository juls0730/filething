<script lang="ts" setup>
let year = new Date().getFullYear()
</script>

<template>
    <footer class="flex flex-col gap-2 sm:flex-row py-6 w-full shrink-0 items-center px-4 md:px-6 border-t">
        <div class="flex flex-row gap-x-2">
            <p>Privacy policy</p>
            <p>Other policy</p>
        </div>
        <p class="text-xs text-subtle ml-auto">&copy; {{ year }} juls0730. All rights reserved.</p>
    </footer>
</template>