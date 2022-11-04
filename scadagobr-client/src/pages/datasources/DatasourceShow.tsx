import {Col, Modal, Row, Space, Typography} from "antd";
import Column from "antd/es/table/Column";
import AppTable from "../../infra/components/AppTable";
import {AppButton} from "../../components/button/AppButton";
import {Link, useNavigate} from "react-location";
import {ExclamationCircleOutlined} from "@ant-design/icons";
import {openNotificationDeleted} from "../../infra/notification";
import React from "react";
import {useDataSources, useDeleteDatasource} from "../../core/hooks/datasource";

const deleteDataSource = (mutator: any, id: string) => {
    Modal.confirm({
        icon: <ExclamationCircleOutlined/>,
        content: 'Are you sure you want to delete this datasource?',
        async onOk() {
            await mutator(id)
            openNotificationDeleted()
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
                <Column title="Datapoint's" dataIndex="datapoints" key="datapoints"/>
                <Column
                    title="Action"
                    key="action"
                    render={(_: any, record: any) => (
                        <Space size="middle">
                            <Link to={record.id}>Edit</Link>
                            <a onClick={() => deleteDataSource(mutate, record.id)}>Delete</a>
                        </Space>
                    )}
                />
            </AppTable>
        </>

    )
}

export default DatasourceShow
