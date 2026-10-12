<script setup>
const props = defineProps({
    status: { type: String, default: null },
});
console.log(status);
</script>

<template>
    <div class="flex flex-row align-items-center">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
            <template v-if="status === 'pending'">
                <circle cx="12" cy="12" r="8" fill="#F59E0B" />
            </template>
            <template v-else-if="status === 'accepted'">
                <circle cx="12" cy="12" r="8" fill="#10B981" />
            </template>
            <template v-else-if="status === 'denied'">
                <circle cx="12" cy="12" r="8" fill="#EF4444" />
            </template>
            <template v-else-if="status === 'left'">
                <circle cx="12" cy="12" r="8" fill="#6B7280" />
            </template>
            <template v-else-if="status === 'removed'">
                <circle cx="12" cy="12" r="8" fill="#7C3AED" />
            </template>
        </svg>
    </div>
</template>

<style scoped></style>
