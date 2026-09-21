import { reactive, ref } from 'vue';

let ws = null;
let lastUserId = null;
let reconnectTimer = null;
const listeners = [];
const unread = reactive(new Map());
const activeChat = ref(null);
const sound = new Audio('/sounds/message_notification.mp3');

function emit(type, payload) {
    listeners.filter((l) => l.type === type).forEach((l) => l.fn(payload));
}

function cacheUnreadMessage(msg) {
    if (!unread.has(msg.SenderId)) unread.set(msg.SenderId, new Set());
    unread.get(msg.SenderId).add(msg.Id);
    emit('unread_change');
}

export function useChat() {
    function ensureWs(userId) {
        if (!userId) return;
        lastUserId = userId;
        if (ws && (ws.readyState === WebSocket.CONNECTING || ws.readyState === WebSocket.OPEN)) return;
        if (ws && ws.readyState === WebSocket.CLOSING) {
            ws.onclose = null;
            ws.close();
        }
        const socket = new WebSocket(`ws://${window.location.host}/ws`);
        ws = socket;
        socket.onopen = async () => {
            const data = await fetch('/api/chat/unread').then((r) => r.json());
            if (data && !data.Error?.Has && data.User?.ChatMessages) {
                data.User.ChatMessages.forEach(cacheUnreadMessage);
            }
        };
        socket.onclose = () => {
            if (ws === socket) ws = null;
            if (lastUserId && !reconnectTimer) {
                reconnectTimer = setTimeout(() => {
                    reconnectTimer = null;
                    ensureWs(lastUserId);
                }, 3000);
            }
        };
        socket.onerror = () => {};
        socket.onmessage = (e) => {
            const envelope = JSON.parse(e.data);
            switch (envelope.type) {
            case 'chat_message': {
                const msg = envelope.payload;
                const isForCurrent =
                    activeChat.value !== null &&
                    (msg.SenderId === activeChat.value || msg.RecipientId === activeChat.value);
                if (!isForCurrent) {
                    cacheUnreadMessage(msg);
                    playNotificationTone();
                }
                emit('chat_message', { msg, isForCurrent });
                break;
            }
            case 'user_status':
                emit('user_status');
                break;
            }
        };
    }

    function disconnectWS() {
        lastUserId = null;
        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null;
        }
        if (ws) {
            ws.onclose = null;
            ws.close();
            ws = null;
        }
    }

    function sendWS(data) {
        if (ws && ws.readyState === WebSocket.OPEN) ws.send(data);
    }

    function sendChatMessage(recipientId, body) {
        sendWS(JSON.stringify({ type: 'chat-message', payload: { recipientId, body } }));
    }

    function notifyMessageRead(id) {
        sendWS(JSON.stringify({ type: 'message-read', payload: { Id: id } }));
    }

    function markMessageAsRead(msg) {
        const ids = unread.get(msg.SenderId);
        if (!ids) return;
        ids.delete(msg.Id);
        if (ids.size === 0) unread.delete(msg.SenderId);
        emit('unread_change');
    }

    function hasUnread(senderId) {
        return (unread.get(senderId)?.size ?? 0) > 0;
    }

    function totalUnread() {
        let total = 0;
        for (const ids of unread.values()) total += ids.size;
        return total;
    }

    function playNotificationTone() {
        sound.volume = 0.5;
        sound.play().catch(() => {});
    }

    function onChatEvent(type, fn) {
        listeners.push({ type, fn });
        return () => {
            const index = listeners.findIndex((l) => l.fn === fn && l.type === type);
            if (index !== -1) listeners.splice(index, 1);
        };
    }

    function setActiveChat(id) {
        activeChat.value = id;
    }

    return {
        activeChat,
        ensureWs,
        disconnectWS,
        sendChatMessage,
        notifyMessageRead,
        markMessageAsRead,
        hasUnread,
        totalUnread,
        playNotificationTone,
        onChatEvent,
        setActiveChat,
    };
}