export function ShowError(data) {
    return `
<fieldset class="error">
    <legend>Error</legend>
    <p>${data.Error.Message}</p>
</fieldset>
`
}

