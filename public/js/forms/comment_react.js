import { commentReactRequest } from '../fetchers/comments.js';

export function commentReactForm(post_id, comment, logged_in){
    return `
                <form method="POST">
                    <input type="hidden" name="post-id" value="${post_id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit"
                            name="action"
                            value="like"
                            ${!logged_in ? 'disabled' : '' }
                            >${comment.Likes} ${comment.Liked ? '👍' : '👍🏻' }</button>
                </form>
                <form method="POST">
                    <input type="hidden" name="post-id" value="${post_id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit"
                            name="action"
                            value="dislike"
                            ${!logged_in ? 'disabled' : '' }
                            >${comment.Dislikes} ${comment.Disliked ? '👎' : '👎🏻' }</button>
                </form>
            </div>
`;
}

export function attachCommentReactionListener() {
    const forms = document.querySelector('.comment').querySelector('.reactions').querySelector('form');
    console.log("jc");
    console.log(forms);
    if (forms) {
        forms.forEach(form => {
            form.addEventListener('submit', commentReactRequest);
        });
    }
}
