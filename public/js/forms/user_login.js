export function showUserLogin() {
    return `
<div class="container">
    <div class="box">
        <form method="POST" action="/api/user/login">
            <fieldset>
                <legend>User Login</legend>
                <input
                        type="hidden"
                        name="action"
                        value="login"/>
                <input
                        type="email"
                        placeholder="Email"
                        name="email"
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
</div>
`
}
