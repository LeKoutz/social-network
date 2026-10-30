<script setup>
import UserAvatar from '@/components/UserAvatar.vue';

defineProps({
    people: {
        type: Array,
        default: () => [],
    },
});
</script>

<template>
    <ul class="people-list">
        <li v-for="person in people" :key="person.UserId">
            <span class="user">
                <UserAvatar :url="person.AvatarURL" :name="person.Username" size="small" />
                <span data-type="username">{{ person.Username }}</span>
            </span>
            <slot name="action" :person="person"></slot>
        </li>
    </ul>
</template>

<style scoped>
.people-list {
    display: flex;
    flex-direction: column;
    list-style: none;
    padding: 0;
    gap: 5px;
}

.people-list li {
    display: flex;
    justify-content: space-between;
}

.user {
    display: flex;
    gap: 6px;
}
</style>
