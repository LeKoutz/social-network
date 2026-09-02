<script setup>
import { useAlerts } from '../composables/useAlerts.js';

const { alert } = useAlerts();
</script>

<template>
    <template v-if="alert">
        <fieldset class="alert">
            <legend>{{ alert.type }}</legend>
            <p>{{ alert.message }}</p>
        </fieldset>
    </template>
</template>

<style scoped></style>
