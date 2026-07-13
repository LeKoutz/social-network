import { ShowMessage } from './message.js';
import { apiFetch } from '../fetchers/api.js';
import { disconnectWS } from '../ws.js';

export function userLogout(data) {
    return `
    <div class="container">
    ${data.Message.Has ? ShowMessage(data) : ''}
    <p>See ya later!</p>
</div>
`;
}

export async function logoutRoute() {
    const data = await apiFetch('/api/user/logout');
    if (data) {
        document.querySelector('.content').innerHTML = userLogout(data);
        disconnectWS();
        setTimeout(() => window.location.hash = '', 1000);
    }
}