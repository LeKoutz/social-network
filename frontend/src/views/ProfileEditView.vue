<script setup>
import { ref, onMounted } from 'vue';
import { apiGet } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';
import ProfileEditForm from '@/components/ProfileEditForm.vue';

const profile = ref(null);
const loading = ref(true);
const error = ref('');

onMounted(async () => {
    try {
        const data = await apiGet('/api/profile/edit');
        profile.value = firstOrNull(data?.Profiles);

        if (!profile.value) {
            error.value = 'Could not load your profile.';
        }
    } catch {
        error.value = 'Could not load your profile. Please try again.';
    } finally {
        loading.value = false;
    }
});
</script>

<template>
    <div class="container padding-sides-1vw">
        <h2>Edit your profile</h2>
        <p v-if="loading">Loading…</p>
        <p v-else-if="error">{{ error }}</p>
        <ProfileEditForm v-else-if="profile" :profile="profile" />
    </div>
</template>
