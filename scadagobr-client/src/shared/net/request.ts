export enum PathsV1 {
    Login = 'api/v1/auth/login',
    RefreshToken = 'api/v1/auth/refresh-token',
    WhoAmI = 'api/v1/auth/who-am-i',
    // User
    UserGet = 'api/v1/user',
    UserCreate = 'api/v1/user',
    UserDelete = 'api/v1/user/',
    UserUpdate = 'api/v1/user/',
    UserGetById = 'api/v1/user/',
    // Views
    ViewGet = 'api/v1/view',
    ViewCreate = 'api/v1/view',
    ViewDelete = 'api/v1/view/',
    ViewUpdate = 'api/v1/view/',
    ViewGetById = 'api/v1/view/',
}

export const Request = {
    V1: PathsV1,
} as const

export enum PathsWsV1 {
    DataPoint = "/api/v1/datapoints/ws/",
}

