export enum PathsV1 {
    Login = 'api/v1/auth/login',
    RefreshToken = 'api/v1/auth/refresh-token',
    WhoAmI = 'api/v1/auth/who-am-i',
}

export const Request = {
    V1: PathsV1,
} as const

