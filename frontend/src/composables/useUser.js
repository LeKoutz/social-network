import { ref } from 'vue';
import { apiFetch } from '@/utils/api.js';
import { useRouter } from 'vue-router';

const user = ref({ LoggedIn: false });

export function useUser() {
    const router = useRouter();

    function setUser(newUser) {
        user.value = newUser;
    }

    async function logoutUser() {
        const data = await apiFetch('/api/user/logout');
        if (data) {
            user.value = { LoggedIn: false };
            router.push('/')
            // TODO: Update Websocket
        }
    }

    return { user, setUser, logoutUser };
}
