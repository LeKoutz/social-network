<script setup>
import { useUser } from '@/composables/useUser.js';
import { useRouter } from 'vue-router';
import { apiPost } from '@/utils/api.js';

const router = useRouter();
const { user } = useUser();

async function submitForm(e) {
    const data = await apiPost('/api/group/create', new FormData(e.target));
    if (data) {
        router.push('/');
    }
}
</script>

<template>
    <template v-if="user.LoggedIn">
        <div class="box">
            <form id="group-create" @submit.prevent="submitForm">
                <fieldset>
                    <legend>Create Group</legend>
                    <input placeholder="Title" type="text" name="title" required maxlength="127" />
                    <textarea placeholder="Description" name="description"></textarea>
                    <input type="submit" value="Create Group" name="submit" />
                </fieldset>
            </form>
        </div>
    </template>
    <p v-else>
        You must be logged in to create a group.
        <router-link to="/user/login">Log in</router-link>
    </p>
</template>

<style scoped></style>