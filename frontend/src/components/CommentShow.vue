<script setup>
import { ref } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { useAlerts } from '@/composables/useAlerts.js';
import { apiPost } from '@/utils/api.js';
import { DateToLocale } from '@/utils/utils.js';
import CommentEditForm from '@/components/CommentEditForm.vue';
import CommentDeleteForm from '@/components/CommentDeleteForm.vue';
import CommentReactForm from '@/components/CommentReactForm.vue';

const props = defineProps({
    post: {
        type: Object,
        required: true,
    },
    comment: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(['updated']);

const { user } = useUser();
const { setAlert } = useAlerts();

const editing = ref(false);
const body = ref(props.comment.Body);

const isOwner = () =>
    user.value.LoggedIn && String(user.value.Id) === String(props.comment.UserId);

function toggleEdit() {
    editing.value = !editing.value;
    if (editing.value) body.value = props.comment.Body;
}

async function saveEdit() {
    if (!body.value.trim()) {
        setAlert({ Error: { Has: true, Message: 'Comment must not be empty' } });
        return;
    }
    const data = await apiPost(
        '/api/comment/edit',
        new URLSearchParams({
            'post-id': props.post.Id,
            'comment-id': props.comment.Id,
            'save-comment': '1',
            comment: body.value,
        })
    );
    if (data) {
        editing.value = false;
        emit('updated');
    }
}

async function remove() {
    const data = await apiPost(
        '/api/comment/delete',
        new URLSearchParams({
            'post-id': props.post.Id,
            'comment-id': props.comment.Id,
        })
    );
    if (data) emit('updated');
}

async function react(action) {
    const data = await apiPost(
        '/api/comment/react',
        new URLSearchParams({
            'post-id': props.post.Id,
            'comment-id': props.comment.Id,
            action,
        })
    );
    if (data) emit('updated');
}
</script>

<template>
    <div class="comment" :id="`comment-${comment.Id}`">
        <span>{{ comment.Username }} ({{ DateToLocale(comment.Timestamp) }})</span>
        <div class="manage-comment">
            <template v-if="isOwner()">
                <CommentEditForm @edit="toggleEdit" />
                <CommentDeleteForm @delete="remove" />
            </template>
        </div>
        <form v-if="editing" class="comment-edit" @submit.prevent="saveEdit">
            <textarea v-model="body" name="comment" required></textarea>
            <button type="submit">Save Comment</button>
            <button type="button" @click="editing = false">Cancel</button>
        </form>
        <pre v-else>{{ comment.Body }}</pre>
        <CommentReactForm :comment="comment" @react="react" />
    </div>
</template>

<style scoped></style>
