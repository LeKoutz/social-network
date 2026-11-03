<script setup>
import { ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { apiPost } from '@/utils/api.js';
import PostVisibilityPanel from '@/components/PostVisibilityPanel.vue';
import ImagePreview from './ImagePreview.vue';
import GenericButton from './GenericButton.vue';
import ImageIcon from './Icons/ImageIcon.vue';

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

const initialCategoryIds = (props.post.Categories ?? []).map((category) => category.Id);

const visibility = ref(props.post.Visibility ?? 'public');
const allowed = ref([]);
const image = ref(null);
const fileInput = useTemplateRef('file-input');
const submitting = ref(false);

function onImage(e) {
    image.value = e.target.files[0] ?? null;
    e.target.value = '';
}

async function submitForm(e) {
    submitting.value = true;
    const data = new FormData(e.target);
    if (image.value) data.append('image', image.value);

    try {
        const res = await apiPost(`/api/post/edit/${props.post.Id}`, data);
        if (res && !res.Error?.Has) {
            router.push(`/post/view/${props.post.Id}`);
        }
    } finally {
        submitting.value = false;
    }
}
</script>

<template>
    <form @submit.prevent="submitForm">
        <fieldset>
            <p>Edit post</p>
            <input type="text" name="title" :value="post.Title" placeholder="Your title here" required />
            <textarea name="body" :value="post.Body" placeholder="Your post here" required></textarea>

            <ImagePreview v-if="image" v-model="image" />
            <img v-else-if="post.ImagePath" :src="`/${post.ImagePath}`" alt="" class="current-image" />

            <div v-if="!post.GroupId" class="category-selector flex flex-wrap gap-10">
                <div
                    v-for="category in categories"
                    :key="category.Id"
                    class="flex align-items-center gap-5 text-color-secondary"
                >
                    <input
                        type="checkbox"
                        :id="`category-${category.Id}`"
                        :name="`category-${category.Id}`"
                        :checked="initialCategoryIds.includes(category.Id)"
                    />
                    <label :for="`category-${category.Id}`">{{ category.Name }}</label>
                </div>
            </div>
            <small v-else>Visibility is set to <strong>Public</strong> for group posts.</small>

            <input type="hidden" name="post-id" :value="post.Id" />
            <input type="hidden" name="save-post" value="1" />
            <input type="hidden" name="visibility" :value="post.GroupId ? 'public' : visibility" />
            <template v-if="!post.GroupId && visibility === 'limited'">
                <input v-for="id in allowed" :key="id" type="hidden" name="allowed-user" :value="id" />
            </template>

            <div class="flex flex-wrap justify-content-space-between gap-6">
                <span class="flex flex-wrap align-items-center gap-6">
                    <PostVisibilityPanel
                        v-if="!post.GroupId"
                        v-model="visibility"
                        v-model:allowed="allowed"
                        :disabled="submitting"
                    />
                    <GenericButton type="button" variant="unfilled" @click="fileInput.click()">
                        <template #leading><ImageIcon /></template>
                        {{ post.ImagePath ? 'Replace image' : 'Image/GIF' }}
                    </GenericButton>
                    <input
                        ref="file-input"
                        type="file"
                        accept="image/jpeg,image/png,image/gif"
                        hidden
                        @change="onImage"
                    />
                </span>
                <div class="flex align-items-center gap-6">
                    <GenericButton type="button" variant="unfilled" @click="router.push(`/post/view/${post.Id}`)">
                        Cancel
                    </GenericButton>
                    <GenericButton type="submit" :disabled="submitting">Save post</GenericButton>
                </div>
            </div>
        </fieldset>
    </form>
</template>

<style scoped>
.current-image {
    max-width: 100%;
    max-height: 200px;
    border-radius: 2px;
    width: fit-content;
}

.category-selector input[type="checkbox"] {
    appearance: none;
    width: 14px;
    height: 14px;
    border-radius: 2px;
    background-color: var(--bg-input);
}

.category-selector input:checked {
    background-color: var(--accent);
}

.category-selector input:checked + label {
    color: var(--text-primary);
}
</style>
