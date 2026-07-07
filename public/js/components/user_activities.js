import {ShowActivityComment} from "./activity_comment.js";
import {ShowActivityPost} from "./activity_post.js";

function parseActivities(activities) {
    return activities !== null ? activities.map(activity=>parseActivity(activity)).join(''):'';
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

export async function ShowUserActivity() {
    const data = await fetch('/api/user/activity').then(response=>response.json());
    console.log(data);
    return `
<div class="container">
    <div class="user-activity">
        <h2>Recent Activities</h2>
        ${parseActivities(data.User.Activities, data)}
    </div>
</div>`;
}

export async function ShowUserLikes() {
    const data = await fetch('/api/user/likes').then(r=>r.json());
    return `
<div class="container">
    <div class="user-activity">
        <h2>My Likes</h2>
        ${parseActivities(data.User.Activities, data)}
    </div>
</div>`;
}

export async function ShowUserPosts() {
    const data = await fetch('/api/user/posts').then(r=>r.json());
    return `
<div class="container">
    <div class="user-activity">
        <h2>My Posts</h2>
        ${parseActivities(data.User.Activities, data)}
    </div>
</div>`;
}
