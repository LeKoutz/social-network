import { ShowError } from '../components/error.js';
import {SetAlertsInner} from '../partials/alerts.js';
import { apiFetch } from '../fetchers/api.js';

export function showPostCreateView(data) {
    return `
    <div class="container">
    <form method="POST" id="post-create" action="/api/post/create" enctype="multipart/form-data">
        <fieldset>
            <legend>New post</legend>
            <input type="text" name="title" placeholder="Your title here" value="" required/>
                ${data.Categories.map(category => `
                <div class="inline">
                    <input type="checkbox" id="category-${category.Id}" name="category-${category.Id}" ${category.Selected ? 'checked' : ''}/>
                    <label for="category-${category.Id}">${category.Name}</label>
                </div>
                `).join('')}
            <textarea name="body" placeholder="Your post here" required></textarea>
            <input type="file" name="image" accept="image/jpeg,image/png,image/gif"/>
            <input type="submit" value="Create post"/>
        </fieldset>
    </form>
</div>
`;
}

export function attachPostCreateListener() {
    const form = document.querySelector('#post-create');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const url = '/api/post/create';
            // TODO: This pattern is repeated maybe it could be a function requestPOST(url) returns data? But it would require an option to use URLSearchParams if the fetch uses ParseForm or not if it uses ParseMultipartForm
            const response = await fetch(url, {
                method: 'POST',
                body: new FormData(e.target)
            });
            const data = await response.json();
            if (data.Error && data.Error.Has) {
                SetAlertsInner(ShowError(data));
                return;
            }
            window.location.hash = `/post/view/${data.Posts[0].Id}`;
        });
    }
}

export async function postCreateRoute() {
    const data = await apiFetch(`/api/post/create`);
    if (data) {
        document.querySelector('.content').innerHTML = showPostCreateView(data);
        attachPostCreateListener();
    }
}
