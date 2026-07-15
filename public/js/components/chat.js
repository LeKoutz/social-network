import { sendWS } from '../ws.js';

const unreadMessages = new Map();
const message_notification = new Audio('/sounds/message_notification.mp3');

export function showChatMessages(data) {
    return data.User.ChatMessages.map(message => `
        <div class="chat-message">
            <span class="sender">${message.SenderUsername}</span>
            <p>${message.Body}</p>
            <span class="timestamp">${message.TimestampString}</span>
        </div>
    `).join('');
}

export function cacheUnreadMessage(msg) {
    if (!unreadMessages.has(msg.SenderId)) {
        unreadMessages.set(msg.SenderId, new Set());
    }
    unreadMessages.get(msg.SenderId).add(msg.Id);
}

function setUnreadMessageBadge() {
    if (unreadMessages.size === 0) return;
    for (const senderId of unreadMessages.keys()) {
        const sender = document.querySelector(`[data-user-id="${senderId}"]`);
        if (sender) sender.insertAdjacentHTML('afterbegin', '<span class="new-message">●</span>');
    }
    const panel = document.querySelector('.users-panel-header');
    if (panel) panel.insertAdjacentHTML('afterbegin', '<span class="new-message">●</span>');
}

export function updateUnreadMessageBadges() {
    document.querySelectorAll('.new-message').forEach(el => el.remove());
    setUnreadMessageBadge();
}

export function markMessageAsRead(message) {
    const unreadIds = unreadMessages.get(message.SenderId);
    if (!unreadIds) return;
    unreadIds.delete(message.Id);
    if (unreadIds.size === 0) unreadMessages.delete(message.SenderId);
    updateUnreadMessageBadges();
}

export function playNotificationTone() {
    message_notification.volume = 0.5;
    message_notification.play();
}

export function notifyMessageRead(msg) {
    sendWS(JSON.stringify({ type: "message-read", payload: { id: msg.id }}));
}