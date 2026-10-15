<script setup>
import { ref, onMounted } from 'vue';
import { apiGet, apiPost } from '@/utils/api.js';
import { useUser } from '@/composables/useUser.js';
import { DateToLocale } from '@/utils/utils.js';
import StatusIcon from '@/components/StatusIcon.vue';

const { user } = useUser();
const followers = ref([]);
const following = ref([]);
const loading = ref(true);

async function loadFollows() {
    const data = await apiGet('/api/follows');
    if (data) {
        const follows = data.Follows ?? {};
        followers.value = follows.Followers ?? [];
        following.value = follows.Following ?? [];
    }
    loading.value = false;
}

async function followRequest(follow) {
    const payload = new FormData();
    payload.append('user-id', follow);
    await apiPost('/api/follow', payload);
    await loadFollows();
}

async function acceptRequest(follow) {
    await apiPost(`/api/follow/accept/${follow.Id}`);
    await loadFollows();
}

async function declineRequest(follow) {
    await apiPost(`/api/follow/deny/${follow.Id}`);
    await loadFollows();
}

async function unfollow(follow) {
    await apiPost(`/api/unfollow/${follow.ToUserId}`);
    await loadFollows();
}

onMounted(loadFollows);
</script>

<template>
    <div class="flex flex-column gap-15 padding-sides-1vw overflow-y-auto follows-view">
        <h2>My follows</h2>
        <template v-if="loading">Loading...</template>
        <template v-else>
            <h3>Followers</h3>
            <p v-if="followers.length === 0">No followers or incoming requests.</p>
            <template v-else>
                <div
                    v-for="follow in followers"
                    :key="follow.Id"
                    class="flex flex-row flex-wrap gap-8"
                >
                    <StatusIcon :status="follow.Status" />
                    <span>{{ follow.Status }}</span>
                    <router-link :to="`/profile/view/${follow.FromUserId}`">
                        <strong>{{ follow.Username }}</strong>
                    </router-link>
                    <span>{{ DateToLocale(follow.Timestamp) }}</span>
                    <div>
                        <template v-if="follow.Status === 'pending'">
                            <button type="button" @click="acceptRequest(follow)">Accept</button>
                            <button type="button" @click="declineRequest(follow)">Decline</button>
                        </template>
                        <template v-else-if="follow.Status === 'accepted'">
                            <button type="button" @click="followRequest(follow.FromUserId)">
                                Follow back
                            </button>
                        </template>
                    </div>
                </div>
            </template>
            <h3>Following</h3>
            <p v-if="following.length === 0">No follows or outgoing requests.</p>
            <template v-else>
                <div v-for="follow in following" :key="follow.Id" class="flex flex-row gap-8">
                    <StatusIcon :status="follow.Status" />
                    <span>{{ follow.Status }}</span>
                    <router-link :to="`/profile/view/${follow.ToUserId}`">
                        <strong>{{ follow.Username }}</strong>
                    </router-link>
                    <span>{{ DateToLocale(follow.Timestamp) }}</span>
                    <div>
                        <template v-if="follow.Status === 'accepted'">
                            <button type="button" @click="unfollow(follow)">Unfollow</button>
                        </template>
                        <template
                            v-else-if="follow.Status === 'unfollowed' || follow.Status === 'denied'"
                        >
                            <button type="button" @click="followRequest(follow.ToUserId)">
                                Re-follow
                            </button>
                        </template>
                    </div>
                </div>
            </template>
        </template>
    </div>
</template>

<style scoped>
.follow-section {
    margin-bottom: 1.5rem;
}

.follow-section h3 {
    border-bottom: 1px solid #ddd;
    padding-bottom: 0.25rem;
}

.follow-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 0;
}

.follow-row span {
    color: #666;
}

.follow-row .status {
    font-style: italic;
}

.follow-row button {
    margin-left: auto;
}

.empty {
    color: #888;
    font-style: italic;
}
</style>
