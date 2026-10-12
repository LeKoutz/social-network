import { useAlerts } from '@/composables/useAlerts';

export async function apiGet(url) {
    const response = await fetch(url);
    const data = await response.json();
    if (data.Error?.Has) {
        const { setAlert } = useAlerts();
        setAlert(data);
        // return null;
    }
    return data;
}

export async function apiPost(url, body) {
    const response = await fetch(url, {
        method: 'POST',
        body,
    });
    const data = await response.json();
    if (data.Error?.Has) {
        const { setAlert } = useAlerts();
        setAlert(data);
        // return null;
    }
    return data;
}

export async function apiPatch(url, body) {
    const response = await fetch(url, {
        method: 'PATCH',
        body,
    });
    const data = await response.json();
    if (data.Error?.Has) {
        const { setAlert } = useAlerts();
        setAlert(data);
        // return null;
    }
    return data;
}

export async function apiDelete(url, body) {
    const response = await fetch(url, {
        method: 'DELETE',
        body,
    });
    const data = await response.json();
    if (data.Error?.Has) {
        const { setAlert } = useAlerts();
        setAlert(data);
        // return null;
    }
    return data;
}
