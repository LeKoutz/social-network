import { ShowError } from '../components/error.js';

export async function ShowUserMenu() {
    const response = await fetch('/api/user');
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        document.querySelector('.alerts').innerHTML = ShowError(data);
        return '';
    }
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
