import { commentReactRequest } from '../fetchers/comments.js';

export function commentReactForm(post_id, comment, logged_in){
    return `
                <form class="comment-reaction" method="POST">
                    <input type="hidden" name="post-id" value="${post_id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <input type="hidden" name="action" value="like"/>
                    <button type="submit"
                            ${!logged_in ? 'disabled' : '' }
                            >${comment.Likes} ${comment.Liked ? '👍' : '👍🏻' }</button>
                </form>
                <form class="comment-reaction" method="POST">
                    <input type="hidden" name="post-id" value="${post_id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <input type="hidden" name="action" value="dislike"/>
                    <button type="submit"
                            ${!logged_in ? 'disabled' : '' }
                            >${comment.Dislikes} ${comment.Disliked ? '👎' : '👎🏻' }</button>
                </form>
            </div>
`;
}

export function attachCommentReactionListener() {
    const forms = document.querySelectorAll('.comment-reaction');
    if (forms) {
        forms.forEach(form => {
            form.addEventListener('submit', commentReactRequest);
        });
    }
}
