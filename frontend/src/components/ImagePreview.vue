<script setup>
import { ref, watch, onUnmounted } from 'vue';
import GenericButton from './GenericButton.vue';
import CloseIcon from './Icons/CloseIcon.vue';

const image = defineModel();
const previewUrl = ref(null);

watch(
    image,
    (file) => {
        if (previewUrl.value) URL.revokeObjectURL(previewUrl.value);
        previewUrl.value = file ? URL.createObjectURL(file) : null;
    },
    { immediate: true },
);

onUnmounted(() => {
    if (previewUrl.value) URL.revokeObjectURL(previewUrl.value);
});
</script>

<template>
    <div v-if="previewUrl" class="image-preview-wrapper">
        <img :src="previewUrl" alt="" class="image-preview" />
        <GenericButton type="button" class="clear-button" @click="image = null">
            <template #leading><CloseIcon /></template>
        </GenericButton>
        <small class="text-color-muted">{{ image.name }}</small>
    </div>
</template>

<style scoped>
.image-preview-wrapper {
    position: relative;
    display: flex;
    flex-direction: column;
    width: fit-content;
}

.image-preview {
    max-width: 100%;
    max-height: 200px;
    border-radius: 2px;
}

.clear-button {
    position: absolute;
    top: 4px;
    right: 4px;
}
</style>
