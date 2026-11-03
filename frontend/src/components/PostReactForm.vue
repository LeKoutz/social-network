<script setup>
import { useUser } from '@/composables/useUser.js';
import GenericButton from '@/components/GenericButton.vue';
import ThumbsUpIcon from '@/components/Icons/ThumbsUpIcon.vue';
import ThumbsDownIcon from '@/components/Icons/ThumbsDownIcon.vue';

defineProps({
    post: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(['react']);

const { user } = useUser();

function react(action) {
    if (!user.value.LoggedIn) return;
    emit('react', action);
}
</script>

<template>
    <div class="flex align-items-center gap-6">
        <GenericButton
            type="button"
            :variant="post.Liked ? 'default' : 'unfilled'"
            :aria-pressed="Boolean(post.Liked)"
            :aria-label="`Like (${post.Likes ?? 0})`"
            :disabled="!user.LoggedIn"
            @click="react('like')"
        >
            <template #leading><ThumbsUpIcon /></template>
            {{ post.Likes ?? 0 }}
        </GenericButton>
        <GenericButton
            type="button"
            :variant="post.Disliked ? 'danger' : 'unfilled'"
            :aria-pressed="Boolean(post.Disliked)"
            :aria-label="`Dislike (${post.Dislikes ?? 0})`"
            :disabled="!user.LoggedIn"
            @click="react('dislike')"
        >
            <template #leading><ThumbsDownIcon /></template>
            {{ post.Dislikes ?? 0 }}
        </GenericButton>
    </div>
</template>

<style scoped></style>
