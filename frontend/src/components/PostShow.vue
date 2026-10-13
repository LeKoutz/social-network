<script setup>
import { ref, watch } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiGet, apiPost } from '@/utils/api.js';
import { firstOrNull, DateToLocale } from '@/utils/utils.js';

import CommentCreateForm from '@/components/CommentCreateForm.vue';
import CommentShow from '@/components/CommentShow.vue';
import PostReactForm from '@/components/PostReactForm.vue';
import PostDeleteForm from '@/components/PostDeleteForm.vue';
import PostHeader from '@/components/PostHeader.vue';

const props = defineProps({
    post: { type: Object, default: null },
});

const { user } = useUser();

const post = ref(props.post);

watch(
    () => props.post,
    (newPost) => {
        post.value = newPost;
    },
);

async function refresh() {
    if (!post.value) return;
    const data = await apiGet(`/api/post/view/${post.value.Id}`);
    if (data) {
        post.value = firstOrNull(data.Posts);
    }
}

async function react(action) {
    const data = await apiPost(
        '/api/post/react',
        new URLSearchParams({
            'post-id': post.value.Id,
            action,
        }),
    );
    if (data) {
        post.value = firstOrNull(data.Posts);
    }
}

const isOwner = () => user.value.LoggedIn && String(user.value.Id) === String(post.value?.User?.Id);
</script>
<template>
    <div v-if="post" class="flex flex-column gap-20 box-shadow post" :id="`${post.Id}`">
        <PostHeader :post="post" />
        <div class="post-content">
            <pre>{{ post.Body }}</pre>
            <img
                v-if="post.ImagePath"
                :src="`/${post.ImagePath}`"
                alt="Post image"
                style="max-width: 100%"
            />
            <PostReactForm :post="post" @react="react" />
        </div>
        <div class="flex flex-column gap-20">
            <CommentCreateForm v-if="user.LoggedIn" :post="post" @created="refresh" />
            <CommentShow
                class="box-shadow"
                v-for="comment in post.Comments"
                :key="comment.Id"
                :post="post"
                :comment="comment"
                @updated="refresh"
            />
        </div>
    </div>
</template>

<style scoped></style>
