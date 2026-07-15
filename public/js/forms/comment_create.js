import {commentCreateRequest, commentEditRequest} from '../fetchers/comments.js';
import {ShowError} from '../components/error.js';

export function showCommentCreate(post) {
    return `
<form method="POST" id="create-comment" action="/api/comment/create">
    <fieldset>
        <legend>Leave a comment</legend>
        <textarea type="text" name="comment" placeholder="Enter your comment" required></textarea>
        <input type="hidden" name="post-id" value="${post.Id}"/>
        <input type="submit" name="submit" value="Post comment"/>
    </fieldset>
</form>
`;
}

export function attachCommentCreateListener() {
    const form = document.querySelector('#create-comment');
    if (form) {
        form.addEventListener('submit', commentCreateRequest);
    }
}

export function attachCommentEditListener() {
    const forms = document.querySelectorAll('.comment-edit');
    if (forms) {
        forms.forEach(form => {
            form.addEventListener('submit', commentEditRequest);
        });
    }
}

export function attachCommentDeleteListener() {
    const forms = document.querySelectorAll('.comment-delete');
    if (forms) {
        forms.forEach(form => {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const commentId = e.target.querySelector('[name="comment-id"]').value;
                const url = `/api/comment/delete`;
                const response = await fetch(url, {
                    method: 'POST',
                    body: new URLSearchParams(new FormData(e.target))
                });
                const data = await response.json();
                if (data.Error && data.Error.Has) {
                    document.querySelector('.alerts').innerHTML = ShowError(data);
                    return;
                }
                document.querySelector(`#comment-${commentId}`).remove();
            });
        });
    }
}
