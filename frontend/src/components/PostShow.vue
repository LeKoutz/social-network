<script setup>
import { ref, watch } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiGet, apiPost } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';

import CommentCreateForm from '@/components/CommentCreateForm.vue';
import CommentShow from '@/components/CommentShow.vue';
import PostReactForm from '@/components/PostReactForm.vue';
import PostHeader from '@/components/PostHeader.vue';
import PlainTag from '@/components/Tags/PlainTag.vue';
import GenericButton from '@/components/GenericButton.vue';
import ChevronDownIcon from '@/components/Icons/ChevronDownIcon.vue';

const props = defineProps({
    post: { type: Object, default: null },
    withComments: { type: Boolean, default: false },
});

const { user } = useUser();

const post = ref(props.post);
const showComments = ref(props.withComments);
const commentsLoaded = ref(props.withComments);

watch(
    () => props.post,
    (newPost) => {
        post.value = newPost;
    },
);

async function refresh() {
    if (!post.value) return;
    const data = await apiGet(`/api/post/view/${post.value.Id}`);
    if (data && !data.Error?.Has) {
        post.value = firstOrNull(data.Posts);
        commentsLoaded.value = true;
    }
}

async function toggleComments() {
    showComments.value = !showComments.value;
    if (showComments.value && !commentsLoaded.value) await refresh();
}

async function react(action) {
    const data = await apiPost(
        '/api/post/react',
        new URLSearchParams({
            'post-id': post.value.Id,
            action,
        }),
    );
    if (data && !data.Error?.Has) {
        post.value = firstOrNull(data.Posts);
    }
}
</script>

<template>
    <article v-if="post" :id="`${post.Id}`" class="post-card flex flex-column gap-12">
        <PostHeader :post="post" />

        <p class="post-body">{{ post.Body }}</p>

        <img v-if="post.ImagePath" :src="`/${post.ImagePath}`" alt="" class="post-image" />

        <div v-if="post.Categories?.length" class="flex flex-wrap gap-6">
            <PlainTag v-for="category in post.Categories" :key="category.Id">
                {{ category.Name }}
            </PlainTag>
        </div>

        <hr class="divider" />

        <div class="flex align-items-center justify-content-space-between gap-6">
            <PostReactForm :post="post" @react="react" />
            <GenericButton
                type="button"
                variant="no-border"
                :aria-expanded="showComments"
                @click="toggleComments"
            >
                {{ showComments ? 'Hide comments' : 'Comments' }}
                <template #trailing>
                    <ChevronDownIcon :class="{ 'transform-rotate-180deg': showComments }" />
                </template>
            </GenericButton>
        </div>

        <div v-if="showComments" class="flex flex-column gap-12">
            <CommentCreateForm v-if="user.LoggedIn" :post="post" @created="refresh" />
            <CommentShow
                v-for="comment in post.Comments"
                :key="comment.Id"
                :post="post"
                :comment="comment"
                @updated="refresh"
            />
        </div>
    </article>
</template>

<style scoped>
.post-card {
    padding: 18px 20px 14px;
    background-color: var(--bg-elevated);
    border: 1px solid var(--accent-line);
    border-radius: 2px;
}

.post-body {
    margin: 0;
    font-size: 14px;
    line-height: 1.5;
    color: var(--text-primary);
    white-space: pre-wrap;
    overflow-wrap: anywhere;
}

.post-image {
    max-width: 100%;
    max-height: 400px;
    width: fit-content;
    border: 1px solid var(--border);
    border-radius: 2px;
}

.divider {
    width: 100%;
    margin: 0;
    border: none;
    border-top: 1px solid var(--accent-line);
}
</style>
