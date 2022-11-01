import {ReactLocation, Router, Navigate, Outlet} from 'react-location'
import Home from "../../pages/Home";
import AuthRequired from "./AuthRequired";
import AppMenu from "../../components/AppMenu";
import DatasourceShow from "../../pages/datasources/DatasourceShow";
import UserShow from "../../pages/users/UserShow";
import NotFound from "../../components/NotFound";
import DatasourceInput from "../../pages/datasources/DatasourceInput";

const reactLocation = new ReactLocation()

export default function AppRouter() {
    return (
        <>
            <AuthRequired>
                <Router
                    location={reactLocation}
                    routes={[
                        {
                            path: '/',
                            element: <Home/>,
                        },
                        {
                            path: 'datasource',
                            children: [
                                {
                                    path: '/',
                                    element: <DatasourceShow/>,
                                },
                                {
                                    path: '/create',
                                    element: <DatasourceInput/>,
                                },
                                {
                                    path: '/:id/edit',
                                    element: <DatasourceInput/>,
                                },
                            ]
                        },
                        {
                            path: '/user/',
                            element: <UserShow/>,
                        },
                        {
                            element: <NotFound/>,
                        },
                    ]}>
                    <AppMenu>
                        <Outlet/>
                    </AppMenu>
                </Router>
            </AuthRequired>
        </>
    )
}
