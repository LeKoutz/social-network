import { TopBar } from './partials/topbar.js';
import { renderFooter } from './partials/footer.js';
import { routeSelect } from './router.js';
import { showWelcome } from './components/welcome.js';

const response = await fetch('/api/');
const data = await response.json();
const body = document.querySelector('body');
const topbar = document.createElement('div');
topbar.classList.add('topbar');
topbar.innerHTML = TopBar(data);
const alerts = document.createElement('div');
alerts.classList.add('alerts');
const content = document.createElement('div');
content.innerHTML = showWelcome(data);
content.classList.add('content');
const footer = document.createElement('div');
footer.classList.add('footer');
footer.innerHTML = renderFooter(data);
body.append(topbar, alerts, content, footer);

routeSelect(data, content);

window.addEventListener('hashchange', () => {
    alerts.innerHTML = '';
    routeSelect(data, content);
});
