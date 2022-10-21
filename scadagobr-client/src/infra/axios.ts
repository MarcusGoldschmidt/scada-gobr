import {__API_URL__} from './env'
import Axios, {AxiosRequestConfig} from 'axios'
import {PathsV1} from "./endpoints";
import {JwtToken} from "./response_types";
import {userStore} from "../core/stores/userStore";

export const axios = Axios.create({
    baseURL: __API_URL__,
    headers: {
        'Content-Type': 'application/json;charset=utf-8',
        'Access-Control-Allow-Origin': '*'
    }
})

const config: AxiosRequestConfig = {
    baseURL: __API_URL__,
    headers: {
        'Content-Type': 'application/json'
    },
}

const axiosJwt = Axios.create(config)

axiosJwt.interceptors.request.use(async (config) => {
    if (
        !config.url?.endsWith(PathsV1.Login) &&
        !config.url?.endsWith(PathsV1.RefreshToken)
    ) {
        const staticAuthStore = userStore.getState().user

        const userTokenExpiration = new Date(+staticAuthStore.tokenExpiration * 1000);
        const today = new Date();

        if (today > userTokenExpiration) {

            const refreshToken = staticAuthStore.refreshToken;

            const {data} = await axiosJwt.post<JwtToken>(PathsV1.RefreshToken, {
                refreshToken
            })

            userStore.getState().updateJwt(data)
        }

        const userToken = staticAuthStore.token;
        // @ts-ignore
        config.headers.Authorization = `Bearer ${userToken}`;
    }

    return config;
});

export default axiosJwt;