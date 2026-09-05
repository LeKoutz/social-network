<script setup>
import { useAlerts } from '../composables/useAlerts.js';

const { alert, clear } = useAlerts();
</script>

<template>
    <template v-if="alert">
        <fieldset class="alert">
            <legend>{{ alert.type }}</legend>
            <p>{{ alert.message }}</p>
            <button @click="clear">Dismiss</button>
        </fieldset>
    </template>
</template>

<style scoped></style>
