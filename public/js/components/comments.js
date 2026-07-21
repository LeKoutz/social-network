import { commentReactForm } from '../forms/comment_react.js';
import { CommentEditForm } from '../forms/comment_edit.js';
import { CommentDeleteForm } from '../forms/comment_delete.js';

export function showPostComments(data) {
    const post = data.Posts[0];
    return post.Comments.map(comment => `
        <div class="comment" id="comment-${comment.Id}">
            <span>${comment.Username} (${comment.TimestampString})</span>
            <div class="manage-comment">
                ${ data.User.Id === comment.UserId ? `
                ${ CommentEditForm(post.Id, comment.Id) }
                ${ CommentDeleteForm(post.Id, comment.Id) }` : '' }
            </div>
            ${ data.EditCommentId === comment.Id ? `
            <form method="POST" class="comment-edit" data-type="save" action="/api/comment/edit">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <input type="hidden" name="comment-id" value="${comment.Id}"/>
                <input type="hidden" name="save-comment" value="1"/>
                <textarea name="comment" required>${comment.Body}</textarea>
                <button type="submit">Save Comment</button>
            </form>` : `<pre>${comment.Body}</pre>` }
            <div class="reactions">
                ${commentReactForm(post.Id, comment, data.User.LoggedIn)}
            </div>
        </div>
    `).join('');
}
