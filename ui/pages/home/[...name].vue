<script lang="ts" setup>
import { useUser } from '~/composables/useUser'
import type { File } from '~/types/file';
import type { FileUpload } from '~/types/user';
const { getUser } = useUser()

definePageMeta({
    middleware: "auth"
});

const user = await getUser()

let { data: usageBytes } = await useFetch<{ usage: number }>('/api/user/usage')
let { data: files } = await useFetch<[File]>('/api/files')

const route = useRoute();
let folder = ref("");
let uploadPaneClosed = ref(true);

if (typeof route.params.name == "object") {
    folder.value = route.params.name.join("/");
}

let recentFiles = ref([]);

const fileInput: Ref<HTMLInputElement | null> = ref(null);

const uploadingFiles: Ref<Array<FileUpload>> = ref([]);

const handleFileChange = (event: Event) => {
    const files = (<HTMLInputElement>event.target).files;
    if (!files) {
        return;
    }

    for (let i = 0; i < files.length; i++) {
        uploadFile(files[i])
    }

    if (!fileInput.value) {
        return
    }

    if (fileInput.value.files.length > 0) {
        fileInput.value.value = "";
    }
}

const uploadFile = (file: File) => {
    const xhr = new XMLHttpRequest();
    const startTime = Date.now();
    let id = `${file.name}-${Math.floor(Math.random() * 1000)}`;

    let uploading_file: FileUpload = {
        id,
        uploading: true,
        controller: xhr,
        startTime,
        file: file,
        length: {},
        status: {}
    }

    uploadingFiles.value.push(uploading_file)

    if (uploadPaneClosed.value === true) {
        uploadPaneClosed.value = false;
    }

    xhr.open('POST', '/api/upload', true);

    xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
            let file = uploadingFiles.value.find(upload => upload.id === id);
            if (!file) {
                throw new Error("Upload is progressing but file is missing!")
            }


            const currentTime = Date.now();
            const timeElapsed = (currentTime - file.startTime) / 1000;

            file.length = { loaded: event.loaded, total: event.total };

            const uploadedBytes = event.loaded;
            const totalBytes = event.total;
            const uploadSpeed = uploadedBytes / timeElapsed;
            const remainingBytes = totalBytes - uploadedBytes;
            const remainingTime = remainingBytes / uploadSpeed;

            file.speed = uploadSpeed;
            file.remainingTime = remainingTime;
        }
    };

    xhr.onload = () => {
        let data = JSON.parse(xhr.response)
        usageBytes.value.usage = data.usage
        files.value?.push(data.file)

        let file = uploadingFiles.value.find(upload => upload.id === id);
        if (!file) {
            throw new Error("Upload has finished but file is missing!")
        }

        if (xhr.status >= 200 && xhr.status < 300) {
            file.uploading = false;

            file.status = {
                error: false,
                aborted: false,
                code: xhr.status,
                message: xhr.statusText
            };
        } else {
            file.uploading = false;

            file.status = {
                error: true,
                aborted: false,
                code: xhr.status,
                message: xhr.statusText
            };
        }
    };

    xhr.onerror = () => {
        let file = uploadingFiles.value.find(upload => upload.id === id);
        if (!file) {
            throw new Error("Upload has errored but file is missing!")
        }

        file.uploading = false;

        file.status = {
            error: true,
            aborted: false,
            code: xhr.status,
            message: xhr.statusText
        };
    };

    xhr.onabort = () => {
        let file = uploadingFiles.value.find(upload => upload.id === id);
        if (!file) {
            throw new Error("Upload has been aborted but file is missing!")
        }

        file.uploading = false;

        file.status = {
            error: true,
            aborted: true,
            code: 0,
            message: "aborted"
        };
    };

    const formData = new FormData();
    formData.append('file', file);

    xhr.send(formData);
};

const openFilePicker = () => {
    fileInput.value?.click();
}
</script>

<template>
    <div class="flex relative min-h-[100dvh]">
        <div class="fixed md:relative -translate-x-full md:translate-x-0">
            <FileNav :usageBytes="usageBytes?.usage" />
        </div>
        <UploadPane :closed="uploadPaneClosed" v-on:update:closed="(newValue) => uploadPaneClosed = newValue"
            :uploadingFiles="uploadingFiles" />
        <div class="w-full">
            <Nav />
            <div class="pt-6 pl-12 overflow-auto max-h-[calc(100vh-var(--nav-height))]">
                <div class="flex gap-x-4 flex-col">
                    <div class="py-5 flex flex-row gap-x-4">
                        <input type="file" ref="fileInput" @change="handleFileChange" multiple class="hidden" />
                        <button v-on:click="openFilePicker"
                            class="rounded-xl border-2 border-surface flex flex-col gap-y-2 px-2 py-3 w-40 justify-center items-center hover:bg-muted/10 active:bg-muted/20 transition-bg">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2">
                                    <path d="M7 18a4.6 4.4 0 0 1 0-9a5 4.5 0 0 1 11 2h1a3.5 3.5 0 0 1 0 7h-1" />
                                    <path d="m9 15l3-3l3 3m-3-3v9" />
                                </g>
                            </svg>
                            Upload
                        </button>
                        <button
                            class="rounded-xl border-2 border-surface flex flex-col gap-y-2 px-2 py-3 w-40 justify-center items-center hover:bg-muted/10 active:bg-muted/20 transition-bg">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M12 19H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h4l3 3h7a2 2 0 0 1 2 2v3.5M16 19h6m-3-3v6" />
                            </svg>
                            New folder
                        </button>
                    </div>
                    <div v-if="recentFiles.length > 0">
                        <h2 class="font-semibold text-2xl">Recent</h2>
                    </div>
                    <div>
                        <h3 class="font-semibold text-xl">
                            <Breadcrumbs :path="route.path" />
                        </h3>
                        <table class="w-full text-sm mt-2">
                            <thead class="border-b">
                                <tr class="flex flex-row h-10 group pl-[30px] -ml-7 relative items-center">
                                    <th class="left-0 absolute">
                                        <div>
                                            <input class="w-4 h-4 hidden group-hover:block" type="checkbox" />
                                        </div>
                                    </th>
                                    <th class="flex-grow text-start">
                                        Name
                                    </th>
                                    <th class="min-w-40 text-start">
                                        Size
                                    </th>
                                    <th class="min-w-40 text-start sm:block hidden">
                                        Last modified
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="block">
                                <tr class="flex flex-row h-10 group items-center border-b hover:bg-muted/10 transition-bg"
                                    v-for="file in files">
                                    <td class="-ml-7 pr-3.5">
                                        <div class="w-4 h-4">
                                            <input class="w-4 h-4 hidden group-hover:block" type="checkbox" />
                                        </div>
                                    </td>
                                    <td class="flex-grow text-start">
                                        <div class="flex items-center">
                                            <svg class="mr-2" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                                viewBox="0 0 24 24">
                                                <g fill="none" stroke="currentColor" stroke-linecap="round"
                                                    stroke-linejoin="round" stroke-width="2">
                                                    <path d="M14 3v4a1 1 0 0 0 1 1h4" />
                                                    <path
                                                        d="M17 21H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h7l5 5v11a2 2 0 0 1-2 2M9 9h1m-1 4h6m-6 4h6" />
                                                </g>
                                            </svg>
                                            {{ file.name }}
                                        </div>
                                    </td>
                                    <td class="min-w-40 text-start">
                                        {{ formatBytes(file.size) }}
                                    </td>
                                    <td class="min-w-40 text-start sm:block hidden">
                                        {{ file.last_modified }}
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
