<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiGet } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';
import CategoryShow from '@/components/CategoryShow.vue';
import PostShow from '@/components/PostShow.vue';
import { useCategory } from '@/composables/useCategory.js';

const route = useRoute();
const category = ref(true);
const loading = ref(true);
const posts = ref(true);
const error = ref(null);

onMounted(async () => {
    loading.value = true;
    error.value = null;
    const data = await apiGet(`/api/category/view/${route.params.id}`);
    if (data) {
        category.value = firstOrNull(data.Categories);
        if (!category.value) error.value = 'Category not found';
        posts.value = data.Posts ?? [];
    }
    loading.value = false;
});

watch(
    () => route.params.id,
    (newId) => {
        if (newId) loadCategory(newId);
    },
);
</script>

<template>
    <div class="flex flex-column gap-20 height-100 overflow-y-auto">
        <div class="category" :id="`${category.Id}`">
            <h2>{{ category.Name }}</h2>
            <p>{{ category.Description }}</p>
            <p v-if="posts.length === 0">This category is empty</p>
            <template v-else v-for="post in posts">
                <PostShow :post="post" />
            </template>
        </div>
    </div>
</template>

<style scoped></style>
