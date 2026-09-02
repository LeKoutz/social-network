<!-- App.vue -->
<script setup>
import { ref, onMounted } from 'vue';
import Topbar from './components/Topbar.vue';
import Footer from './components/Footer.vue';
import { apiFetch } from './utils/api.js';
import Alerts from './components/Alerts.vue';
import { useAppData } from './composables/useAppData.js';

const loaded = ref(false);
const { setAppData } = useAppData();

onMounted(async () => {
    const data = await apiFetch('/api');
    if (data) {
        setAppData(data);
        loaded.value = true
    }
});
</script>

<template>
    <template v-if="loaded">
        <Topbar />
        <Alerts />
        <router-view />
        <Footer />
    </template>
</template>

<style scoped></style>
