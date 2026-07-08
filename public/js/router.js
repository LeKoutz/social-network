import { showCategoryView } from './components/categories.js';
import { showUserLogin, attachLoginListener } from './forms/user_login.js';
import { showUserRegister, attachRegisterListener } from './forms/user_register.js';
import { showPostCreateView, attachPostCreateListener } from './forms/post_create.js';
import { ShowUserActivity, ShowUserLikes, ShowUserPosts } from './components/user_activities.js';
import { userLogout } from './components/logout.js';
import { showPostView } from './components/posts.js';
import { ShowUserMenu } from './components/user.js';
import { showWelcome } from './components/welcome.js';
import { attachCommentCreateListener, attachCommentEditListener } from './forms/comment_create.js';
import { ShowError } from './components/error.js';
import { TopBar } from './partials/topbar.js';

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
        setTimeout(() => window.location.hash = '',1000);
        break;
    case '#/user/likes':
        content.innerHTML = await ShowUserLikes();
        break;
    case '#/user/posts':
        content.innerHTML = await ShowUserPosts();
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
        content.innerHTML = await showPostCreateView();
        attachPostCreateListener();
        break;
    case `#/post/edit/${id}`:
        content.innerHTML = await showPostCreateView({ editing: true, postId: id });
        attachPostCreateListener({ editing: true, postId: id });
        break;
    case '#/':
    case '':
    case '#': {
        const response = await fetch('/api/');
        const data = await response.json();
        content.innerHTML = showWelcome(data);
        document.querySelector('.topbar').innerHTML = TopBar(data);
        document.querySelector('.alerts').innerHTML = '';
        break;
    }
    default:
        ShowError({Error:{Message:"Soft 404"}});
        content.innerHTML = showWelcome(data);
        break;
    }
}
