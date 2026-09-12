<script setup>
import { useUser } from '@/composables/useUser.js';

defineProps({
    comment: {
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
    <div class="reactions">
        <button
            type="button"
            class="comment-reaction"
            :disabled="!user.LoggedIn"
            @click="react('like')"
        >
            {{ comment.Likes }} {{ comment.Liked ? '👍' : '👍🏻' }}
        </button>
        <button
            type="button"
            class="comment-reaction"
            :disabled="!user.LoggedIn"
            @click="react('dislike')"
        >
            {{ comment.Dislikes }} {{ comment.Disliked ? '👎' : '👎🏻' }}
        </button>
    </div>
</template>

<style scoped></style>
