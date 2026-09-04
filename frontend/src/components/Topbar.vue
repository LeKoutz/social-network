<script setup>
import { useAppData } from '@/composables/useAppData.js';

const { websiteName } = useAppData();
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
                <!-- TODO: Expand with v-if user.LoggedIn case -->
                <li><router-link to="/user/login">Log in</router-link></li>
                <li><router-link to="/user/register">Register</router-link></li>
            </ul>
        </div>
    </div>
</template>

<style scoped></style>
