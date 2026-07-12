import { Categories } from './categories.js';
import { TopBar } from '../partials/topbar.js';
import { apiFetch } from '../fetchers/api.js';
import { renderFooter } from '../partials/footer.js';
import { connectWS } from '../ws.js';
import { UsersPanel, addUsersPanelButtonListener } from './users_panel.js';

export function showWelcome(data) {
    return `
        <div class="welcome">
            <h2>Welcome ${(data.User.LoggedIn)? data.User.Username:''}!</h2>
            <p>Below, you will see the post categories available as well as a brief
            description when applicable. Click any of them to navigate to its list
            of posts.</p>
                ${(!data.User.LoggedIn)?
        `<p>If you want to upload posts or comment to other ones or just react
                (like or dislike) comments or posts, you have to <a href="#/user/login">login</a>.</p>
                <p>In case you don't have an account, you can create one by clicking
                <a href="#/user/register">here</a>  or via menu at the top of the page.</p>`:''}
        </div>
        ${Categories(data)}
`;
}

export async function indexRoute() {
    const data = await apiFetch('/api/');
    if (data) {
        document.querySelector('.content').innerHTML = showWelcome(data);
        document.querySelector('.topbar').innerHTML = TopBar(data);
        document.querySelector('.footer').innerHTML = renderFooter(data);
    }
    if (data && data.User.LoggedIn) {
        const usersData = await apiFetch('/api/users');
        if (usersData) {
            document.querySelector('.users-panel').innerHTML = UsersPanel(usersData);
        }
        addUsersPanelButtonListener();
        connectWS();
    } else if (data && !data.User.LoggedIn) {
        document.querySelector('.users-panel').innerHTML = '';
    }
}
