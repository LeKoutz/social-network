<script setup>
import { useUser } from '@/composables/useUser.js';

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
    <div class="reactions">
        <button
            type="button"
            class="post-reaction"
            :disabled="!user.LoggedIn"
            @click="react('like')"
        >
            {{ post.Likes }} {{ post.Liked ? '👍' : '👍🏻' }}
        </button>
        <button
            type="button"
            class="post-reaction"
            :disabled="!user.LoggedIn"
            @click="react('dislike')"
        >
            {{ post.Dislikes }} {{ post.Disliked ? '👎' : '👎🏻' }}
        </button>
    </div>
</template>

<style scoped></style>
