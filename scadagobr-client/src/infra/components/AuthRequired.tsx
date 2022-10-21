import React, {useState} from 'react'
import {userStore} from "../../core/stores/userStore";
import Login from "../../pages/auth/Login";

function AuthRequired({children}: { children: React.ReactElement }) {

    const isLoggedIn = userStore(e => e.user.isLoggedIn)

    if (isLoggedIn) {
        return children
    }

    return <Login/>
}

export default AuthRequired
