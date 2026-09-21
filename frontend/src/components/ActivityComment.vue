<script setup>
defineProps({
    post: {
        type: Object,
        required: true,
    },
    comment: {
        type: Object,
        required: true,
    },
});
</script>

<template>
    <div class="container">
        <router-link :to="`/post/view/${post.Id}`">
                <h3>{{ post.Title }}</h3>
        </router-link>
        <div class="activity-highlight">
            <div class="comment" id="comment-{{comment.Id}}">
                <router-link :to="`/post/view/{{post.Id}}#comment-{{comment.Id}}`">
                    <pre>{{comment.Body}}</pre>
                </router-link>
                <img v-if="comment.ImagePath" :src="`/${comment.ImagePath}`" alt="Comment image" style="max-width: 100%"/>
                <div class="reactions">
                    <button disabled>{{comment.Likes}} {{comment.Liked?"👍":"👍🏻"}}</button>
                    <button disabled>{{comment.Dislikes}} {{comment.Disliked?"👎":"👎🏻"}}</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped></style>
