<script setup>
import { ref, watch, nextTick } from 'vue';
import ChatMessage from '@/components/ChatMessage.vue';
import { throttle } from '@/utils/utils.js';

const props = defineProps({
    messages: { type: Array, required: true },
    currentUserId: { type: Number, required: true },
});

const emit = defineEmits(['load-older']);

const container = ref(null);

function scrollToBottom() {
    if (container.value) container.value.scrollTop = container.value.scrollHeight;
}

const handleScroll = throttle(() => {
    if (container.value && container.value.scrollTop === 0) emit('load-older');
}, 300);

watch(
    [() => props.messages[0]?.Id, () => props.messages.at(-1)?.Id],
    async ([first, last], [oldFirst, oldLast]) => {
        const firstVisible = container.value?.firstElementChild;
        await nextTick();
        if (last !== oldLast) {
            scrollToBottom();
        } else if (first !== oldFirst) {
            firstVisible?.scrollIntoView();
        }
    },
);
</script>

<template>
    <div class="chat-messages" ref="container" @scroll="handleScroll">
        <ChatMessage
            v-for="m in messages"
            :key="m.Id"
            :message="m"
            :is-own="m.SenderId === currentUserId"
        />
    </div>
</template>

<style scoped></style>