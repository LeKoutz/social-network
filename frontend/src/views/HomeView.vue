<script setup>
import { useUser } from '@/composables/useUser';
import Categories from '@/components/Categories.vue';
import GroupsBrowser from '@/components/GroupsBrowser.vue';

const { user } = useUser();
</script>

<template>
    <div class="flex flex-column gap-20 padding-sides-1vw overflow-y-auto">
        <h2>Welcome {{ user.LoggedIn ? user.Username : '' }}!</h2>
        <p>
            Below, you will see the post categories available as well as a brief description when
            applicable. Click any of them to navigate to its list of posts.
        </p>
        <template v-if="!user.LoggedIn">
            <p>
                If you want to upload posts or comment on other ones or just react (like or dislike)
                comments or posts, you have to
                <router-link to="/user/login">login</router-link>
                .
            </p>
            <p>
                In case you don't have an account, you can create one by clicking
                <router-link to="/user/register">here</router-link>
                or via the menu at the top of the page.
            </p>
        </template>
        <div v-else class="flex flex-row flex-wrap gap-20">
            <Categories />
            <GroupsBrowser />
        </div>
    </div>
</template>

<style scoped></style>
