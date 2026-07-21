export function ShowActivityPost(activity) {
    return `
<div class="container">
    <div class="activity-highlight">
        <a href="/#/post/view/${activity.Post.Id}"><h3>${activity.Post.Title}</h3></a>
            <pre>${activity.Post.Body}</pre>
        <div class="reactions">
            <button disabled>${activity.Post.Likes} ${activity.Post.Liked?"👍":"👍🏻"}</button>
            <button disabled>${activity.Post.Dislikes} ${activity.Post.Disliked?"👎":"👎🏻"}</button>
        </div>
    </div>
</div>
`;
}
