import { TopBar } from './partials/topbar.js';
import { renderFooter } from './partials/footer.js';
import { routeSelect } from './router.js';

const response = await fetch('/api/');
const data = await response.json();
const body = document.querySelector('body');
const topbar = document.createElement('div');
topbar.classList.add('topbar');
topbar.innerHTML = TopBar(data);
const alerts = document.createElement('div');
alerts.classList.add('alerts');
const container = document.createElement('div');
container.classList.add('container');
const content = document.createElement('div');
content.classList.add('content');
const usersPanel = document.createElement('div');
usersPanel.classList.add('users-panel');
container.append(content, usersPanel);
const footer = document.createElement('div');
footer.classList.add('footer');
footer.innerHTML = renderFooter(data);
body.append(topbar, alerts, container, footer);

routeSelect();

window.addEventListener('hashchange', () => {
    alerts.innerHTML = '';
    routeSelect();
});
