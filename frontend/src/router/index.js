import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import CategoryView from '@/views/CategoryView.vue';
import LoginView from '@/views/LoginView.vue';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', component: HomeView },
        { path: '/category/view/', component: CategoryView },
        { path: '/user/login', component: LoginView },
    ],
});

export default router;
