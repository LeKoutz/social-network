<script setup>
import { ref, onMounted, watch } from 'vue';
import FollowersList from '@/components/FollowersList.vue';
import { useProfile } from '@/composables/useProfile.js';

const props = defineProps({
    group: { type: Object, default: null },
    user: { type: Object, default: null },
});

const { profile, loadProfile } = useProfile();
const followers = ref(null);

const closeDialog = (e) => {
    e.target.parentNode.parentNode.open = false;
};

const openDialog = (e) => {
    e.target.nextSibling.open = true;
    if (props.user && props.user.Id) {
        loadProfile(props.user.Id);
    }
};

watch(profile, (newProfile) => {
    if (newProfile && newProfile.Followers) {
        followers.value = newProfile.Followers;
    } else {
        followers.value = null;
    }
});

onMounted(async () => {
    if (props.user && props.user.Id) {
        await loadProfile(props.user.Id);
    }
});
</script>

<template>
    <button @click="openDialog">Invite your followers</button>
    <dialog class="background-color-base text-color-primary" open>
        <div class="flex flex-row justify-content-space-between">
            <h3>Invite users</h3>
            <button @click="closeDialog">Close</button>
        </div>

        <p>
            Send invites to your followers so they can join the
            <strong>{{ group.Title }}</strong>
            group.
        </p>
        <FollowersList :followers="followers" :group="group" />
    </dialog>
</template>

<style scoped></style>
