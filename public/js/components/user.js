import { apiFetch } from '../fetchers/api.js';

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

export async function userMenuRoute() {
    const data = await apiFetch('/api/user');
    if (data && data.User.LoggedIn) {
        document.querySelector('.content').innerHTML = ShowUserMenu();
    }
}