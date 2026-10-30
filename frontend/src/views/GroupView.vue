<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiGet, apiPost } from '@/utils/api.js';
import { firstOrNull, DateToLocale } from '@/utils/utils.js';

import { useUser } from '@/composables/useUser.js';

import PostCreateForm from '@/components/PostCreateForm.vue';
import PostShow from '@/components/PostShow.vue';
import GroupAddIcon from '@/components/GroupAddIcon.vue';
import GroupHeader from '@/components/GroupHeader.vue';
import PostsList from '@/components/PostsList.vue';
import GroupInviteModal from '../components/GroupInviteModal.vue';

const route = useRoute();
const { user } = useUser();
const group = ref(null);
const posts = ref(true);
const loading = ref(true);
const error = ref(null);

const owner = ref(false);

async function loadGroup(id) {
    loading.value = true;
    error.value = null;
    try {
        const data = await apiGet(`/api/group/view/${id}`);
        if (data) {
            group.value = firstOrNull(data.Groups);
            posts.value = data.Posts;
            error.value = data.Error;

            // if (!group.value) {
            //     error.value = 'Group not found';
            // }
        }
    } catch (e) {
        console.log(e);
    }
    loading.value = false;
}

async function join(user, group) {
    // const group = firstOrNull(group.value)
    const payload = new FormData();
    payload.append('group-id', group.Id);
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
    <div class="flex flex-column gap-20 overflow-y-auto padding-sides-1vw">
        <template v-if="loading">Loading…</template>
        <template v-else-if="group">
            <GroupHeader :group="group" />
            <GroupInviteModal
                :group="group"
                :user="user"
                v-if="group.Member || user.Id === group.OwnerId"
            />
            <div class="flex flex-column gap-20" v-if="group.Member || user.Id === group.OwnerId">
                <PostCreateForm :group="group" />
                <PostsList :posts="posts" />
            </div>
            <template v-else class="group-locked">
                This group is private. Only its members can see the posts and join the conversation.
                <button
                    class="flex flex-row align-items-center width-fit-content"
                    @click="join(user, group)"
                >
                    <GroupAddIcon />
                    Join group
                </button>
            </template>
        </template>
    </div>
</template>

<style scoped></style>
