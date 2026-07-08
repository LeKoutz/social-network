import {commentCreateRequest, commentEditRequest} from '../fetchers/comments.js';

export function showCommentCreate(post) {
    return `
<form method="POST" id=create-comment action="/api/comment/create">
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
