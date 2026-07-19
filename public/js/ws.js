import { apiFetch } from './fetchers/api.js';
import { UsersPanel, addUsersPanelButtonListener } from './components/users_panel.js';
import { showChatMessages, cacheUnreadMessage, updateUnreadMessageBadges, playNotificationTone, notifyMessageRead } from './components/chat.js';

let ws = null;
// Function to refresh the users panel by fetching the latest user data and updating the panel's HTML
async function refreshUsersPanel() {
    const panel = document.querySelector('.users-panel');

    if (!panel) return;

    // Θυμόμαστε αν το panel ήταν κλειστό.
    const wasCollapsed =
        panel.querySelector('.users-panel-inner')?.hidden ?? false;

    const usersData = await apiFetch('/api/users');

    if (!usersData) return;

    // Αντικαθιστούμε μόνο το users panel, όχι ολόκληρη τη σελίδα.
    panel.innerHTML = UsersPanel(usersData);

    const panelInner = panel.querySelector('.users-panel-inner');

    if (panelInner) {
        panelInner.hidden = wasCollapsed;
    }

    // Επαναφέρουμε badges και click listener επειδή δημιουργήθηκε νέο HTML.
    updateUnreadMessageBadges();
    addUsersPanelButtonListener();
}

export function connectWS(userId) {
    if (ws) disconnectWS();
    ws = new WebSocket(`ws://${window.location.host}/ws`);
    // debug prints for testing
    ws.onopen = async () => {
        console.log('WS connected');
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
                
                // Refresh the users panel to reflect the latest message timestamps and user statuses
                await refreshUsersPanel();
                break;
            }
            case 'user_status': {
                await refreshUsersPanel();
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