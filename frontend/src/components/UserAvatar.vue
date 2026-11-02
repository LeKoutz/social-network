<script setup>
import { computed } from 'vue';

const props = defineProps({
    url: { type: String, default: '' },
    name: { type: String, default: '' },
    size: { type: String, default: 'small' },
    variant: {
        type: String,
        default: 'default',
        validator: (v) => ['default', 'muted'].includes(v),
    },
});

const classes = computed(() => ['user-avatar', `user-avatar-${props.size}`, props.variant]);
</script>

<template>
    <img v-if="props.url" :src="`/${props.url}`" alt="" :class="classes" />
    <span v-else-if="$slots.default" :class="classes">
        <slot></slot>
    </span>
    <span v-else :class="classes">
        {{ props.name.charAt(0).toUpperCase() }}
    </span>
</template>

<style scoped>
.user-avatar {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    border-radius: 50%;
    font-weight: 500;
    object-fit: cover;
}

.user-avatar :slotted(svg) {
    width: 60%;
    height: 60%;
}

/* Sizes */
.user-avatar-small {
    /* Topbar */
    width: 28px;
    height: 28px;
    font-size: 13px;
}

/* Variants */
.default {
    background-color: var(--accent-dim);
    color: var(--bg-base);
}
.muted {
    background-color: var(--bg-elevated);
    color: var(--text-muted);
}
</style>
