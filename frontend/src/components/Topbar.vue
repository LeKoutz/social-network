<script setup>
import { useAppData } from '@/composables/useAppData.js';
import { useUser } from '@/composables/useUser.js';
import { useNotifications } from '@/composables/useNotifications.js';
import NotificationsPanel from './NotificationsPanel.vue';
import UserAvatar from './UserAvatar.vue';
import BurgerIcon from './Icons/BurgerIcon.vue';
import BellIcon from './Icons/BellIcon.vue';
import ChevronDownIcon from './Icons/ChevronDownIcon.vue';

const { websiteName } = useAppData();
const { user, logoutUser } = useUser();
const { unreadNotificationsCount } = useNotifications();
</script>

<template>
    <div class="topbar">
        <div class="title">
            <button class="menu-button"><BurgerIcon /></button>
            <router-link to="/">
                <h1>{{ websiteName }}</h1>
            </router-link>
        </div>
        <ul class="menu">
            <template v-if="user.LoggedIn">
                <li>
                    <details>
                        <summary>
                            <BellIcon />
                            <span class="label">Notifications</span>
                            <span v-if="unreadNotificationsCount > 0" class="count">
                                {{ unreadNotificationsCount }}
                            </span>
                        </summary>
                        <NotificationsPanel />
                    </details>
                </li>
                <li>
                    <div class="user-menu">
                        <UserAvatar :url="user.AvatarURL" :name="user.Username" size="small" />
                        <span class="username">{{ user.Username }}</span>
                        <ChevronDownIcon class="chevron"/>
                    </div>
                    <ul class="user-dropdown">
                        <li>
                            <router-link :to="`/profile/view/${user.Id}`">My profile</router-link>
                        </li>
                        <li><router-link to="/post/create">Create post</router-link></li>
                        <li><router-link to="/group/create">Create group</router-link></li>
                        <li><router-link to="/groups">Browse groups</router-link></li>
                        <li><router-link to="/user/posts">My posts</router-link></li>
                        <li><router-link to="/user/likes">My likes</router-link></li>
                        <li><router-link to="/follows">My follows</router-link></li>
                        <li class="logout"><button @click="logoutUser">Log out</button></li>
                    </ul>
                </li>
            </template>
            <template v-else>
                <li><router-link to="/user/login">Log in</router-link></li>
                <li><router-link to="/user/register">Register</router-link></li>
            </template>
        </ul>
    </div>
</template>

<style scoped>
.topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
    padding: 0 24px;
    line-height: 1.2;
    background-color: var(--bg-base);
    border-bottom: 1px solid var(--accent-line);
}

.topbar a {
    text-decoration: none;
}

.title {
    display: flex;
    align-items: center;
    gap: 12px;
}

.title h1 {
    margin: 0;
    font-size: 20px;
    font-weight: 700;
    color: var(--accent);
}

.menu-button,
.bell {
    display: none;
}

.menu {
    display: flex;
    align-items: center;
    gap: 24px;
    margin: 0;
    padding: 0;
    list-style: none;
}

.menu > li {
    position: relative;
}

.menu > li:hover > .user-dropdown {
    display: flex;
}

summary {
    display: flex;
    align-items: center;
    gap: 6px;
    list-style: none;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
}

.count {
    padding: 1px 6px;
    border-radius: 2px;
    background-color: var(--accent-line);
    font-size: 12px;
    font-weight: 700;
}

.user-menu {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 3px 7px 3px 3px;
    border: 1px solid var(--border);
    border-radius: 2px;
    font-size: 14px;
    font-weight: 500;
}

.chevron {
    color: var(--text-secondary);
}

.user-dropdown {
    display: none;
    flex-direction: column;
    position: absolute;
    top: 100%;
    right: 0;
    z-index: 1;
    box-sizing: border-box;
    width: 200px;
    margin: 0;
    padding: 5px;
    list-style: none;
    background-color: var(--bg-elevated);
    border: 1px solid var(--accent-line);
    border-radius: 2px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.user-dropdown a,
.user-dropdown button {
    display: block;
    width: 100%;
    box-sizing: border-box;
    padding: 8px 10px;
    border: none;
    border-radius: 2px;
    background: none;
    font: inherit;
    font-size: 14px;
    text-align: left;
    text-decoration: none;
    color: var(--text-primary);
}

.user-dropdown a:hover,
.user-dropdown button:hover {
    background-color: var(--bg-input);
    color: var(--accent);
}

.user-dropdown .logout {
    border-top: 1px solid var(--border);
}

.user-dropdown .logout button {
    color: var(--danger);
}

@media (max-width: 768px) {
    .topbar {
        height: 56px;
        padding: 0 16px 0 12px;
    }

    .title h1 {
        font-size: 17px;
    }

    .menu {
        gap: 12px;
    }

    .menu-button {
        display: flex;
        padding: 6px;
        border: 1px solid var(--border);
        border-radius: 2px;
        background: none;
        color: var(--text-primary);
    }

    .bell {
        display: block;
    }

    .label,
    .username,
    .chevron {
        display: none;
    }

    summary {
        gap: 4px;
    }

    .count {
        padding: 1px 5px;
        font-size: 11px;
    }

    .user-menu {
        padding: 3px;
    }
}
</style>
