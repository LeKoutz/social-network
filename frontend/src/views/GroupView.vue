<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch } from '@/utils/api.js';
import { firstOrNull, DateToLocale } from '@/utils/utils.js';
import PostCreateForm from '@/components/PostCreateForm.vue';
import PostShow from '@/components/PostShow.vue';

const route = useRoute();
const group = ref(null);
const posts = ref(null);
const loading = ref(true);
const error = ref(null);

async function loadGroup(id) {
    loading.value = true;
    error.value = null;
    const data = await apiFetch(`/api/group/view/${id}`);
    if (data) {
        group.value = firstOrNull(data.Groups);
        posts.value = data.Posts;
        if (!group.value) error.value = 'Group not found';
    }
    loading.value = false;
}

onMounted(() => loadGroup(route.params.id));

watch(
    () => route.params.id,
    (newId) => {
        if (newId) loadGroup(newId);
    }
);
</script>

<template>
    <div class="container">
        <div v-if="loading">Loading…</div>
        <div v-else-if="error">{{ error }}</div>
        <div v-else-if="group">
            <h2>{{ group.Title }}</h2>
            <p class="group-description">{{ group.Description }}</p>
            <p class="group-owner">
                Owned by 
                <router-link :to="`/profile/view/${group.OwnerUserId}`">
                    <strong>{{ group.OwnerUsername }}</strong>
                </router-link>
                on <em>({{ DateToLocale(group.Timestamp) }})</em>
            </p>
            <template v-if="group.Member">
                <PostCreateForm :group="group" />
                <div v-if="posts && posts.length">
                    <PostShow v-for="post in posts" :key="post.Id" :post="post" />
                </div>
                <p v-else>No posts in this group yet. Be the first to post!</p>
            </template>
            <p v-else class="group-locked">
                This group is private. Only its members can see the posts and join the
                conversation.
            </p>
        </div>
    </div>
</template>

<style scoped></style>
