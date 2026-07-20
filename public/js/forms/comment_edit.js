export function CommentEditForm(post_id, comment_id) {
    return `
        <form method="POST" class="comment-edit" data-type="edit" action="/api/comment/edit">
            <input type="hidden" name="post-id" value="${post_id}"/>
            <input type="hidden" name="comment-id" value="${comment_id}"/>
            <button type="submit">Edit Comment</button>
        </form>`;
}
