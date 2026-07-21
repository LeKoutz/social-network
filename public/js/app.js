import { renderTopBar } from './partials/topbar.js';
import { renderFooter } from './partials/footer.js';
import { routeSelect } from './router.js';
import { connectWS } from './ws.js';
import { UsersPanel, addUsersPanelButtonListener } from './components/users_panel.js';
import { apiFetch } from './fetchers/api.js';
import {SetAlertsInner} from './partials/alerts.js';

function buildDOM() {
    const body = document.querySelector('body');
    const topbar = document.createElement('div');
    topbar.classList.add('topbar');
    const alerts = document.createElement('div');
    alerts.classList.add('alerts');
    const content = document.createElement('div');
    content.classList.add('content');
    const usersPanel = document.createElement('div');
    usersPanel.classList.add('users-panel');
    const container = document.createElement('div');
    container.classList.add('container');
    container.append(content, usersPanel);
    const footer = document.createElement('div');
    footer.classList.add('footer');
    body.append(topbar, alerts, container, footer);
}

export async function initUserFeatures(data) {
    connectWS(data.User.Id);
    const usersData = await apiFetch('/api/users');
    if (usersData) {
        document.querySelector('.users-panel').innerHTML = UsersPanel(usersData);
        addUsersPanelButtonListener();
    }
}

async function init() {
    const data = await apiFetch('/api/');
    if (!data) return;
    buildDOM();
    renderTopBar(data);
    document.querySelector('.footer').innerHTML = renderFooter(data);
    if (data.User.LoggedIn) await initUserFeatures(data);
    window.addEventListener('hashchange', () => {
        SetAlertsInner('');
        document.querySelector('.alerts').innerHTML = '';
        routeSelect();
    });
    routeSelect();
}

init();
