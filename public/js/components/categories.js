import { ShowMessage } from './message.js';
import { ShowError } from './error.js';
import { displayPosts } from './posts.js';

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

export async function showCategoryView(id) {
    const response = await fetch(`/api/category/view/${id}`);
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        document.querySelector('.alerts').innerHTML = ShowError(data);
        return '';
    }
    const category = createCategories(data.Categories);
    document.querySelector('.alerts').innerHTML = data.Posts ? '' : ShowMessage(data, {Type: 'Oops', Content: 'This category is empty'});
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
