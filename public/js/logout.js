import { ShowMessage } from './message.js'
import { ShowError } from './error.js'

export async function userLogout() {
    const response = await fetch('/api/user/logout')
    const data = await response.json()
    return `
    <div class="container">
    ${(data.Error.Has) ? ShowError(data) : ''}
    ${(data.Message.Has) ? ShowMessage(data) : ''}
    <p>See ya later!</p>
</div>
`
}