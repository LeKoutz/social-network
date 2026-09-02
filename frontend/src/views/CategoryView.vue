<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch } from '@/utils/api.js';

const route = useRoute();
const category = ref(null);
const posts = ref(null);

onMounted(async () => {
    const data = apiFetch(`/api/category/view/${route.params.id}`);
    if (data) {
        posts.value = data.Posts ?? [];
        category.value = data.Categories[0];
    }
});
</script>

<template>
    <div class="container">
        <div class="category" :id="`${category.Id}`">
            <h2>{{ category.Name }}</h2>
            <p>{{ category.Description }}</p>
            <p v-if="posts.length === 0">This category is empty</p>
            <!-- TODO: posts rendering -->
        </div>
    </div>
</template>

<style scoped></style>
