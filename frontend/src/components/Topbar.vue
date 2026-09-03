<script setup>
import { useAppData } from '@/composables/useAppData.js';
import { useUser } from '@/composables/useUser.js';
import { useNotifications } from '@/composables/useNotifications.js';
import NotificationsPanel from './NotificationsPanel.vue';

const { websiteName } = useAppData();
const { user, logoutUser } = useUser();
const { unreadNotificationsCount } = useNotifications();
</script>

<template>
    <div class="topbar">
        <div class="title">
            <router-link to="/">
                <h1>{{ websiteName }}</h1>
            </router-link>
        </div>
        <div class="menu">
            <ul>
                <template v-if="user.LoggedIn">
                    <li>
                        <router-link to="/user">Welcome back, {{ user.Username }}</router-link>
                    </li>
                    <li>
                        <ul>
                            <li><router-link to="/post/create">Create post</router-link></li>
                            <li><router-link to="/user/activity">My activity</router-link></li>
                            <li><router-link to="/user/posts">My posts</router-link></li>
                            <li><router-link to="/user/likes">My likes</router-link></li>
                            <li>
                                <details>
                                    <summary>
                                        Notifications 🔔 {{ unreadNotificationsCount }}
                                    </summary>
                                    <NotificationsPanel v-if="unreadNotificationsCount > 0" />
                                </details>
                            </li>
                            <li><button @click="logoutUser">Log out</button></li>
                        </ul>
                    </li>
                </template>
                <template v-else>
                    <li><router-link to="/user/login">Log in</router-link></li>
                    <li><router-link to="/user/register">Register</router-link></li>
                </template>
            </ul>
        </div>
    </div>
</template>

<style scoped></style>
