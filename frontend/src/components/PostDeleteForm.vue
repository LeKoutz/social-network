<script setup>
import { useRouter } from 'vue-router';
import { usePost } from '@/composables/usePost.js';

const props = defineProps({
    post: {
        type: Object,
        required: true,
    },
});

const router = useRouter();

const deletePost = async (id) => {
    const data = await usePost().deletePost(id);
    if (data) router.push('/');
};
</script>

<template>
    <button type="button" @click="deletePost(post.Id)">Delete Post</button>
</template>

<style scoped></style>
