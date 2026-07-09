import { ShowMessage } from './message.js';
import { displayPosts } from './posts.js';
import { apiFetch } from '../fetchers/api.js';
import { attachCommentCreateListener } from '../forms/comment_create.js';

export function Categories(data) {
    return `
    <div class="categories">
        <h2>Categories</h2>
        ${createCategories(data.Categories)}
    </div>
    `;
}

function createCategories(categories) {
    return categories.map(category => `
        <div class="category" id="${category.Id}">
            <a href="#/category/view/${category.Id}"><h2>${category.Name}</h2></a>
            <p class="category-description">${category.Description}</p>
        </div>
    `).join('');
}

export function showCategoryView(data) {
    const category = createCategories(data.Categories);
    if (!data.Posts) document.querySelector('.alerts').innerHTML = ShowMessage(data, {Type: 'Oops', Content: 'This category is empty'});
    return `
    <div class="container">
        ${category}
        ${data.Posts ? displayPosts(data) : ''}
    </div>
    `;
}

export function showPostCategories(post) {
    return post.Categories.map(category => `
        <a href="#/category/view/${category.Id}">${category.Name}</a>
    `).join('');
}

export async function categoryRoute(id) {
    const data = await apiFetch(`/api/category/view/${id}`);
    if (data) {
        document.querySelector('.content').innerHTML = showCategoryView(data);
        attachCommentCreateListener();
    }
}