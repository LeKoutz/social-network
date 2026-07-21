export function CommentDeleteForm(post_id, comment_id) {
    return `
        <form method="POST" class="comment-delete" action="/api/comment/delete">
            <input type="hidden" name="post-id" value="${post_id}"/>
            <input type="hidden" name="comment-id" value="${comment_id}"/>
            <button type="submit">Delete Comment</button>
        </form>`;
}
