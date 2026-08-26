import {ShowError} from "../components/error.js";
import {attachPostReactionListener, postReactionForm} from "../forms/post_react.js";
import {SetAlertsInner} from '../partials/alerts.js';

export async function postReactRequest(e) {
    e.preventDefault();
    const postId = e.target.querySelector('[name="post-id"]').value;
    const content = document.querySelector(`[id="${postId}"] .reactions`);
    const response = await fetch('/api/post/react', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target))
    });
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        SetAlertsInner(ShowError(data));
        return;
    }
    content.innerHTML = postReactionForm(data.Posts[0], data.User.LoggedIn);
    attachPostReactionListener();
}

