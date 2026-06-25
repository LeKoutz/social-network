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
            <h2>${category.Name}</h2>
            <p>${category.Description}</p>
        </div>
    `).join('')
}