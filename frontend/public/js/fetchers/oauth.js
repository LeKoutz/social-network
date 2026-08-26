import { ShowError } from '../components/error.js';
import { renderTopBar } from '../partials/topbar.js';
import { initUserFeatures } from '../app.js';

function handleOAuthPopup(url) {
    const popup = window.open(url);
    const interval = setInterval(() => {
        try {
            if (popup.closed) {
                clearInterval(interval);
                return;
            }
            if (popup.location.origin !== window.location.origin) return;
            if (popup.location.pathname !== '/auth/github/callback' &&
                popup.location.pathname !== '/auth/google/callback') return;
            const data = JSON.parse(popup.document.body.innerText);
            clearInterval(interval);
            popup.close();
            if (data.Error?.Has) {
                document.querySelector('.alerts').innerHTML = ShowError(data);
                return;
            }
            renderTopBar(data);
            initUserFeatures(data);
            window.location.hash = '';
        } catch (e) {
            console.log(e);
        }
    }, 100);
}

export function attachGoogleAuthListener() {
    document.querySelector('#google-auth')?.addEventListener('click', (e) => {
        e.preventDefault();
        handleOAuthPopup('/auth/google');
    });
}

export function attachGithubAuthListener() {
    document.querySelector('#github-auth')?.addEventListener('click', (e) => {
        e.preventDefault();
        handleOAuthPopup('/auth/github');
    });
}
