import { displayPosts } from './posts.js'

export function Categories(data) {
    return `
    <div class="categories">
        <h2>Categories</h2>
        ${createCategories(data.Categories)}
    </div>
    `
}

function createCategories(categories) {
    return categories.map(category => `
        <div class="category" id="${category.Id}">
             <a href="#/category/view/${category.Id}"><h2>${category.Name}</h2></a>
            <p>${category.Description}</p>
        </div>
    `).join('')
}

export async function showCategoryView(id) {
    const response = await fetch(`/api/category/view/${id}`)
    const data = await response.json()
    return `
    <div class="container">
        ${data.Categories[0].Name}
        ${displayPosts(data)}
    </div>
    `
}
