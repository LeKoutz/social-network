<script setup>
import { useCategories } from '@/composables/useCategories.js';
import GenericButton from './GenericButton.vue';

const selected = defineModel({ type: Number, default: null });

const { categories } = useCategories();
</script>

<template>
    <div class="flex flex-wrap gap-6" role="group">
        <GenericButton
            type="button"
            :variant="selected === null ? 'accent' : 'unfilled'"
            @click="selected = null"
        >
            All
        </GenericButton>
        <GenericButton
            v-for="category in categories"
            :key="category.Id"
            type="button"
            :variant="selected === category.Id ? 'accent' : 'unfilled'"
            @click="selected = category.Id"
        >
            {{ category.Name }}
        </GenericButton>
    </div>
</template>

<style scoped></style>
