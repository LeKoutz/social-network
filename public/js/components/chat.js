export function showChatMessages(data) {
    return data.User.ChatMessages.map(message => `
        <div class="chat-message">
            <span class="sender">${message.SenderUsername}</span>
            <p>${message.Body}</p>
            <span class="timestamp">${message.TimestampString}</span>
        </div>
    `).join('');
}