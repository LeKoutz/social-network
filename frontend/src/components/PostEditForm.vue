<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { apiPost } from '@/utils/api.js';

const props = defineProps({
    post: {
        type: Object,
        required: true,
    },
    categories: {
        type: Array,
        required: true,
    },
});

const router = useRouter();

const title = ref(props.post.Title);
const body = ref(props.post.Body);
const image = ref(null);
const selected = ref(props.post.Categories.map((category) => category.Id));
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
    data.append('post-id', props.post.Id);
    data.append('save-post', '1');
    data.append('title', title.value);
    data.append('body', body.value);
    for (const id of selected.value) {
        data.append(`category-${id}`, 'on');
    }
    if (image.value) data.append('image', image.value);
    const res = await apiPost(`/api/post/edit/${props.post.Id}`, data);
    submitting.value = false;
    if (res) {
        router.push(`/post/view/${props.post.Id}`);
    }
}
</script>

<template>
    <div class="container">
        <form class="post-edit" @submit.prevent="submit">
            <fieldset>
                <legend>Edit Post</legend>
                <input
                    type="text"
                    v-model="title"
                    name="title"
                    placeholder="Your title here"
                    required
                />
                <template v-if="!post.GroupId">
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
                <img
                    v-if="post.ImagePath"
                    :src="`/${post.ImagePath}`"
                    alt="Post image"
                    style="max-width: 100%"
                />
                <textarea v-model="body" name="body" placeholder="Your post here" required></textarea>
                <input
                    type="file"
                    name="image"
                    accept="image/jpeg,image/png,image/gif"
                    @change="onImage"
                />
                <input type="submit" :disabled="submitting" value="Save post" />
            </fieldset>
        </form>
    </div>
</template>

<style scoped></style>
