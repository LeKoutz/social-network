<script setup>
const emit = defineEmits(['edit']);
</script>

<template>
    <button type="button" class="comment-edit" @click="emit('edit')">Edit Comment</button>
</template>

<style scoped></style>
