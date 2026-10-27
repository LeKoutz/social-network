import { ref } from 'vue';

const websiteName = ref('');
const version = ref('');
const mobileMediaQuery = window.matchMedia('(max-width: 768px)');
const isMobileView = ref(mobileMediaQuery.matches);

mobileMediaQuery.addEventListener('change', (e) => {
    isMobileView.value = e.matches;
});

export function useAppData() {
    function setAppData(data) {
        websiteName.value = data.WebsiteName;
        version.value = data.Version;
    }

    return { websiteName, version, isMobileView, setAppData };
}
