<script setup>
const visibility = defineModel({
    type: String,
    default: 'public',
})

defineProps({
    disabled: {
        type: Boolean,
        default: false,
    },
})

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
    }
]
</script>

<template>
    <fieldset class="post-visibility" :disabled="disabled">
        <legend>Who can see this post?</legend>

        <label
            v-for="option in options"
            :key="option.value"
            class="visibility-option"
            :class="{ selected: visibility === option.value }"
        >
            <input
                v-model="visibility"
                type="radio"
                name="visibility"
                :value="option.value"
                :disabled="option.disabled"
            />
            <span>
                <strong>{{ option.label }}</strong>
                <small>{{ option.description }}</small>
            </span>
        </label>

        <small>You can always see your own post.</small>
    </fieldset>
</template>

<style scoped>
.post-visibility {
    gap: 10px;
    margin: 12px 0;
    box-shadow: none;
}

.visibility-option {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    padding: 12px;
    border: 1px solid var(--border);
    cursor: pointer;
}

.visibility-option.selected {
    border-color: var(--accent);
    background-color: var(--bg-input);
}

.visibility-option span {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.visibility-option small {
    color: var(--text-secondary);
}

.visibility-option input {
    accent-color: var(--accent);
}
</style>