import { ShowMessage } from './message.js';
import { ShowError } from './error.js';

export async function userLogout() {
    const response = await fetch('/api/user/logout');
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        document.querySelector('.alerts').innerHTML = ShowError(data);
        return '';
    }
    return `
    <div class="container">
    ${data.Message.Has ? (() => { setTimeout(() => window.location.hash = '', 1000); return ShowMessage(data); })() : ''}
    <p>See ya later!</p>
</div>
`;
}