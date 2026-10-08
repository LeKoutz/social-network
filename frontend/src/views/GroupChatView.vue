<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useUser } from '@/composables/useUser.js';
import { useChat } from '@/composables/useChat.js';
import { apiFetch } from '@/utils/api.js';
import ChatMessageList from '@/components/ChatMessageList.vue';

const route = useRoute();
const { user } = useUser();
const { ensureWs, onChatEvent, sendGroupMessage, notifyGroupMessageRead, markGroupAsRead, setActiveChat, myGroups } = useChat();

const messages = ref([]);
const body = ref('');
const historyEnded = ref(false);
const canChat = ref(null);

const groupTitle = computed(() => myGroups.value.find((g) => g.Id === groupId())?.Title ?? '');

function groupId() {
    return Number(route.params.id);
}

async function loadHistory(offset) {
    const data = await apiFetch(`/api/group/chat/${groupId()}?offset=${offset}`);
    if (!data) {
        canChat.value = false;
        return [];
    }
    canChat.value = data.Groups?.[0]?.Member ?? false;
    return data.Groups?.[0]?.ChatMessages ?? [];
}

async function initChat() {
    setActiveChat(`g:${groupId()}`);
    historyEnded.value = false;
    messages.value = await loadHistory(0);
    markGroupAsRead(groupId());
}

watch(() => route.params.id, initChat);

const unsubscribes = [];

onMounted(() => {
    ensureWs(user.value.Id);
    initChat();
    unsubscribes.push(
        onChatEvent('group_message', ({ msg, isForCurrent }) => {
            if (!isForCurrent) return;
            messages.value.push(msg);
            if (msg.SenderId !== user.value.Id) notifyGroupMessageRead(msg.Id, msg.GroupId);
        })
    );
});

onUnmounted(() => {
    setActiveChat(null);
    unsubscribes.forEach((unsubscribe) => unsubscribe());
});

async function loadOlder() {
    if (historyEnded.value) return;
    const older = await loadHistory(messages.value.length);
    if (!older.length) {
        historyEnded.value = true;
        return;
    }
    messages.value = [...older, ...messages.value];
}

function send() {
    if (!canChat.value) return;
    if (!body.value.trim()) return;
    sendGroupMessage(groupId(), body.value);
    body.value = '';
}
</script>

<template>
    <div class="container padding-sides-1vw group-chat">
        <h2><router-link :to="`/group/view/${groupId()}`">{{ groupTitle }}</router-link></h2>
        <ChatMessageList
            :messages="messages"
            :current-user-id="user.Id"
            @load-older="loadOlder"
        />
        <form id="chat-message" @submit.prevent="send">
            <input v-model="body" name="body" placeholder="Type a message..." required :hidden="!canChat"/>
            <input type="submit" value="Send" :hidden="!canChat"/>
            <p :hidden="canChat !== false">You are not a member of this group</p>
        </form>
    </div>
</template>

<style scoped>
.group-chat {
    width: 100%;
    max-width: 700px;
    margin-inline: auto;
}

#chat-message {
    margin-inline: auto;
}
</style>