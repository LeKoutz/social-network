<script setup>
import { ref, onMounted, watch } from 'vue';
import FollowersList from '@/components/FollowersList.vue';
import CloseIcon from '@/components/Icons/CloseIcon.vue';
import PlusIcon from '@/components/Icons/PlusIcon.vue';
import Spacer from '@/components/Spacer.vue';
import { useProfile } from '@/composables/useProfile.js';
import GenericButton from '@/components/GenericButton.vue';

const props = defineProps({
    group: { type: Object, default: null },
    user: { type: Object, default: null },
});

const { profile, loadProfile } = useProfile();
const followers = ref(null);
const isModalOpen = ref(false);

const closeDialog = (e) => {
    isModalOpen.value = false;
    // document.querySelector('#invite-to-group-modal').hidden = true;
};

const openDialog = (e) => {
    isModalOpen.value = true;
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
    <GenericButton @click="openDialog">
        <template v-slot:icon>
            <PlusIcon />
        </template>
        Invite followers
    </GenericButton>
    <div
        v-if="isModalOpen"
        id="invite-to-group-modal"
        class="position-fixed display-flex flex-direction-row justify-content-center align-content-center align-items-center top-50 left-50 z-index-20"
    >
        <div
            class="position-fixed width-340px height-auto background-color-elevated border-radius-2px box-shadow-0px-8px-24px-0px-00000066 display-flex align-items-stretch align-content-stretch justify-content-start row-gap-2px padding-16px flex-direction-column flex-wrap-nowrap"
        >
            <div
                class="width-100 height-auto display-flex align-items-center align-content-stretch justify-content-start padding-block-end-4px flex-direction-row flex-wrap-nowrap"
            >
                <span
                    class="text-color-primary font-size-16px text-decoration-none font-family-Inter font-style-normal font-weight-700"
                >
                    Invite followers
                </span>
                <Spacer />
                <div
                    class="background-color-bg-input border-radius-2px flex align-items-center align-content-stretch justify-content-start padding-5px flex-direction-row flex-wrap-nowrap"
                    @click="closeDialog"
                >
                    <CloseIcon />
                </div>
            </div>
            <span
                class="text-color-muted font-size-12px text-decoration-none font-family-Inter font-style-normal font-weight-400"
            >
                Your followers who are not in {{ group.Title }} yet.
            </span>
            <div>
                <FollowersList :followers="followers" :group="group" />
            </div>
        </div>
    </div>
</template>
<style scoped></style>
