<script setup>
import { ref } from 'vue';
import { DateToLocale } from '@/utils/utils.js';

const props = defineProps({
    group: { type: Object, default: null },
});

const group = ref(props.group);
</script>

<template>
    <h2>{{ group.Title }}</h2>
    <template v-if="group.Description.length === 0">
        <p class="text-color-muted">No description avaliable.</p>
    </template>
    <template v-else>
        <p class="group-description">{{ group.Description }}</p>
    </template>
    <p class="group-owner">
        Created by
        <router-link :to="`/profile/view/${group.OwnerUserId}`">
            <strong>{{ group.OwnerUsername }}</strong>
        </router-link>
        on
        <em>({{ DateToLocale(group.Timestamp) }})</em>
    </p>
</template>

<style scoped></style>
