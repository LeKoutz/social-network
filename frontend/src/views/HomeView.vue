<script setup>
import { ref, onMounted, computed } from 'vue';
import { apiGet } from '@/utils/api.js';
import PostCreateForm from '@/components/PostCreateForm.vue';
import PostsList from '@/components/PostsList.vue';
import CategoryFilter from '@/components/CategoryFilter.vue';

const posts = ref(null);
const selectedCategory = ref(null);
const visiblePosts = computed(() =>
    selectedCategory.value === null
        ? posts.value
        : posts.value?.filter((post) =>
              post.Categories?.some((c) => c.Id === selectedCategory.value),
          ),
);

async function loadPosts() {
    const data = await apiGet('/api/posts');
    if (data && !data.Error?.Has) {
        // TODO: remove once the backend sorts (ORDER BY timestamp DESC)
        posts.value = (data.Posts ?? []).sort((a, b) => Number(b.Timestamp) - Number(a.Timestamp));
    }
}

onMounted(loadPosts);
</script>

<template>
    <div
        class="flex justify-content-center padding-block-start-28px padding-block-end-32px padding-inline-start-32px padding-inline-end-32px"
    >
        <div class="feed flex flex-column gap-16">
            <PostCreateForm @created="loadPosts" />
            <CategoryFilter v-model="selectedCategory" />
            <p v-if="posts === null" class="text-color-muted">Loading…</p>
            <PostsList v-else :posts="visiblePosts" />
        </div>
    </div>
</template>

<style scoped>
.feed {
    max-width: 680px;
    width: 100%;
}
</style>
