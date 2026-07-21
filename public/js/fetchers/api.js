import { ShowError } from "../components/error.js";
import {SetAlertsInner} from '../partials/alerts.js';

export async function apiFetch(url) {
    const response = await fetch(url);
    const data = await response.json();
    if (data.Error?.Has) {
        SetAlertsInner(ShowError(data));
        return null;
    }
    return data;
}
