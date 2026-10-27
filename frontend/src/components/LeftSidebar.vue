<script setup>
import { onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useUser } from '@/composables/useUser.js';
import { useChat } from '@/composables/useChat.js';
import DoubleChevronsLeftIcon from './Icons/DoubleChevronsLeftIcon.vue';
import DoubleChevronsRightIcon from './Icons/DoubleChevronsRightIcon.vue';
import { useLeftSidebar } from '@/composables/useLeftSidebar.js';
import { useAppData } from '@/composables/useAppData.js';

const route = useRoute();
const { user } = useUser();
const { ensureWs, myGroups } = useChat();
const { sidebarOpen, toggleLeftSidebar } = useLeftSidebar();
const { isMobileView } = useAppData();

onMounted(() => {
    ensureWs(user.value.Id);
});

watch(
    () => route.fullPath,
    () => (sidebarOpen.value = !isMobileView.value),
);
</script>

<template>
    <div :class="['sidebar-menu', { collapsed: !sidebarOpen }]">
        <div class="header">
            <span v-if="sidebarOpen" class="section-title">Menu</span>
            <button class="toggle" @click="toggleLeftSidebar">
                <DoubleChevronsLeftIcon v-if="sidebarOpen" />
                <DoubleChevronsRightIcon v-else />
            </button>
        </div>
        <template v-if="sidebarOpen">
            <router-link to="/" class="link">Home</router-link>
            <router-link to="/chats" class="link">Chats</router-link>
            <router-link :to="`/profile/view/${user.Id}`" class="link">My profile</router-link>

            <span class="section-title section">Posts</span>
            <router-link to="/post/create" class="link">Create post</router-link>

            <span class="section-title section">Groups</span>
            <router-link to="/groups" class="link">Browse groups</router-link>
            <router-link to="/group/create" class="link">Create group</router-link>
            <template v-if="myGroups.length">
                <span class="subsection">Joined</span>
                <router-link
                    v-for="group in myGroups"
                    :key="group.Id"
                    :to="`/group/view/${group.Id}`"
                    class="group-link"
                >
                    {{ group.Title }}
                </router-link>
            </template>

            <span class="section-title section">Events</span>
            <router-link to="/TODO" class="link">Browse events</router-link>
            <router-link to="/TODO" class="link">Create event</router-link>
        </template>
    </div>
</template>

<style scoped>
.sidebar-menu {
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    gap: 2px;
    box-sizing: border-box;
    width: 220px;
    padding: 16px 16px 24px 24px;
    overflow-y: auto;
    line-height: 1.2;
    background-color: var(--bg-base);
    border-right: 1px solid var(--accent-line);
}

.sidebar-menu.collapsed {
    width: 48px;
    padding: 16px 0;
    align-items: center;
}

.header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 0 8px 10px;
}

.collapsed .header {
    padding: 0;
}

.toggle {
    display: flex;
    padding: 4px;
    border: 1px solid var(--border);
    border-radius: 2px;
    background-color: var(--bg-input);
    color: var(--accent);
    cursor: pointer;
}

.section-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
}

.section {
    padding: 16px 10px 6px;
}

.link {
    padding: 6px 9px;
    border: 1px solid transparent;
    border-radius: 2px;
    font-size: 14px;
    color: var(--text-primary);
    text-decoration: none;
}

.link.router-link-exact-active {
    border-color: var(--accent-line);
    background-color: var(--bg-elevated);
    color: var(--accent);
    font-weight: 600;
}

.subsection {
    padding: 10px 10px 4px;
    font-size: 11px;
    font-weight: 600;
    color: var(--text-muted);
}

.group-link {
    padding: 5px 10px 5px 22px;
    font-size: 13px;
    color: var(--text-secondary);
    text-decoration: none;
}

@media (max-width: 768px) {
    .sidebar-menu.collapsed {
        display: none;
    }
}
</style>
