import { apiFetch } from '@/utils/api.js';
import { computed, ref } from 'vue';

const notifications = ref([])

export function useNotifications() {
    // TODO: Replace backends UnreadNotificationsCount?
    const unreadNotificationsCount = computed(() => {
        notifications.value.filter((n) => !n.Read.length)
    })

    async function getNotifications() {
        const data = await apiFetch('/api')
        if (data && !data.Error.Has) {
            const newNotifications = data.User.Notifications
            setNotifications(newNotifications)
        }
    }

    function setNotifications(newNotifications) {
        notifications.value = newNotifications ?? []
    }

    async function markNotificationsAsRead(e) {
        const response = await fetch('/api/user/notifications', {
            method: 'POST',
            body: new URLSearchParams(new FormData(e.target))
        });
        const data = await response.json();
        if (data && !data.Error.Has) {
            notifications.value.forEach((notification) => {
                notification.Read = true
            })
        }
    }

    return { notifications, unreadNotificationsCount, getNotifications, setNotifications, markNotificationsAsRead }
}