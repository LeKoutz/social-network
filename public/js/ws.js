import { apiFetch } from './fetchers/api.js';
import { UsersPanel, addUsersPanelButtonListener } from './components/users_panel.js';
import { showChatMessages, cacheUnreadMessage, updateUnreadMessageBadges, playNotificationTone, notifyMessageRead } from './components/chat.js';

let ws = null;

export function connectWS(userId) {
    ws = new WebSocket(`ws://${window.location.host}/ws`);
    // debug prints for testing
    ws.onopen = async () => {
        const data = await apiFetch('/api/chat/unread');
        if (data && data.User.ChatMessages) {
            data.User.ChatMessages.forEach(message => cacheUnreadMessage(message));
            updateUnreadMessageBadges();
        }
    };
    ws.onclose = (e) => console.log('WebSocket closed', e.code, e.reason);
    ws.onerror = (e) => console.log('WebSocket error', e);
    ws.onmessage = async (e) => {
        const envelope = JSON.parse(e.data);
        switch (envelope.type) {
        case 'chat_message': {
            const msg = envelope.payload;
            handleIncomingMessage(userId, msg);
            break;
        }
        case 'user_status': {
            const usersData = await apiFetch('/api/users');
            if (usersData) {
                document.querySelector('.users-panel').innerHTML = UsersPanel(usersData);
                updateUnreadMessageBadges();
            }
            addUsersPanelButtonListener();
            break;
        }
        }
    };
}

export function sendWS(data) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(data);
    }
}

export function disconnectWS() {
    if (ws) {
        ws.close();
        ws = null;
    }
}

function handleIncomingMessage(userId, msg) {
    const chatMessages = document.querySelector('.chat-messages');
    const currentId = parseInt(window.location.hash.split('/').at(-1));
    if (!chatMessages || (msg.SenderId !== currentId && msg.RecipientId !== currentId)) {
        cacheUnreadMessage(msg);
        updateUnreadMessageBadges();
        playNotificationTone();
        return;
    };
    chatMessages.insertAdjacentHTML('beforeend', showChatMessages({User: { ChatMessages: [msg] }}));
    const insertedMessage = chatMessages.lastElementChild;
    if (insertedMessage) insertedMessage.scrollIntoView();
    if (userId === msg.RecipientId) notifyMessageRead(msg);
}