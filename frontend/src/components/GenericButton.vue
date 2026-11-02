<script setup>
defineProps({
    variant: {
        type: String,
        default: 'default',
        validator: (v) => ['default', 'accent', 'danger', 'filled', 'unfilled', 'muted'].includes(v),
    },
});
</script>

<template>
    <button
        class="flex flex-row align-items-center gap-6 border-radius-2px font-weight-500 font-size-13px padding-block-start-6px padding-block-end-6px padding-inline-start-10px padding-inline-end-10px"
        :class="variant"
    >
        <slot name="leading"></slot>
        <span v-if="$slots.default"><slot></slot></span>
        <slot name="trailing"></slot>
    </button>
</template>

<style scoped>
.default {
    background-color: var(--bg-input);
    color: var(--accent);
}
.accent {
    background-color: var(--bg-input);
    border-color: var(--accent);
    color: var(--accent);
}
.danger {
    border-color: var(--border);
    color: var(--danger);
}
.filled {
    background-color: var(--bg-elevated);
    border-color: var(--border);
    color: var(--text-primary);
}
.unfilled {
    border-color: var(--border);
    color: var(--text-primary);
}
.muted {
    border-color: var(--border);
    color: var(--text-muted);
}
</style>
