export function UsersPanel(data) {
    const onlineUsers = data.Users.filter(user => user.LoggedIn);
    const offlineUsers = data.Users.filter(user => !user.LoggedIn);
    return `
    <div class="panel-header">
        <h3>Users</h3>
        <button id="collapse-panel">−</button>
    </div>
    <h4>Online</h4>
    <ul>
        ${onlineUsers.map(user => `${data.User.Id !== user.Id ? `<li><a href="#/chat/${user.Id}">${user.Username}</a></li>` : ''}`)
        .join('')}
    </ul>
    <h4>Offline</h4>
    <ul>
        ${offlineUsers.map(user => `${data.User.Id !== user.Id ? `<li class="offline"><a href="#/chat/${user.Id}">${user.Username}</a></li>` : ''}`)
        .join('')}
    </ul>
    `;
}