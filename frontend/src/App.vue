<script setup>
import { ref, onMounted } from 'vue';
import Topbar from './components/Topbar.vue';
import Footer from './components/Footer.vue';
import { apiGet } from './utils/api.js';
import Alerts from './components/Alerts.vue';
import { useAppData } from './composables/useAppData.js';
import { useNotifications } from './composables/useNotifications.js';
import { useUser } from './composables/useUser.js';
import { useCategories } from './composables/useCategories.js';
import LeftSidebar from './components/LeftSidebar.vue';
import ChatPanel from './components/ChatPanel.vue';

const loaded = ref(false);
const { setAppData } = useAppData();
const { setNotifications } = useNotifications();
const { user, setUser, resolveUserLoaded } = useUser();
const { setCategories } = useCategories();

onMounted(async () => {
    try {
        const data = await apiGet('/api');
        if (data) {
            setAppData(data);
            setCategories(data.Categories);
            loaded.value = true;
            if (data.User.LoggedIn) {
                setUser(data.User);
                setNotifications(data.User.Notifications);
            }
        }
    } finally {
        resolveUserLoaded();
    }
});
</script>

<template>
    <template v-if="loaded">
        <Topbar />
        <Alerts />
        <div class="flex flex-row overflow-y-auto height-100">
            <LeftSidebar v-if="user.LoggedIn" />
            <router-view class="flex-grow-1" />
            <ChatPanel v-if="user.LoggedIn" />
        </div>
        <Footer />
    </template>
</template>

<style scoped></style>
