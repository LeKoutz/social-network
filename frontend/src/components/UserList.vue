<script setup>
import { ref } from 'vue';

const props = defineProps({
    group: { type: Object, default: null },
});

const group = ref(props.group);

</script>

<template>
</template>

<style scoped></style>
