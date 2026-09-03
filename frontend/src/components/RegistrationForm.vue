<script setup>
import { useAlerts } from '@/composables/useAlerts.js';
import { useUser } from '@/composables/useUser.js';
import { useRouter } from 'vue-router';

const router = useRouter();
const { user } = useUser();
const { setAlert } = useAlerts();

async function submitForm(e) {
    const response = await fetch('/api/user/register', {
        method: 'POST',
        body: new URLSearchParams(new FormData(e.target)),
    });
    const data = await response.json();
    setAlert(data);
    if (data && !data.Error.Has) {
        router.push('/user/login');
    }
}
</script>

<template>
    <template v-if="!user.LoggedIn">
        <div class="box">
            <form id="register" @submit.prevent="submitForm">
                <fieldset>
                    <legend>User Registration</legend>
                    <input type="hidden" name="action" value="register" />
                    <input placeholder="Email" type="email" name="email" required />
                    <input placeholder="Nickname" type="text" name="username" required />
                    <input placeholder="First Name" type="first_name" name="first_name" required />
                    <input placeholder="Last Name" type="last_name" name="last_name" required />
                    <input placeholder="Age" type="age" name="age" required />
                    <select name="gender" required>
                        <option value="">Gender</option>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                        <option value="other">Other</option>
                    </select>
                    <input placeholder="Password" type="password" name="password1" required />
                    <input
                        placeholder="Confirm password"
                        type="password"
                        name="password2"
                        required
                    />
                    <input type="submit" value="Register" name="submit" />
                </fieldset>
            </form>
            <div class="oauth-buttons">
                <router-link id="google-auth" to="/auth/google">Continue with Google</router-link>
                <router-link id="github-auth" to="/auth/github">Continue with GitHub</router-link>
            </div>
            <p>
                Already have an account?
                <router-link to="/user/login">Login!</router-link>
            </p>
        </div>
    </template>
</template>

<style scoped></style>
