import { showCategoryView } from './components/categories.js';
import { loginRoute } from './forms/user_login.js';
import { registerRoute } from './forms/user_register.js';
import { showPostCreateView, attachPostCreateListener } from './forms/post_create.js';
import { ShowUserActivity, ShowUserLikes, ShowUserPosts } from './components/user_activities.js';
import { logoutRoute } from './components/logout.js';
import { showPostView } from './components/posts.js';
import { userMenuRoute } from './components/user.js';
import { showWelcome } from './components/welcome.js';
import { attachCommentCreateListener, attachCommentEditListener } from './forms/comment_create.js';
import { ShowError } from './components/error.js';

export async function routeSelect() {
    const hash = window.location.hash;
    const id = parseInt(hash.split('/').at(-1)) || '';
    const content = document.querySelector('.content');
    switch (hash) {
    case '#/user':
        userMenuRoute();
        break;
    case '#/user/login':
        loginRoute();
        break;
    case '#/user/register':
        registerRoute();
        break;
    case '#/user/logout':
        logoutRoute();
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
        content.innerHTML = await showWelcome();
        break;
    }
    default:
        document.querySelector('.alerts').innerHTML = ShowError({Error:{Message:"Soft 404"}});
        content.innerHTML = await showWelcome();
        break;
    }
}
