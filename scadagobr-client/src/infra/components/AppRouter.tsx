import {ReactLocation, Router, Navigate, Outlet} from 'react-location'
import Home from "../../pages/Home";
import AuthRequired from "./AuthRequired";

const reactLocation = new ReactLocation()

export default function AppRouter() {
    return (
        <Router
            location={reactLocation}
            routes={[
                {
                    path: '/',
                    element: <AuthRequired><Home/></AuthRequired>,
                }
            ]}>
            <Outlet/>
        </Router>
    )
}
