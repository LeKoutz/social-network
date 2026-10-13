import { ref } from 'vue';

import { apiPost } from '@/utils/api.js';

const post = ref(null);

export function usePost() {
    function setPost(data) {
        if (data.Error?.Has) {
            post.value = { type: 'Error', message: data.Error.Message };
        } else if (data.Message?.Has) {
            post.value = { type: data.Message.Type, message: data.Message.Content };
        }
    }

    function clear() {
        post.value = null;
    }

    async function deletePost(id) {
        const data = await apiPost('/api/post/delete', new URLSearchParams({ 'post-id': id }));
        return data;
    }

    return { post, setPost, deletePost };
}
