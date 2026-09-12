<script setup>
import { ref, watch } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { apiFetch } from '@/utils/api.js';
import { firstOrNull, DateToLocale } from '@/utils/utils.js';

import CommentCreateForm from '@/components/CommentCreateForm.vue';
import CommentShow from '@/components/CommentShow.vue';

const props = defineProps({
    post: { type: Object, default: null },
});

const { user } = useUser();

const post = ref(props.post);

watch(
    () => props.post,
    (newPost) => {
        post.value = newPost;
    }
);

async function refresh() {
    if (!post.value) return;
    const data = await apiFetch(`/api/post/view/${post.value.Id}`);
    if (data) {
        post.value = firstOrNull(data.Posts);
    }
}

const isOwner = () =>
    user.value.LoggedIn &&
    String(user.value.Id) === String(post.value?.User?.Id);

</script>
<template>
    <div class="container">
        <div v-if="post" class="post" :id="`${post.Id}`">
            <div class="post-content">
                <router-link :to="`/post/view/${post.Id}`"><h3>{{post.Title}}</h3></router-link>
                <div class="post-details">
                    <p>Categories:
                    <span v-for="category in post.Categories" :key="category.Id">
                        <router-link :to="`/category/view/${category.Id}`">
                            {{ category.Name }}
                        </router-link>
                    </span>
                    </p>
                    <p>Posted by
                    <router-link :to="`/user/${post.User.Id}`">
                        <strong>{{post.User.Username}}</strong>
                    </router-link>
                    on
                    <em>({{DateToLocale(post.Timestamp)}})</em>
                    </p>
                </div>
                <div v-if="isOwner()" class="manage-post">
                    <router-link :to="`/post/edit/${post.Id}`">
                        <button type="button">Edit Post</button>
                    </router-link>
                    <router-link :to="`/post/delete/${post.Id}`">
                        <button type="button">Delete Post</button>
                    </router-link>
                </div>
                <pre>{{post.Body}}</pre>
                <img v-if="post.ImagePath" :src="`/${post.ImagePath}`" alt="Post image" style="max-width: 100%;"/>
                <div class="comments">
                    <CommentCreateForm v-if="user.LoggedIn" :post="post" @created="refresh" />
                    <CommentShow
                        v-for="comment in post.Comments"
                        :key="comment.Id"
                        :post="post"
                        :comment="comment"
                        @updated="refresh"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped></style>
