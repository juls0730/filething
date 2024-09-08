<script lang="ts" setup>
definePageMeta({
    middleware: "unauth"
});
</script>

<template>
    <div class="flex flex-col min-h-[100dvh]">
        <Nav />
        <main class="flex-1">
            <section
                class="flex justify-center w-full py-12 md:py-24 lg:py-32 xl:py-48 min-h-[calc(100dvh-var(--nav-height))]">
                <div class="container px-4 md:px-6 h-fit">
                    <div class="grid gap-6 lg:grid-cols-[1fr_400px] lg:gap-12 xl:grid-cols-[1fr_600px]">
                        <img src="/placeholder.svg" width="550" height="310" alt="Hero"
                            class="mx-auto aspect-video overflow-hidden rounded-xl object-cover sm:w-full lg:order-last" />
                        <div class="flex flex-col justify-center space-y-4">
                            <div class="space-y-2">
                                <h1 class="text-3xl font-bold tracking-tight sm:text-5xl xl:text-6xl/none">
                                    Effortless Open Source File Hosting
                                </h1>
                                <p class="max-w-[600px] text-subtle md:text-xl">
                                    Store, share, and access your files from anywhere with our powerful file hosting
                                    platform.
                                </p>
                            </div>
                            <div class="flex flex-col gap-2 min-[400px]:flex-row">
                                <NuxtLink to="/signup"
                                    class="inline-flex h-10 items-center justify-center rounded-md boder px-8 text-sm font-medium transition-bg border-love/40 bg-love/10 text-love hover:bg-love/15 active:bg-love/25 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring">
                                    Sign Up
                                </NuxtLink>
                                <NuxtLink to="/login"
                                    class="inline-flex h-10 items-center justify-center rounded-md border px-8 text-sm font-medium transition-bg hover:bg-muted/10 active:bg-muted/15 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring">
                                    Log In
                                </NuxtLink>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
            <section class="flex justify-center w-full py-12 md:py-24 lg:py-32">
                <div class="container px-4 md:px-6">
                    <div class="flex flex-col items-center justify-center space-y-4 text-center">
                        <div class="space-y-2">
                            <div class="inline-block rounded-lg bg-subtle/20 px-3 py-1 text-sm">
                                Secure File Hosting
                            </div>
                            <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">Store Your Files Safely</h2>
                            <p
                                class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                                Our platform uses industry-leading encryption to keep your files secure and protected.
                                Access your
                                data from any device, anywhere.
                            </p>
                        </div>
                    </div>
                    <div class="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-2 lg:gap-12">
                        <img src="/placeholder.svg" width="550" height="310" alt="Image"
                            class="mx-auto aspect-video overflow-hidden rounded-xl object-cover object-center sm:w-full lg:order-last" />
                        <div class="flex flex-col justify-center space-y-4">
                            <ul class="grid gap-6">
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Secure Storage</h3>
                                        <p class="text-muted-foreground">
                                            Your files are encrypted and stored in a secure environment.
                                        </p>
                                    </div>
                                </li>
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Easy Sharing</h3>
                                        <p class="text-muted-foreground">
                                            Share files with your team or clients with just a few clicks.
                                        </p>
                                    </div>
                                </li>
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Cross-Device Access</h3>
                                        <p class="text-muted-foreground">Access your files from any device, anywhere
                                            in the world.</p>
                                    </div>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </section>
            <section class="flex justify-center w-full py-12 md:py-24 lg:py-32 bg-subtle/10">
                <div class="container px-4 md:px-6">
                    <div class="flex flex-col items-center justify-center space-y-4 text-center">
                        <div class="space-y-2">
                            <div class="inline-block rounded-lg bg-subtle/20 px-3 py-1 text-sm">
                                File Management
                            </div>
                            <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">Effortless File Management</h2>
                            <p
                                class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                                Our intuitive file browser makes it easy to upload, download, and manage your files.
                                Organize your data with folders and tags, and share files with your team or clients.
                            </p>
                        </div>
                    </div>
                    <div class="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-2 lg:gap-12">
                        <img src="/placeholder.svg" width="550" height="310" alt="Image"
                            class="mx-auto aspect-video overflow-hidden rounded-xl object-cover object-center sm:w-full lg:order-last" />
                        <div class="flex flex-col justify-center space-y-4">
                            <ul class="grid gap-6">
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Intuitive File Browser</h3>
                                        <p class="text-muted-foreground">
                                            Easily navigate and manage your files with our user-friendly interface.
                                        </p>
                                    </div>
                                </li>
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Folder Organization</h3>
                                        <p class="text-muted-foreground">
                                            Create and organize your files into folders for better structure.
                                        </p>
                                    </div>
                                </li>
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Sharing and Collaboration</h3>
                                        <p class="text-muted-foreground">
                                            Easily share files with your team or clients
                                            and collaborate on projects.
                                        </p>
                                    </div>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </section>
            <section class="flex justify-center w-full py-12 md:py-24 lg:py-32">
                <div class="container px-4 md:px-6">
                    <div class="flex flex-col items-center justify-center space-y-4 text-center">
                        <div class="space-y-2">
                            <div class="inline-block rounded-lg bg-subtle/20 px-3 py-1 text-sm">
                                Storage Plans
                            </div>
                            <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">Store 10GB For Free</h2>
                            <p
                                class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                                Our storage plans give you flexibility to store whatever you need, whether it's
                                presentations and documents, or pictures and vides, there's a plan for you.
                            </p>
                        </div>
                    </div>
                    <div class="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-2 lg:gap-12">
                        <img src="/placeholder.svg" width="550" height="310" alt="Image"
                            class="mx-auto aspect-video overflow-hidden rounded-xl object-cover object-center sm:w-full lg:order-last" />
                        <div class="flex flex-col justify-center space-y-4">
                            <ul class="grid gap-6">
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Lorem Ipsum</h3>
                                        <p class="text-muted-foreground">
                                            Lorem ipsum amet dolar sit.
                                        </p>
                                    </div>
                                </li>
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Lorem Ipsum</h3>
                                        <p class="text-muted-foreground">
                                            Lorem ipsum amet dolar sit.
                                        </p>
                                    </div>
                                </li>
                                <li>
                                    <div class="grid gap-1">
                                        <h3 class="text-xl font-bold">Lorem Ipsum</h3>
                                        <p class="text-muted-foreground">Lorem ipsum amet dolar sit.</p>
                                    </div>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </section>
        </main>
        <Footer />
    </div>
</template>


<style></style>