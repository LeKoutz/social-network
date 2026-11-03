<script setup>
import { ref, useTemplateRef, watch, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';

const props = defineProps({
    variant: {
        type: String,
        default: 'input',
        validator: (v) => ['elevated', 'input'].includes(v),
    },
    openOn: {
        type: String,
        default: 'click',
        validator: (v) => ['click', 'hover'].includes(v),
    },
    align: {
        type: String,
        default: 'left',
        validator: (v) => ['left', 'right'].includes(v),
    },
});

const HOVER_CLOSE_DELAY = 150;

const open = ref(false);
const root = useTemplateRef('root');
const route = useRoute();
let closeTimer = null;
let openedByHover = false;

function close() {
    clearTimeout(closeTimer);
    open.value = false;
    openedByHover = false;
}

function onTriggerClick() {
    if (openedByHover && open.value) return;
    open.value = !open.value;
}

function onPointerEnter(e) {
    if (props.openOn !== 'hover' || e.pointerType !== 'mouse') return;
    clearTimeout(closeTimer);
    open.value = true;
    openedByHover = true;
}

function onPointerLeave(e) {
    if (props.openOn !== 'hover' || e.pointerType !== 'mouse') return;
    closeTimer = setTimeout(close, HOVER_CLOSE_DELAY);
}

function onDocumentClick(e) {
    if (open.value && !root.value.contains(e.target)) close();
}

function onKeydown(e) {
    if (e.key === 'Escape') close();
}

onMounted(() => {
    document.addEventListener('click', onDocumentClick);
    document.addEventListener('keydown', onKeydown);
});

onUnmounted(() => {
    clearTimeout(closeTimer);
    document.removeEventListener('click', onDocumentClick);
    document.removeEventListener('keydown', onKeydown);
});

watch(() => route.fullPath, close);
</script>

<template>
    <div
        ref="root"
        class="position-relative"
        @pointerenter="onPointerEnter"
        @pointerleave="onPointerLeave"
    >
        <div @click="onTriggerClick">
            <slot name="trigger" :open="open"></slot>
        </div>
        <div
            v-if="open"
            class="position-absolute top-100 z-index-20 border-radius-2px min-width-100 flex flex-column padding-5px line-height-1-2"
            :class="[variant, align === 'right' ? 'right-0px' : 'left-0px']"
        >
            <slot :close="close"></slot>
        </div>
    </div>
</template>

<style scoped>
.elevated {
    background-color: var(--bg-elevated);
    border: 1px solid var(--accent-line);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.input {
    background-color: var(--bg-input);
    border: 1px solid var(--border);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}
</style>
