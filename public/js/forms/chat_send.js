import { apiFetch } from '../fetchers/api.js';
import { sendWS } from '../ws.js';
import { markMessageAsRead, showChatMessages } from '../components/chat.js';
import { throttle } from '../utils/utils.js';

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
            sendWS(JSON.stringify({ type: "chat-message", payload: { recipientId: id, body: body }}));
            form.reset();
        });
    }
}

export async function chatRoute(id, offset = 0) {
    const data = await apiFetch(`/api/chat/${id}?offset=${offset}`);
    if (data) {
        document.querySelector('.content').innerHTML = showChat(data);
        const chatMessages = document.querySelector('.chat-messages');
        chatMessages.scrollTop = chatMessages.scrollHeight;
        data.User.ChatMessages.forEach(message => {
            markMessageAsRead(message);
        });
        attachChatListener(id);
        attachScrollListener(id, offset);
    }
}

function attachScrollListener(id, offset) {
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages) return;
    let currentOffset = offset;
    const scrollHandler = throttle(async () => {
        if (chatMessages.scrollTop > 0) return;
        currentOffset += 10;
        const olderMessages = await fetchOlderMessages(id, currentOffset);
        if (olderMessages) {
            olderMessages.User.ChatMessages.forEach(message => {
                markMessageAsRead(message);
            });
            chatMessages.insertAdjacentHTML('afterbegin', showChatMessages(olderMessages));
        } else {
            chatMessages.removeEventListener('scroll', scrollHandler);
        }
    }, 300);
    chatMessages.addEventListener('scroll', scrollHandler);
}

async function fetchOlderMessages(id, offset) {
    const data = await apiFetch(`/api/chat/${id}?offset=${offset}`);
    if (data && data.User.ChatMessages && data.User.ChatMessages.length > 0) {
        return data;
    }
    return null;
}