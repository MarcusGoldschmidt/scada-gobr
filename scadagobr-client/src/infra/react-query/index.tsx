import axios from 'axios'
import React, {FC} from 'react'
import {
    QueryCache,
    QueryClient,
    QueryClientProvider,
    MutationCache
} from 'react-query'
import {openNotificationWithIcon} from "../notification";

interface Error {
    errors: Record<string, string>[]
    type: string
    title: string
    status: number
    detail: string
}

export const getErrorDetails = (error: Error) => {
    if (error?.status === 400) {
        return error?.errors
    }

    return 'Server error'
}

export const reactQueryClient = new QueryClient({
    mutationCache: new MutationCache({
        onError: error => {
            if (axios.isAxiosError(error)) {
                const errorDetails = getErrorDetails(error?.response?.data as Error)
                if (typeof errorDetails === 'string') {
                    openNotificationWithIcon({
                        message: errorDetails,
                    }, 'error')
                } else {
                    errorDetails.forEach(error => {
                        openNotificationWithIcon({
                            message: error.error,
                        }, 'error')
                    })
                }
            }
        }
    }),
    queryCache: new QueryCache({
        onError: (error: unknown) => {
            if (axios.isAxiosError(error)) {
                const errorDetails = getErrorDetails(error?.response?.data as Error)
                if (typeof errorDetails === 'string') {
                    openNotificationWithIcon({
                        message: errorDetails,
                    }, 'error')
                } else {
                    errorDetails.forEach(error => {
                        errorDetails.forEach(error => {
                            openNotificationWithIcon({
                                message: error.error,
                            }, 'error')
                        })
                    })
                }
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
