import { TopBar } from './topbar.js';
import { showUserLogin } from './forms/user_login.js';
import { showUserRegister, attachRegisterListener } from './forms/user_register.js'
import { showPostCreateView, attachPostCreateListener, showPostEditView } from './forms/post_create.js'
import { ShowError } from './error.js';
import { Categories, showCategoryView } from './categories.js';
import { ShowMessage } from './message.js'
import { userLogout } from './logout.js';
import { showPostView, displayPost, } from './posts.js';
import { attachCommentCreateListener } from './forms/comment_create.js'

function mainViewLayout(data) {
    return `
        ${TopBar(data)}
        ${(data.Error.Has) ? ShowError(data) : ''}
        ${(data.Message.Has) ? ShowMessage(data) : ''}
        <div class="container">
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
        </div>
        <div class="footer">
            <p>Version: ${data.Version}</p>
        </div>`
}

const response = await fetch('/api/')
const data = await response.json()
const body = document.querySelector('body');
body.innerHTML = mainViewLayout(data);
const container = document.querySelector('.container');

async function routeSelect() {
    const hash = window.location.hash
    const id = parseInt(hash.split('/').at(-1)) || ''
    switch (hash) {
        case '#/user/login':
            container.innerHTML = showUserLogin();
            break;
        case '#/user/register':
            container.innerHTML = showUserRegister();
            attachRegisterListener(data, container);
            break;
        case '#/user/logout':
            container.innerHTML = await userLogout()
            break
        case `#/category/view/${id}`:
            container.innerHTML = await showCategoryView(id);
            attachCommentCreateListener()
            break
        case `#/post/view/${id}`:
            container.innerHTML = await showPostView(id);
            attachCommentCreateListener()
            break
        case `#/post/create`:
            container.innerHTML = await showPostCreateView(data);
            attachPostCreateListener()
            break
        case `#/post/edit/${id}`:
            container.innerHTML = await showPostEditView(id);
            attachPostCreateListener({ editing: true })
            break
        case '#/' || '' || '#':
            body.innerHTML = mainViewLayout(data);
            break;
    }
}

routeSelect();

window.addEventListener('hashchange', () => {
    routeSelect()
});
