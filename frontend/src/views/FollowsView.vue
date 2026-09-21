<script setup>
import { ref, onMounted } from 'vue';
import { apiFetch, apiPost } from '@/utils/api.js';
import { DateToLocale } from '@/utils/utils.js';

const followers = ref([]);
const following = ref([]);
const pendingIncoming = ref([]);
const pendingOutgoing = ref([]);
const loading = ref(true);

async function loadFollows() {
    const data = await apiFetch('/api/follows');
    if (data) {
        const follows = data.Follows ?? {};
        followers.value = follows.Followers ?? [];
        following.value = follows.Following ?? [];
        pendingIncoming.value = follows.PendingIncoming ?? [];
        pendingOutgoing.value = follows.PendingOutgoing ?? [];
    }
    loading.value = false;
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
    <div class="follows-view">
        <h2>My follows</h2>

        <template v-if="!loading">
            <section class="follow-section pending-incoming">
                <h3>Follow requests</h3>
                <template v-if="pendingIncoming.length > 0">
                    <div v-for="follow in pendingIncoming" :key="'pending-in-' + follow.Id" class="follow-row">
                        <router-link :to="`/user/${follow.FromUserId}`">
                            <strong>{{ follow.Username }}</strong>
                        </router-link>
                        <button type="button" @click="acceptRequest(follow)">Accept</button>
                        <button type="button" @click="declineRequest(follow)">Decline</button>
                    </div>
                </template>
                <p v-else class="empty">No pending follow requests.</p>
            </section>

            <section class="follow-section followers">
                <h3>Followers</h3>
                <template v-if="followers.length > 0">
                    <div v-for="follow in followers" :key="'follower-' + follow.Id" class="follow-row">
                        <router-link :to="`/user/${follow.FromUserId}`">
                            <strong>{{ follow.Username }}</strong>
                        </router-link>
                        <span>{{ DateToLocale(follow.Timestamp) }}</span>
                    </div>
                </template>
                <p v-else class="empty">No followers yet.</p>
            </section>

            <section class="follow-section following">
                <h3>Following</h3>
                <template v-if="following.length > 0">
                    <div v-for="follow in following" :key="'following-' + follow.Id" class="follow-row">
                        <router-link :to="`/user/${follow.ToUserId}`">
                            <strong>{{ follow.Username }}</strong>
                        </router-link>
                        <span>{{ DateToLocale(follow.Timestamp) }}</span>
                        <button type="button" @click="unfollow(follow)">Unfollow</button>
                    </div>
                </template>
                <p v-else class="empty">You are not following anyone yet.</p>
            </section>

            <section class="follow-section pending-outgoing">
                <h3>Outgoing requests</h3>
                <template v-if="pendingOutgoing.length > 0">
                    <div v-for="follow in pendingOutgoing" :key="'pending-out-' + follow.Id" class="follow-row">
                        <router-link :to="`/user/${follow.ToUserId}`">
                            <strong>{{ follow.Username }}</strong>
                        </router-link>
                        <span class="status">Awaiting response</span>
                    </div>
                </template>
                <p v-else class="empty">No pending outgoing requests.</p>
            </section>
        </template>
        <p v-else>Loading...</p>
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