import { ref } from 'vue';

const categories = ref([]);

export function useCategories() {
    function setCategories(newCategories) {
        categories.value = newCategories ?? [];
    }

    return { categories, setCategories };
}
