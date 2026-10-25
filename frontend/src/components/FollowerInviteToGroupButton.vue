<script setup>
import { useGroupInvite } from '@/composables/useGroupInvite.js';
import Spacer from '@/components/Spacer.vue';
import { ref } from 'vue';
import CheckIcon from '@/components/Icons/CheckIcon.vue';
import AvatarUsernameRow from '@/components/AvatarUsernameRow.vue';
import GroupInviteSuccessButton from '@/components/GroupInviteSuccessButton.vue';
import GroupInviteButton from '@/components/GroupInviteButton.vue';

const { loading, error, success, sendInvite } = useGroupInvite();

const props = defineProps({
    follower: { type: Object, default: null },
    group: { type: Object, default: null },
});

async function handleInviteClick(userId) {
    if (!props.group?.Id || !userId) return;
    await sendInvite(props.group.Id, userId);
}
</script>

<template>
    <div
        class="position-relative width-100 height-auto display-flex align-items-center align-content-stretch justify-content-start column-gap-10px padding-block-start-6px padding-block-end-6px flex-direction-row flex-wrap-nowrap"
    >
        <AvatarUsernameRow :username="follower.Username" />
        <Spacer />
        <template v-if="loading">Sending...</template>
        <template v-else-if="success"><GroupInviteSuccessButton /></template>
        <template v-else>
            <GroupInviteButton @click="handleInviteClick(follower.UserId)" :disabled="loading" />
        </template>
    </div>
</template>
<style scoped></style>
