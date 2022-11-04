import {
    DatabaseOutlined,
    DesktopOutlined,
    ExclamationCircleOutlined,
    FundViewOutlined,
    LogoutOutlined,
    PieChartOutlined,
    TeamOutlined
} from '@ant-design/icons';
import type {MenuProps} from 'antd';
import {Layout, Menu, Modal} from 'antd';
import React from 'react';
import {Colors} from "../core/colors";
import styled from "styled-components";
import {useUserStore} from "../core/stores/userStore";
import {openNotificationWithIcon} from "../infra/notification"
import {useNavigate} from 'react-location';
import {useMenuStore} from "../core/stores/menuStore";

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
            navigate({to: '/view'});
            break;
        case '5':
            navigate({to: '/user'});
            break;
        case '9':
            Modal.confirm({
                icon: <ExclamationCircleOutlined/>,
                content: 'Are you sure you want to logout?',
                onOk() {
                    useUserStore.getState().unSetUser();
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
    const navigate = useNavigate();

    const isLoggedIn = useUserStore(e => e.user.isLoggedIn)

    const {collapsed, show} = useMenuStore(e => e.data)

    const updateMenu = useMenuStore(e => e.setData)

    if (!isLoggedIn || !show) {
        return <>{children}</>
    }

    return (
        <Layout style={{minHeight: '100vh'}}>
            <Sider collapsible collapsed={collapsed} onCollapse={value => updateMenu({collapsed: value})}>
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