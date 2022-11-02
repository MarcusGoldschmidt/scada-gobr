import create from 'zustand'
import {getLocalStorage, removeLocalStorage, setLocalStorage} from "../../infra/storage";
import {JwtToken} from "../../infra/response_types";

export interface User {
    token: string
    tokenExpiration: string
    refreshToken: string
    isLoggedIn: boolean
    name: string,
}

const defaultUser: User = {
    token: '',
    refreshToken: '',
    tokenExpiration: '',
    isLoggedIn: false,
    name: '',
}

export interface UserStore {
    user: User
    setUser: (user: User) => void
    updateJwt: (jwt: JwtToken) => void
    unSetUser: () => void
}

const AuthStoreKey = "CURRENT_USER"

export const useUserStore = create<UserStore>((set) => {
    return {
        user: getLocalStorage<User>(AuthStoreKey) || defaultUser,
        setUser: (user: User) => {
            setLocalStorage(AuthStoreKey, user)
            set({user})
        },
        unSetUser: () => {
            removeLocalStorage(AuthStoreKey)
            set({user: defaultUser})
        },
        updateJwt: (jwt: JwtToken) => {
            set((state) => {
                const newData = {
                    ...state.user,
                    ...jwt
                }

                setLocalStorage(AuthStoreKey, newData)
                return {...state, user: newData}
            })
        }
    }
})