export function ShowUserMenu() {
    return `
<div class="container">
    <div class="user-section">
        <p>You can do the following</p>
        <ul>
        <li><a href="#/post/create">Create post</a></li>
        <li><a href="#/user/activity">My activity</a></li>
        <li><a href="#/user/posts">My posts</a></li>
        <li><a href="#/user/likes">My likes</a></li>
        <li><a href="#/user/logout">Log out</a></li>
        </ul>
    </div>
</div>`;
}
