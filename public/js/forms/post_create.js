import { showPostView } from '../posts.js';
import { ShowError } from '../error.js';

export function showPostCreateView(data, options = { editing: false }) {
    const post = options.editing ? data.Posts[0] : null;
    return `
    <div class="container">
    <form method="POST" id="post-create" action="${options.editing ? `/api/post/edit/${post.Id}` : `/api/post/create`}" enctype="multipart/form-data">
        <fieldset>
            <legend>${options.editing ? `Edit Post` : `New post`}</legend>
            ${options.editing ? `
                <input type="hidden" name="post-id" value="${post.Id}"/>
                <input type="hidden" name="save-post" value="1"/>
                `: '' }
            <input type="text" name="title" placeholder="Your title here" value="${options.editing ? `${post.Title}` : ''}" required/>
                ${data.Categories.map(category => `
                <div class="inline">
                    <input type="checkbox" id="category-${category.Id}" name="category-${category.Id}" ${category.Selected ? 'checked' : ''}/>
                    <label for="category-${category.Id}">${category.Name}</label>
                </div>
                `).join('')}
            <textarea name="body" placeholder="Your post here" required>${options.editing ? `${post.Body}` : ''}</textarea>
            ${options.editing && post.ImagePath ? `<img src="/${post.ImagePath}" alt="Post image" style="max-width: 100%;"/>` : ''}
            <input type="file" name="image" accept="image/jpeg,image/png,image/gif"/>
            <input type="submit" value="${options.editing ? `Save post` : `Create post`}"/>
        </fieldset>
    </form>
</div>
`;
}

export function attachPostCreateListener(options = { editing: false }) {
        const form = document.querySelector('#post-create');
        const container = document.querySelector('.container');
        if (form) {
            
            form.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    const url = options.editing ? `/api/post/edit/${form.querySelector('[name="post-id"]').value}` : '/api/post/create';
                    // TODO: This pattern is repeated maybe it could be a function requestPOST(url) returns data? But it would require an option to use URLSearchParams if the fetch uses ParseForm or not if it uses ParseMultipartForm
                    const response = await fetch(url, {
                            method: 'POST',
                            body: new FormData(e.target)
                    });
                    const data = await response.json();
                    if (data.Error && data.Error.Has) {
                        container.innerHTML = ShowError(data);
                        return;
                    }
                    container.innerHTML = await showPostView(data.Posts[0].Id);
            });
        }
}

export async function showPostEditView(id) {
    const response = await fetch(`/api/post/edit?post-id=${id}`);
    const data = await response.json();
    return showPostCreateView(data, { editing: true });
}
