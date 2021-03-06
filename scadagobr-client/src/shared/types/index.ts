export interface JwtToken {
    token: string
    tokenExpiration: string
    refreshToken: string
}

export interface User {
    id: string
    email: string
}

export type CurrentUser = {
    user: User
    jwt: JwtToken
} | undefined

export type LoginHandler = (email: string, password: string) => Promise<JwtToken | undefined>

export type WhoIam = (token: JwtToken) => Promise<User>

export enum LoadingStatus {
    Fetching,
    Error,
    Success,
}

export enum ViewComponentType {
    TimeSeries,
    Graphical,
    Text,
}