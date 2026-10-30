<script setup>
import { computed, onMounted } from 'vue';
import GlobeIcon from '@/components/Icons/GlobeIcon.vue';
import LockIcon from '@/components/Icons/LockIcon.vue';
import UsersIcon from '@/components/Icons/UsersIcon.vue';
import ChevronDownIcon from '@/components/Icons/ChevronDownIcon.vue';
import { useUser } from '@/composables/useUser.js';
import { useProfile } from '@/composables/useProfile.js';
import PeopleList from '@/components/PeopleList.vue';

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
        disabled: false,
    },
    {
        value: 'private',
        label: 'Private',
        description: 'Only your accepted followers can see this post.',
        disabled: false,
    },
    {
        value: `limited`,
        label: `Selected Followers`,
        description: `Only selected followers can see this post.`,
        disabled: false,
    },
];

const selected = computed(() => options.find((option) => option.value === visibility.value));

onMounted(async () => {
    if (user.value.Id) await loadProfile(user.value.Id);
});
</script>

<template>
    <div class="visibility-panel">
        <label class="visibility-select" :title="selected.description">
            <GlobeIcon v-if="visibility === 'public'" />
            <LockIcon v-else-if="visibility === 'private'" />
            <UsersIcon v-else />

            <span>{{ selected.label }}</span>
            <ChevronDownIcon />

            <select v-model="visibility" name="visibility" :disabled="disabled">
                <option v-for="option in options" :key="option.value" :value="option.value">
                    {{ option.label }}
                </option>
            </select>
        </label>
        <div v-if="visibility === 'limited'" class="followers-dropdown">
            <small v-if="!loading && !followers.length">You have no followers yet</small>
            <small v-else>Choose followers</small>
            <PeopleList :people="followers">
                <template #action="{ person }">
                    <label class="row-toggle">
                        <input type="checkbox" v-model="allowed" :value="person.UserId" />
                    </label>
                </template>
            </PeopleList>
        </div>
    </div>
</template>

<style scoped>
.visibility-panel {
    position: relative;
}

.visibility-select {
    display: flex;
    align-items: center;
    position: relative;
    gap: 8px;
    padding: 4px 10px;
    color: var(--text-secondary);
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: 2px;
    width: fit-content;
}

.visibility-select > select {
    position: absolute;
    inset: 0;
    opacity: 0;
    appearance: none;
    background: transparent;
    border: none;
    color: inherit;
    font-size: 13px;
    field-sizing: content;
}

.visibility-select:focus-within {
    border-color: white;
}

.followers-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    z-index: 10;
    width: 200px;
    max-height: 300px;
    padding: 6px 10px;
    overflow-y: auto;
    overflow-x: hidden;
    background: var(--bg-input);
}

.followers-dropdown :deep(li) {
    position: relative;
    padding: 4px 6px;
}

.followers-dropdown :deep([data-type='username']) {
    color: var(--text-secondary);
}

.followers-dropdown :deep(li:has(input:checked) [data-type='username']) {
    color: var(--text-primary);
}

.row-toggle::after {
    content: '';
    position: absolute;
    inset: 0;
    cursor: pointer;
}

.followers-dropdown input[type='checkbox'] {
    appearance: none;
    width: 14px;
    height: 14px;
    margin: 0;
    border-radius: 2px;
    background-color: var(--bg-base);
}

.followers-dropdown input:checked {
    background-color: var(--accent);
}
</style>
