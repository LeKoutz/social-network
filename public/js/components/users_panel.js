export function UsersPanel(data) {
    return `
    <div class="panel-header">
        <h3>Users</h3>
        <button id="collapse-panel">−</button>
    </div>
    <ul>
        ${data.Users.map(user => `${data.User.Id !== user.Id ? `<li><a href="#/chat/${user.Id}">${user.Username}</a></li>` : ''}`)
        .join('')}
    </ul>
    `;
}