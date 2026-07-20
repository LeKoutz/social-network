import { postReactRequest } from '../fetchers/posts.js';

export function postReactionForm(post, logged_in) {
    return `
        <form method="POST">
            <input type="hidden" name="post-id" value="${post.Id}"/>
            <input type="hidden" name="action" value="like"/>
            <button type="submit"
                    ${logged_in ? '' : 'disabled'}
                    >${post.Likes} ${post.Liked ? '👍': '👍🏻'}</button>
        </form>
        <form method="POST">
            <input type="hidden" name="post-id" value="${post.Id}"/>
            <input type="hidden" name="action" value="dislike"/>
            <button type="submit"
                    ${logged_in ? '' : 'disabled'}
                    >${post.Dislikes} ${post.Disliked ? '👎' : '👎🏻' }</button>
        </form>
        `;
}

export function attachPostReactionListener() {
    const forms = document.querySelector('.reactions')?.querySelectorAll('form');
    if (forms) {
        forms.forEach(form => {
            form.addEventListener('submit', postReactRequest);
        });
    }
}
