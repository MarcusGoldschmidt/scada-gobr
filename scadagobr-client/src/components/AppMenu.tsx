import {
    DesktopOutlined,
    FileOutlined,
    PieChartOutlined,
    TeamOutlined,
    DatabaseOutlined,
    LogoutOutlined,
    FundViewOutlined,
    ExclamationCircleOutlined
} from '@ant-design/icons';
import type {MenuProps} from 'antd';
import {Layout, Menu} from 'antd';
import React, {useState} from 'react';
import {Colors} from "../core/colors";
import styled from "styled-components";
import {userStore} from "../core/stores/userStore";
import {Button, Modal} from 'antd';
import {openNotificationWithIcon} from "../infra/notification"
import {BuildNextOptions, DefaultGenerics, useNavigate} from 'react-location';
import {NavigateOptions} from "react-location/src";

const {Header, Content, Sider} = Layout;

type MenuItem = Required<MenuProps>['items'][number];

function getItem(
    label: React.ReactNode,
    key: React.Key,
    icon?: React.ReactNode,
    children?: MenuItem[],
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
    } as MenuItem;
}

const items: MenuItem[] = [
    getItem('Home', '1', <PieChartOutlined/>),
    getItem('Views', '2', <DesktopOutlined/>),
    getItem('DataSources', '3', <DatabaseOutlined/>),
    getItem('Views', '4', <FundViewOutlined/>),
    getItem('Users', '5', <TeamOutlined/>),
    getItem('Logout', '9', <LogoutOutlined/>),
];

const onMenuClick = (item: MenuItem, navigate: any) => {
    switch (item?.key) {
        case '1':
            navigate({to: '/'});
            break;
        case '2':
            navigate({to: '/view'});
            break;
        case '3':
            navigate({to: '/datasource'});
            break;
        case '4':
            navigate({to: '/'});
            break;
        case '5':
            navigate({to: '/user'});
            break;
        case '9':
            Modal.confirm({
                icon: <ExclamationCircleOutlined/>,
                content: 'Are you sure you want to logout?',
                onOk() {
                    userStore.getState().unSetUser();
                    openNotificationWithIcon({message: "Successfully Logged out"}, `info`);
                }
            });
            break;
    }
}

const Logo = styled.h1`
  color: ${Colors.Secondary};
  text-align: center;
  font-size: 3rem;
  transition-delay: 1s;
  font-family: "3270Medium Nerd Font Mono", serif;
  max-width: 100%;
`

const AppMenu: React.FC<{ children: React.ReactNode }> = ({children}) => {
    const [collapsed, setCollapsed] = useState(false);

    const navigate = useNavigate();

    return (
        <Layout style={{minHeight: '100vh'}}>
            <Sider collapsible collapsed={collapsed} onCollapse={value => setCollapsed(value)}>
                {collapsed ? <Logo>GO</Logo> : <Logo>GOBR</Logo>}
                <Menu selectable={false} onClick={(e) => onMenuClick(e, navigate)} theme="dark" mode="inline"
                      items={items}/>
            </Sider>
            <Layout className="site-layout" style={{backgroundColor: Colors.Light}}>
                <Content style={{margin: '0.5vw 1vh'}}>
                    {children}
                </Content>
            </Layout>
        </Layout>
    );
};

export default AppMenu;