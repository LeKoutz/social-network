<script setup>
import { useAppData } from '@/composables/useAppData';

const { version } = useAppData();
</script>

<template>
    <div class="footer">
        <p>Version: {{ version }}</p>
    </div>
</template>

<style scoped>
.footer {
    background-color: var(--bg-base);
    display: flex;
    align-items: center;
    padding: 24px;
    box-sizing: border-box;
    font-size: 13px;
    line-height: 1.2;
    color: var(--text-muted);
    border-top: 1px solid var(--accent-line);
    max-height: 55px;
    height: 55px;
}
</style>
