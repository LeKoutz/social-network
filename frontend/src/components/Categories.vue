<script setup>
import { useCategories } from '@/composables/useCategories';

const { categories } = useCategories();
</script>

<template>
    <div class="categories">
        <h2>Categories</h2>
        <div class="category" v-for="category in categories" :key="category.Id">
            <router-link :to="`/category/view/${category.Id}`">
                <h2>{{ category.Name }}</h2>
            </router-link>
            <p class="category-description">{{ category.Description }}</p>
        </div>
    </div>
</template>

<style scoped></style>
