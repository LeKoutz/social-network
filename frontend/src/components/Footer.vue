<script setup>
import { useAppData } from '@/composables/useAppData';

const { version } = useAppData();
</script>

<template>
    <div class="footer">
        <p>Version: {{ version }}</p>
    </div>
</template>

<style scoped></style>
