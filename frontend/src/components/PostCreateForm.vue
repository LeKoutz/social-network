<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useCategories } from '@/composables/useCategories.js';
import { apiPost } from '@/utils/api.js';

const router = useRouter();
const { categories } = useCategories();

const title = ref('');
const body = ref('');
const image = ref(null);
const selected = ref([]);
const submitting = ref(false);

function toggle(categoryId) {
    const index = selected.value.indexOf(categoryId);
    if (index === -1) {
        selected.value.push(categoryId);
    } else {
        selected.value.splice(index, 1);
    }
}

function onImage(e) {
    image.value = e.target.files[0] ?? null;
}

async function submit() {
    submitting.value = true;
    const data = new FormData();
    data.append('title', title.value);
    data.append('body', body.value);
    for (const id of selected.value) {
        data.append(`category-${id}`, 'on');
    }
    if (image.value) data.append('image', image.value);
    const res = await apiPost('/api/post/create', data);
    submitting.value = false;
    if (res) {
        router.push(`/post/view/${res.Posts[0].Id}`);
    }
}
</script>

<template>
    <div class="container">
        <form class="post-create" @submit.prevent="submit">
            <fieldset>
                <legend>New post</legend>
                <input type="text" v-model="title" name="title" placeholder="Your title here" required />
                <template v-if="categories.length">
                    <div class="inline" v-for="category in categories" :key="category.Id">
                        <input
                            type="checkbox"
                            :id="`category-${category.Id}`"
                            :name="`category-${category.Id}`"
                            :checked="selected.includes(category.Id)"
                            @change="toggle(category.Id)"
                        />
                        <label :for="`category-${category.Id}`">{{ category.Name }}</label>
                    </div>
                </template>
                <p v-else>No categories available</p>
                <textarea v-model="body" name="body" placeholder="Your post here" required></textarea>
                <input
                    type="file"
                    name="image"
                    accept="image/jpeg,image/png,image/gif"
                    @change="onImage"
                />
                <input type="submit" :disabled="submitting" value="Create post" />
            </fieldset>
        </form>
    </div>
</template>

<style scoped></style>
