<script setup>
import { useNotifications } from '@/composables/useNotifications';

const { notifications, markNotificationsAsRead } = useNotifications();

</script>

<template>
    <div class="notifications">
        <form id="mark-as-read" @submit.prevent="markNotificationsAsRead">
            <button type="submit">Mark all as read</button>
        </form>
        <ul>
            <li :class="{ unread: !notification.Read }" v-for="notification in notifications">
                <router-link v-if="notification.Type === 'comment'" :to="`/post/view/${notification.PostId}#comment-${notification.CommentId}`">
                    {{ notification.Actor.Username }} commented on your post
                </router-link>
                <router-link v-if="notification.Type === 'like'" :to="`/post/view/${notification.PostId}`">
                    {{ notification.Actor.Username }} liked on post
                </router-link>
                <router-link v-if="notification.Type === 'dislike'" :to="`/post/view/${notification.PostId}`">
                    {{ notification.Actor.Username }} disliked your post
                </router-link>
                <router-link v-if="notification.Type === 'commentLike'" :to="`/post/view/${notification.PostId}#comment-${notification.CommentId}`">
                    {{ notification.Actor.Username }} liked your comment
                </router-link>
                <router-link v-if="notification.Type === 'commentDislike'" :to="`/post/view/${notification.PostId}#comment-${notification.CommentId}`">
                    {{ notification.Actor.Username }} disliked your comment
                </router-link>
            </li>
        </ul>
    </div>
</template>

<style scoped>
</style>
