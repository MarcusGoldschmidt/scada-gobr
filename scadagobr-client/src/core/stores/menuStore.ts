import create from "zustand";

export interface MenuData {
    collapsed: boolean
    show: boolean
}

export interface MenuStore {
    data: MenuData
    setData: (user: MenuData | Partial<MenuData>) => void
}

export const useMenuStore = create<MenuStore>((set) => {
    return {
        data: {
            collapsed: false,
            show: true,
        },
        setData: (data: MenuData | Partial<MenuData>) => {
            set(s => {
                return {...s, data: {...s.data, ...data}}
            })
        }
    }
})