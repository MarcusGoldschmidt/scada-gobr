export const __DEV__ = import.meta.env.DEV

export const __BASE_URL__ = import.meta.env.REACT_APP_BASE_URL

export const __EMBEDDED__ = import.meta.env.VITE_EMBEDDED ?? false

export const __API_URL__ = __EMBEDDED__ ? window.location.href : import.meta.env.VITE_API_URL
