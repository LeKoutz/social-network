<script setup>
import { useGroupInvite } from '@/composables/useGroupInvite.js';
import FollowerInviteToGroupButton from '@/components/FollowerInviteToGroupButton.vue';
import { ref } from 'vue';

const { loading, error, success, sendInvite } = useGroupInvite();

const props = defineProps({
    followers: { type: Array, default: null },
    group: { type: Object, default: null },
});

async function handleInviteClick(userId) {
    if (!props.group?.Id || !userId) return;
    console.log(props.group.Id, userId);
    await sendInvite(props.group.Id, userId);
}
</script>

<template>
    <div v-if="followers" class="flex flex-column">
        <div v-if="error">{{ error }}</div>
        <div v-else v-for="follower in followers" :key="follower.Id" class="flex flex-column">
            <FollowerInviteToGroupButton :follower="follower" :group="group" />
        </div>
    </div>
</template>
<style scoped></style>
