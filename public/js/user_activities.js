import {ShowActivityComment} from "./activity_comment.js";
import {ShowActivityPost} from "./activity_post.js";

function parseActivities(activities) {

    return activities !== null ? activities.map(activity=>parseActivity(activity)).join(''):'';
}

function parseActivity(activity) {
    return `
    <div class="activity">
    ${activity.TimestampString}
    <strong><a href="/#/user">${activity.User.Username}}</a></strong> 

    ${ (activity.Type === "post") ? 
            ` posted:
            ${ShowActivityPost(activity)}
            `:

            (activity.Type === "comment") ?
            ` commented:
            ${ShowActivityComment(activity)}
            `:

            (activity.Type === "postLike") ?
            ` liked a post:
            ${ShowActivityPost(activity)}
            `:

            (activity.Type === "postDislike") ?
            `  disliked a post:
            ${ShowActivityPost(activity)}
            `:
            (activity.Type ===  "commentLike") ?
            `  liked a comment:
            ${ShowActivityComment(activity)}
            `:

            (activity.Type === "commentDislike") ?
            `  disliked a comment:
            ${ShowActivityComment(activity)}
            `:''}
    </div>
    `
}

export async function ShowUserActivity() {
    const data = await fetch('/api/').then(response=>response.json());
    console.log(data.User.Activities);
    return `
<div class="container">
    <div class="user-activity">
        <h2>Recent Activities</h2>
        ${parseActivities(data.User.Activities)}
    </div>
</div>`;
}
