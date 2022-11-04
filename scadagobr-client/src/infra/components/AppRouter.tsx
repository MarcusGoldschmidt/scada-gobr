import {Outlet, ReactLocation, Router} from 'react-location'
import Home from "../../pages/Home";
import AuthRequired from "./AuthRequired";
import AppMenu from "../../components/AppMenu";
import DatasourceShow from "../../pages/datasources/DatasourceShow";
import UserShow from "../../pages/users/UserShow";
import NotFound from "../../components/NotFound";
import DatasourceInput from "../../pages/datasources/DatasourceInput";
import Login from "../../pages/auth/Login";
import React from "react";
import DatasourceEdit from "../../pages/datasources/DatasourceEdit";

const reactLocation = new ReactLocation()

const useAuth = (page: React.ReactElement) => <AuthRequired children={page}/>

export default function AppRouter() {
    return (
        <>
            <Router
                location={reactLocation}
                routes={[
                    {
                        path: '/',
                        element: useAuth(<Home/>),
                    },
                    {
                        path: '/login/',
                        element: <Login/>,
                    },
                    {
                        path: '/datasource',
                        children: [
                            {
                                path: '/',
                                element: useAuth(<DatasourceShow/>),
                            },
                            {
                                path: '/create',
                                element: useAuth(<DatasourceInput/>),
                            },
                            {
                                path: '/:id',
                                element: async (e) => useAuth(<DatasourceEdit datasourceId={e.params.id}/>),
                            },
                        ]
                    },
                    {
                        path: '/user/',
                        element: useAuth(<UserShow/>),
                    },
                    {
                        element: useAuth(<NotFound/>),
                    },
                ]}>
                <AppMenu>
                    <Outlet/>
                </AppMenu>
            </Router>
        </>
    )
}
