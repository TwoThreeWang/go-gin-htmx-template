document.addEventListener('htmx:configRequest', function(event) {
    event.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]')?.content;
});

document.addEventListener('htmx:afterSwap', function(event) {
    console.log('Content swapped:', event.detail.target);
});

document.addEventListener('htmx:responseError', function(event) {
    console.error('Error:', event.detail.xhr.responseText);
    alert('An error occurred. Please try again.');
});

function formatJSON(str) {
    try {
        const obj = JSON.parse(str);
        return JSON.stringify(obj, null, 2);
    } catch (e) {
        return str;
    }
}

document.body.addEventListener('htmx:afterSwap', function(event) {
    const apiResult = document.getElementById('api-result');
    if (apiResult && apiResult.textContent && !apiResult.textContent.startsWith('{')) {
        apiResult.textContent = formatJSON(apiResult.textContent);
    }
});
