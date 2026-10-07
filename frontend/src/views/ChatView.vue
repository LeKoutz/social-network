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
const canChat = ref(null);

function chatUserId() {
    return Number(route.params.id);
}

async function loadHistory(offset) {
    const data = await apiFetch(`/api/chat/${chatUserId()}?offset=${offset}`);
    if (!data) return [];
    canChat.value = data.Users?.[0]?.CanChat ?? false;
    recipientUsername.value = data.Users?.[0]?.Username ?? '';
    const msgs = data.User?.ChatMessages ?? [];
    msgs.forEach((m) => {
        if (m.RecipientId === user.value.Id) markMessageAsRead(m);
    });
    return msgs;
}

async function initChat() {
    setActiveChat(chatUserId());
    historyEnded.value = false;
    messages.value = await loadHistory(0);
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
    if (canChat.value) {
        if (!body.value.trim()) return;
        sendChatMessage(chatUserId(), body.value);
        body.value = '';
    }
}
</script>

<template>
    <div class="container padding-sides-1vw">
        <h2>Chat with <span><router-link :to="`/profile/view/${chatUserId()}`">{{ recipientUsername }}</router-link></span></h2>
        <div class="chat-messages" ref="container" @scroll="handleScroll">
            <div v-for="message in messages" :key="message.Id" class="chat-message">
                <span class="timestamp">{{ DateToLocale(message.Timestamp) }}</span>
                <span class="sender"><router-link :to="`/profile/view/${message.SenderId}`">{{ message.SenderUsername }}</router-link></span>
                <span>{{ message.Body }}</span>
            </div>
        </div>
        <form id="chat-message" @submit.prevent="send">
            <input v-model="body" name="body" placeholder="Type a message..." required :hidden="!canChat"/>
            <input type="submit" value="Send" :hidden="!canChat"/>
            <p :hidden="canChat !== false">You can't chat with this user. Try following them first</p>
        </form>
    </div>
</template>

<style scoped></style>
