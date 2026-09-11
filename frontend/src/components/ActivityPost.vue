<script setup>
defineProps({
    post: {
        type: Object,
        required: true,
    },
});
</script>

<template>
    <div class="container">
        <div class="activity-highlight">
            <router-link :to="`/post/view/${post.Id}`">
                <h3>{{ post.Title }}</h3>
            </router-link>
            <pre>{{ post.Body }}</pre>
            <div class="reactions">
                <button disabled>
                    {{ post.Likes }} {{ post.Liked ? '👍' : '👍🏻' }}
                </button>
                <button disabled>
                    {{ post.Dislikes }} {{ post.Disliked ? '👎' : '👎🏻' }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped></style>
