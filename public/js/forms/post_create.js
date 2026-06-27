import { showPostView } from '../posts.js'
import { ShowError } from '../error.js'

export function showPostCreateView(data) {
    return `
    <div class="container">
    <form method="POST" id="post-create" action="/api/post/create" enctype="multipart/form-data">
        <fieldset>
            <legend>New post</legend>
            <input type="text" name="title" placeholder="Your title here" required/>
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
`
}

export function attachPostCreateListener() {
        const form = document.querySelector('#post-create')
        const container = document.querySelector('.container');
        if (form) {
                form.addEventListener('submit', async (e) => {
                        e.preventDefault()
                        const response = await fetch('/api/post/create', { // TODO: This pattern is repeated maybe it could be a function requestPOST(url) returns data? But it would require an option to use URLSearchParams if the fetch uses ParseForm or not if it uses ParseMultipartForm
                                method: 'POST',
                                body: new FormData(e.target)
                        })
                        console.log(response.status)
                        const data = await response.json()
                        console.log(data)
                        if (data.Error && data.Error.Has) {
                            container.innerHTML = ShowError(data)
                            return
                        }
                        container.innerHTML = await showPostView(data.postId)
                })
        }
}