import { showPostView } from "./posts.js"
import { ShowError } from './error.js';

export function showPostComments(data) {
    const post = data.Posts[0]
    return post.Comments.map(comment => `
        <div class="comment" id="comment-${comment.Id}">
            <span>${comment.Username} (${comment.TimestampString})</span>
            <div class="manage-comment">
                ${ data.User.Id === comment.UserId ? `
                <form method="POST" action="/comment/edit">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit">Edit Comment</button>
                </form>
                <form method="POST" action="/comment/delete">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit">Delete Comment</button>
                </form>
                ` :
                ''
                }
            </div>
            ${ data.EditCommentId === comment.Id ? `
            <form method="POST" action="/comment/edit">
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
            <div class="reactions">
                <form method="POST" action="/comment/react">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit"
                            name="action"
                            value="like"
                            ${!data.User.LoggedIn ? 'disabled' : '' }
                            >${comment.Likes} ${comment.Liked ? '👍' : '👍🏻' }</button>
                </form>
                <form method="POST" action="/comment/react">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <input type="hidden" name="comment-id" value="${comment.Id}"/>
                    <button type="submit"
                            name="action"
                            value="dislike"
                            ${!data.User.LoggedIn ? 'disabled' : '' }
                            >${comment.Dislikes} ${comment.Disliked ? '👎' : '👎🏻' }</button>
                </form>
            </div>
        </div>
    `).join('')
}

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
`
}

export function attachCommentCreateListener() {
        const form = document.querySelector('#create-comment')
        const container = document.querySelector('.container');
        if (form) {
                form.addEventListener('submit', async (e) => {
                        e.preventDefault()
                        const response = await fetch('/api/comment/create', {
                                method: 'POST',
                                body: new URLSearchParams(new FormData(e.target))
                        })
                        const data = await response.json()
                        console.log(data)
                        if (data.Error && data.Error.Has) {
                            container.innerHTML = ShowError(data)
                            return
                        }
                        container.innerHTML = await showPostView(data.postId)
                        document.querySelector(`#comment-${data.commentId}`).scrollIntoView()
                })
        }
}