<script setup lang="ts">
definePageMeta({
    middleware: ["auth", "admin"],
    layout: "admin"
});

let {data: systemStatusData, refresh} = await useFetch("/api/admin/status")

const calculateTimeSince = (time) => {
    const now = new Date();
    const date = new Date(time);
    const diffInSeconds = Math.floor((now - date) / 1000);

    const days = Math.floor(diffInSeconds / (3600 * 24));
    const hours = Math.floor((diffInSeconds % (3600 * 24)) / 3600);
    const minutes = Math.floor((diffInSeconds % 3600) / 60);
    const seconds = diffInSeconds % 60;

    // Constructing the output based on non-zero values
    const timeParts = [];
    if (days > 0) timeParts.push(`${days} days`);
    if (hours > 0 || days > 0) timeParts.push(`${hours} hours`);
    if (minutes > 0 || hours > 0 || days > 0) timeParts.push(`${minutes} minutes`);
    timeParts.push(`${seconds} seconds`);

    return timeParts.join(', ');
}

let uptime = ref(calculateTimeSince(systemStatusData.value.uptime));
let lastGcTime = ref(calculateTimeSince(systemStatusData.value.last_gc_time));

let systemStatusInterval;
let timeInterval;

const updateTime = () => {
    uptime.value = calculateTimeSince(systemStatusData.value.uptime);
    lastGcTime.value = calculateTimeSince(systemStatusData.value.last_gc_time)
};

onMounted(() => {
    updateTime();

    systemStatusInterval = setInterval(async () => {
        refresh()
    }, 5000);

    timeInterval = setInterval(updateTime, 1000);
});

onUnmounted(() => {
    clearInterval(systemStatusInterval);
    clearInterval(timeInterval);
});
</script>

<template>
    <div class="w-full overflow-hidden rounded-md border h-fit text-[15px]">
        <h4 class="bg-surface px-3.5 py-3 border-b">System Status</h4>
        <div class="p-3.5 text-sm">
            <dl class="flex-wrap">
                <dt>Server Uptime</dt>
                <dd>{{ uptime }}</dd>
                <dt>Current Goroutine</dt>
                <dd>{{ systemStatusData.num_goroutine }}</dd>
                <hr />
                <dt>Current Memory Usage</dt>
                <dd>{{ systemStatusData.cur_mem_usage }}</dd>
                <dt>Total Memory Allocated</dt>
                <dd>{{ systemStatusData.total_mem_usage }}</dd>
                <dt>Memory Obtained</dt>
                <dd>{{ systemStatusData.mem_obtained }}</dd>
                <dt>Pointer Lookup Times</dt>
                <dd>{{ systemStatusData.ptr_lookup_times }}</dd>
                <dt>Memory Allocations</dt>
                <dd>{{ systemStatusData.mem_allocations }}</dd>
                <dt>Memory Frees</dt>
                <dd>{{ systemStatusData.mem_frees }}</dd>
                <hr />
                <dt>Current Heap Usage</dt>
                <dd>{{ systemStatusData.cur_heap_usage }}</dd>
                <dt>Heap Memory Obtained</dt>
                <dd>{{ systemStatusData.heap_mem_obtained }}</dd>
                <dt>Heap Memory Idle</dt>
                <dd>{{ systemStatusData.heap_mem_idle }}</dd>
                <dt>Heap Memory In Use</dt>
                <dd>{{ systemStatusData.heap_mem_inuse }}</dd>
                <dt>Heap Memory Released</dt>
                <dd>{{ systemStatusData.heap_mem_release }}</dd>
                <dt>Heap Objects</dt>
                <dd>{{ systemStatusData.heap_objects }}</dd>
                <hr />
                <dt>Bootstrap Stack Usage</dt>
                <dd>{{ systemStatusData.bootstrap_stack_usage }}</dd>
                <dt>Stack Memory Obtained</dt>
                <dd>{{ systemStatusData.stack_mem_obtained }}</dd>
                <dt>MSpan Structures Usage</dt>
                <dd>{{ systemStatusData.mspan_structures_usage }}</dd>
                <dt>MSpan Structures Obtained</dt>
                <dd>{{ systemStatusData.mspan_structures_obtained }}</dd>
                <dt>MCache Structures Usage</dt>
                <dd>{{ systemStatusData.mcache_structures_usage }}</dd>
                <dt>MCache Structures Obtained</dt>
                <dd>{{ systemStatusData.mcache_structures_obtained }}</dd>
                <dt>Profiling Bucket Hash Table Obtained</dt>
                <dd>{{ systemStatusData.buck_hash_sys }}</dd>
                <dt>GC Metadata Obtained</dt>
                <dd>{{ systemStatusData.gc_sys }}</dd>
                <dt>Other System Allocation Obtained</dt>
                <dd>{{ systemStatusData.other_sys }}</dd>
                <hr />
                <dt>Next GC Recycle</dt>
                <dd>{{ systemStatusData.next_gc }}</dd>
                <dt>Since Last GC Time</dt>
                <dd>{{ lastGcTime }}</dd>
                <dt>Total GC Pause</dt>
                <dd>{{ systemStatusData.pause_total_ns }}</dd>
                <dt>Last GC Pause</dt>
                <dd>{{ systemStatusData.pause_ns }}</dd>
                <dt>GC Times</dt>
                <dd>{{ systemStatusData.num_gc }}</dd>
            </dl>
        </div>
    </div>
</template>

<style>
dl {
    display: flex;
    flex-wrap: wrap;
}

dt {
    font-weight: 600;
    width: 300px;
    max-width: calc(100% - 100px - 1em);
    padding-top: 5px;
    padding-bottom: 5px;
}

dd {
    padding-top: 5px;
    padding-bottom: 5px;
    width: calc(100% - 300px);
    min-width: 100px;
}

hr {
    width: 100%;
    margin-top: 4px;
    margin-bottom: 4px;
}
</style>