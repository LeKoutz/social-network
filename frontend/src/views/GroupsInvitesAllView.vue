<script setup>
import { ref, onMounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiFetch, apiPost } from '@/utils/api.js';
import { DateToLocale, firstOrNull } from '@/utils/utils.js';
import JoinStatusIcon from '@/components/JoinStatusIcon.vue';

const { user } = useUser();

const invites = ref([]);

onMounted(async () => {
    const data = await apiFetch('/api/groups/invites');
    if (data) {
        invites.value = data.GroupInvitations ?? [];
    }
});

async function acceptInvite(invite){
    console.log(invite);
    const payload = new FormData();
    //payload.append('group-id', group.Id)
    payload.append('to-user-id', user.Id)
    // const data = await apiPost('/api/group/join', payload);
    // console.log(data);
}
</script>

<template>
    <div class="flex flex-column gap-5 padding-sides-1vw height-100 overflow-y-auto">
        <h2>Groups Invites</h2>
        <template v-if="user.LoggedIn">
            <div class="flex flex-row flex-wrap align-items-center gap-8" v-for="invite in invites" :key="invite.Id" :id="invite.Id">
                <span><JoinStatusIcon :status="`${invite.Status}`"/></span>
                <span><router-link :to="`/profile/view/${invite.InvitedBy}`">{{ invite.Username }}</router-link></span>
                <span>invited you to join:</span>
                <span>
                    <router-link :to="`/group/view/${invite.GroupId}`">
                        {{ invite.GroupTitle }}
                    </router-link>
                </span>
                <span>on {{ DateToLocale(invite.Timestamp) }}</span>
                <template v-if="invite.Status === 'pending'">
                    <button @click="acceptInvite(invite.Id)">
                        ✅ Accept
                    </button>
                    <button>
                        ❌ Deny
                    </button>
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

