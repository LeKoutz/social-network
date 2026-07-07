import { showCategoryView } from './categories.js';
import { showUserLogin, attachLoginListener } from './forms/user_login.js';
import { showUserRegister, attachRegisterListener } from './forms/user_register.js';
import { showPostCreateView, attachPostCreateListener, showPostEditView } from './forms/post_create.js';
import { ShowUserActivity } from './user_activities.js';
import { userLogout } from './logout.js';
import { showPostView } from './posts.js';
import { ShowUserMenu } from './user.js';
import { showWelcome } from './welcome.js';
import { attachCommentCreateListener, attachCommentEditListener } from './forms/comment_create.js';
import { ShowError } from './error.js';

export async function routeSelect(data, content) {
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
    case '#/user/activity':
        content.innerHTML = await ShowUserActivity();
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
    case '#/':
    case '':
    case '#':
        content.innerHTML = showWelcome(data);
        break;
    default:
        ShowError({Error:{Message:"Soft 404"}});
        content.innerHTML = showWelcome(data);
        break;
    }
}
