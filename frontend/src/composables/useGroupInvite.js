import { ref } from 'vue';
import { apiPost, apiPatch } from '@/utils/api.js';

export function useGroupInvite() {
    const loading = ref(false);
    const error = ref(null);
    const success = ref(false);

    async function sendInvite(groupId, userId) {
        loading.value = true;
        error.value = null;
        success.value = false;
        console.log(groupId, userId);

        try {
            const payload = new FormData();
            payload.append('group-id', groupId);
            payload.append('to-user-id', userId);

            const data = await apiPost('/api/group/invite', payload);
            if (data.Error?.Has) throw new Error(data.Error.Message);

            success.value = true;
            return true;
        } catch (err) {
            console.error('Failed to send invite:', err);
            error.value = err.message || 'An error occurred while sending the invite.';
            return false;
        } finally {
            loading.value = false;
        }
    }

    async function acceptInvite(inviteId) {
        loading.value = true;
        error.value = null;
        success.value = false;

        try {
            const payload = new FormData();
            payload.append('id', inviteId);

            const data = await apiPatch('/api/group/invite/accept', payload);
            if (data.Error?.Has) throw new Error(data.Error.Message);

            success.value = true;
            return true;
        } catch (err) {
            console.error('Failed to accept invite:', err);
            error.value = err.message || 'An error occurred while accepting the invite.';
            return false;
        } finally {
            loading.value = false;
        }
    }

    async function declineInvite(inviteId) {
        loading.value = true;
        error.value = null;
        success.value = false;

        try {
            const payload = new FormData();
            payload.append('id', inviteId);

            const data = await apiPatch('/api/group/invite/decline', payload);
            if (data.Error?.Has) throw new Error(data.Error.Message);

            success.value = true;
            return true;
        } catch (err) {
            console.error('Failed to decline invite:', err);
            error.value = err.message || 'An error occurred while declining the invite.';
            return false;
        } finally {
            loading.value = false;
        }
    }


    return {
        loading,
        error,
        success,
        sendInvite,
        acceptInvite,
        declineInvite,
    };
}
