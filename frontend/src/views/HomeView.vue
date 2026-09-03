<script setup>
import { useUser } from '@/composables/useUser';

const { user } = useUser();
</script>

<template>
    <div class="welcome">
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
        <!-- TODO: Show Categories if LoggedIn -->
    </div>
</template>

<style scoped></style>
