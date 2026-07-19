import { ShowError } from '../components/error.js';
import { displayPost, showPostView } from '../components/posts.js';
import { attachCommentEditListener } from '../forms/comment_create.js';

export async function commentCreateRequest(e) {
    e.preventDefault();
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
    content.innerHTML = await showPostView(data.postId);
    document.querySelector(`#comment-${data.commentId}`).scrollIntoView();
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
}

export async function commentReactRequest(e) {
    e.preventDefault();
    const content = document.querySelector('.content');
    const response = await fetch('/api/comment/react', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target))
    });
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        content.innerHTML = ShowError(data);
        return;
    }
    if (form.dataset.type === 'save') {
        content.innerHTML = await showPostView(data.Posts[0].Id);
    } else {
        content.innerHTML = displayPost(data);
    }
    attachCommentReactListener();
}
