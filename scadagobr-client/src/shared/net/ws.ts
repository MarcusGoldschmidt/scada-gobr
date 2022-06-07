import env from "../env";
import {get} from "svelte/store";
import {authStore} from "../stores/user";


export function createWs(path: string) : WebSocket {
    return new WebSocket(env.WS_BASE_URL + path);
}