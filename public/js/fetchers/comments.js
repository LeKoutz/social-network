import { ShowError } from '../components/error.js';
import { displayPost, showPostView } from '../components/posts.js';
import { attachCommentEditListener, attachCommentDeleteListener } from '../forms/comment_create.js';
import {commentReactForm, attachCommentReactionListener} from '../forms/comment_react.js';
import {SetAlertsInner} from '../partials/alerts.js';

export async function commentCreateRequest(e) {
    e.preventDefault();
    const postId = e.target.querySelector('[name="post-id"]').value;
    const content = document.querySelector('.content');
    const response = await fetch('/api/comment/create', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target))
    });
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        content.innerHTML = ShowError(data);
        return;
    }
    window.location.hash = `/post/view/${postId}#comment-${data.commentId}`;
}

export async function commentEditRequest(e) {
    e.preventDefault();
    const content = document.querySelector('.content');
    const response = await fetch('/api/comment/edit', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target))
    });
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        content.innerHTML = ShowError(data);
        return;
    }
    if (e.target.dataset.type === 'save') {
        content.innerHTML = await showPostView(data.Posts[0].Id);
    } else {
        content.innerHTML = displayPost(data);
    }
    attachCommentEditListener();
    attachCommentDeleteListener();
}

export async function commentReactRequest(e) {
    e.preventDefault();
    const commentId = e.target.querySelector('[name="comment-id"]').value
    const content = document.querySelector(`[id="comment-${commentId}"] .reactions`);
    const response = await fetch('/api/comment/react', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target))
    });
    const data = await response.json();
    console.log(data)
    if (data.Error && data.Error.Has) {
        SetAlertsInner(ShowError(data));
        return;
    }
    content.innerHTML = commentReactForm(data.Posts[0].Id, data.Posts[0].Comments[0], data.User.LoggedIn);
    attachCommentReactionListener();
}
