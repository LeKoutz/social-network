export function displayPosts(data) {
    return data.Posts.map(post => `
    <div class="post" id="${post.Id}">
        <a href="/post/view/${post.Id}"><h3>${post.Title}</h3></a>
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
        ${data.User.LoggedIn ? /*TODO: showCommentCreate()*/ '' : ''}
        </div>
    </div>
    `).join('')
}
