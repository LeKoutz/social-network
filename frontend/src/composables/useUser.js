import { ref } from 'vue';
import { apiGet, apiPost } from '@/utils/api.js';
import router from '@/router/index.js';
import { useChat } from '@/composables/useChat.js';

const user = ref({ LoggedIn: false });
let resolveUserLoaded;
const userLoaded = new Promise((resolve) => {
    resolveUserLoaded = resolve;
});

export function useUser() {
    const loading = ref(false);
    const alert = ref({ Error: { Has: false, Message: '' } });

    function setUser(newUser) {
        user.value = newUser;
    }

    async function logoutUser() {
        try {
            const data = await apiGet('/api/user/logout');
            if (data) {
                user.value = { LoggedIn: false };
                router.push('/');
                useChat().disconnectWS();
            }
            return true;
        } catch (e) {
            console.log('Logout failed', e);
            return false;
        }
    }

    async function registerUser(payload) {
        loading.value = true;
        alert.value = { Error: { Has: false, Message: '' } };

        try {
            const data = await apiPost('/api/user/register', payload);
            alert.value = data;
            return data;
        } catch (error) {
            alert.value = { Error: { Has: true, Message: 'Network error. Please try again.' } };
            return null;
        } finally {
            loading.value = false;
        }
    }

    return {
        loading,
        alert,
        user,
        setUser,
        logoutUser,
        userLoaded,
        resolveUserLoaded,
        registerUser,
    };
}
