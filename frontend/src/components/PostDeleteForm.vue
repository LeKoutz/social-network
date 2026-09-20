<script setup>
import { useRouter } from 'vue-router';
import { apiPost } from '@/utils/api.js';

const props = defineProps({
    post: {
        type: Object,
        required: true,
    },
});

const router = useRouter();

async function remove() {
    const data = await apiPost(
        '/api/post/delete',
        new URLSearchParams({ 'post-id': props.post.Id })
    );
    if (data) router.push('/');
}
</script>

<template>
    <button type="button" class="post-delete" @click="remove">Delete Post</button>
</template>

<style scoped></style>
