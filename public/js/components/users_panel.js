function formatLastMessageTimestamp(timestamp) {
    if (!timestamp) return '';

    const date = new Date(timestamp * 1000);
    const day = date.getDate();
    const month = date.getMonth() + 1;
    const year = String(date.getFullYear()).slice(-2);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${day}-${month}-${year}, ${hours}:${minutes}`;
}

function renderUser(user, currentUserId, offline = false) {
    if (user.Id === currentUserId) return '';

    const lastMessage = formatLastMessageTimestamp(user.LastMessageTimestamp);
    const lastMessageLabel = lastMessage
        ? `<span class="last-message-time">- last msg ${lastMessage}</span>`
        : '';
    const className = offline ? ' class="offline"' : '';

    return `<li${className} data-user-id="${user.Id}">
        <a href="#/chat/${user.Id}">${user.Username}</a>
        ${lastMessageLabel}
    </li>`;
}

export function UsersPanel(data) {
    const sortedUsers = [...(data.Users || [])].sort((a, b) => {
        const aHasMessages = a.LastMessageTimestamp > 0;
        const bHasMessages = b.LastMessageTimestamp > 0;
        if (aHasMessages && bHasMessages) {
            const timestampDifference =
                b.LastMessageTimestamp - a.LastMessageTimestamp;
            if (timestampDifference !== 0) {
                return timestampDifference;
            }
        }
        if (aHasMessages && !bHasMessages) {
            return -1;
        }
        if (!aHasMessages && bHasMessages) {
            return 1;
        }
        return a.Username.localeCompare(b.Username);
    });
    const onlineUsers = sortedUsers.filter(user => user.LoggedIn);
    const offlineUsers = sortedUsers.filter(user => !user.LoggedIn);
    return `
    <div class="users-panel-header">
        <h3>Users</h3>
        <button id="collapse-panel">−</button>
    </div>
    <div class="users-panel-inner">
        <h4>Online</h4>
        <ul>
            ${onlineUsers.map(user => renderUser(user, data.User.Id)).join('')}
        </ul>
        <h4>Offline</h4>
        <ul>
            ${offlineUsers.map(user => renderUser(user, data.User.Id, true)).join('')}
        </ul>
    </div>
    `;
}

export function addUsersPanelButtonListener() {
    const collapse_button = document.querySelector('#collapse-panel');
    const upi = document.querySelector('.users-panel-inner');
    collapse_button.addEventListener("click", ()=>{
        upi.hidden = !upi.hidden;
    });
}
