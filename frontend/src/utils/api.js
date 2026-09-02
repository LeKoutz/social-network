import { useAlerts } from "@/composables/useAlerts";

export async function apiFetch(url) {
    const response = await fetch(url);
    const data = await response.json();
    if (data.Error?.Has) {
        const { setAlert } = useAlerts()
        setAlert(data)
        return null;
    }
    return data;
}
