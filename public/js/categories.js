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