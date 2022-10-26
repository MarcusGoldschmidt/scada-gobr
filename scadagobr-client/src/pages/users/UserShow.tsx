import {Button, Col, Modal, Row, Space, Typography} from "antd";
import Column from "antd/es/table/Column";
import AppTable from "../../infra/components/AppTable";
import {AppButton} from "../../components/button/AppButton";
import useUsers from "../../core/hooks/useUsers";
import {
    CloseOutlined,
    CheckOutlined, ExclamationCircleOutlined,
} from '@ant-design/icons';
import {
    Link
} from "react-location";
import {userStore} from "../../core/stores/userStore";
import {openNotificationWithIcon} from "../../infra/notification";
import React from "react";

const removeUser = (id: number) => {
    Modal.confirm({
        icon: <ExclamationCircleOutlined/>,
        content: 'Are you sure you want to delete this user?',
        onOk() {
            openNotificationWithIcon({message: "Successfully deleted"}, `info`);
        }
    })
}

function DatasourceShow() {
    const data = useUsers();

    return (
        <>
            <Row>
                <Col span={12}>
                    <Typography.Title>Users</Typography.Title>
                </Col>
                <Col span={12} style={{textAlign: `right`}}>
                    <AppButton type="primary" size="middle">
                        Add User
                    </AppButton>
                </Col>
            </Row>
            <br/>
            <AppTable
                {...data}
            >
                <Column title="Name" dataIndex="name" key="name"/>
                <Column title="Home Url" dataIndex="homeUrl" key="homeUrl"/>
                <Column
                    title="Admin"
                    key="name"
                    render={(_: any, record: any) => (
                        (record.administrator ? <CheckOutlined/> : <CloseOutlined/>)
                    )}
                />
                <Column
                    title="Action"
                    key="action"
                    render={(_: any, record: any) => (
                        <Space size="middle">
                            <Link to={`/user/${record.id}`}>
                                Edit
                            </Link>
                            <a onClick={() => removeUser(record.id)}>Delete</a>
                        </Space>
                    )}
                />
            </AppTable>
        </>

    )
}

export default DatasourceShow
