<script setup>
const props = defineProps({
    url: { type: String, default: '' },
    name: { type: String, default: '' },
    size: { type: String, default: 'small' },
});
</script>

<template>
    <img
        v-if="props.url"
        :src="`/${props.url}`"
        alt=""
        :class="['user-avatar', `user-avatar-${props.size}`]"
    />
    <span v-else :class="['user-avatar', `user-avatar-${props.size}`]">
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
    background-color: var(--accent-dim);
    color: var(--bg-base);
    font-weight: 700;
    object-fit: cover;
}

.user-avatar-small {
    /* Topbar */
    width: 28px;
    height: 28px;
    font-size: 13px;
}
</style>
