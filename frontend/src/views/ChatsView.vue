<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useUser } from '@/composables/useUser.js';
import { useChat } from '@/composables/useChat.js';
import { apiFetch } from '@/utils/api.js';
import { DateToLocale, firstOrNull } from '@/utils/utils.js';

const { user } = useUser();
const { ensureWs, onChatEvent, hasUnread, myGroups } = useChat();

const chats = ref([]);
const loaded = ref(false);

async function refresh() {
    const data = await apiFetch('/api/chats');
    if (data) {
        chats.value = data.Users ?? [];
    }
    loaded.value = true;
}

const unsubscribes = [];

onMounted(async () => {
    ensureWs(user.value.Id);
    await refresh();
    unsubscribes.push(onChatEvent('chat_message', refresh));
});

onUnmounted(() => {
    unsubscribes.forEach((unsubscribe) => unsubscribe());
    unsubscribes.length = 0;
});

function lastMessage(u) {
    return firstOrNull(u.ChatMessages);
}

function previewPrefix(u) {
    return lastMessage(u)?.SenderId === user.value.Id ? 'You: ' : '';
}
</script>

<template>
    <div class="container padding-sides-1vw">
        <h2>Chats</h2>
        <p v-if="loaded && chats.length === 0 && myGroups.length === 0">No conversations yet.</p>
        <template v-else>
            <template v-if="chats.length">
                <h3>Direct messages</h3>
                <ul class="chats-list">
                    <li v-for="u in chats" :key="u.Id">
                        <router-link :to="`/chat/${u.Id}`">
                            <strong>{{ u.Username }}</strong>
                            <span v-if="hasUnread(u.Id)" class="new-message">●</span>
                            <span class="timestamp">{{ DateToLocale(lastMessage(u)?.Timestamp) }}</span>
                            <span class="preview">{{ previewPrefix(u) }}{{ lastMessage(u)?.Body }}</span>
                        </router-link>
                    </li>
                </ul>
            </template>
            <template v-if="myGroups.length">
                <h3>Groups</h3>
                <ul class="chats-list">
                    <li v-for="g in myGroups" :key="g.Id">
                        <router-link :to="`/group/${g.Id}/chat`">
                            <strong>{{ g.Title }}</strong>
                            <span v-if="hasUnread(`g:${g.Id}`)" class="new-message">●</span>
                        </router-link>
                    </li>
                </ul>
            </template>
        </template>
    </div>
</template>

<style scoped>
.chats-list .preview {
    display: block;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
}
</style>
