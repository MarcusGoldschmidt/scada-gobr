import {writable} from 'svelte/store';
import type {CurrentUser} from "../types";

const AuthStoreKey = "CURRENT_USER"

const storedUser = JSON.parse(localStorage.getItem(AuthStoreKey));
export const authStore = writable<CurrentUser>(storedUser);

authStore.subscribe(value => {
    localStorage.setItem(AuthStoreKey, JSON.stringify(value));
});
