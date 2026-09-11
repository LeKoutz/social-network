import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import CategoryView from '@/views/CategoryView.vue';
import LoginView from '@/views/LoginView.vue';
import RegisterView from '@/views/RegisterView.vue';
import ProfileView from '@/views/ProfileView.vue';
import UserActivityView from '@/views/UserActivityView.vue';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', component: HomeView },
        { path: '/category/view/', component: CategoryView },
        { path: '/user/login', component: LoginView },
        { path: '/user/register', component: RegisterView },
        { path: '/profile', component: ProfileView },
        { path: '/user/activity', component: UserActivityView },
    ],
});

export default router;
