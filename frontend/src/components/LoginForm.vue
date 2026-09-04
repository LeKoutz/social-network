<script setup>
import { useUser } from '@/composables/useUser.js';
import { useAlerts } from '@/composables/useAlerts';
import { useRouter } from 'vue-router';

const router = useRouter();
const { setAlert } = useAlerts();
const { user, setUser } = useUser();

if (user.value.LoggedIn) {
    setAlert({ Error: { Has: true, Message: 'You are already logged-in' } });
}

async function submitForm(e) {
    const response = await fetch('/api/user/login', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target)),
    });
    const data = await response.json();
    setAlert(data);
    if (data && !data.Error.Has) {
        setUser(data.User);
        router.push('/');
        // TODO: Handle Websocket
    }
}
</script>
<template>
    <template v-if="!user.LoggedIn">
        <div class="box">
            <form id="login" @submit.prevent="submitForm">
                <fieldset>
                    <legend>User Login</legend>
                    <input type="hidden" name="action" value="login" />
                    <input
                        type="identifier"
                        placeholder="Email or Nickname"
                        name="identifier"
                        required
                    />
                    <input type="password" placeholder="Password" name="password" required />
                    <input type="submit" value="Login" name="submit" />
                </fieldset>
            </form>
            <div class="oauth-buttons">
                <router-link id="google-auth" to="/auth/google">Continue with Google</router-link>
                <router-link id="github-auth" to="/auth/github">Continue with GitHub</router-link>
            </div>
            <p>
                You don't have an account?
                <router-link to="/user/register">Register!</router-link>
            </p>
        </div>
    </template>
</template>
<style scoped></style>
