<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiGet } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';
import ProfileShow from '@/components/ProfileShow.vue';
import { useProfile } from '@/composables/useProfile.js';

const route = useRoute();

const { profile, loading, error, loadProfile } = useProfile();

onMounted(() => {
    loadProfile(route.params.id);
});

watch(
    () => route.params.id,
    (newId) => {
        if (newId) loadProfile(newId);
    },
);
</script>

<template>
    <div class="flex flex-column gap-20 overflow-y-auto padding-sides-1vw">
        <div v-if="loading">Loading…</div>
        <div v-else-if="error">{{ error }}</div>
        <ProfileShow v-else-if="profile" :profile="profile" />
    </div>
</template>

<style scoped></style>
