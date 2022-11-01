import axios from 'axios'
import React, {FC} from 'react'
import {
    QueryCache,
    QueryClient,
    QueryClientProvider,
    MutationCache
} from 'react-query'
import {openNotificationWithIcon} from "../notification";

interface multipleErrors {
    errors: Record<string, string>
}

interface singleError {
    error: string
}

type Error = multipleErrors | singleError

export const notifyError = (error: any) => {
    const errorDetails = error as Error

    if ('errors' in errorDetails) {
        Object.keys(errorDetails.errors).forEach(key => {
            openNotificationWithIcon({message: errorDetails.errors[key]}, 'error')
        })
    }

    if ('error' in errorDetails) {
        openNotificationWithIcon({message: errorDetails.error}, 'error')
    }
}

export const reactQueryClient = new QueryClient({
    mutationCache: new MutationCache({
        onError: error => {
            if (axios.isAxiosError(error)) {
                notifyError(error?.response?.data)
            }
        }
    }),
    queryCache: new QueryCache({
        onError: (error: unknown) => {
            if (axios.isAxiosError(error)) {
                notifyError(error?.response?.data)
            }
        }
    }),
    defaultOptions: {
        queries: {
            refetchOnWindowFocus: false
        }
    }
})

export const ReactQueryProvider: FC<{ children: React.ReactNode }> = props => (
    <QueryClientProvider client={reactQueryClient}>
        {props.children}
    </QueryClientProvider>
)
