import { ref } from 'vue';
import { apiFetch } from '@/utils/api.js';

const user = ref({ LoggedIn: false });

export function useUser() {
    function setUser(newUser) {
        user.value = newUser;
    }

    async function logoutUser() {
        try {
            const data = await apiFetch('/api/user/logout');
            if (data) {
                user.value = { LoggedIn: false };
                router.push('/')
                // TODO: Update Websocket
            }
            return true;
        } catch (e) {
            console.log('Logout failed', e);
            return false;
        }
    }

    return { user, setUser, logoutUser };
}
