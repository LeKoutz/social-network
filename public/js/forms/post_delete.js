import { ShowError } from "../components/error.js";

export function attachPostDeleteListener() {
    const form = document.querySelector('#post-delete');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const url = `/api/post/delete`;
            const response = await fetch(url, {
                method: 'POST',
                body: new URLSearchParams(new FormData(e.target))
            });
            const data = await response.json();
            if (data.Error && data.Error.Has) {
                document.querySelector('.alerts').innerHTML = ShowError(data);
                return;
            }
            window.location.hash = `#/user/posts`;
        });
    }
}