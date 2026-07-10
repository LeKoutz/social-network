import { routeSelect } from './router.js';

const body = document.querySelector('body');
const topbar = document.createElement('div');
topbar.classList.add('topbar');
const alerts = document.createElement('div');
alerts.classList.add('alerts');
const content = document.createElement('div');
content.classList.add('content');
const footer = document.createElement('div');
footer.classList.add('footer');
body.append(topbar, alerts, content, footer);

routeSelect();

window.addEventListener('hashchange', () => {
    alerts.innerHTML = '';
    routeSelect();
});
