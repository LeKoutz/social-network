import { ref, watch } from 'vue';
import { useAppData } from '@/composables/useAppData.js';

const { isMobileView } = useAppData();
const sidebarOpen = ref(!isMobileView.value);

watch(isMobileView, (mobile) => {
    sidebarOpen.value = !mobile;
});

export function useLeftSidebar() {
    function toggleLeftSidebar() {
        sidebarOpen.value = !sidebarOpen.value;
    }

    return { sidebarOpen, toggleLeftSidebar };
}
