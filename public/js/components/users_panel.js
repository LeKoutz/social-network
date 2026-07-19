export function UsersPanel(data) {

    //create a new array of users sorted by username
    const sortedUsers = [...data.Users].sort((a, b) =>
        a.Username.localeCompare(b.Username)
    );
    //filter the users into online and offline arrays
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
            ${onlineUsers.map(user => `${data.User.Id !== user.Id ? `<li data-user-id="${user.Id}"><a href="#/chat/${user.Id}">${user.Username}</a></li>` : ''}`)
        .join('')}
        </ul>
        <h4>Offline</h4>
        <ul>
            ${offlineUsers.map(user => `${data.User.Id !== user.Id ? `<li class="offline" data-user-id="${user.Id}"><a href="#/chat/${user.Id}">${user.Username}</a></li>` : ''}`)
        .join('')}
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
