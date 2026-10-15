<script setup>
import { ref, computed } from 'vue';
import { useUser } from '@/composables/useUser.js';
import ActivityList from '@/components/ActivityList.vue';
import { apiPost } from '@/utils/api.js';

const props = defineProps({
    profile: {
        type: Object,
        required: true,
    },
});

const { user } = useUser();
const activeTab = ref('activity');

const isOwnProfile = computed(
    () => user.value.LoggedIn && String(user.value.Id) === String(props.profile.UserId),
);

async function followRequest(id) {
    const payload = new FormData();
    payload.append('user-id', id);
    await apiPost('/api/follow', payload);
}
</script>

<template>
    <header class="profile-header">
        <div class="profile-identity">
            <img
                :src="profile.AvatarURL ? `/${profile.AvatarURL}` : '/default-avatar.png'"
                alt="Avatar"
                class="profile-avatar"
            />
            <h2>{{ profile.Username }}</h2>
            <h4 v-if="profile.CanView && profile.Nickname">({{ profile.Nickname }})</h4>
        </div>

        <div v-if="!profile.CanView" class="profile-locked">
            User profile is private, request follow to be able to see more info
            <button @click="followRequest(profile.UserId)">Follow</button>
        </div>
        <div v-else class="profile-details">
            <h3 v-if="profile.FirstName || profile.LastName" class="profile-fullname">
                {{ profile.FirstName }} {{ profile.LastName }}
            </h3>

            <router-link v-if="isOwnProfile" to="/profile/edit">Edit profile</router-link>
            <button v-else @click="followRequest(profile.UserId)">Follow</button>

            <p class="profile-counts">
                <span>Followers {{ profile.Followers.length }} |</span>
                <span>Following {{ profile.Following.length }}</span>
            </p>

            <p class="profile-info">
                <span v-if="profile.Email">{{ profile.Email }}</span>
                <span v-if="profile.DateOfBirth">
                    {{ profile.DateOfBirth }}
                </span>
                <span v-if="profile.About">{{ profile.About }}</span>
            </p>
        </div>

    </header>

    <div v-if="profile.CanView" class="profile-body">
        <nav class="profile-tabs">
            <button
                type="button"
                :class="{ active: activeTab === 'activity' }"
                @click="activeTab = 'activity'"
            >
                Activity
            </button>
            <button
                type="button"
                :class="{ active: activeTab === 'followers' }"
                @click="activeTab = 'followers'"
            >
                Followers
            </button>
            <button
                type="button"
                :class="{ active: activeTab === 'following' }"
                @click="activeTab = 'following'"
            >
                Following
            </button>
        </nav>

        <section v-if="activeTab === 'activity'">
            <ActivityList :activities="profile.Activities" />
        </section>
        <section v-else-if="activeTab === 'followers'">
            <h2>Followers</h2>
            <ul v-for="follower in profile.Followers">
                <li>
                    <router-link :to="`/profile/view/${follower.UserId}`" :key="follower.UserId">
                        {{ follower.Username }}
                    </router-link>
                </li>
            </ul>
        </section>
        <section v-else-if="activeTab === 'following'">
            <h2>Following</h2>
            <ul v-for="followee in profile.Following">
                <li>
                    <router-link :to="`/profile/view/${followee.UserId}`" :key="followee.UserId">
                        {{ followee.Username }}
                    </router-link>
                </li>
            </ul>
        </section>
    </div>
</template>

<style scoped>
.profile-header {
    gap: 10%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.profile-identity {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.profile-info {
    display: flex;
    flex-direction: column;
}

.profile-tabs {
    display: flex;
    justify-content: center;
}
.profile-avatar {
    max-height: 128px;
    max-width: 128px;
}
</style>
