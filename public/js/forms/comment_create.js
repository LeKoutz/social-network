import { showPostView } from '../posts.js';
import { ShowError } from '../error.js';
import { displayPost } from '../posts.js';


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
        const container = document.querySelector('.container');
        if (form) {
                form.addEventListener('submit', async (e) => {
                        e.preventDefault();
                        const response = await fetch('/api/comment/create', {
                                method: 'POST',
                                body: new URLSearchParams(new FormData(e.target))
                        });
                        const data = await response.json();
                        if (data.Error && data.Error.Has) {
                            container.innerHTML = ShowError(data);
                            return;
                        }
                        container.innerHTML = await showPostView(data.postId);
                        document.querySelector(`#comment-${data.commentId}`).scrollIntoView();
                });
        }
}

export function attachCommentEditListener() {
        const forms = document.querySelectorAll('.comment-edit');
        const container = document.querySelector('.container');
        if (forms) {
            forms.forEach(form => {
                form.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    const response = await fetch('/api/comment/edit', {
                        method: 'POST',
                        body: new URLSearchParams(new FormData(e.target))
                    });
                    const data = await response.json();
                    if (data.Error && data.Error.Has) {
                        container.innerHTML = ShowError(data);
                        return;
                    }
                    if (form.dataset.type === 'save') {
                        container.innerHTML = await showPostView(data.Posts[0].Id);
                    } else {
                        container.innerHTML = displayPost(data);
                    }
                    attachCommentEditListener();
                });
            });
        }
}