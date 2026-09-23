<script setup>
import { ref, computed } from 'vue';
import { useUser } from '@/composables/useUser.js';
import ActivityList from '@/components/ActivityList.vue';

const props = defineProps({
	profile: {
		type: Object,
		required: true,
	},
});

const { user } = useUser();
const activeTab = ref('posts');

const isOwnProfile = computed(
	() => user.value.LoggedIn && String(user.value.Id) === String(props.profile.UserId)
);
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
            <h4 v-if="profile.CanView" >({{ profile.Nickname }})</h4>
        </div>
        <div v-if="profile.CanView" class="profile-details">
            <h3>
                <span v-if="profile.CanView && (profile.FirstName || profile.LastName)" class="profile-fullname">
                    {{ profile.FirstName }} {{ profile.LastName }}
                </span>
            </h3>
            <button v-if="isOwnProfile" type="button" disabled title="Edit profile — coming soon">
                Edit
            </button>
            <button v-else type="button">Follow</button>
            <p class="profile-counts">
                <span>Followers | </span>
                <span>Following</span>
            </p>
            <p v-if="profile.CanView" class="profile-info">
                <span v-if="profile.Email">{{ profile.Email }}</span>
                <span v-if="profile.DateOfBirth"> · {{ profile.DateOfBirth }}</span>
                <span v-if="profile.About"> · {{ profile.About }}</span>
            </p>
		</div>
        <p v-else class="profile-locked">
            User profile is private, request follow to be able to see more info
        </p>
	</header>

    <div class="profile-body">
        <template v-if="profile.CanView">
            
                <nav class="profile-tabs">
                    <button type="button" :class="{ active: activeTab === 'activity' }" @click="activeTab = 'activity'">
                        Activity
                    </button>
                    <button type="button" :class="{ active: activeTab === 'followers' }" @click="activeTab = 'followers'">
                        Followers
                    </button>
                    <button type="button" :class="{ active: activeTab === 'following' }" @click="activeTab = 'following'">
                        Following
                    </button>
                </nav>

                <section v-if="activeTab === 'activity'">
                    <ActivityList :activities="profile.Activities" />
                </section>
                <section v-else-if="activeTab === 'followers'">
                    <h2>Followers</h2>
                </section>
                <section v-else-if="activeTab === 'following'">
                    <h2>Following</h2>
                </section>
        </template>
    </div>
</template>

<style scoped>
header {
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
</style>