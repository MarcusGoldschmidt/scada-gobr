export function getLocalStorage<T>(key: string): T | null {
    // @ts-ignore
    return JSON.parse(localStorage.getItem(key));
}

export const setLocalStorage = (key: string, value: object) => localStorage.setItem(key, JSON.stringify(value));

export const removeLocalStorage = (key: string) => localStorage.removeItem(key);