export interface JwtToken {
    token: string
    refreshToken: string
}

export interface User {
    id: string
    email: string
}

export type CurrentUser = {
    user: User
    jts: JwtToken
} | undefined

export type LoginHandler = (email: string, password: string) => Promise<JwtToken | undefined>

export type WhoIam = (token: JwtToken) => Promise<User>