import { ShowMessage } from '../message.js';
import { ShowError } from '../error.js';

export function showUserRegister() {
    return `
    <div class="box">
        <form id="register" method="POST" action="/api/user/register">
            <fieldset>
                <legend>User Registration</legend>
                <input
                        type="hidden"
                        name="action"
                        value="register"/>
                <input
                        placeholder="Email"
                        type="email"
                        name="email"
                        required/>
                <input
                        placeholder="Nickname"
                        type="text"
                        name="username"
                        required/>
                <input
                        placeholder="First Name"
                        type="first_name"
                        name="first_name"
                        required/>
                <input
                        placeholder="Last Name"
                        type="last_name"
                        name="last_name"
                        required/>
                <input
                        placeholder="Age"
                        type="age"
                        name="age"
                        required/>
                <select name="gender" required>
                        <option value="">Gender</option>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                        <option value="other">Other</option>
                </select>
                <input
                        placeholder="Password"
                        type="password"
                        name="password1"
                        required/>
                <input
                        placeholder="Confirm password"
                        type="password"
                        name="password2"
                        required/>
                <input
                        type="submit"
                        value="Register"
                        name="submit"/>
            </fieldset>
        </form>
        <div class="oauth-buttons">
            <a href="/auth/google">Continue with Google</a>
            <a href="/auth/github">Continue with GitHub</a>
        </div>
        <p>Already have an account? <a href="#/user/login">Login!</a></p>
    </div>
`;
}

export function attachRegisterListener() {
        const form = document.querySelector('#register');
        if (form) {
                form.addEventListener('submit', async (e) => {
                        e.preventDefault();
                        const response = await fetch('/api/user/register', {
                                method: 'POST',
                                body: new URLSearchParams(new FormData(e.target))
                        });
                        const data = await response.json();
                        document.querySelector('.alerts').innerHTML = data.Error.Has ? ShowError(data) : ShowMessage(data);
                });
        }
}