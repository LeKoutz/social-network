<script setup>
import { computed, onMounted } from 'vue';
import GlobeIcon from '@/components/Icons/GlobeIcon.vue';
import LockIcon from '@/components/Icons/LockIcon.vue';
import UsersIcon from '@/components/Icons/UsersIcon.vue';
import ChevronDownIcon from '@/components/Icons/ChevronDownIcon.vue';
import { useUser } from '@/composables/useUser.js';
import { useProfile } from '@/composables/useProfile.js';
import PeopleList from '@/components/PeopleList.vue';
import GenericDropdown from './GenericDropdown.vue';
import GenericButton from './GenericButton.vue';
import DropdownItem from './DropdownItem.vue';

const allowed = defineModel('allowed', {
    type: Array,
    default: () => [],
});

const visibility = defineModel({
    type: String,
    default: 'public',
});

const { user } = useUser();
const { profile, loadProfile, loading } = useProfile();

const followers = computed(() => profile.value?.Followers ?? []);

defineProps({
    disabled: {
        type: Boolean,
        default: false,
    },
});

const options = [
    {
        value: 'public',
        label: 'Public',
        description: 'Anyone can see this post.',
    },
    {
        value: 'private',
        label: 'Private',
        description: 'Only your accepted followers can see this post.',
    },
    {
        value: `limited`,
        label: `Selected Followers`,
        description: `Only selected followers can see this post.`,
    },
];

const selected = computed(() => options.find((option) => option.value === visibility.value));

onMounted(async () => {
    if (user.value.Id) await loadProfile(user.value.Id);
});
</script>

<template>
    <GenericDropdown open-on="click" variant="input">
        <template #trigger="{ open }">
            <GenericButton type="button" variant="unfilled" :disabled="disabled">
                <template #leading>
                    <GlobeIcon v-if="visibility === 'public'" />
                    <LockIcon v-else-if="visibility === 'private'" />
                    <UsersIcon v-else />
                </template>
                {{ selected.label }}
                <template #trailing>
                    <ChevronDownIcon :class="{ 'transform-rotate-180deg': open }" />
                </template>
            </GenericButton>
        </template>
        <template #default="{ close }">
            <DropdownItem
                :title="option.description"
                v-for="option in options"
                :key="option.value"
                :active="option.value === visibility"
                @click="
                    visibility = option.value;
                    close();
                "
            >
                {{ option.label }}
            </DropdownItem>
        </template>
    </GenericDropdown>
    <GenericDropdown v-if="visibility === 'limited'" variant="input">
        <template #trigger="{ open }">
            <GenericButton type="button" variant="unfilled" :disabled="disabled">
                Choose followers ({{ allowed.length }})
                <template #trailing>
                    <ChevronDownIcon :class="{ 'transform-rotate-180deg': open }" />
                </template>
            </GenericButton>
        </template>
        <small v-if="!loading && !followers.length">You have no followers yet</small>
        <PeopleList v-else class="followers" :people="followers">
            <template #action="{ person }">
                <label class="row-toggle">
                    <input type="checkbox" v-model="allowed" :value="person.UserId" />
                </label>
            </template>
        </PeopleList>
    </GenericDropdown>
</template>

<style scoped>
.followers {
    min-width: 180px;
    max-height: 300px;
    margin: 0;
    overflow-y: auto;
}

.followers :deep(li) {
    position: relative;
    padding: 4px 6px;
}

.followers :deep([data-type='username']) {
    color: var(--text-secondary);
}

.followers :deep(li:has(input:checked) [data-type='username']) {
    color: var(--text-primary);
}

.row-toggle::after {
    content: '';
    position: absolute;
    inset: 0;
    cursor: pointer;
}

.followers input[type='checkbox'] {
    appearance: none;
    width: 14px;
    height: 14px;
    margin: 0;
    border-radius: 2px;
    background-color: var(--bg-base);
}

.followers input:checked {
    background-color: var(--accent);
}
</style>
