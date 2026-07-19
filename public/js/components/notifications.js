import { ShowError } from '../components/error.js';
import { apiFetch } from '../fetchers/api.js';
import { renderTopBar } from '../partials/topbar.js';

export function expandNotificationsPanel(data) {
    return `<div class="notifications">
    <form id="mark-as-read">
        <button type="submit">Mark all as read</button>
    </form>
    <ul>
        ${showNotificationsList(data)}
    </ul>
</div>`;
}

function showNotificationsList(data) {
    return data.User.Notifications ? data.User.Notifications.map(notification => {
        const strongStart = notification.Read ? '' : '<strong>';
        const strongEnd = notification.Read ? '' : '</strong>';
        switch (notification.Type) {
        case 'comment':
            return `${strongStart}<li><a href="#/post/view/${notification.PostId}#comment-${notification.CommentId}">${notification.Actor.Username} commented on your post</a></li>${strongEnd}`;
        case 'like':
            return `${strongStart}<li><a href="#/post/view/${notification.PostId}">${notification.Actor.Username} liked your post</a></li>${strongEnd}`;
        case 'dislike':
            return `${strongStart}<li><a href="#/post/view/${notification.PostId}">${notification.Actor.Username} disliked your post</a></li>${strongEnd}`;
        case 'commentLike':
            return `${strongStart}<li><a href="#/post/view/${notification.PostId}#comment-${notification.CommentId}">${notification.Actor.Username} liked your comment</a></li>${strongEnd}`;
        case 'commentDislike':
            return `${strongStart}<li><a href="#/post/view/${notification.PostId}#comment-${notification.CommentId}">${notification.Actor.Username} disliked your comment</a></li>${strongEnd}`;
        default:
            return '';
        }
    }).join('') : '';
}

export async function attachNotificationsButtonListener() {
    const form = document.querySelector('#mark-as-read');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const response = await fetch('/api/user/notifications', {
                method: 'POST',
                body: new URLSearchParams(new FormData(e.target))
            });
            const data = await response.json();
            data.Error.Has ? document.querySelector('.alerts').innerHTML = ShowError(data) : '';
            renderTopBar(data);
        });
    }
}