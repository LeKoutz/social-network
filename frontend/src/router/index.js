import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import CategoryView from '@/views/CategoryView.vue';
import LoginView from '@/views/LoginView.vue';
import RegisterView from '@/views/RegisterView.vue';
import ProfileView from '@/views/ProfileView.vue';
import UserActivityView from '@/views/UserActivityView.vue';
import PostView from '@/views/PostView.vue';
import PostCreateView from '@/views/PostCreateView.vue';
import PostEditView from '@/views/PostEditView.vue';
import { useUser } from '@/composables/useUser.js';
import GroupCreateView from '@/views/GroupCreateView.vue';
import GroupsBrowserView from '@/views/GroupsBrowserView.vue';
import GroupView from '@/views/GroupView.vue';
import ChatView from '@/views/ChatView.vue';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', component: HomeView },
        { path: '/category/view/', component: CategoryView },
        { path: '/user/login', component: LoginView },
        { path: '/user/register', component: RegisterView },
        { path: '/profile', component: ProfileView },
        { path: '/user/activity', component: UserActivityView },
        { path: '/post/view/:id', component: PostView },
        {
            path: '/post/create',
            component: PostCreateView,
            meta: { requiresAuth: true },
        },
        {
            path: '/post/edit/:id',
            component: PostEditView,
            meta: { requiresAuth: true },
        },
        { path: '/group/create', component: GroupCreateView },
        { path: '/groups', component: GroupsBrowserView },
        { path: '/group/view/:id', component: GroupView },
        {
            path: '/chat/:id',
            component: ChatView,
            meta: { requiresAuth: true },
        },
    ],
});

router.beforeEach((to) => {
    if (to.meta.requiresAuth) {
        const { user } = useUser();
        if (!user.value.LoggedIn) return '/user/login';
    }
});

export default router;
