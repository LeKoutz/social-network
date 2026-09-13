export function throttle(func, wait) {
    let timer;
    let callPending = false;
    return function (...args) {
        if (timer) {
            callPending = true;
            return;
        }
        func(...args);
        timer = setTimeout(function run() {
            timer = null;
            if (callPending) {
                callPending = false;
                func(...args);
                timer = setTimeout(run, wait);
            }
        }, wait);
    };
}

export function DateToLocale(s) {
    return new Date(parseInt(s) * 1000).toLocaleString();
}

export function firstOrNull(arr) {
    return Array.isArray(arr) && arr.length ? arr[0] : null;
}
