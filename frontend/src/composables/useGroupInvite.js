import { ref } from 'vue';
import { apiGet, apiPost, apiDelete, apiPatch } from '@/utils/api.js';

export function useGroupInvite() {
    const loading = ref(false);
    const error = ref(null);
    const success = ref(false);
    const groupJoinRequests = ref(false);
    const groupInvitations = ref(false);

    async function sendInvite(groupId, userId) {
        loading.value = true;
        error.value = null;
        success.value = false;
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

    async function sendJoinRequest(groupId) {
        loading.value = true;
        error.value = null;
        success.value = false;
        try {
            const payload = new FormData();
            payload.append('group-id', groupId);
            const data = await apiPost('/api/group/join', payload);
            if (data.Error?.Has) throw new Error(data.Error.Message);
            success.value = true;
            return true;
        } catch (err) {
            console.error('Failed to send join request:', err);
            error.value = err.message || 'An error occurred while sending the join request.';
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

    async function withdrawJoinRequest(id) {
        loading.value = true;
        error.value = null;
        success.value = false;
        try {
            const payload = new FormData();
            payload.append('id', id);
            const data = await apiDelete('/api/group/join/withdraw', payload);
            if (data.Error?.Has) throw new Error(data.Error.Message);
            success.value = true;
            return true;
        } catch (err) {
            console.error('Failed to withdraw group join request:', err);
            error.value = err.message || 'An error occurred while withdrawing group join request.';
            return false;
        } finally {
            loading.value = false;
        }
    }

    async function getGroupJoinRequests() {
        loading.value = true;
        error.value = null;
        success.value = false;
        try {
            const data = await apiGet('/api/groups/joins');
            if (data.Error?.Has) throw new Error(data.Error.Message);
            groupJoinRequests.value = data.GroupInvitations;
        } catch (err) {
            console.error('Failed to get group join requests:', err);
            error.value = err.message || 'An error occurred while getting group join requests.';
        } finally {
            loading.value = false;
        }
    }

    async function getGroupInvitations() {
        loading.value = true;
        error.value = null;
        success.value = false;
        try {
            const data = await apiGet('/api/groups/invites');
            if (data.Error?.Has) throw new Error(data.Error.Message);
            groupInvitations.value = data.GroupInvitations;
        } catch (err) {
            console.error('Failed to get group invites:', err);
            error.value = err.message || 'An error occurred while getting group invites.';
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
        sendJoinRequest,
        declineInvite,
        withdrawJoinRequest,
        getGroupJoinRequests,
        getGroupInvitations,
        groupInvitations,
        groupJoinRequests,
    };
}
