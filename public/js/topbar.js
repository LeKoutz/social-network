export function TopBar(data) {
    return `
<div class="topbar">
    <div class="title">
        <a href="/"><h1>${data.WebsiteName}</h1></a>
    </div>
    <div class="menu">
        <ul>
            ${data.User.LoggedIn ? `
            <li>
                <a href="/user">${data.User.Username}</a>
            </li>
            <li>
                <ul>
                    <li><a href="/post/create">Create post</a></li>
                    <li><a href="/user/activity">My activity</a></li>
                    <li><a href="/user/posts">My posts</a></li>
                    <li><a href="/user/likes">My likes</a></li>
                    <li><details>
                        <summary>Notifications 🔔 ${data.User.UnreadNotificationsCount > 0 ? `(${data.user.unreadNotificationsCount})` : ''}</summary>
                        ${data.User.Notifications.length > 0 ? NotificationsList(data.User.Notifications) : ''}
                    </details></li>
                    <li><a href="/user/logout">Log out</a></li>
                </ul>
            </li>
            ` : `
            <li><a href="/user/login">Log in</a></li>
            <li><a href="/user/register">Register</a></li>
            `}
        </ul>
    </div>
</div>
`
}