import { apiFetch } from './fetchers/api.js';
import { UsersPanel, addUsersPanelButtonListener } from './components/users_panel.js';
import { showChatMessages } from './components/chat.js';

let ws = null;

export function connectWS() {
    ws = new WebSocket(`ws://${window.location.host}/ws`);
    // debug prints for testing
    ws.onopen = () => console.log('WebSocket connection established');
    ws.onclose = (e) => console.log('WebSocket closed', e.code, e.reason);
    ws.onerror = (e) => console.log('WebSocket error', e);
    ws.onmessage = async (e) => {
        const envelope = JSON.parse(e.data);
        switch (envelope.type) {
        case 'chat_message':
            handleIncomingMessage(envelope.payload);
            break;
        case 'user_status': {
            const usersData = await apiFetch('/api/users');
            if (usersData) {
                document.querySelector('.users-panel').innerHTML = UsersPanel(usersData);
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

function handleIncomingMessage(msg) {
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages)  return; // TODO: Create notification if user is not in chat page or something like that
    const currentId = parseInt(window.location.hash.split('/').at(-1));
    if (msg.SenderId !== currentId && msg.RecipientId !== currentId) return;
    chatMessages.insertAdjacentHTML('beforeend', showChatMessages({User: { ChatMessages: [msg] }}));
}