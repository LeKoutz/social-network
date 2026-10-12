<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { apiPost } from '@/utils/api.js';
import PostVisibilityPanel from '@/components/PostVisibilityPanel.vue';
import { useCategories } from '@/composables/useCategories.js';

const props = defineProps({
    group: {
        type: Object,
        default: null,
    },
});

const { categories } = useCategories();
const visibility = ref('public');

const router = useRouter();

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
    if (props.group) data.append('group-id', props.group.Id);
    if (image.value) data.append('image', image.value);
    data.append('visibility', props.group ? 'public' : visibility.value);
    const res = await apiPost('/api/post/create', data);
    submitting.value = false;
    if (res) {
        const id = props.group ? props.group.Id : res.Posts[0].Id;
        router.push(props.group ? `/group/view/${id}` : `/post/view/${id}`);
    }
}
</script>

<template>
    <div class="flex flex-column gap-20">
        <form class="post-create" @submit.prevent="submit">
            <fieldset>
                <legend>New post</legend>
                <PostVisibilityPanel
                    v-if="!group"
                    v-model="visibility"
                    :disabled="submitting"
                />
                <p v-else>Only group members can see this post.</p>
                <template v-if="!group">
                    <input
                        type="text"
                        v-model="title"
                        name="title"
                        placeholder="Your title here"
                        required
                    />
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
                </template>
                <input
                    v-else
                    type="text"
                    v-model="title"
                    name="title"
                    placeholder="Your title here"
                    required
                />
                <textarea
                    v-model="body"
                    name="body"
                    placeholder="Your post here"
                    required
                ></textarea>
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
