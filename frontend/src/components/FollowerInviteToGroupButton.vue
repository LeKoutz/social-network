<script setup>
import { useGroupInvite } from '@/composables/useGroupInvite.js';
import { ref } from 'vue';

const { loading, error, success, sendInvite } = useGroupInvite();

const props = defineProps({
    follower: { type: Object, default: null },
    group: { type: Object, default: null },
});

async function handleInviteClick(userId) {
    if (!props.group?.Id || !userId) return;
    console.log(props.group.Id, userId);
    await sendInvite(props.group.Id, userId);
}
</script>

<template>
    <button @click="handleInviteClick(follower.UserId)" :disabled="loading" class="invite-button">
        <span v-if="loading">Sending...</span>
        <span v-else-if="success">Neat</span>
        <span v-else>Invite {{ follower.Username }}</span>
    </button>
</template>
