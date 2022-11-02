import useDataSources from "../../core/hooks/useDataSources";
import {Col, Modal, Row, Space, Typography} from "antd";
import Column from "antd/es/table/Column";
import AppTable from "../../infra/components/AppTable";
import {AppButton} from "../../components/button/AppButton";
import {useNavigate} from "react-location";
import {ExclamationCircleOutlined} from "@ant-design/icons";
import {useUserStore} from "../../core/stores/userStore";
import {openNotificationWithIcon} from "../../infra/notification";
import React from "react";
import useDeleteDatasource from "../../core/hooks/useDeleteDatasource";

const deleteDataSource = (mutator: any, id: string) => {
    Modal.confirm({
        icon: <ExclamationCircleOutlined/>,
        content: 'Are you sure you want to delete this datasource?',
        onOk() {
            mutator(id)
        }
    });
}

function DatasourceShow() {
    const data = useDataSources();
    const navigate = useNavigate();

    const {mutate, isLoading} = useDeleteDatasource();

    return (
        <>
            <Row>
                <Col span={12}>
                    <Typography.Title>Datasources</Typography.Title>
                </Col>
                <Col span={12} style={{textAlign: `right`}}>
                    <AppButton type="primary" size="middle" onClick={() => navigate({to: "/datasource/create"})}>
                        Add Datasource
                    </AppButton>
                </Col>
            </Row>
            <br/>
            <AppTable
                {...data}
            >
                <Column title="Name" dataIndex="name" key="name"/>
                <Column title="Type" dataIndex="type" key="type"/>
                <Column
                    title="Action"
                    key="action"
                    render={(_: any, record: any) => (
                        <Space size="middle">
                            <a>Invite {record.id}</a>
                            <a onClick={() => deleteDataSource(mutate, record.id)}>Delete</a>
                        </Space>
                    )}
                />
            </AppTable>
        </>

    )
}

export default DatasourceShow
