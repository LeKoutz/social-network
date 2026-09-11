<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { apiFetch } from '@/utils/api.js';
import ActivityPost from '@/components/ActivityPost.vue';
import ActivityComment from '@/components/ActivityComment.vue';
import { DateToLocale } from '@/utils/utils.js';

const route = useRoute();
const activities = ref(null);
const user = ref(null);

function describeActivity(activity) {
    switch(activity.Type) {
    case "post":
        return ' posted: ';
    case "comment":
        return ' commented: ';

    case "postLike":
        return ' liked a post: ';

    case "postDislike":
        return ' disliked a post: ';

    case  "commentLike":
        return ' liked a comment: ';

    case "commentDislike":
        return ' disliked a comment: ';
    default:
        return '';
    }
}

function isPostActivity(activity) {
	return ['post', 'postLike', 'postDislike'].includes(activity.Type);
}

function isCommentActivity(activity) {
	return ['comment', 'commentLike', 'commentDislike'].includes(activity.Type);
}

onMounted(async () => {
    const data = await apiFetch(`/api/user/activity`);
    if (data) {
		user.value = data.User ?? null;
		activities.value = data.User?.Activities ?? [];
    }
});
</script>

<template>
	<div class="container">
		<div class="user-activity">
			<h2>Recent Activities</h2>
			<div
                v-for="(activity, index) in activities"
                :key="index"
				class="activity"
            >
                {{ DateToLocale(activity.Timestamp) }}
				<strong>
                    <router-link :to="`/user/${user?.Id}`">{{ user?.Username }}</router-link>
				</strong>
                {{ describeActivity(activity) }}
                <ActivityPost
                    v-if="isPostActivity(activity)"
                    :post="activity.Post"
                />
                <ActivityComment
                    v-else-if="isCommentActivity(activity)"
                    :post="activity.Post"
                    :comment="activity.Comment"
                />
			</div>
		</div>
	</div>
</template>

<style scoped></style>

