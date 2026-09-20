<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';
import PostEditForm from '@/components/PostEditForm.vue';

const route = useRoute();
const post = ref(null);
const categories = ref([]);
const loading = ref(true);
const error = ref(null);

async function loadPost(id) {
    loading.value = true;
    error.value = null;
    const data = await apiFetch(`/api/post/edit/${id}`);
    if (data) {
        post.value = firstOrNull(data.Posts);
        categories.value = data.Categories ?? [];
        if (!post.value) error.value = 'Post not found';
    }
    loading.value = false;
}

onMounted(() => loadPost(route.params.id));

watch(
    () => route.params.id,
    (newId) => {
        if (newId) loadPost(newId);
    }
);
</script>

<template>
    <div class="container">
        <div v-if="loading">Loading…</div>
        <div v-else-if="error">{{ error }}</div>
        <PostEditForm v-else-if="post" :post="post" :categories="categories" />
    </div>
</template>

<style scoped></style>