export function postReactionForm(post, logged_in) {
    return `
        <div class="reactions">
            <form method="POST" action="/api/post/react">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <button type="submit"
                        name="action"
                        value="like"
                        ${logged_in ? '' : 'disabled'}
                        >${post.Likes} ${post.Liked ? '👍': '👍🏻'}</button>
            </form>
            <form method="POST" action="/api/post/react">
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <button type="submit"
                        name="action"
                        value="dislike"
                        ${logged_in ? '' : 'disabled'}
                        >${post.Dislikes} ${post.Disliked ? '👎' : '👎🏻' }</button>
            </form>
        </div>
        `;
}
