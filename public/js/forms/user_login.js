import { ShowMessage } from '../components/message.js';
import { ShowError } from '../components/error.js';

export async function showUserLogin() {
    const response = await fetch ('/api/user/login');
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        document.querySelector('.alerts').innerHTML = ShowError(data);
        return '';
    }
    return `
<div class="box">
<form method="POST" id="login" action="/api/user/login">
        <fieldset>
        <legend>User Login</legend>
        <input
                type="hidden"
                name="action"
                value="login"/>
        <input
                type="identifier"
                placeholder="Email or Nickname"
                name="identifier"
                required
                />
        <input
                type="password"
                placeholder="Password"
                name="password"
                required
                />
        <input
                type="submit"
                value="Login"
                name="submit"/>
        </fieldset>
</form>
<div class="oauth-buttons">
        <a href="/auth/google">Continue with Google</a>
        <a href="/auth/github">Continue with GitHub</a>
</div>
<p>You don't have an account? <a href="#/user/register">Register!</a></p>
</div>
`;
}

export function attachLoginListener() {
    const form = document.querySelector('#login');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const response = await fetch('/api/user/login', {
                method: 'POST',
                body: new URLSearchParams(new FormData(e.target))
            });
            const data = await response.json();
            document.querySelector('.alerts').innerHTML = data.Error.Has ? ShowError(data) : ShowMessage(data);
            data.User.LoggedIn ? setTimeout(() => window.location.hash = '', 1000) : '';
        });
    }
}
