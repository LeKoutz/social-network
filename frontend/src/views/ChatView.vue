<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useRoute } from 'vue-router';
import { useUser } from '@/composables/useUser.js';
import { useChat } from '@/composables/useChat.js';
import { apiFetch } from '@/utils/api.js';
import { throttle, DateToLocale } from '@/utils/utils.js';

const route = useRoute();
const { user } = useUser();
const { ensureWs, onChatEvent, sendChatMessage, notifyMessageRead, markMessageAsRead, setActiveChat } = useChat();

const messages = ref([]);
const recipientUsername = ref('');
const body = ref('');
const historyEnded = ref(false);
const container = ref(null);

function chatUserId() {
    return Number(route.params.id);
}

async function loadHistory(offset) {
    const data = await apiFetch(`/api/chat/${chatUserId()}?offset=${offset}`);
    if (!data) return [];
    const msgs = data.User?.ChatMessages ?? [];
    msgs.forEach((m) => {
        if (m.RecipientId === user.value.Id) markMessageAsRead(m);
    });
    return msgs;
}

async function loadRecipientUsername() {
    const data = await apiFetch('/api/users');
    const found = (data?.Users ?? []).find((u) => u.Id === chatUserId());
    if (found) recipientUsername.value = found.Username;
}

async function initChat() {
    setActiveChat(chatUserId());
    historyEnded.value = false;
    messages.value = await loadHistory(0);
    loadRecipientUsername();
    await nextTick();
    scrollToBottom();
}

watch(() => route.params.id, initChat);

const unsubscribes = [];

onMounted(() => {
    ensureWs(user.value.Id);
    setActiveChat(chatUserId());
    initChat();
    unsubscribes.push(
        onChatEvent('chat_message', ({ msg, isForCurrent }) => {
            if (!isForCurrent) return;
            messages.value.push(msg);
            nextTick(scrollToBottom);
            if (msg.RecipientId === user.value.Id) notifyMessageRead(msg.Id);
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
    const firstVisible = container.value?.firstElementChild;
    messages.value = [...older, ...messages.value];
    await nextTick();
    firstVisible?.scrollIntoView();
}

const handleScroll = throttle(() => {
    if (container.value && container.value.scrollTop === 0) loadOlder();
}, 300);

function scrollToBottom() {
    if (container.value) container.value.scrollTop = container.value.scrollHeight;
}

function send(e) {
    e.preventDefault();
    if (!body.value.trim()) return;
    sendChatMessage(chatUserId(), body.value);
    body.value = '';
}
</script>

<template>
    <div class="chat-container">
        <h2>Chat with <span>{{ recipientUsername }}</span></h2>
        <div class="chat-messages" ref="container" @scroll="handleScroll">
            <div v-for="message in messages" :key="message.Id" class="chat-message">
                <span class="timestamp">{{ DateToLocale(message.Timestamp) }}</span>
                <span class="sender">&lt;{{ message.SenderUsername }}&gt;</span>
                <p>{{ message.Body }}</p>
            </div>
        </div>
        <form id="chat-message" @submit.prevent="send">
            <input v-model="body" name="body" placeholder="Type a message..." required />
            <input type="submit" value="Send" />
        </form>
    </div>
</template>

<style scoped></style>