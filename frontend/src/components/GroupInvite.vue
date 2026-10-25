<script setup>
import { ref, onMounted, watch } from 'vue';
import FollowersList from '@/components/FollowersList.vue';
import CloseIcon from '@/components/Icons/CloseIcon.vue';
import PlusIcon from '@/components/Icons/PlusIcon.vue';
import Spacer from '@/components/Spacer.vue';
import { useProfile } from '@/composables/useProfile.js';

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
    <button @click="openDialog">Invite followers</button>
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
                    class="color-rgba-212-224-200-1 text-transform-none line-break-auto overflow-wrap-initial white-space-pre font-size-16px text-rendering-geometricPrecision caret-color-rgba-212-224-200-1 text-decoration-none letter-spacing-0px font-family-Inter font-style-normal font-weight-700"
                >
                    Invite followers
                </span>
                <Spacer />
                <div
                    class="position-relative width-auto height-auto background-1a1d1aFF border-radius-2px display-flex align-items-center align-content-stretch justify-content-start padding-5px flex-direction-row flex-wrap-nowrap"
                    @click="closeDialog"
                >
                    <CloseIcon />
                </div>
            </div>
            <!-- text: Your followers who are not in -->
            <div class="width-100 flex-shrink-0 margin-block-end-6px">
                <div>
                    <div class="display-flex white-space-break-spaces align-items-flex-start">
                        <div
                            class="display-inline-flex flex-direction-column justify-content-inherit min-width-100 margin-right-1px vertical-align-top"
                        >
                            <p
                                class="font-size-0 line-height-1-5 margin-0 text-align-left"
                                dir="ltr"
                            >
                                <span
                                    class="color-rgba-107-114-100-1 text-transform-none line-break-auto overflow-wrap-initial white-space-break-spaces font-size-12px text-rendering-geometricPrecision caret-color-rgba-107-114-100-1 text-decoration-none letter-spacing-0px font-family-Inter font-style-normal font-weight-400"
                                >
                                    Your followers who are not in {{ group.Title }} yet.
                                </span>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
            <FollowersList :followers="followers" :group="group" />
        </div>
    </div>
</template>
<style scoped></style>
