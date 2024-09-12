<script setup lang="ts">
import { formatBytes } from '~/utils/formatBytes';
import type { FileUpload } from '~/types/file';

const props = defineProps({
    uploadingFiles: {
        type: Array<FileUpload>,
        required: true
    },
    closed: Boolean,
})
defineEmits(['update:closed'])

const abortUpload = (id: string) => {
    let file = props.uploadingFiles.find(upload => upload.id === id);
    if (!file) {
        throw new Error("Upload cannot be aborted file is missing!")
    }

    const controller = file.controller;
    if (controller) {
        controller.abort();
    }
};

const formatRemainingTime = (seconds: number): string => {
    if (seconds < 60) {
        return `${Math.floor(seconds)} second${Math.floor(seconds) === 1 ? '' : 's'} left`;
    }

    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) {
        return `${minutes} minute${minutes === 1 ? '' : 's'} left`;
    }

    const hours = Math.floor(minutes / 60);
    if (hours < 24) {
        return `${hours} hour${hours === 1 ? '' : 's'} left`;
    }

    const days = Math.floor(hours / 24);
    return `${days} day${days === 1 ? '' : 's'} left`;
};

const truncateFilenameToFitWidth = (filename: string, maxWidthPx: number, font = '18px ui-sans-serif,system-ui,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji') => {
    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');

    if (context === null) {
        return
    }

    context.font = font;

    const name = filename.substring(0, filename.lastIndexOf('.'));
    const extension = filename.substring(filename.lastIndexOf('.'));

    function getTextWidth(text: string): number {
        if (context === null) {
            return 0
        }

        return context.measureText(text).width;
    }

    if (getTextWidth(filename) <= maxWidthPx) {
        return filename;
    }

    let truncatedName = name;
    let charsToRemove = 4;
    while (getTextWidth(truncatedName + extension) > maxWidthPx && truncatedName.length > charsToRemove) {
        const start = Math.ceil((truncatedName.length - charsToRemove) / 2);
        const end = Math.floor((truncatedName.length + charsToRemove) / 2);

        truncatedName = truncatedName.substring(0, start) + '...' + truncatedName.substring(end);
        charsToRemove++;
    }

    canvas.remove()

    return truncatedName + extension;
}

let collapsed = ref(false);
let closeable = computed(() => props.uploadingFiles.filter(x => x.uploading === true).length === 0);
let overallRemaining = computed(() => {
    if (closeable.value) {
        return
    }
    const uploadingFiles = props.uploadingFiles.filter(x => x.uploading === true);

    return uploadingFiles.reduce((max, item) => item.remainingTime > max.remainingTime ? item : max).remainingTime
});
let overallPercentage = computed(() => {
    const uploadingFiles = props.uploadingFiles.filter(x => x.uploading === true);

    const totalLoaded = uploadingFiles.reduce((acc, file) => acc + file.length.loaded, 0);
    const totalSize = uploadingFiles.reduce((acc, file) => acc + file.length.total, 0);

    if (totalSize === 0) return 0; // Avoid division by zero

    return (totalLoaded / totalSize) * 100; // Return percentage
})

let uploadedSuccessfully = computed(() => props.uploadingFiles.filter(x => x.status.error === false));
let uploadFailed = computed(() => props.uploadingFiles.filter(x => x.status.error === true));
</script>

<template>
    <div class="absolute bottom-0 right-0 m-3 rounded-2xl border flex flex-col sm:w-[440px] w-[calc(100%-24px)] shadow-md bg-surface z-20"
        :class="{ 'h-[510px]': !collapsed, 'hidden': closed }">
        <div class="flex flex-row justify-between h-14 items-center mb-3 px-4" :class="{ 'hidden': collapsed }">
            <h3 class="text-xl font-semibold">Upload</h3>
            <div class="flex flex-row gap-x-2">
                <button v-on:click="collapsed = !collapsed"
                    class="p-1 border h-fit rounded-md hover:bg-muted/10 active:bg-muted/20 transition-bg">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                        <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                            stroke-width="2" d="m6 9l6 6l6-6" />
                    </svg>
                </button>
                <button v-on:click="$emit('update:closed', true)" v-if="closeable"
                    class="p-1 border h-fit rounded-md hover:bg-muted/10 active:bg-muted/20 transition-bg">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                        <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                            stroke-width="2" d="M18 6L6 18M6 6l12 12" />
                    </svg>
                </button>
            </div>
        </div>
        <div class="flex-grow px-4 overflow-y-auto max-h-[358px]" :class="{ 'hidden': collapsed }">
            <div v-for="(upload, index) in uploadingFiles" :key="index" :id="`file-upload-${upload.id}`">
                <div class="flex flex-row gap-x-2 py-2 w-full">
                    <div>
                        <svg v-if="upload.uploading" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                stroke-width="2">
                                <path d="M7 18a4.6 4.4 0 0 1 0-9a5 4.5 0 0 1 11 2h1a3.5 3.5 0 0 1 0 7h-1" />
                                <path d="m9 15l3-3l3 3m-3-3v9" />
                            </g>
                        </svg>
                        <svg v-else-if="upload.status.aborted" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                stroke-width="2"
                                d="M13 18.004H6.657C4.085 18 2 15.993 2 13.517s2.085-4.482 4.657-4.482c.393-1.762 1.794-3.2 3.675-3.773c1.88-.572 3.956-.193 5.444 1c1.488 1.19 2.162 3.007 1.77 4.769h.99c1.37 0 2.556.8 3.117 1.964M22 22l-5-5m0 5l5-5" />
                        </svg>
                        <svg v-else-if="upload.status.error" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                stroke-width="2"
                                d="M15 18.004H6.657C4.085 18 2 15.993 2 13.517s2.085-4.482 4.657-4.482c.393-1.762 1.794-3.2 3.675-3.773c1.88-.572 3.956-.193 5.444 1c1.488 1.19 2.162 3.007 1.77 4.769h.99c1.374 0 2.562.805 3.121 1.972M19 16v3m0 3v.01" />
                        </svg>
                        <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                stroke-width="2"
                                d="M11 18.004H6.657C4.085 18 2 15.993 2 13.517s2.085-4.482 4.657-4.482c.393-1.762 1.794-3.2 3.675-3.773c1.88-.572 3.956-.193 5.444 1c1.488 1.19 2.162 3.007 1.77 4.769h.99c1.388 0 2.585.82 3.138 2.007M15 19l2 2l4-4" />
                        </svg>
                    </div>
                    <div class="px-2 flex-grow">
                        <div class="flex flex-col">
                            <span
                                class="font-medium overflow-hidden overflow-ellipsis whitespace-nowrap inline-block max-w-[220px]">{{
                                    truncateFilenameToFitWidth(upload.file.name, 220) }}</span>
                            <div class="flex flex-row">
                                <div
                                    class="font-medium uppercase rounded-full bg-overlay text-[10px] px-2 py-0.5 -ml-1 mr-2 w-fit max-w-20 overflow-hidden overflow-ellipsis max-h-[19px] whitespace-nowrap inline-block">
                                    {{ upload.file.name.split(".")[upload.file.name.split(".").length - 1] }}
                                </div>
                                <div class="flex text-[10px] items-end text-subtle">
                                    <span
                                        class="h-min overflow-hidden overflow-ellipsis max-w-56 whitespace-nowrap inline-block"
                                        v-if="upload.uploading">
                                        Uploading - {{ formatBytes(upload.length.loaded, 1) }} / {{
                                            formatBytes(upload.length.total, 1) }} - {{
                                            formatRemainingTime(upload.remainingTime) }}
                                    </span>
                                    <span v-else-if="upload.status.code >= 200 && upload.status.code < 300"
                                        class="h-min overflow-hidden overflow-ellipsis max-w-56 whitespace-nowrap inline-block">
                                        Uploaded to upload path
                                    </span>
                                    <span
                                        class="h-min overflow-hidden overflow-ellipsis max-w-56 whitespace-nowrap inline-block"
                                        v-else-if="upload.status.aborted">
                                        Canceled
                                    </span>
                                    <span
                                        class="h-min overflow-hidden overflow-ellipsis max-w-56 whitespace-nowrap inline-block"
                                        v-else-if="upload.status.error">
                                        {{ upload.status.message }}
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="flex items-center" v-if="upload.uploading">
                        <button v-on:click="abortUpload(upload.id)"
                            class="h-fit p-1 border rounded-md hover:bg-love/10 active:bg-love/20 hover:text-love transition-[background-color,color] text-sm py-1 px-2">
                            Cancel
                        </button>
                    </div>
                </div>
                <div v-if="upload.length.loaded !== undefined && upload.status.code === undefined"
                    class="w-full rounded-full h-1 bg-foam/20 relative -mt-1">
                    <div class="bg-foam rounded-full absolute left-0 top-0 bottom-0 transition-[width]"
                        :style="'width: ' + Math.round((upload.length.loaded / upload.length.total) * 100) + '%;'">
                    </div>
                </div>
            </div>
        </div>
        <div class="m-3 rounded-md bg-overlay border bottom-2 flex flex-row">
            <div class="flex p-3 w-fit rounded-md">
                <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24">
                    <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                        stroke-width="2" d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2M7 9l5-5l5 5m-5-5v12" />
                </svg>
            </div>
            <div class="flex flex-row flex-grow items-center">
                <div v-if="!closeable" class="p-2 flex flex-col flex-grow relative">
                    <span class="font-medium font-pine">Uploading Files</span>
                    <span class="text-xs items-end text-subtle" v-if="overallRemaining">{{
                        formatRemainingTime(overallRemaining) }}</span>
                    <div class="bg-pine/25 absolute left-0 bottom-0 top-0"
                        :style="'width: ' + overallPercentage + '%;'">
                    </div>
                </div>
                <div v-else-if="uploadFailed.length === 0" class="p-2 flex flex-col flex-grow">
                    <span class="font-medium">Successfully Uploaded all files</span>
                </div>
                <div v-else-if="uploadedSuccessfully.length === 0" class="p-2 flex flex-col flex-grow">
                    <span class="font-medium">Failed to Uploaded all files</span>
                </div>
                <div v-else class="p-2 flex flex-col flex-grow">
                    <span class="font-medium">Successfully Uploaded some files</span>
                </div>
                <button v-if="collapsed" v-on:click="collapsed = !collapsed"
                    class="p-1 border h-fit rounded-md hover:bg-muted/10 active:bg-muted/20 transition-bg mr-4">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
                        <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                            stroke-width="2" d="m6 15l6-6l6 6" />
                    </svg>
                </button>
            </div>
        </div>
    </div>
</template>