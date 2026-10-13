<script setup>
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { apiPost } from '@/utils/api.js';
import { useAlerts } from '@/composables/useAlerts.js';

const props = defineProps({
    profile: {
        type: Object,
        required: true,
    },
});

const router = useRouter();
const { setAlert } = useAlerts();
const submitting = ref(false);
const removeAvatar = ref(false);

const form = reactive({
    first_name: props.profile.FirstName ?? '',
    last_name: props.profile.LastName ?? '',
    nickname: props.profile.Nickname ?? '',
    date_of_birth: props.profile.DateOfBirth ?? '',
    about: props.profile.About ?? '',
    private_profile: props.profile.PrivateProfile,
});

async function submit(event) {
    if (submitting.value) return;

    const data = new FormData(event.currentTarget);
    const avatar = data.get('avatar');

    if (!avatar?.name || removeAvatar.value) {
        data.delete('avatar');
    }

    if (!removeAvatar.value && avatar?.size > 20 * 1024 * 1024) {
        setAlert({
            Error: {
                Has: true,
                Message: 'Avatar must be no larger than 20 MB.',
            },
        });
        return;
    }

    submitting.value = true;

    try {
        const response = await apiPost('/api/profile/edit', data);
        if (!response) return;

        setAlert(response);
        await router.push(`/profile/view/${props.profile.UserId}`);
    } catch {
        setAlert({
            Error: {
                Has: true,
                Message: 'Could not save your profile. Please try again.',
            },
        });
    } finally {
        submitting.value = false;
    }
}
</script>

<template>
    <form @submit.prevent="submit">
        <fieldset :disabled="submitting">
            <legend>Edit profile</legend>

            <label>
                First name
                <input v-model="form.first_name" name="first_name" required />
            </label>

            <label>
                Last name
                <input v-model="form.last_name" name="last_name" required />
            </label>

            <label>
                Nickname
                <input v-model="form.nickname" name="nickname" />
            </label>

            <label>
                Date of birth
                <input
                    v-model="form.date_of_birth"
                    type="date"
                    name="date_of_birth"
                />
            </label>

            <label>
                About
                <textarea v-model="form.about" name="about"></textarea>
            </label>

            <label>
                <input
                    v-model="form.private_profile"
                    type="checkbox"
                    name="private_profile"
                    value="on"
                />
                Private profile
            </label>

            <img
                v-if="profile.AvatarURL && !removeAvatar"
                :src="`/${profile.AvatarURL}`"
                alt="Current avatar"
                class="current-avatar"
            />

            <label>
                Replace avatar — JPEG, PNG or GIF, up to 20 MB
                <input
                    type="file"
                    name="avatar"
                    accept="image/jpeg,image/png,image/gif"
                    :disabled="removeAvatar"
                />
            </label>

            <label v-if="profile.AvatarURL">
                <input
                    v-model="removeAvatar"
                    type="checkbox"
                    name="remove_avatar"
                    value="on"
                />
                Remove current avatar
            </label>

            <button type="submit">
                {{ submitting ? 'Saving…' : 'Save profile' }}
            </button>

            <button
                type="button"
                @click="router.push(`/profile/view/${profile.UserId}`)"
            >
                Cancel
            </button>
        </fieldset>
    </form>
</template>

<style scoped>
label {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.current-avatar {
    max-width: 150px;
    max-height: 150px;
    object-fit: contain;
}
</style>