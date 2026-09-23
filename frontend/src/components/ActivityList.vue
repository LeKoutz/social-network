<script setup>
import { DateToLocale } from '@/utils/utils.js';
import ActivityPost from '@/components/ActivityPost.vue';
import ActivityComment from '@/components/ActivityComment.vue';

defineProps({
	activities: {
		type: Array,
		default: () => [],
	},
});

function describeActivity(activity) {
	switch (activity.Type) {
		case 'post': return ' posted: ';
		case 'comment': return ' commented: ';
		case 'postLike': return ' liked a post: ';
		case 'postDislike': return ' disliked a post: ';
		case 'commentLike': return ' liked a comment: ';
		case 'commentDislike': return ' disliked a comment: ';
		default: return '';
	}
}
function isPostActivity(activity) {
	return ['post', 'postLike', 'postDislike'].includes(activity.Type);
}
function isCommentActivity(activity) {
	return ['comment', 'commentLike', 'commentDislike'].includes(activity.Type);
}
</script>

<template>
	<div v-if="activities && activities.length">
        <h2>Recent Activity</h2>
		<div v-for="(activity, index) in activities" :key="index" class="activity">
			{{ DateToLocale(activity.Timestamp) }}
			{{ describeActivity(activity) }}
			<ActivityPost v-if="isPostActivity(activity)" :post="activity.Post" />
			<ActivityComment
				v-else-if="isCommentActivity(activity)"
				:post="activity.Post"
				:comment="activity.Comment"
			/>
		</div>
	</div>
	<p v-else>No activity yet.</p>
</template>

<style scoped></style>