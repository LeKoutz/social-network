<script setup>
import { ref } from 'vue';
import { DateToLocale } from '@/utils/utils.js';
import GroupIcon from '@/components/Icons/GroupIcon.vue';
import HighLightedTag from '@/components/Tags/HighLightedTag.vue';

const props = defineProps({
    group: { type: Object, default: null },
    user: { type: Object, default: null },
});

const group = ref(props.group);
</script>

<template>
    <div class="flex flex-row gap-16">
        <div
            class="width-72px height-72px align-items-center flex justify-content-center background-color-accent-line"
        >
            <GroupIcon />
        </div>
        <div class="flex flex-column gap-6 justify-content-flex-start">
            <div class="flex flex-row align-items-center gap-10">
                <span class="font-size-24px font-family-Inter font-style-normal font-weight-700">
                    {{ group.Title }}
                </span>
                <HighLightedTag v-if="group.Owner">Owner</HighLightedTag>
                <PlainTag v-else-if="group.Member">Member</PlainTag>
            </div>
            <span class="text-color-muted font-size-13px">
                by
                <router-link :to="`/profile/view/${group.OwnerUserId}`">
                    {{ group.OwnerUsername }}
                </router-link>
                ·
                {{ DateToLocale(group.Timestamp) }}
            </span>
        </div>
    </div>
    <div>
        <template v-if="group.Description.length === 0">
            <p class="text-color-muted">No description avaliable.</p>
        </template>
        <template v-else>
            <span class="text-color-secondary font-weight-400">{{ group.Description }}</span>
        </template>
    </div>
</template>

<style scoped></style>
