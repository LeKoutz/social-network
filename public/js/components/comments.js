import {commentReactForm} from '../forms/comment_react.js';
export function showPostComments(data) {
    const post = data.Posts[0];
    return post.Comments.map(comment => `
        <div class="comment" id="comment-${comment.Id}">
            <span>${comment.Username} (${comment.TimestampString})</span>
            <div class="manage-comment">
                ${ data.User.Id === comment.UserId ? `
                <form method="POST" class="comment-edit" data-type="edit" action="/api/comment/edit">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit">Edit Comment</button>
                </form>
                <form method="POST" class="comment-delete" action="/api/comment/delete">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit">Delete Comment</button>
                </form>
                ` :
        ''
}
            </div>
            ${ data.EditCommentId === comment.Id ? `
            <form method="POST" class="comment-edit" data-type="save" action="/api/comment/edit">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <input type="hidden" name="comment-id" value="${comment.Id}"/>
                <input type="hidden" name="save-comment" value="1"/>
                <textarea name="comment" required>${comment.Body}</textarea>
                <button type="submit">Save Comment</button>
            </form>
            `
        :
        `<pre>${comment.Body}</pre>`
}
        ${commentReactForm(post.Id, comment, data.User.LoggedIn)}
        </div>
    `).join('');
}
