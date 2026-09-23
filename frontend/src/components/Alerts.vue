<script setup>
import { useAlerts } from '../composables/useAlerts.js';

const { alert, clear } = useAlerts();
</script>

<template>
    <div v-if="alert" class="alert-space">
        <div :class="['alert', `alert-type-${alert.type}`]">
            <span>{{ alert.message }}</span>
            <button @click="clear">✗</button>
        </div>
    </div>
</template>

<style scoped></style>
