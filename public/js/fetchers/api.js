import { ShowError } from "../components/error.js";

export async function apiFetch(url) {
    const response = await fetch(url);
    const data = await response.json();
    if (data.Error?.Has) {
        document.querySelector('.alerts').innerHTML = ShowError(data);
        return null;
    }
    return data;
}