<script setup>
defineProps({
    to: { type: [String, Object], default: null },
    variant: {
        type: String,
        default: 'default',
        validator: (v) => ['default', 'danger'].includes(v),
    },
    active: { type: Boolean, default: false },
});
</script>

<template>
    <router-link v-if="to" :to="to" class="dropdown-item" :class="[variant, { active }]">
        <slot></slot>
    </router-link>
    <button v-else type="button" class="dropdown-item" :class="[variant, { active }]">
        <slot></slot>
    </button>
</template>

<style scoped>
.dropdown-item {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 8px 10px;
    box-sizing: border-box;
    border: none;
    border-radius: 2px;
    background: none;
    font: inherit;
    font-size: 13px;
    text-align: left;
    text-decoration: none;
    white-space: nowrap;
    cursor: pointer;
}

.dropdown-item:hover,
.dropdown-item.active {
    background-color: var(--bg-input);
}

.default {
    color: var(--text-primary);
}
.default:hover,
.default.active {
    color: var(--accent);
}

.danger {
    color: var(--danger);
}
</style>
