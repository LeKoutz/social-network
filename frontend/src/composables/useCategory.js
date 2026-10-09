import { ref } from 'vue';

const category = ref([]);

export function useCategory() {
    function setCategory(newCategory) {
        category.value = newCategory ?? [];
    }

    return { category, setCategory };
}
