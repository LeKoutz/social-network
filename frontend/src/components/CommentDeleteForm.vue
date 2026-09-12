<script setup>
const emit = defineEmits(['delete']);
</script>

<template>
    <button type="button" class="comment-delete" @click="emit('delete')">Delete Comment</button>
</template>

<style scoped></style>
