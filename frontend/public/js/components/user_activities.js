import { apiFetch } from '../fetchers/api.js';
import { displayPosts } from "./posts.js";
import { attachPostReactionListener } from '../forms/post_react.js';
import { DateToLocale } from '../utils/utils.js';

export function ShowUserLikes(data) {
    return `
<div class="container">
    <div class="user-activity">
        <h2>My Likes</h2>
        ${displayPosts(data)}
    </div>
</div>`;
}

export function ShowUserPosts(data) {
    return `
<div class="container">
    <div class="user-activity">
        <h2>My Posts</h2>
        ${displayPosts(data)}
    </div>
</div>`;
}

export async function userPostsRoute() {
    const data = await apiFetch('/api/user/posts');
    if (data) {
        document.querySelector('.content').innerHTML = ShowUserPosts(data);
        attachPostReactionListener();
    }
}

export async function userLikesRoute() {
    const data = await apiFetch('/api/user/likes');
    if (data) {
        document.querySelector('.content').innerHTML = ShowUserLikes(data);
        attachPostReactionListener();
    }
}
