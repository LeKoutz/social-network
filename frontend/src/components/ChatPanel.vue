<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAppData } from '@/composables/useAppData.js';
import { useUser } from '@/composables/useUser.js';
import { useChat } from '@/composables/useChat.js';
import { apiGet } from '@/utils/api.js';
import ChevronDownIcon from './Icons/ChevronDownIcon.vue';
import DoubleChevronsLeftIcon from './Icons/DoubleChevronsLeftIcon.vue';
import DoubleChevronsRightIcon from './Icons/DoubleChevronsRightIcon.vue';
import ChatBubbleIcon from './Icons/ChatBubbleIcon.vue';

const route = useRoute();
const { user } = useUser();
const { ensureWs, onChatEvent, hasUnread, totalUnread, myGroups } = useChat();
const { isMobileView } = useAppData();
const users = ref([]);
const loaded = ref(false);

const collapsed = ref(isMobileView.value);

const list = computed(() => {
    const sorted = [...users.value].sort((a, b) => {
        const aHasMessages = a.LastMessageTimestamp > 0;
        const bHasMessages = b.LastMessageTimestamp > 0;
        if (aHasMessages && bHasMessages) {
            const timestampDifference = b.LastMessageTimestamp - a.LastMessageTimestamp;
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
    const lastMessage = new Date(u.LastMessageTimestamp * 1000);
    const today = new Date();
    const yesterday = new Date();
    yesterday.setDate(today.getDate() - 1);
    const day = lastMessage.getDate();
    const month = lastMessage.getMonth() + 1;
    const year = String(lastMessage.getFullYear()).slice(-2);
    const hours = String(lastMessage.getHours()).padStart(2, '0');
    const minutes = String(lastMessage.getMinutes()).padStart(2, '0');
    switch (lastMessage.toDateString()) {
        case today.toDateString():
            return `${hours}:${minutes}`;
        case yesterday.toDateString():
            return 'yesterday';
        default:
            return `${day}-${month}-${year}`;
    }
}

async function refresh() {
    const data = await apiGet('/api/users');
    if (data) {
        users.value = data.Users ?? [];
    }
    loaded.value = true;
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

watch(isMobileView, (mobile) => {
    collapsed.value = mobile;
});
watch(
    () => route.fullPath,
    () => (collapsed.value = isMobileView.value),
);
</script>

<template>
    <button v-if="collapsed && isMobileView" class="chat-bar" @click="collapsed = false">
        <span class="bar-title">Chat</span>
        <span v-if="totalUnread()" class="dot dot-title"></span>
        <span class="online-count">{{ list.online.length }} online</span>
        <span class="bar-chevron"><ChevronDownIcon /></span>
    </button>
    <div v-else :class="['chat-panel', { collapsed }]">
        <div class="header">
            <template v-if="!collapsed">
                <h3 class="title">Chat</h3>
                <span v-if="totalUnread()" class="dot dot-title"></span>
                <router-link to="/chats" class="see-all">See all chats</router-link>
            </template>
            <button class="toggle" @click="collapsed = !collapsed">
                <DoubleChevronsRightIcon v-if="!collapsed" />
                <DoubleChevronsLeftIcon v-else />
            </button>
            <div v-if="collapsed" class="chat-icon">
                <ChatBubbleIcon />
                <span v-if="totalUnread()" class="dot dot-corner"></span>
            </div>
        </div>

        <template v-if="!collapsed">
            <span class="section-title">Online</span>
            <ul>
                <li
                    v-for="u in list.online"
                    :key="u.Id"
                    :class="['row', { unread: hasUnread(u.Id) }]"
                >
                    <span class="status"></span>
                    <router-link :to="`/chat/${u.Id}`">{{ u.Username }}</router-link>
                    <span v-if="hasUnread(u.Id)" class="dot"></span>
                    <span v-if="u.LastMessageTimestamp" class="time">
                        {{ lastMessageLabel(u) }}
                    </span>
                </li>
            </ul>

            <span class="section-title">Offline</span>
            <ul>
                <li
                    v-for="u in list.offline"
                    :key="u.Id"
                    :class="['row', 'offline', { unread: hasUnread(u.Id) }]"
                >
                    <span class="status"></span>
                    <router-link :to="`/chat/${u.Id}`">{{ u.Username }}</router-link>
                    <span v-if="hasUnread(u.Id)" class="dot"></span>
                    <span v-if="u.LastMessageTimestamp" class="time">
                        {{ lastMessageLabel(u) }}
                    </span>
                </li>
            </ul>

            <template v-if="myGroups.length">
                <span class="section-title">Groups</span>
                <ul>
                    <li
                        v-for="g in myGroups"
                        :key="g.Id"
                        :class="['row', { unread: hasUnread(`g:${g.Id}`) }]"
                    >
                        <span class="hash">#</span>
                        <router-link :to="`/group/${g.Id}/chat`">{{ g.Title }}</router-link>
                        <span v-if="hasUnread(`g:${g.Id}`)" class="dot"></span>
                    </li>
                </ul>
            </template>
        </template>
    </div>
</template>

<style scoped>
.chat-panel {
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    gap: 2px;
    box-sizing: border-box;
    width: 220px;
    padding: 12px 16px 16px;
    overflow-y: auto;
    line-height: 1.2;
    background-color: var(--bg-base);
    border-left: 1px solid var(--accent-line);
}

.chat-panel.collapsed {
    width: 48px;
    padding: 12px 0 16px;
    align-items: center;
}

.header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-bottom: 8px;
}

.collapsed .header {
    flex-direction: column;
    gap: 14px;
    padding: 0;
}

.title {
    margin: 0;
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
}

.see-all {
    margin-left: auto;
    font-size: 13px;
    font-weight: 500;
    color: var(--accent);
    text-decoration: none;
}

.section-title {
    padding: 12px 0 4px;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
}

ul {
    margin: 0;
    padding: 0;
    list-style: none;
}

.row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 5px 0;
}

.row a {
    font-size: 14px;
    color: var(--accent);
    text-decoration: none;
}

.row.offline a {
    color: var(--text-secondary);
}

.row.unread a {
    font-weight: 600;
}

.status {
    flex-shrink: 0;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--accent);
}

.offline .status {
    background-color: var(--text-muted);
}

.hash {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-muted);
}

.dot {
    flex-shrink: 0;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: var(--danger);
}

.dot-title {
    width: 8px;
    height: 8px;
}

.time {
    margin-left: auto;
    font-size: 11px;
    white-space: nowrap;
    color: var(--text-muted);
}

.chat-icon {
    position: relative;
    width: 22px;
    height: 22px;
    color: var(--text-secondary);
}

.dot-corner {
    position: absolute;
    top: -2px;
    right: -2px;
    width: 9px;
    height: 9px;
    box-shadow: 0 0 0 2px var(--bg-elevated);
}

@media (max-width: 768px) {
    .chat-panel {
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        z-index: 10;
        width: auto;
        max-height: 60vh;
        border-left: none;
        border-top: 1px solid var(--accent-line);
    }

    .chat-bar {
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        z-index: 10;
        display: flex;
        align-items: center;
        gap: 8px;
        height: 52px;
        padding: 0 16px;
        border: none;
        border-top: 1px solid var(--accent-line);
        background-color: var(--bg-base);
        font: inherit;
        cursor: pointer;
    }

    .bar-title {
        font-size: 16px;
        font-weight: 700;
        color: var(--text-primary);
    }

    .online-count {
        margin-left: auto;
        font-size: 13px;
        color: var(--text-muted);
    }

    .bar-chevron {
        display: flex;
        padding: 4px 6px;
        border: 1px solid var(--border);
        border-radius: 2px;
        background-color: var(--bg-input);
        color: var(--accent);
        transform: rotate(180deg);
    }
}
</style>
