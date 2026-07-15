export function UsersPanel(data) {
    // The spread operator creates a new array containing the same users.
    // This is important because sort() changes the array on which it runs;
    // copying data.Users prevents us from changing the original API data. 
    // localeCompare() compares the Username strings alphabetically.
    const sortedUsers = [...data.Users].sort((a, b) =>
        a.Username.localeCompare(b.Username)
    );
    // Keep the logged-in users from the alphabetically sorted array.
    const onlineUsers = sortedUsers.filter(user => user.LoggedIn);
    // Keep the logged-out users from the same alphabetically sorted array.
    // The ! operator changes false to true, selecting offline users only.
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
