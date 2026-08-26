export function ShowMessage(data, message = {Type: '', Content: ''}) {
    const type = message.Content ? message.Type : data.Message.Type;
    const content = message.Content ? message.Content : data.Message.Content;
    return `
<fieldset class="error">
    <legend>${type}</legend>
    <p>${content}</p>
</fieldset>
`;
}
