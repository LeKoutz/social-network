<script setup>
import { useAlerts } from '@/composables/useAlerts.js';
import { useUser } from '@/composables/useUser.js';
import { useRouter } from 'vue-router';

const router = useRouter();
const { user } = useUser();
const { setAlert } = useAlerts();

const MAX_FILE_SIZE = 20 * 1024 * 1024;
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'];
const MIN_DIMENSION = 400;
const MAX_DIMENSION = 4000;

function getImageDimensions(file) {
    return new Promise((resolve, reject) => {
        const img = new Image();
        const url = URL.createObjectURL(file);
        img.onload = () => {
            URL.revokeObjectURL(url);
            resolve({ width: img.naturalWidth, height: img.naturalHeight });
        };
        img.onerror = () => {
            URL.revokeObjectURL(url);
            reject(new Error('Could not read image'));
        };
        img.src = url;
    });
}

async function validateAvatar(file) {
    if (!file) return null;

    if (!ALLOWED_TYPES.includes(file.type)) {
        return 'Avatar must be a JPEG, PNG, or GIF image';
    }
    if (file.size > MAX_FILE_SIZE) {
        return 'Avatar must be under 20MB';
    }

    const { width, height } = await getImageDimensions(file);
    if (width < MIN_DIMENSION || height < MIN_DIMENSION) {
        return `Avatar must be at least ${MIN_DIMENSION}x${MIN_DIMENSION}px`;
    }
    if (width > MAX_DIMENSION || height > MAX_DIMENSION) {
        return `Avatar must be at most ${MAX_DIMENSION}x${MAX_DIMENSION}px`;
    }
    return null;
}

async function submitForm(e) {
    const avatarFile = e.target.avatar.files[0]
    const avatarError = await validateAvatar(avatarFile)
    if (avatarError) {
        setAlert({ Error: { Has: true, Message: avatarError } })
        return
    }

    const response = await fetch('/api/user/register', {
        method: 'POST',
        body: new FormData(e.target),
    })
    const data = await response.json()
    setAlert(data)
    if (data && !data.Error.Has) {
        router.push('/user/login')
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
                    <input placeholder="Username" type="text" name="username" required />
                    <input placeholder="First Name" type="text" name="first_name" required />
                    <input placeholder="Last Name" type="text" name="last_name" required />
                    <input
                        placeholder="Date of Birth"
                        type="date"
                        name="date_of_birth"
                        required
                    />
                    <input placeholder="Password" type="password" name="password1" required />
                    <input
                        placeholder="Confirm password"
                        type="password"
                        name="password2"
                        required
                    />
                    <fieldset>
                        <legend>Optional</legend>
                        <input placeholder="Nickname" type="text" name="nickname" />
                        <textarea placeholder="About Me" name="about_me"></textarea>
                        <input type="file" name="avatar" accept="image/jpeg,image/png,image/gif"/>
                        <input type="checkbox" name="private_profile" value="on" />
                    </fieldset>
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
