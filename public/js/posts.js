import { showPostComments, showCommentCreate, attachCommentCreateListener } from "./comments.js"
import { showPostCategories } from "./categories.js"

export function displayPosts(data) {
    return data.Posts.map(post => `
    <div class="post" id="${post.Id}">
        <a href="#/post/view/${post.Id}"><h3>${post.Title}</h3></a>
        <p>Posted on <em>${post.TimestampString}</em>.</p>
        <pre>${post.Body}</pre>
        ${post.ImagePath ? `<img src="/${post.ImagePath}" alt="Post image" style="max-width: max-content;"/>` : ''}
        <div class="reactions">
            <form method="POST" action="/post/react">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <button type="submit"
                        name="action"
                        value="like"
                        ${data.User.LoggedIn ? '' : 'disabled'}
                        >${post.Likes} ${post.Liked ? '👍': '👍🏻'}</button>
            </form>
            <form method="POST" action="/post/react">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <button type="submit"
                        name="action"
                        value="dislike"
                        ${data.User.LoggedIn ? '' : 'disabled'}
                        >${post.Dislikes} ${post.Disliked ? '👎' : '👎🏻' }</button>
            </form>
        </div>
        <div class="comments">
        ${data.User.LoggedIn ? showCommentCreate(post) : ''}
        </div>
    </div>
    `).join('')
}

export function displayPost(data) {
    const post = data.Posts[0]
    return `
    <div class="container">
    <div class="post" id="${post.Id}">
        <div class="post-content">
            <a href="#/post/view/${post.Id}"><h3>${post.Title}</h3></a>
            <div class="post-details">
                <p>Categories:
                ${showPostCategories(post)}
                </p>
                <p>Posted by <strong>${post.User.Username}</strong> on <em>(${post.TimestampString})</em></p>
            </div>
            <div class=manage-post>
                ${data.User.Id === post.User.Id ? `
                <form method="GET" action="/post/edit">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <button type="submit">Edit Post</button>
                </form>
                <form method="POST" action="/post/delete">
                    <input type="hidden" name="post-id" value="${post.Id}"/>
                    <button type="submit">Delete Post</button>
                </form>
                ` 
                : 
                ''
            }
            </div>
            <pre>${post.Body}</pre>
            ${post.ImagePath ? `<img src="/${post.ImagePath}" alt="Post image" style="max-width: 100%;"/>` : '' }
            <div class="manage-post">
        </div>
        <div class="reactions">
            <form method="POST" action="/post/react">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <button type="submit"
                        name="action"
                        value="like"
                        ${!data.User.LoggedIn ? 'disabled' : '' }
                        >${post.Likes} ${post.Liked ? '👍' : '👍🏻' }</button>
            </form>
            <form method="POST" action="/post/react">
                <input type="hidden" name="post-id" value="{{$post.Id}}"/>
                <button type="submit"
                        name="action"
                        value="dislike"
                        ${!data.User.LoggedIn ? 'disabled' : '' }
                        >${post.Dislikes} ${post.Disliked ? '👎' : '👎🏻' }</button>
            </form>
        </div>
        <div class="comments">
            ${data.User.LoggedIn ? showCommentCreate(post) : ''}
            ${post.Comments ? showPostComments(data) : ''}
        </div>
    </div>
</div>
</div>
    `
}

export async function showPostView(id) {
    const response = await fetch(`/api/post/view/${id}`)
    const data = await response.json()
    return `${displayPost(data)}`
}
