<script setup>
import { ref, useTemplateRef } from 'vue';
import { apiPost } from '@/utils/api.js';
import PostVisibilityPanel from '@/components/PostVisibilityPanel.vue';
import { useCategories } from '@/composables/useCategories.js';
import ImagePreview from './ImagePreview.vue';
import GenericButton from './GenericButton.vue';
import ImageIcon from './Icons/ImageIcon.vue';

const props = defineProps({
    group: {
        type: Object,
        default: null,
    },
});

const emit = defineEmits(['created']);

const { categories } = useCategories();

const visibility = ref('public');
const allowed = ref([]);
const image = ref(null);
const fileInput = useTemplateRef('file-input');
const submitting = ref(false);

function onImage(e) {
    image.value = e.target.files[0] ?? null;
    e.target.value = '';
}

async function submitForm(e) {
    const form = e.target;
    submitting.value = true;
    const data = new FormData(form);
    if (image.value) data.append('image', image.value);

    try {
        const res = await apiPost('/api/post/create', data);
        if (res && !res.Error?.Has) {
            emit('created', res.Posts[0]);
            resetForm(form);
        }
    } finally {
        submitting.value = false;
    }
}

function resetForm(form) {
    form.reset();
    image.value = null;
    visibility.value = 'public';
    allowed.value = [];
}
</script>

<template>
    <form @submit.prevent="submitForm">
        <fieldset>
            <p>New post</p>
            <input type="text" name="title" placeholder="Your title here" required />
            <textarea name="body" placeholder="Your post here" required></textarea>
            <ImagePreview v-model="image" />

            <div v-if="!group" class="category-selector flex flex-wrap gap-10">
                <template v-if="categories.length">
                    <div
                        v-for="category in categories"
                        :key="category.Id"
                        class="flex align-items-center gap-5 text-color-secondary"
                    >
                        <input
                            type="checkbox"
                            :id="`category-${category.Id}`"
                            :name="`category-${category.Id}`"
                        />
                        <label :for="`category-${category.Id}`">{{ category.Name }}</label>
                    </div>
                </template>
                <p v-else>No categories available</p>
            </div>
            <small v-else>Only group members can see this post.</small>

            <input v-if="group" type="hidden" name="group-id" :value="group.Id" />
            <input type="hidden" name="visibility" :value="group ? 'public' : visibility" />
            <template v-if="!group && visibility === 'limited'">
                <input
                    v-for="id in allowed"
                    :key="id"
                    type="hidden"
                    name="allowed-user"
                    :value="id"
                />
            </template>

            <div class="flex flex-wrap justify-content-space-between gap-6">
                <span class="flex flex-wrap align-items-center gap-6">
                    <PostVisibilityPanel
                        v-if="!group"
                        v-model="visibility"
                        v-model:allowed="allowed"
                        :disabled="submitting"
                    />
                    <GenericButton type="button" variant="unfilled" @click="fileInput.click()">
                        <template #leading><ImageIcon /></template>
                        Image/GIF
                    </GenericButton>
                    <input
                        ref="file-input"
                        type="file"
                        accept="image/jpeg,image/png,image/gif"
                        hidden
                        @change="onImage"
                    />
                </span>
                <GenericButton type="submit" :disabled="submitting">Create post</GenericButton>
            </div>
        </fieldset>
    </form>
</template>

<style scoped>
.category-selector input[type='checkbox'] {
    appearance: none;
    width: 14px;
    height: 14px;
    border-radius: 2px;
    background-color: var(--bg-input);
}

.category-selector input:checked {
    background-color: var(--accent);
}

.category-selector input:checked label {
    color: var(--text-primary);
}
</style>
