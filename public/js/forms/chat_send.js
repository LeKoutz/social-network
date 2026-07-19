import { apiFetch } from '../fetchers/api.js';
import { sendWS } from '../ws.js';
import { markMessageAsRead, showChatMessages } from '../components/chat.js';
import { throttle } from '../utils/utils.js';

function showChat(data, recipientUsername) {
    return `<div class="chat-container">
        <h2>Chat with <span>${recipientUsername}</span></h2>
        <div class="chat-messages">
            ${data.User.ChatMessages ? showChatMessages(data) : ''}
        </div>
        <form id="chat-message">
            <input type="text" name="body" placeholder="Type a message..." required />
            <input type="submit" value="Send"/>
        </form>
    </div>`;
}

function getChatUsername(id) {
    const userLink = document.querySelector('.users-panel [data-user-id="' + id + '"] a');
    return userLink ? userLink.textContent.trim() : 'user';
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
        document.querySelector('.content').innerHTML = showChat(data, getChatUsername(id));
        const chatMessages = document.querySelector('.chat-messages');
        chatMessages.scrollTop = chatMessages.scrollHeight;
        if (data.User.ChatMessages) data.User.ChatMessages.forEach(message => markMessageAsRead(message));
        attachChatListener(id);
        attachScrollListener(id);
    }
}

function attachScrollListener(id) {
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages) return;
    const scrollHandler = throttle(async () => {
        if (chatMessages.scrollTop > 0) return;
        const lastViewedMessage = chatMessages.firstElementChild;
        const currentOffset = chatMessages.children.length;
        const olderMessages = await fetchOlderMessages(id, currentOffset);
        if (olderMessages) {
            olderMessages.User.ChatMessages.forEach(message => {
                markMessageAsRead(message);
            });
            chatMessages.insertAdjacentHTML('afterbegin', showChatMessages(olderMessages));
            lastViewedMessage.scrollIntoView();
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