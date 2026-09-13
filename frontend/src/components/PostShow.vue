<script setup>
import { useRouter } from 'vue-router';
import { useUser } from '@/composables/useUser.js';
import { DateToLocale } from '@/utils/utils.js';

const props = defineProps({
    post: { type: Object, default: null },
});

const { user } = useUser();
const router = useRouter();

const isOwner = () =>
    user.value.LoggedIn &&
    String(user.value.Id) === String(props.post.User?.Id);

</script>
<template>
    <div class="container">
        <div v-if="post" class="post" :id="`${post.Id}`">
            <div class="post-content">
                <router-link :to="`/post/view/${post.Id}`"><h3>{{post.Title}}</h3></router-link>
                <div class="post-details">
                    <p>Categories:
                    <div v-for="category in post.Categories" :key="category.Id">
                        <router-link :to="`/category/view/${category.Id}`">
                            {{ category.Name }}
                        </router-link>
                    </div>
                    </p>
                    <p>Posted by
                        <router-link :to="`/user/${post.User.Id}`">
                            <strong>{{post.User.Username}}</strong>
                        </router-link>
                        on
                        <em>({{DateToLocale(post.Timestamp)}})</em>
                    </p>
                </div>
                <div v-if="isOwner()" class="manage-post">
                    <!-- TODO:
                        extract both to actions
                        ${postDeleteForm(post.Id)}`:''} -->
                    <router-link :to="`/post/edit/${post.Id}`">
                        <button type="button">Edit Post</button>
                    </router-link>
                    <router-link :to="`/post/delete/${post.Id}`">
                        <button type="button">Delete Post</button>
                    </router-link>
                </div>
                <pre>{{post.Body}}</pre>
                <img v-if="post.ImagePath" :src="`/${post.ImagePath}`" alt="Post image" style="max-width: 100%;"/>
                <!-- TODO: Add the following
                <div class="reactions">
                    ${postReactionForm(post,data.User.LoggedIn)}
                </div>
                <div class="comments">
                    ${data.User.LoggedIn ? showCommentCreate(post) : ''}
                    ${post.Comments ? showPostComments(data) : ''}
                </div>
                -->
                <div v-for="(comment,index) in post.Comments">
                    {{ comment.PostId }}
                    {{ comment.UserId }}
                    {{ comment.Username }}
                    {{ comment.Body }}
                    {{ comment.Timestamp }}
                    {{ comment.Likes }}
                    {{ comment.Liked }}
                    {{ comment.Dislikes }}
                    {{ comment.Disliked }}
                </div>
            </div>
        </div>
    </div>
</template>

