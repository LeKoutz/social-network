import {ShowActivityComment} from "./activity_comment.js";
import {SetAlertsInner} from '../partials/alerts.js';
import {ShowActivityPost} from "./activity_post.js";
import { apiFetch } from '../fetchers/api.js';
import { displayPosts } from "./posts.js";
import { attachPostReactionListener } from '../forms/post_react.js'

function parseActivities(activities, data) {
    return activities !== null ? activities.map(activity=>parseActivity(activity, data)).join(''):'';
}

function parseActivity(activity, data) {
    let inner = '';
    switch(activity.Type) {
    case "post":
        inner = ` posted: ${ShowActivityPost(activity)}`;
        break;
    case "comment":
        inner = ` commented: ${ShowActivityComment(activity)}`;
        break;

    case "postLike":
        inner = ` liked a post: ${ShowActivityPost(activity)}`;
        break;

    case "postDislike":
        inner = ` disliked a post: ${ShowActivityPost(activity)}`;
        break;

    case  "commentLike":
        inner = ` liked a comment: ${ShowActivityComment(activity)}`;
        break;

    case "commentDislike":
        inner = ` disliked a comment: ${ShowActivityComment(activity)}`;
        break;
    }
    return `
    <div class="activity">
    ${activity.TimestampString}
    <strong><a href="/#/user">${data?.User?.Username}</a></strong> 
    ${inner}
    </div>
    `;
}

export function ShowUserActivity(data) {
    return `
<div class="container">
    <div class="user-activity">
        <h2>Recent Activities</h2>
        ${parseActivities(data.User.Activities, data)}
    </div>
</div>`;
}

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

export async function userActivityRoute() {
    const data = await apiFetch('/api/user/activity');
    if (data) {
        document.querySelector('.content').innerHTML = ShowUserActivity(data);
    }
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
