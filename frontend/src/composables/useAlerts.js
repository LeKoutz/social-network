import { ref } from 'vue';

const alert = ref(null)

export function useAlerts() {
    function createAlert(data) {
        if (data.Error.Has) {
            alert.value = { type: 'Error', message: data.Error.Message }
        } else if (data.Message.Has) {
            alert.value = { type: data.Message.Type, message: data.Message.Content }
        }
    }

    function clear() {
        alert.value = null
    }
    
    return { alert, createAlert, clear }
}