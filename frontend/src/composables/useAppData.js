import { ref } from 'vue';

const websiteName = ref('');
const version = ref('');

export function useAppData() {
    function setAppData(data) {
        websiteName.value = data.WebsiteName;
        version.value = data.Version;
    }

    return { websiteName, version, setAppData };
}
