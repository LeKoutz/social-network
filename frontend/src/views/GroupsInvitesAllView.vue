<script setup>
import { ref, onMounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiGet, apiPost } from '@/utils/api.js';
import { DateToLocale, firstOrNull } from '@/utils/utils.js';
import StatusIcon from '@/components/StatusIcon.vue';
import { useGroupInvite } from '@/composables/useGroupInvite.js';

const { user } = useUser();

const invites = ref([]);

onMounted(async () => {
    const data = await apiGet('/api/groups/invites');
    if (data) {
        invites.value = data.GroupInvitations ?? [];
    }
});

const { loading, error, success, acceptInvite } = useGroupInvite();

async function handleAcceptGroupInviteButton(inviteId) {
    await acceptInvite(inviteId);
}
</script>

<template>
    <div class="flex flex-column gap-5 padding-sides-1vw overflow-y-auto">
        <h2>Groups Invites</h2>
        <template v-if="user.LoggedIn">
            <div
                class="flex flex-row flex-wrap align-items-center gap-8"
                v-for="invite in invites"
                :key="invite.Id"
                :id="invite.Id"
            >
                <span><StatusIcon :status="`${invite.Status}`" /></span>
                <span>
                    <router-link :to="`/profile/view/${invite.InvitedBy}`">
                        {{ invite.Username }}
                    </router-link>
                </span>
                <span>invited you to join:</span>
                <span>
                    <router-link :to="`/group/view/${invite.GroupId}`">
                        {{ invite.GroupTitle }}
                    </router-link>
                </span>
                <span>on {{ DateToLocale(invite.Timestamp) }}</span>
                <template v-if="invite.Status === 'pending'">
                    <button @click="handleAcceptGroupInviteButton(invite.Id)">✅ Accept</button>
                    <button>❌ Deny</button>
                </template>
            </div>
        </template>
        <p v-else>
            You must be logged in to browse groups.
            <router-link to="/user/login">Log in</router-link>
        </p>
    </div>
</template>

<style scoped></style>
