export function later(delay: number) {
    return new Promise(function(resolve) {
        setTimeout(resolve, delay);
    });
}

export default {
    later
}