import { categoryRoute } from './components/categories.js';
import { loginRoute } from './forms/user_login.js';
import { registerRoute } from './forms/user_register.js';
import { postCreateRoute } from './forms/post_create.js';
import { postEditRoute } from './forms/post_edit.js';
import { userActivityRoute, ShowUserLikes, userPostsRoute } from './components/user_activities.js';
import { logoutRoute } from './components/logout.js';
import { postRoute } from './components/posts.js';
import { userMenuRoute } from './components/user.js';
import { indexRoute } from './components/welcome.js';
import { ShowError } from './components/error.js';

export async function routeSelect() {
    const hash = window.location.hash;
    const id = parseInt(hash.split('/').at(-1)) || '';
    const content = document.querySelector('.content');
    switch (hash) {
    case '#/user':
        await userMenuRoute();
        break;
    case '#/user/login':
        await loginRoute();
        break;
    case '#/user/register':
        await registerRoute();
        break;
    case '#/user/logout':
        await logoutRoute();
        break;
    case '#/user/likes':
        content.innerHTML = await ShowUserLikes();
        break;
    case '#/user/posts':
        await userPostsRoute();
        break;
    case '#/user/activity':
        userActivityRoute();
        break;
    case `#/category/view/${id}`:
        await categoryRoute(id);
        break;
    case `#/post/view/${id}`:
        await postRoute(id);
        break;
    case `#/post/create`:
        await postCreateRoute();
        break;
    case `#/post/edit/${id}`:
        await postEditRoute(id);
        break;
    case '#/':
    case '':
    case '#': {
        await indexRoute();
        break;
    }
    default:
        document.querySelector('.alerts').innerHTML = ShowError({Error:{Message:"Soft 404"}});
        await indexRoute();
        break;
    }
}
