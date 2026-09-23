<script setup>
import { ref, onMounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiFetch } from '@/utils/api.js';
import { DateToLocale } from '@/utils/utils.js';

const { user } = useUser();

const groups = ref([]);

onMounted(async () => {
    const data = await apiFetch('/api/groups');
    if (data) {
        groups.value = data.Groups ?? [];
    }
});
</script>

<template>
    <div class="groups-browser">
        <h2>Groups</h2>
        <template v-if="user.LoggedIn">
            <div class="group" v-for="group in groups" :key="group.Id">
                <router-link :to="`/group/view/${group.Id}`"><h3>{{ group.Title }}</h3></router-link>
                <p class="group-description">{{ group.Description }}</p>
                <p class="group-owner">
                Owned by 
                <router-link :to="`/profile/view/${group.OwnerUserId}`">
                    <strong>{{ group.OwnerUsername }}</strong>
                </router-link>
                on<em>({{ DateToLocale(group.Timestamp) }})</em>
            </p>
            </div>
            <p v-if="groups.length === 0">No groups yet. Create the first one!</p>
        </template>
        <p v-else>
            You must be logged in to browse groups.
            <router-link to="/user/login">Log in</router-link>
        </p>
    </div>
</template>

<style scoped></style>
