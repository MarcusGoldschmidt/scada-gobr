export enum PathsV1 {
    Login = 'api/v1/auth/login',
    RefreshToken = 'api/v1/auth/refresh-token',
    WhoAmI = 'api/v1/auth/who-am-i',
    UserGet = 'api/v1/user',
    UserCreate = 'api/v1/user',
    UserDelete = 'api/v1/user/',
    UserUpdate = 'api/v1/user/',
    UserGetById = 'api/v1/user/',
}

export const Request = {
    V1: PathsV1,
} as const

