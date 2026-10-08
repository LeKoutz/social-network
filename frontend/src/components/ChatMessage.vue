<script setup>
import { DateToLocale } from '@/utils/utils.js';

defineProps({
    message: { type: Object, required: true },
    isOwn: { type: Boolean, default: false },
});
</script>

<template>
    <div class="chat-message" :class="{ own: isOwn }">
        <img
            :src="message.AvatarURL ? `/${message.AvatarURL}` : '/default-avatar.png'"
            alt="Avatar"
            class="chat-avatar"
        />
        <div class="chat-content">
            <div class="chat-header">
                <span class="sender">
                    <router-link :to="`/profile/view/${message.SenderId}`">{{ message.SenderUsername }}</router-link>
                </span>
                <span class="timestamp">{{ DateToLocale(message.Timestamp) }}</span>
            </div>
            <span class="body">{{ message.Body }}</span>
        </div>
    </div>
</template>

<style scoped>
.chat-message.own {
    display: flex;
    flex-direction: row-reverse;
    margin-left: auto;
}

.chat-message.own .chat-content {
    align-items: flex-end;
}

.chat-message {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    padding: 0.5rem 0;
    width: 80%;
}

.chat-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
}

.chat-content {
    display: flex;
    flex-direction: column;
    min-width: 0;
}

.chat-header {
    display: flex;
    align-items: baseline;
    gap: 0.5rem;
}

.timestamp {
    color: var(--text-muted);
    font-size: 0.8em;
}

.body {
    overflow-wrap: anywhere;
}
</style>