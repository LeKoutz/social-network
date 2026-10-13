<script setup>
import { ref, watch } from 'vue';
import { useUser } from '@/composables/useUser.js';

import { DateToLocale } from '@/utils/utils.js';
import PostDeleteForm from '@/components/PostDeleteForm.vue';

const { user } = useUser();
const props = defineProps({
    post: { type: Object, default: null },
});
const post = ref(props.post);

const isOwner = () => user.value.LoggedIn && String(user.value.Id) === String(post.value?.User?.Id);
</script>

<template>
    <router-link :to="`/post/view/${post.Id}`">
        <h3>{{ post.Title }}</h3>
    </router-link>
    <div class="post-details">
        <template v-if="post.GroupId">
            <p>
                Group:
                <router-link :to="`/group/view/${post.GroupId}`">
                    <strong>{{ post.GroupTitle }}</strong>
                </router-link>
            </p>
        </template>
        <template v-if="post.Categories">
            <p>
                Categories:
                <span v-for="category in post.Categories" :key="category.Id">
                    <router-link :to="`/category/view/${category.Id}`">
                        {{ category.Name }}
                    </router-link>
                </span>
            </p>
        </template>
        <p>
            Posted by
            <router-link :to="`/profile/view/${post.User.Id}`">
                <strong>{{ post.User.Username }}</strong>
            </router-link>
            on
            <em>({{ DateToLocale(post.Timestamp) }})</em>
        </p>
    </div>
    <div v-if="isOwner()" class="manage-post">
        <router-link :to="`/post/edit/${post.Id}`">
            <button type="button">Edit Post</button>
        </router-link>
        <PostDeleteForm :post="post" />
    </div>
</template>

<style scoped></style>
