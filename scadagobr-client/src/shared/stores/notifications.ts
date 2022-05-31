import {writable} from 'svelte/store';

export enum NotificationType {
    Danger= "danger",
    Success= "success",
    Warning = "warning",
    Info = "info",
}

export interface Notification {
    id: string
    title: string
    body: string
    type: NotificationType,
    timeout: number
}

export const notificationStore = writable<Notification>();

export function sendNotification(
    title: string,
    body: string = "",
    type: NotificationType = NotificationType.Info,
    timeout: number = 3000) {
    notificationStore.set({
        id: Math.random().toString(36).replace(/[^a-z]+/g, ''),
        title,
        body,
        type,
        timeout
    })
}