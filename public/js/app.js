import { TopBar } from './topbar.js';
import { showUserLogin, attachLoginListener } from './forms/user_login.js';
import { showUserRegister, attachRegisterListener } from './forms/user_register.js'
import { showPostCreateView, attachPostCreateListener, showPostEditView } from './forms/post_create.js'
import { ShowError } from './error.js';
import { Categories, showCategoryView } from './categories.js';
import { ShowMessage } from './message.js'
import { userLogout } from './logout.js';
import { showPostView, displayPost, } from './posts.js';
import { attachCommentCreateListener, attachCommentEditListener } from './forms/comment_create.js'
import { renderFooter } from './footer.js'

function showWelcome(data) {
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
`
}

const response = await fetch('/api/')
const data = await response.json()
const body = document.querySelector('body')
const topbar = document.createElement('div')
topbar.classList.add('topbar')
topbar.innerHTML = TopBar(data)
const alerts = document.createElement('div')
alerts.classList.add('alerts')
const content = document.createElement('div')
content.innerHTML = showWelcome(data);
content.classList.add('content')
const footer = document.createElement('div')
footer.classList.add('footer')
footer.innerHTML = renderFooter(data)
body.append(topbar, alerts, content, footer)

async function routeSelect() {
    const hash = window.location.hash
    const id = parseInt(hash.split('/').at(-1)) || ''
    switch (hash) {
        case '#/user/login':
            content.innerHTML = showUserLogin();
            attachLoginListener();
            break;
        case '#/user/register':
            content.innerHTML = showUserRegister();
            attachRegisterListener();
            break;
        case '#/user/logout':
            content.innerHTML = await userLogout()
            break
        case `#/category/view/${id}`:
            content.innerHTML = await showCategoryView(id);
            attachCommentCreateListener()
            break
        case `#/post/view/${id}`:
            content.innerHTML = await showPostView(id);
            attachCommentCreateListener()
            attachCommentEditListener()
            break
        case `#/post/create`:
            content.innerHTML = await showPostCreateView(data);
            attachPostCreateListener()
            break
        case `#/post/edit/${id}`:
            content.innerHTML = await showPostEditView(id);
            attachPostCreateListener({ editing: true })
            break
        case '#/' || '' || '#':
            content.innerHTML = showWelcome(data);
            break;
    }
}

routeSelect();

window.addEventListener('hashchange', () => {
    routeSelect()
});
