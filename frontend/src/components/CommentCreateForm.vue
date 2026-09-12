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
const submitting = ref(false);

async function submit() {
    if (!body.value.trim()) return;
    submitting.value = true;
    const data = await apiPost(
        '/api/comment/create',
        new URLSearchParams({
            'post-id': props.post.Id,
            comment: body.value,
        })
    );
    submitting.value = false;
    if (data) {
        body.value = '';
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
            <input type="submit" :disabled="submitting" value="Post comment" />
        </fieldset>
    </form>
</template>

<style scoped></style>
