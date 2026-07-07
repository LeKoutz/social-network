import { TopBar } from './topbar.js';
import { showUserLogin, attachLoginListener } from './forms/user_login.js';
import { showUserRegister, attachRegisterListener } from './forms/user_register.js';
import { showPostCreateView, attachPostCreateListener, showPostEditView } from './forms/post_create.js';
import { ShowError } from './error.js';
import { Categories, showCategoryView } from './categories.js';
import { ShowMessage } from './message.js';
import { userLogout } from './logout.js';
import { showPostView, displayPost, } from './posts.js';
import { ShowUserMenu } from './user.js';
import { attachCommentCreateListener, attachCommentEditListener } from './forms/comment_create.js';
import { renderFooter } from './footer.js';
import { showWelcome } from './welcome.js';

const response = await fetch('/api/');
const data = await response.json();
const body = document.querySelector('body');
const topbar = document.createElement('div');
topbar.classList.add('topbar');
topbar.innerHTML = TopBar(data);
const alerts = document.createElement('div');
alerts.classList.add('alerts');
const content = document.createElement('div');
content.innerHTML = showWelcome(data);
content.classList.add('content');
const footer = document.createElement('div');
footer.classList.add('footer');
footer.innerHTML = renderFooter(data);
body.append(topbar, alerts, content, footer);

async function routeSelect() {
    const hash = window.location.hash;
    const id = parseInt(hash.split('/').at(-1)) || '';
    switch (hash) {
    case '#/user':
        content.innerHTML = ShowUserMenu();
        break;
    case '#/user/login':
        content.innerHTML = showUserLogin();
        attachLoginListener();
        break;
    case '#/user/register':
        content.innerHTML = showUserRegister();
        attachRegisterListener();
        break;
    case '#/user/logout':
        content.innerHTML = await userLogout();
        break;
    case `#/category/view/${id}`:
        content.innerHTML = await showCategoryView(id);
        attachCommentCreateListener();
        break;
    case `#/post/view/${id}`:
        content.innerHTML = await showPostView(id);
        attachCommentCreateListener();
        attachCommentEditListener();
        break;
    case `#/post/create`:
        content.innerHTML = showPostCreateView(data);
        attachPostCreateListener();
        break;
    case `#/post/edit/${id}`:
        content.innerHTML = await showPostEditView(id);
        attachPostCreateListener({ editing: true });
        break;
    case '#/' || '' || '#':
        content.innerHTML = showWelcome(data);
        break;
    default:
        content.innerHTML = showWelcome(data);
        break;
    }
}

routeSelect();

window.addEventListener('hashchange', () => {
    routeSelect();
});
