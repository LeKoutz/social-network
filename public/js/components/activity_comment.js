export function ShowActivityComment(activity) {
    return `
<div class="container">
    <h3>${activity.Post.Title}</h3>
    <pre>${activity.Post.Body}</pre>
    <div class="activity-highlight">
        <div class="comment" id="comment-${activity.Comment.Id}">
            <a href="#/post/view/${activity.Post.Id}#comment-${activity.Comment.Id}">
            <pre>${activity.Comment.Body}</pre>
            </a>
            <div class="reactions">
                <button disabled>${activity.Comment.Likes} ${activity.Comment.Liked?"👍":"👍🏻"}</button>
                <button disabled>${activity.Comment.Dislikes} ${activity.Comment.Disliked?"👎":"👎🏻"}</button>
            </div>
        </div>
    </div>
</div>
`;
}
