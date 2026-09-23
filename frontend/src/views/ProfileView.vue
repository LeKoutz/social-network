<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';
import ProfileShow from '@/components/ProfileShow.vue';

const route = useRoute();
const profile = ref(null);
const loading = ref(true);
const error = ref(null);

async function loadProfile(id) {
	loading.value = true;
	error.value = null;
	const data = await apiFetch(`/api/profile/view/${id}`);
	if (data) {
		profile.value = firstOrNull(data.Profiles);
		if (!profile.value) error.value = 'Profile not found';
	}
	loading.value = false;
}

onMounted(() => loadProfile(route.params.id));
watch(() => route.params.id, (newId) => { if (newId) loadProfile(newId); });
</script>

<template>
	<div class="container">
		<div v-if="loading">Loading…</div>
		<div v-else-if="error">{{ error }}</div>
		<ProfileShow v-else-if="profile" :profile="profile" />
	</div>
</template>

<style scoped></style>