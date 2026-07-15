import { categoryRoute } from './components/categories.js';
import { loginRoute } from './forms/user_login.js';
import { registerRoute } from './forms/user_register.js';
import { postCreateRoute } from './forms/post_create.js';
import { postEditRoute } from './forms/post_edit.js';
import { userActivityRoute, userLikesRoute, userPostsRoute } from './components/user_activities.js';
import { logoutRoute } from './components/logout.js';
import { postRoute } from './components/posts.js';
import { userMenuRoute } from './components/user.js';
import { indexRoute } from './components/welcome.js';
import { ShowError } from './components/error.js';
import { chatRoute } from './forms/chat_send.js';

export async function routeSelect() {
    const hash = window.location.hash?.split('#');
    if (hash.length < 2) {
        await indexRoute();
        return;
    }
    const id = parseInt(hash[1].split('/').at(-1)) || '';
    const content = document.querySelector('.content');
    switch (hash[1]) {
    case '/user':
        await userMenuRoute();
        break;
    case '/user/login':
        await loginRoute();
        break;
    case '/user/register':
        await registerRoute();
        break;
    case '/user/logout':
        await logoutRoute();
        break;
    case '/user/likes':
        content.innerHTML = await ShowUserLikes();
        break;
    case '/user/posts':
        content.innerHTML = await ShowUserPosts();
        break;
    case '/user/activity':
        userActivityRoute();
        break;
    case `/category/view/${id}`:
        await categoryRoute(id);
        break;
    case `/post/view/${id}`:
        if ( hash.length === 3 ) {
            await postRoute(id, hash[2]);
        }
        await postRoute(id);
        break;
    case `/post/create`:
        await postCreateRoute();
        break;
    case `/post/edit/${id}`:
        await postEditRoute(id);
        break;
    case `/chat/${id}`:
        await chatRoute(id);
        break;
    case '/':
    case '': {
        await indexRoute();
        break;
    }
    default:
        document.querySelector('.alerts').innerHTML = ShowError({Error:{Message:"Soft 404"}});
        await indexRoute();
        break;
    }
}
