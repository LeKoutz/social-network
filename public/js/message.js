export function ShowMessage(data) {
	return `
<fieldset class="error">
    <legend>${data.Message.Type}</legend>
    <p>${data.Message.Content}</p>
</fieldset>
`
}