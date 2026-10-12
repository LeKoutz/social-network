<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useUser } from '@/composables/useUser.js';
import { apiGet, apiPost, apiDelete } from '@/utils/api.js';
import { DateToLocale, firstOrNull } from '@/utils/utils.js';
import JoinStatusIcon from '@/components/JoinStatusIcon.vue';

const route = useRoute();
const { user } = useUser();

const joins = ref([]);

onMounted(async () => {
    const data = await apiGet('/api/groups/joins');
    if (data) {
        joins.value = data.GroupInvitations ?? [];
    }
    console.log(data);
});

async function withdraw(id) {
    const payload = new FormData();
    payload.append('id', id);
    const data = await apiDelete('/api/group/join/withdraw', payload);
    if (data) router.push('/groups/joins');
}
</script>

<template>
    <div class="flex flex-column gap-5 padding-sides-1vw height-100 overflow-y-auto">
        <template v-if="user.LoggedIn">
            <h2>Join requests</h2>
            <div class="flex flex-row gap-5" v-for="jr in joins" :key="jr.Id">
                <template v-if="jr.InvitedBy === 0">
                    <span><JoinStatusIcon :status="`${jr.Status}`" /></span>
                    <span>Requested to join</span>
                    <span>
                        <router-link :to="`/group/view/${jr.GroupId}`">
                            {{ jr.GroupTitle }}
                        </router-link>
                    </span>
                    <span>on {{ DateToLocale(jr.Timestamp) }}</span>
                    <template v-if="jr.Status === 'pending'">
                        <button @click="withdraw(jr.Id)">❌ Withdraw request</button>
                    </template>
                </template>
            </div>
            <p v-if="joins.length === 0">
                No join requests. Find a
                <router-link to="/groups">group</router-link>
                to join!
            </p>
        </template>
        <p v-else>
            You must be logged in to browse groups.
            <router-link to="/user/login">Log in</router-link>
        </p>
    </div>
</template>

<style scoped></style>
