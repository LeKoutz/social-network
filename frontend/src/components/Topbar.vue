<script setup>
import { useAppData } from '@/composables/useAppData.js';
import { useUser } from '@/composables/useUser.js';
import { useNotifications } from '@/composables/useNotifications.js';
import NotificationsPanel from './NotificationsPanel.vue';
import UserAvatar from './UserAvatar.vue';
import BurgerIcon from './Icons/BurgerIcon.vue';
import BellIcon from './Icons/BellIcon.vue';
import ChevronDownIcon from './Icons/ChevronDownIcon.vue';
import { useLeftSidebar } from '../composables/useLeftSidebar.js';
import GenericDropdown from './GenericDropdown.vue';
import DropdownItem from './DropdownItem.vue';
import GenericButton from './GenericButton.vue';

const { websiteName, isMobileView } = useAppData();
const { user, logoutUser } = useUser();
const { unreadNotificationsCount } = useNotifications();
const { toggleLeftSidebar } = useLeftSidebar();
</script>

<template>
    <div class="topbar">
        <div class="title">
            <GenericButton v-if="isMobileView" variant="no-border" @click="toggleLeftSidebar">
                <template #leading><BurgerIcon /></template>
            </GenericButton>
            <router-link to="/">
                <h1>{{ websiteName }}</h1>
            </router-link>
        </div>
        <ul class="menu">
            <template v-if="user.LoggedIn">
                <li>
                    <GenericDropdown open-on="click" align="right" variant="elevated">
                        <template #trigger>
                            <GenericButton variant="no-border">
                                <template #leading v-if="isMobileView"><BellIcon /></template>
                                <template #default v-if="!isMobileView">Notifications</template>
                                <template #trailing v-if="unreadNotificationsCount > 0">
                                    <span class="count">{{ unreadNotificationsCount }}</span>
                                </template>
                            </GenericButton>
                        </template>
                        <template #default>
                            <NotificationsPanel />
                        </template>
                    </GenericDropdown>
                </li>
                <li>
                    <GenericDropdown open-on="hover" align="right" variant="elevated">
                        <template #trigger="{ open }">
                            <GenericButton variant="no-border">
                                <template #leading>
                                    <UserAvatar
                                        :url="user.AvatarURL"
                                        :name="user.Username"
                                        size="small"
                                    />
                                </template>
                                <template v-if="!isMobileView" #default>
                                    {{ user.Username }}
                                </template>
                                <template v-if="!isMobileView" #trailing>
                                    <ChevronDownIcon :class="{ 'transform-rotate-180deg': open }" />
                                </template>
                            </GenericButton>
                        </template>
                        <template #default>
                            <DropdownItem :to="`/profile/view/${user.Id}`">My profile</DropdownItem>
                            <DropdownItem to="/post/create">Create post</DropdownItem>
                            <DropdownItem to="/group/create">Create group</DropdownItem>
                            <DropdownItem to="/groups">Browse groups</DropdownItem>
                            <DropdownItem to="/user/posts">My posts</DropdownItem>
                            <DropdownItem to="/user/likes">My likes</DropdownItem>
                            <DropdownItem to="/follows">My follows</DropdownItem>
                            <hr class="divider" />
                            <DropdownItem variant="danger" @click="logoutUser()">
                                Log out
                            </DropdownItem>
                        </template>
                    </GenericDropdown>
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

.menu {
    display: flex;
    align-items: center;
    gap: 24px;
    margin: 0;
    padding: 0;
    list-style: none;
}

.count {
    padding: 1px 6px;
    border-radius: 2px;
    background-color: var(--accent-line);
    font-size: 12px;
    font-weight: 700;
}

.divider {
    margin: 4px 0;
    border: none;
    border-top: 1px solid var(--border);
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

    .count {
        padding: 1px 5px;
        font-size: 11px;
    }
}
</style>
