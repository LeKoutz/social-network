<script setup>
import { ref } from 'vue';

import { useUser } from '@/composables/useUser.js';
import PostShow from '@/components/PostShow.vue';

const user = useUser();

const props = defineProps({
    posts: { type: Array, default: null },
});
</script>

<template>
    <div class="flex flex-column gap-16" v-if="posts && posts.length">
        <PostShow v-for="post in posts" :key="post.Id" :post="post" />
    </div>
    <p v-else>No posts yet. Be the first to post!</p>
</template>

<style scoped></style>
