import useDataSources from "../../core/hooks/useDataSources";
import {Button, Col, Row, Space, Typography} from "antd";
import Column from "antd/es/table/Column";
import AppTable from "../../infra/components/AppTable";
import {AppButton} from "../../components/button/AppButton";

function DatasourceShow() {
    const data = useDataSources();

    return (
        <>
            <Row>
                <Col span={12}>
                    <Typography.Title>Datasources</Typography.Title>
                </Col>
                <Col span={12} style={{textAlign: `right`}}>
                    <AppButton type="primary" size="middle">
                        Add Datasource
                    </AppButton>
                </Col>
            </Row>
            <br/>
            <AppTable
                {...data}
            >
                <Column title="Name" dataIndex="name" key="name" />
                <Column title="Type" dataIndex="type" key="type" />
                <Column
                    title="Action"
                    key="action"
                    render={(_: any, record: any) => (
                        <Space size="middle">
                            <a>Invite {record.id}</a>
                            <a>Delete</a>
                        </Space>
                    )}
                />
            </AppTable>
        </>

    )
}

export default DatasourceShow
