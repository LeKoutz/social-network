export async function fetchPostCreateData(options = { editing: false, postId: null }) {
    const url = options.editing ? `/api/post/edit/${options.postId}` : '/api/post/create';
    const response = await fetch(url);
    const data = await response.json();
    if (data.Error && data.Error.Has) {
        document.querySelector('.alerts').innerHTML = ShowError(data);
        return null;
    }
    return data
}
