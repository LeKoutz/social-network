import { ref } from 'vue';

const mobileView = window.matchMedia('(max-width: 768px)');
const sidebarOpen = ref(!mobileView.matches);

mobileView.addEventListener('change', (e) => {
    sidebarOpen.value = !e.matches;
});

export function useLeftSidebar() {
    function toggleLeftSidebar() {
        sidebarOpen.value = !sidebarOpen.value;
    }

    return { sidebarOpen, toggleLeftSidebar };
}
