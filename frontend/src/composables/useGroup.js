import { ref } from 'vue';
import { apiGet, apiPost, apiPatch } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';

export function useGroup() {
    const loading = ref(false);
    const error = ref(null);
    const group = ref(null);
    const posts = ref(null);

    async function loadGroup(id) {
        loading.value = true;
        error.value = null;
        try {
            const data = await apiGet(`/api/group/view/${id}`);
            if (data) {
                group.value = firstOrNull(data.Groups);
                posts.value = data.Posts;
                if (!group.value) error.value = 'Group not found';
            }
        } catch (err) {
            console.error('Failed to get group:', err);
            error.value = err.message || 'An error occurred while getting group.';
        } finally {
            loading.value = false;
        }
    }

    return {
        loading,
        error,
        group,
        posts,
        loadGroup,
    };
}
