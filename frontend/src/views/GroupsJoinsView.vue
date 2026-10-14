<script setup>
import { ref, onMounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiGet, apiPost } from '@/utils/api.js';
import { DateToLocale, firstOrNull } from '@/utils/utils.js';

const { user } = useUser();

const groups = ref([]);

onMounted(async () => {
    const data = await apiGet('/api/groups');
    if (data) {
        groups.value = data.Groups ?? [];
    }
});

async function join(user) {
    const group = firstOrNull(groups.value);
    const payload = new FormData();
    payload.append('group-id', group.Id);
    payload.append('to-user-id', user.Id);
    const data = await apiPost('/api/group/join', payload);
    console.log(data);
}
</script>

<template>
    <div class="flex flex-column gap-5 padding-sides-1vw overflow-y-auto">
        <h2>Join requests</h2>
        <template v-if="user.LoggedIn">
            <div
                class="flex flex-row flex-wrap align-items-center gap-2 border-1px-dashed-accent-line padding-sides-1vw piped-elements"
                v-for="group in groups"
                :key="group.Id"
            >
                <div class="flex flex-row align-items-center">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        width="24"
                        height="24"
                    >
                        <circle cx="12" cy="12" r="8" fill="#F59E0B" />
                    </svg>
                    <span>pending</span>
                </div>

                <div class="flex flex-row align-items-center">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        width="24"
                        height="24"
                    >
                        <!-- accepted - green -->
                        <circle cx="12" cy="12" r="8" fill="#10B981" />
                    </svg>
                    <span>accepted</span>
                </div>

                <div class="flex flex-row align-items-center">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        width="24"
                        height="24"
                    >
                        <!-- denied - red -->
                        <circle cx="12" cy="12" r="8" fill="#EF4444" />
                    </svg>
                    <span>denied</span>
                </div>

                <div class="flex flex-row align-items-center">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        width="24"
                        height="24"
                    >
                        <!-- left - gray -->
                        <circle cx="12" cy="12" r="8" fill="#6B7280" />
                    </svg>
                    <span>left</span>
                </div>

                <div class="flex flex-row align-items-center">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        width="24"
                        height="24"
                        alt="removed"
                    >
                        <!-- removed - dark purple -->
                        <circle cx="12" cy="12" r="8" fill="#7C3AED" />
                    </svg>
                    <span>removed</span>
                </div>

                <div>
                    🔴 Denied/Removed 🟠 Pending 🟡 ??? 🟤 Left
                    <button>❌ Deny</button>
                    <button>✅ Accept</button>
                </div>
                <div>
                    <span v-if="group.Owner">
                        🟣
                        <button disabled>Owner</button>
                    </span>
                    <span v-else-if="group.Member">
                        🟢
                        <button @click="whatever">Member</button>
                    </span>
                    <span v-else>
                        ⚫
                        <button @click="join(user)">Join</button>
                    </span>
                </div>
                <router-link class="font-size-larger" :to="`/group/view/${group.Id}`">
                    {{ group.Title }}
                </router-link>
                <p v-if="group.Description.length !== 0" class="group-description">
                    {{ group.Description }}
                </p>
                <p class="group-owner">
                    Owned by
                    <router-link :to="`/profile/view/${group.OwnerUserId}`">
                        <strong>{{ group.OwnerUsername }}</strong>
                    </router-link>
                    on
                    <em>({{ DateToLocale(group.Timestamp) }})</em>
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
