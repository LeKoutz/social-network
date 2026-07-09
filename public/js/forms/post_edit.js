import { ShowError } from '../components/error.js';
import { apiFetch } from '../fetchers/api.js';

export function showPostEditView(data) {
    const post = data.Posts[0];
    return `
    <div class="container">
    <form method="POST" id="post-edit" action="/api/post/edit/${post.Id}" enctype="multipart/form-data">
        <fieldset>
            <legend>Edit Post</legend>
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <input type="hidden" name="save-post" value="1"/>
            <input type="text" name="title" placeholder="Your title here" value="${post.Title}" required/>
                ${data.Categories.map(category => `
                <div class="inline">
                    <input type="checkbox" id="category-${category.Id}" name="category-${category.Id}" ${category.Selected ? 'checked' : ''}/>
                    <label for="category-${category.Id}">${category.Name}</label>
                </div>
                `).join('')}
            <textarea name="body" placeholder="Your post here" required>${post.Body}</textarea>
            ${post.ImagePath ? `<img src="/${post.ImagePath}" alt="Post image" style="max-width: 100%;"/>` : ''}
            <input type="file" name="image" accept="image/jpeg,image/png,image/gif"/>
            <input type="submit" value="Save post"/>
        </fieldset>
    </form>
</div>
`;
}

export function attachPostEditListener(id) {
    const form = document.querySelector('#post-edit');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const url = `/api/post/edit/${id}`;
            const response = await fetch(url, {
                method: 'POST',
                body: new FormData(e.target)
            });
            const data = await response.json();
            if (data.Error && data.Error.Has) {
                document.querySelector('.alerts').innerHTML = ShowError(data);
                return;
            }
            window.location.hash = `/post/view/${id}`;
        });
    }
}

export async function postEditRoute(id) {
    const data = await apiFetch(`/api/post/edit/${id}`);
    if (data) {
        document.querySelector('.content').innerHTML = showPostEditView(data);
        attachPostEditListener(id);
    }
}