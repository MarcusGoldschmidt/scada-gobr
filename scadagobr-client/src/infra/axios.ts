import {__API_URL__} from './env'
import Axios, {AxiosRequestConfig} from 'axios'
import {PathsV1} from "./endpoints";
import {JwtToken} from "./response_types";
import {useUserStore} from "../core/stores/userStore";

const config: AxiosRequestConfig = {
    baseURL: __API_URL__,
    headers: {
        'Content-Type': 'application/json'
    },
}

export const axios = Axios.create(config)

axios.interceptors.request.use(async (config) => {
    if (
        !config.url?.endsWith(PathsV1.Login) &&
        !config.url?.endsWith(PathsV1.RefreshToken)
    ) {
        const staticAuthStore = useUserStore.getState().user
        let userToken = staticAuthStore.token;

        const userTokenExpiration = new Date(+staticAuthStore.tokenExpiration * 1000);
        const today = new Date();

        if (today > userTokenExpiration) {
            const refreshToken = staticAuthStore.refreshToken;

            try {
                const {data} = await axios.post<JwtToken>(PathsV1.RefreshToken, {
                    refreshToken
                })

                useUserStore.getState().updateJwt(data)
                userToken = data.token
            }catch (e) {
                useUserStore.getState().unSetUser()
            }
        }

        // @ts-ignore
        config.headers.Authorization = `Bearer ${userToken}`;
    }

    return config;
});