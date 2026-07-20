import { showPostComments } from "./comments.js";
import { showCommentCreate, attachCommentCreateListener, attachCommentEditListener, attachCommentDeleteListener } from "../forms/comment_create.js";
import { SetAlertsInner } from '../partials/alerts.js';
import { showPostCategories } from "./categories.js";
import { ShowError } from "./error.js";
import { apiFetch } from '../fetchers/api.js';
import { postDeleteForm, attachPostDeleteListener } from '../forms/post_delete.js';
import {attachPostReactionListener, postReactionForm} from "../forms/post_react.js";
import {attachCommentReactionListener} from '../forms/comment_react.js'

export function displayPosts(data) {
    return data.Posts ? data.Posts.map(post => `
    <div class="post" id="${post.Id}">
        <a href="#/post/view/${post.Id}"><h3>${post.Title}</h3></a>
        <p>Posted on <em>${post.TimestampString}</em>.</p>
        <pre>${post.Body}</pre>
        ${post.ImagePath ? `<img src="/${post.ImagePath}" alt="Post image" style="max-width: max-content;"/>` : ''}
        <div class="reactions">
            ${postReactionForm(post,data.User.LoggedIn)}
        </div>
        <div class="comments">
            ${data.User.LoggedIn ? showCommentCreate(post) : ''}
        </div>
    </div>
    `).join('') : '';
}

export function displayPost(data) {
    console.log(data);
    const post = data.Posts[0];
    return `
    <div class="container">
    <div class="post" id="${post.Id}">
        <div class="post-content">
            <a href="#/post/view/${post.Id}"><h3>${post.Title}</h3></a>
            <div class="post-details">
                <p>Categories:
                ${showPostCategories(post)}
                </p>
                <p>Posted by <strong>${post.User.Username}</strong> on <em>(${post.TimestampString})</em></p>
            </div>
            <div class="manage-post">
                ${data.User.Id === post.User.Id ? `
                <a href="#/post/edit/${post.Id}"><button type="button">Edit Post</button></a>
                ${postDeleteForm(post.Id)}
                `
        :
        ''
}
            </div>
            <pre>${post.Body}</pre>
            ${post.ImagePath ? `<img src="/${post.ImagePath}" alt="Post image" style="max-width: 100%;"/>` : '' }
            <div class="manage-post">
        </div>
        <div class="reactions">
        ${postReactionForm(post,data.User.LoggedIn)}
        </div>
        <div class="comments">
            ${data.User.LoggedIn ? showCommentCreate(post) : ''}
            ${post.Comments ? showPostComments(data) : ''}
        </div>
    </div>
</div>
</div>
    `;
}

// This should probably be replaced by postRoute and removed completely
export async function showPostView(id) {
    const response = await fetch(`/api/post/view/${id}`);
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        SetAlertsInner(ShowError(data));
        return '';
    }
    return `${displayPost(data)}`;
}

export async function postRoute(id, arg) {
    const data = await apiFetch(`/api/post/view/${id}`);
    if (data) {
        document.querySelector('.content').innerHTML = displayPost(data);
        attachPostReactionListener();
        attachCommentCreateListener();
        attachCommentEditListener();
        attachCommentDeleteListener();
        attachCommentReactionListener();
        attachPostDeleteListener();
        if (arg !== undefined) {
            const args = arg.split('-');
            if ( args.length === 2 ) {
                if ( args[0] === 'comment' ) {
                    const r = /^\d*$/;
                    if ( r.test(args[1]) ) {
                        console.log(arg);
                        document.querySelector(`#${arg}`).scrollIntoView();
                    }
                }
            }
        }
    }
}
