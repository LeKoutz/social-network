import { apiFetch } from '../fetchers/api.js';
import { sendWS } from '../ws.js';
import { showChatMessages } from '../components/chat.js';

function showChat(data) {
    return `<div class="chat-container">
        <h2>Chat</h2>
        <div class="chat-messages">
            ${data.User.ChatMessages ? showChatMessages(data) : ''}
        </div>
        <form id="chat-message">
            <input type="text" name="body" placeholder="Type a message..." required />
            <input type="submit" value="Send"/>
        </form>
    </div>`;
}

function attachChatListener(id) {
    const form = document.querySelector('#chat-message');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const body = form.querySelector('[name="body"]').value;
            sendWS(JSON.stringify({ recipientId: id, body }));
            form.reset();
        });
    }
}

export async function chatRoute(id) {
    const data = await apiFetch(`/api/chat/${id}`);
    if (data) {
        document.querySelector('.content').innerHTML = showChat(data);
        attachChatListener(id);
    }
}