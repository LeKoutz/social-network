import {ShowError} from "../components/error.js";
// import {displayPost, showPostView} from "../components/posts.js";
import {attachPostReactionListener, postReactionForm} from "../forms/post_react.js";

export async function postReactRequest(e) {
    e.preventDefault();
    const content = document.querySelector('.content');
    const response = await fetch('/api/post/react', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target))
    });
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        content.innerHTML = ShowError(data);
        return;
    }
    e.target.parentElement.innerHTML = postReactionForm(data.Posts[0], data.User.LoggedIn);
    attachPostReactionListener();
}

