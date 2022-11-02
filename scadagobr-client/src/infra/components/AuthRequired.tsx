import React from 'react'
import {useUserStore} from "../../core/stores/userStore";
import {Navigate} from "react-location";

function AuthRequired({children}: { children: React.ReactElement }) {

    const isLoggedIn = useUserStore(e => e.user.isLoggedIn)

    if (isLoggedIn) {
        return children
    }

    return <Navigate to="/login" />
}

export default AuthRequired
