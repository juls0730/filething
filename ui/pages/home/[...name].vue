<script lang="ts" setup>
import { useUser } from '~/composables/useUser'
import type { File } from '~/types/file';
import type { FileUpload } from '~/types/user';
const { getUser } = useUser()

definePageMeta({
    middleware: "auth"
});

const user = await getUser()
const route = useRoute();

let { data: files } = await useFetch<[File]>('/api/files/get/' + route.path.replace(/^\/home/, ''))

const sortedFiles = computed(() => {
    files.value?.forEach(file => file.toggled === undefined ? file.toggled = 'unchecked' : {})

    let folders = files.value?.filter(file => file.is_dir).sort((a, b) => {
        return ('' + a.name).localeCompare(b.name);
    });
    let archives = files.value?.filter(file => !file.is_dir).sort((a, b) => {
        return ('' + a.name).localeCompare(b.name);
    });

    return folders?.concat(archives)
})

let selectAll = ref('unchecked');
let selectedFiles = computed(() => sortedFiles.value?.filter(file => file.toggled === 'checked'))

watch(sortedFiles, (newVal, oldVal) => {
    let checkedFilesLength = newVal?.filter(file => file.toggled === 'checked').length;
    if (checkedFilesLength > 0) {
        if (checkedFilesLength < newVal?.length) {
            selectAll.value = 'some';
        } else {
            selectAll.value = 'checked';
        }
    } else {
        selectAll.value = 'unchecked';
    }
})

watch(selectAll, (newVal, oldVal) => {
    if (newVal === 'some') {
        return
    }

    sortedFiles.value?.forEach(file => {
        file.toggled = newVal
    })
});

let folderName = ref('');
let folder = ref("");
let folderError = ref('');
let popupVisable = ref(false);
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

    xhr.open('POST', '/api/files/upload/' + route.path.replace(/^\/home/, ''), true);

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
        let file = uploadingFiles.value.find(upload => upload.id === id);
        if (!file) {
            throw new Error("Upload has finished but file is missing!")
        }

        if (xhr.status >= 200 && xhr.status < 300) {
            let data = JSON.parse(xhr.response)
            user.usage = data.usage
            files.value?.push(data.file)

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

const createFolder = async () => {
    const { data, error } = await useAsyncData(
        () => $fetch('/api/files/upload' + route.path.replace(/^\/home/, '') + '/' + folderName.value, {
            method: "POST",
            body: {
                files: selectedFiles.value?.map(file => ({ name: file.name }))
            }
        })
    )

    if (error.value != null) {
        folderError.value = error.value.data.message;
    } else {
        user.usage = data.value.usage
        files.value?.push(data.value.file)

        popupVisable.value = false;
        navigateTo(route.path + '/' + folderName.value);
    }
}

const deleteFiles = async () => {
    await $fetch('/api/files/delete' + route.path.replace(/^\/home/, ''), {
        method: "POST",
        body: {
            files: selectedFiles.value?.map(file => ({ name: file.name }))
        }
    })

    files.value = files.value?.filter(file => !selectedFiles.value?.includes(file))
}

const downloadFile = (file) => {
    const anchor = document.createElement('a');
    anchor.href = '/api/files/download/' + file.name;
    anchor.download = file.name;

    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
}

const downloadFiles = async () => {
    let filenames = ""

    selectedFiles.value?.forEach((file, i) => {
        filenames += encodeURIComponent(file.name)
        if (i != selectedFiles.value?.length - 1) {
            filenames += ",";
        }
    })

    let { data, error } = await useAsyncData(
        () => $fetch('/api/files/download', {
            params: {
                "filenames": filenames
            }
        })
    )

    console.log("DATA", data.value)

    if (error.value == null) {
        const anchor = document.createElement('a');
        anchor.href = window.URL.createObjectURL(data.value)
        anchor.download = "filething.zip";

        document.body.appendChild(anchor);
        anchor.click();
        document.body.removeChild(anchor);
    }
}
</script>

<template>
    <div class="flex relative min-h-[100dvh]">
        <div class="fixed md:relative -translate-x-full md:translate-x-0">
            <FileNav :usageBytes="user.usage" />
        </div>
        <UploadPane :closed="uploadPaneClosed" v-on:update:closed="(newValue) => uploadPaneClosed = newValue"
            :uploadingFiles="uploadingFiles" />
        <Popup v-model="popupVisable" header="New Folder">
            <div class="flex flex-col p-2">
                <div class="mb-3 flex flex-col">
                    <label for="folderNameInput" class="text-sm">name</label>
                    <Input id="folderNameInput" v-model="folderName" placeholder="Folder name" />
                    <p class="text-love">{{ folderError }}</p>
                </div>
                <div class="ml-auto flex gap-x-1.5">
                    <button v-on:click="popupVisable = !popupVisable"
                        class=" px-2 py-1 rounded-md text-sm border bg-muted/10 hover:bg-muted/15 active:bg-muted/25 transition-bg">Close</button>
                    <button v-on:click="createFolder" :disabled="folderName === ''"
                        class=" px-2 py-1 rounded-md text-sm
                        disabled:bg-highlight-med/50 bg-highlight-med not:hover:brightness-105 not:active:brightness-110 transition-[background-color,filter] text-surface disabled:cursor-not-allowed">Confirm</button>
                </div>
            </div>
        </Popup>

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
                                    <path stroke="rgb(var(--color-accent))" d="m9 15l3-3l3 3m-3-3v9" />
                                </g>
                            </svg>
                            Upload
                        </button>
                        <button v-on:click="popupVisable = !popupVisable"
                            class="rounded-xl border-2 border-surface flex flex-col gap-y-2 px-2 py-3 w-40 justify-center items-center hover:bg-muted/10 active:bg-muted/20 transition-bg">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2">
                                    <path d="M12 19H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h4l3 3h7a2 2 0 0 1 2 2v3.5M16" />
                                    <path stroke="rgb(var(--color-accent))" d="M16 19h6m-3-3v6" />
                                </g>
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
                        <div class="mt-2">
                            <div class="flex flex-row gap-x-2" v-if="selectedFiles?.length > 0">
                                <button v-on:click="downloadFiles"
                                    class="flex flex-row px-2 py-1 rounded-md transition-bg text-xs border hover:bg-muted/10 active:bg-muted/20 items-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                                        <path fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2"
                                            d="M4 20h16m-8-6V4m0 10l4-4m-4 4l-4-4" />
                                    </svg>
                                    Download
                                </button>
                                <button v-on:click="deleteFiles"
                                    class="flex flex-row px-2 py-1 rounded-md transition-bg text-xs border hover:bg-love/10 active:bg-love/20 hover:text-love active:text-love items-center">
                                    <svg class="mr-1" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                        viewBox="0 0 24 24">
                                        <path fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2"
                                            d="M4 7h16M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2l1-12M9 7V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v3m-5 5l4 4m0-4l-4 4" />
                                    </svg>
                                    Delete
                                </button>
                            </div>
                        </div>
                        <table class="w-full text-sm mt-2 table-fixed">
                            <thead class="border-b">
                                <tr class="flex flex-row h-10 group pl-[30px] -ml-7 relative items-center">
                                    <th class="left-0 absolute">
                                        <div>
                                            <Checkbox :class="{ 'hidden': selectAll === 'unchecked' }"
                                                v-model="selectAll" class="group-hover:flex" type="checkbox" />
                                        </div>
                                    </th>
                                    <th v-on:click="selectAll === 'unchecked' ? selectAll = 'checked' : selectAll = 'unchecked'"
                                        class="flex-grow min-w-40 text-start flex items-center h-full">
                                        Name
                                    </th>
                                    <th class="min-w-32 text-start">
                                        Size
                                    </th>
                                    <th class="min-w-28 text-start sm:block hidden">
                                        Modified
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="block">
                                <tr class="flex border-l-2 flex-row h-10 group items-center border-b active:bg-surface/45 transition-bg relative"
                                    v-for="file in sortedFiles"
                                    :class="file.toggled === 'checked' ? 'bg-accent/20 border-l-accent' : 'border-l-transparent hover:bg-surface'">
                                    <td class="-ml-7 pr-4 flex-shrink-0">
                                        <div class="w-5 h-5">
                                            <Checkbox class="group-hover:flex"
                                                :class="{ 'hidden': file.toggled === 'unchecked' }"
                                                v-model="file.toggled" />
                                        </div>
                                    </td>
                                    <td v-on:click="file.toggled === 'unchecked' ? file.toggled = 'checked' : file.toggled = 'unchecked'"
                                        class="flex-grow text-start flex items-center h-full min-w-40">
                                        <div class="flex items-center min-w-40">
                                            <svg v-if="!file.is_dir" class="mr-2 flex-shrink-0"
                                                xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                                viewBox="0 0 24 24">
                                                <g fill="none" stroke="currentColor" stroke-linecap="round"
                                                    stroke-linejoin="round" stroke-width="2">
                                                    <path d="M14 3v4a1 1 0 0 0 1 1h4" />
                                                    <path
                                                        d="M17 21H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h7l5 5v11a2 2 0 0 1-2 2M9 9h1m-1 4h6m-6 4h6" />
                                                </g>
                                            </svg>
                                            <svg v-else class="mr-2 flex-shrink-0" xmlns="http://www.w3.org/2000/svg"
                                                width="16" height="16" viewBox="0 0 24 24">
                                                <path fill="none" stroke="currentColor" stroke-linecap="round"
                                                    stroke-linejoin="round" stroke-width="2"
                                                    d="M5 4h4l3 3h7a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2" />
                                            </svg>
                                            <span class="overflow-hidden whitespace-nowrap text-ellipsis">
                                                <NuxtLink v-if="file.is_dir" :to="`${route.path}/${file.name}`">
                                                    {{ file.name }}
                                                </NuxtLink>
                                                <span v-else>{{ file.name }}</span>
                                            </span>
                                        </div>
                                    </td>
                                    <td class="min-w-32 text-start">
                                        {{ formatBytes(file.size) }}
                                    </td>
                                    <td class="min-w-28 text-start sm:block hidden">
                                        {{ file.last_modified }}
                                    </td>
                                    <td :class="file.toggled === 'checked' ? 'context-active' : 'context'"
                                        class="absolute pl-6 top-0 bottom-0 right-0 hidden group-hover:flex items-center pr-8">
                                        <button v-on:click="downloadFile(file)"
                                            class="p-2 rounded hover:bg-muted/10 active:bg-muted/20">
                                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                                                viewBox="0 0 24 24">
                                                <path fill="none" stroke="currentColor" stroke-linecap="round"
                                                    stroke-linejoin="round" stroke-width="2"
                                                    d="M4 20h16m-8-6V4m0 10l4-4m-4 4l-4-4" />
                                            </svg>
                                        </button>
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

<style>
td,
th {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.context-active {
    background: linear-gradient(to right, transparent, var(--color-accent-20) 16px, var(--color-accent-20) 100%);
}

.context {
    background: linear-gradient(to right, transparent, rgb(var(--color-surface)) 16px, rgb(var(--color-surface)) 100%);
}
</style>