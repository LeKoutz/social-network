import { TopBar } from './topbar.js';
import { ShowError } from './error.js';

function mainViewLayout(data) {
    return `
        ${TopBar(data)}
        ${(data.Error.Has) ? ShowError(data) : ''}
        <div class="container">
            <div class="welcome">
                <h2>Welcome${(data.User.LoggedIn)? data.User.Username:''}!</h2>
                <p>Below, you will see the post categories available as well as a brief
                description when applicable. Click any of them to navigate to its list
                of posts.</p>
                    ${(!data.User.LoggedIn)?
                            `<p>If you want to upload posts or comment to other ones or just react
                    (like or dislike) comments or posts, you have to <a href="/user/login">login</a>.</p>
                    <p>In case you don't have an account, you can create one by clicking
                    <a href="/user/register">here</a>  or via menu at the top of the page.</p>`:''}
            </div>
            <!-- {{template "categories" .}} -->
        </div>
        <div class="footer">
            <p>Version: ${data.Version}</p>
        </div>`
}

const response = await fetch('/api/')
const data = await response.json()
document.querySelector('body').innerHTML = mainViewLayout(data);
