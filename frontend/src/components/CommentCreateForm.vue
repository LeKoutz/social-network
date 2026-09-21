<script setup>
import { ref } from 'vue';
import { apiPost } from '@/utils/api.js';

const props = defineProps({
    post: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(['created']);

const body = ref('');
const image = ref(null);
const submitting = ref(false);

function onImage(e) {
    image.value = e.target.files[0] ?? null;
}

async function submit() {
    if (!body.value.trim()) return;
    submitting.value = true;
    const data = new FormData();
    data.append('post-id', props.post.Id);
    data.append('comment', body.value);
    if (image.value) data.append('image', image.value);
    const res = await apiPost('/api/comment/create', data);
    submitting.value = false;
    if (res) {
        body.value = '';
        image.value = null;
        emit('created');
    }
}
</script>

<template>
    <form class="comment-create" @submit.prevent="submit">
        <fieldset>
            <legend>Leave a comment</legend>
            <textarea
                v-model="body"
                name="comment"
                placeholder="Enter your comment"
                required
            ></textarea>
            <input
                type="file"
                name="image"
                accept="image/jpeg,image/png,image/gif"
                @change="onImage"
            />
            <input type="submit" :disabled="submitting" value="Post comment" />
        </fieldset>
    </form>
</template>

<style scoped></style>
