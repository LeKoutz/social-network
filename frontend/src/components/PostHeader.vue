<script setup>
import { computed } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { DateToLocale } from '@/utils/utils.js';
import PostDeleteForm from '@/components/PostDeleteForm.vue';

const props = defineProps({
    post: { type: Object, required: true },
});

const { user } = useUser();

const isOwner = computed(
    () => user.value.LoggedIn && String(user.value.Id) === String(props.post.User?.Id),
);
</script>

<template>
    <header class="flex flex-column gap-5">
        <h3 class="post-title">
            <router-link :to="`/post/view/${post.Id}`">{{ post.Title }}</router-link>
        </h3>
        <div class="flex flex-wrap align-items-center gap-5 font-size-12px text-color-muted">
            <template v-if="post.User">
                <span>by</span>
                <router-link :to="`/profile/view/${post.User.Id}`" class="author">
                    {{ post.User.Username }}
                </router-link>
                <span aria-hidden="true">·</span>
            </template>
            <time>{{ DateToLocale(post.Timestamp) }}</time>
            <template v-if="post.GroupId">
                <span aria-hidden="true">·</span>
                <span>in</span>
                <router-link :to="`/group/view/${post.GroupId}`" class="author">
                    {{ post.GroupTitle }}
                </router-link>
            </template>
            <span v-if="isOwner" class="owner-actions flex align-items-center gap-12">
                <router-link :to="`/post/edit/${post.Id}`">Edit</router-link>
                <PostDeleteForm :post="post" class="delete-link" />
            </span>
        </div>
    </header>
</template>

<style scoped>
.post-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
}

.post-title a,
.author,
.owner-actions a {
    text-decoration: none;
}

.author {
    font-weight: 600;
}

.owner-actions {
    margin-left: auto;
}

.delete-link {
    padding: 0;
    border: none;
    background: none;
    font: inherit;
    color: var(--danger);
    cursor: pointer;
}
</style>
