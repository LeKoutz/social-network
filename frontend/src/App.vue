<!-- App.vue -->
<script setup>
import { ref, onMounted } from 'vue';
import Topbar from './components/Topbar.vue';
import Footer from './components/Footer.vue';
import { apiFetch } from './utils/api.js';
import Alerts from './components/Alerts.vue';
import { useAppData } from './composables/useAppData.js';
import { useNotifications } from './composables/useNotifications.js';
import { useUser } from './composables/useUser.js';

const loaded = ref(false);
const { setAppData } = useAppData();
const { setNotifications } = useNotifications();
const { setUser } = useUser();

onMounted(async () => {
    const data = await apiFetch('/api');
    if (data) {
        setAppData(data);
        loaded.value = true;
        if (data.User.LoggedIn) {
            setUser(data.User);
            setNotifications(data.User.Notifications);
        }
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
