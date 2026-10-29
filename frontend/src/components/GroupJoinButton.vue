<script setup>
import GenericButton from '@/components/GenericButton.vue';
import PlusIcon from '@/components/Icons/PlusIcon.vue';
</script>

<template>
    <GenericButton>
        <template v-slot:icon>
            <PlusIcon />
        </template>
        Join group
    </GenericButton>
</template>
