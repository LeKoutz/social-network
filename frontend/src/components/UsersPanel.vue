<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { useChat } from '@/composables/useChat.js';
import { apiFetch } from '@/utils/api.js';

const { user } = useUser();
const { ensureWs, onChatEvent, hasUnread, totalUnread } = useChat();

const users = ref([]);
const collapsed = ref(false);

const list = computed(() => {
    const sorted = [...users.value].sort((a, b) => {
        const aHasMessages = a.LastMessageTimestamp > 0;
        const bHasMessages = b.LastMessageTimestamp > 0;
        if (aHasMessages && bHasMessages) {
            const timestampDifference =
                b.LastMessageTimestamp - a.LastMessageTimestamp;
            if (timestampDifference !== 0) {
                return timestampDifference;
            }
        }
        if (aHasMessages && !bHasMessages) {
            return -1;
        }
        if (!aHasMessages && bHasMessages) {
            return 1;
        }
        return String(a.Username).localeCompare(String(b.Username));
    });
    return {
        online: sorted.filter((u) => u.LoggedIn),
        offline: sorted.filter((u) => !u.LoggedIn),
    };
});

function lastMessageLabel(u) {
    if (!u.LastMessageTimestamp) return '';
    const date = new Date(u.LastMessageTimestamp * 1000);
    const day = date.getDate();
    const month = date.getMonth() + 1;
    const year = String(date.getFullYear()).slice(-2);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${day}-${month}-${year}, ${hours}:${minutes}`;
}

async function refresh() {
    const data = await apiFetch('/api/users');
    if (data) {
        users.value = data.Users ?? [];
    }
}

const unsubscribes = [];

onMounted(async () => {
    ensureWs(user.value.Id);
    await refresh();
    unsubscribes.push(onChatEvent('user_status', refresh));
    unsubscribes.push(onChatEvent('chat_message', refresh));
});

onUnmounted(() => {
    unsubscribes.forEach((unsubscribe) => unsubscribe());
    unsubscribes.length = 0;
});
</script>

<template>
    <div v-if="user.LoggedIn" class="users-panel">
        <div class="users-panel-header">
            <h3>Users <span v-if="totalUnread()" class="new-message">●</span></h3>
            <button id="collapse-panel" @click="collapsed = !collapsed">−</button>
        </div>
        <div class="users-panel-inner" :hidden="collapsed">
            <h4>Online</h4>
            <ul>
                <li
                    v-for="u in list.online"
                    :key="u.Id"
                    :data-user-id="u.Id"
                >
                    <router-link :to="`/chat/${u.Id}`">{{ u.Username }}</router-link>
                    <span v-if="hasUnread(u.Id)" class="new-message">●</span>
                    <span v-if="u.LastMessageTimestamp" class="last-message-time">
                        - last msg {{ lastMessageLabel(u) }}
                    </span>
                </li>
            </ul>
            <h4>Offline</h4>
            <ul>
                <li
                    v-for="u in list.offline"
                    :key="u.Id"
                    :data-user-id="u.Id"
                    class="offline"
                >
                    <router-link :to="`/chat/${u.Id}`">{{ u.Username }}</router-link>
                    <span v-if="hasUnread(u.Id)" class="new-message">●</span>
                    <span v-if="u.LastMessageTimestamp" class="last-message-time">
                        - last msg {{ lastMessageLabel(u) }}
                    </span>
                </li>
            </ul>
        </div>
    </div>
</template>

<style scoped></style>