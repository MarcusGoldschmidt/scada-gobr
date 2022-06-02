import type {AxiosRequestConfig} from "axios";
import * as axios from 'axios'
import env from "../env";
import {authStore} from "../stores/user";
import type {JwtToken} from "../types";
import {get} from "svelte/store";
import {PathsV1} from "./request";
import {NotificationType, sendNotification} from "../stores/notifications";

const config: AxiosRequestConfig = {
    baseURL: env.API_BASE_URL,
    headers: {
        'Content-Type': 'application/json'
    },
}

export const axiosJwt = axios.default.create(config)

axiosJwt.interceptors.request.use(async (config) => {
    if (
        !config.url.endsWith(PathsV1.Login) &&
        !config.url.endsWith(PathsV1.RefreshToken)
    ) {
        const staticAuthStore = get(authStore)

        const userTokenExpiration = new Date(+staticAuthStore.jwt.tokenExpiration * 1000);
        const today = new Date();

        if (today > userTokenExpiration) {

            const refreshToken = staticAuthStore.jwt.refreshToken;

            const {data} = await axiosJwt.post<JwtToken>(PathsV1.RefreshToken, {
                refreshToken
            })

            authStore.set({
                ...staticAuthStore,
                jwt: data
            })
        }

        const userToken = get(authStore).jwt.token;
        config.headers.Authorization = `Bearer ${userToken}`;
    }

    return config;
});

export async function loginUser(name, password) {

    try {
        const {data: jwt} = await axiosJwt.post(PathsV1.Login, {name, password})

        authStore.set({
            user: undefined,
            jwt
        })

        const {data: user} = await axiosJwt.get(PathsV1.WhoAmI)

        authStore.set({
            user,
            jwt
        })

        sendNotification("Successful login", "", NotificationType.Info)

        return true;
    } catch (e) {
        sendNotification("Wrong username or password", "", NotificationType.Danger)
        return false;
    }
}