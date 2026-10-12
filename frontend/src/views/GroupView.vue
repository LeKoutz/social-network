<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch, apiPost } from '@/utils/api.js';
import { firstOrNull, DateToLocale } from '@/utils/utils.js';

import { useUser } from '@/composables/useUser.js';

import PostCreateForm from '@/components/PostCreateForm.vue';
import PostShow from '@/components/PostShow.vue';
import GroupAddIcon from '@/components/GroupAddIcon.vue';

const route = useRoute();
const { user } = useUser();
const group = ref(null);
const posts = ref(null);
const loading = ref(true);
const error = ref(null);

async function loadGroup(id) {
    loading.value = true;
    error.value = null;
    try {
        const data = await apiFetch(`/api/group/view/${id}`);
        if (data) {
            group.value = firstOrNull(data.Groups);
            posts.value = data.Posts;
            error.value = data.Error;
            // if (!group.value) {
            //     error.value = 'Group not found';
            // }
        }
    } catch (e) {
        console.log(e)
    }
    loading.value = false;
}

async function join(user,group){
    // const group = firstOrNull(group.value)
    const payload = new FormData();
    payload.append('id', group.Id)
    const data = await apiPost('/api/group/join', payload);
    console.log(data);
}

onMounted(() => loadGroup(route.params.id));

watch(
    () => route.params.id,
    (newId) => {
        if (newId) loadGroup(newId);
    },
);
</script>

<template>
    <div class="flex flex-column gap-20 overflow-y-auto padding-sides-1vw height-100">
        <template v-if="loading">Loading…</template>
        <template v-else-if="group">
                <h2>{{ group.Title }}</h2>
                <template v-if="group.Description.length === 0">
                    <p class="text-color-muted">No description avaliable.</p>
                </template>
                <template v-else>
                    <p class="group-description">{{ group.Description }}</p>
                </template>
                <p class="group-owner">
                Owned by 
                <router-link :to="`/profile/view/${group.OwnerUserId}`">
                    <strong>{{ group.OwnerUsername }}</strong>
                </router-link>
                on <em>({{ DateToLocale(group.Timestamp) }})</em>
                </p>
                <div class="flex flex-column gap-20" v-if="group.Member">
                    <PostCreateForm :group="group" />
                    <div class="flex flex-column gap-20" v-if="posts && posts.length">
                        <PostShow v-for="post in posts" :key="post.Id" :post="post" />
                    </div>
                    <p v-else>No posts in this group yet. Be the first to post!</p>
                </div>
                <template v-else class="group-locked">
                    This group is private. Only its members can see the posts and join the
                    conversation.
                    <button
                        class="flex flex-row align-items-center width-fit-content"
                        @click="join(user,group)"
                        >
                        <GroupAddIcon/> Join group
                    </button>
                </template>
        </template>
    </div>
</template>

<style scoped></style>
