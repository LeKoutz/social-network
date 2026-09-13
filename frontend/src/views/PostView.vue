<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';
import PostShow from '@/components/PostShow.vue';

const route = useRoute();
const post = ref(null);
const loading = ref(true);
const error = ref(null);

async function loadPost(id) {
    loading.value = true;
    error.value = null;
    try {
        const data = await apiFetch(`/api/post/view/${id}`);
        console.log(id);
        post.value = firstOrNull( data?.Posts );
        if (!post.value) error.value = 'Post not found';
    } catch (e) {
        error.value = e?.message ?? 'Failed to load post';
    } finally {
        loading.value = false;
    }
}

onMounted(()=>loadPost(route.params.id));

watch(() => route.params.id, (newId) => {
    if (newId) loadPost(newId);
});
</script>

<template>
    <div class="container">
        <div v-if="loading">Loading…</div>
        <div v-else-if="error">{{ error }}</div>
        <PostShow v-else-if="post" :post="post" />
    </div>
</template>

<style scoped></style>
