<script setup>
import { ref } from 'vue';
import UserList from '@/components/UserList.vue';
import { useProfile } from '@/composables/useProfile.js';

const props = defineProps({
    group: { type: Object, default: null },
    user: { type: Object, default: null },
});

const group = ref(props.group);
const user = ref(props.user);
const profile = useProfile();

const closeDialog = (e) => {
    e.target.parentNode.parentNode.open = false;
};

const openDialog = (e) => {
    e.target.nextSibling.open = true;
};
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
            <em>{{ group.Title }}</em>
            group.
        </p>
        <UserList :group="group" />
    </dialog>
</template>

<style scoped></style>
