<script setup>
import { ref, onMounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiFetch, apiPost } from '@/utils/api.js';
import { DateToLocale, firstOrNull } from '@/utils/utils.js';

const { user } = useUser();

const groups = ref([]);

onMounted(async () => {
    const data = await apiFetch('/api/groups');
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
    <div class="flex flex-column gap-5 padding-sides-1vw height-100 overflow-y-auto">
        <h2>Groups</h2>
        <template v-if="user.LoggedIn">
            <div
                class="flex flex-row flex-wrap align-items-center gap-2 border-1px-dashed-accent-line padding-sides-1vw piped-elements"
                v-for="group in groups"
                :key="group.Id"
            >
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
                    Created by
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
