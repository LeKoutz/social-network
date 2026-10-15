import { ref } from 'vue';
import { apiGet } from '@/utils/api.js';
import { firstOrNull } from '@/utils/utils.js';

export function useProfile() {
    const profile = ref(null);
    const loading = ref(false);
    const error = ref(null);

    async function loadProfile(id) {
        loading.value = true;
        error.value = null;
        try {
            const data = await apiGet(`/api/profile/view/${id}`);
            if (data) {
                profile.value = firstOrNull(data.Profiles);
                if (!profile.value) error.value = 'Profile not found';
            }
        } catch (err) {
            error.value = err.message || 'Something went wrong';
        } finally {
            loading.value = false;
        }
    }

    return {
        profile,
        loading,
        error,
        loadProfile,
    };
}
