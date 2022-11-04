import {Button, Col, Form, Row, Typography} from "antd";
import {useLocation} from "react-location";
import DatasourceCommonForm from "./components/DatasourceCommonForm";
import React, {useState} from "react";
import {DatasourceType} from "../../core/enums/datasource_type";
import SqlDatasourceForm from "./components/forms/SqlDatasourceForm";
import RandomValueDatasourceForm from "./components/forms/RandomValueDatasourceForm";
import HttpRequestDatasourceForm from "./components/forms/HttpRequestDatasourceForm";
import HttpServerDatasourceForm from "./components/forms/HttpServerDatasourceForm";
import {openNotificationCreated} from "../../infra/notification";
import {useCreateDatasource} from "../../core/hooks/datasource";

const mapComponent = {
    [DatasourceType.Sql]: <SqlDatasourceForm/>,
    [DatasourceType.RandomValue]: <RandomValueDatasourceForm/>,
    [DatasourceType.HttpRequest]: <HttpRequestDatasourceForm/>,
    [DatasourceType.HttpServer]: <HttpServerDatasourceForm/>,
}

function DatasourceInput() {
    const [form] = Form.useForm();
    const {history} = useLocation()

    const [datasourceType, setDatasourceType] = useState<DatasourceType>(DatasourceType.RandomValue);

    const {mutate} = useCreateDatasource({
        onSuccess: () => {
            openNotificationCreated()
            history.back()
        }
    });

    return (
        <>
            <Row>
                <Col span={12}>
                    <Typography.Title> Create Datasource</Typography.Title>
                </Col>
            </Row>
            <Form form={form} layout="vertical" onFinish={(e) => mutate(e)}>
                <Row>
                    <Col xs={24} md={12}>
                        <DatasourceCommonForm
                            onChangeType={setDatasourceType}
                        ></DatasourceCommonForm>
                        {mapComponent[datasourceType] ?? <Typography.Title>Datasource type not found</Typography.Title>}
                        <Form.Item>
                            <Button type="primary" htmlType="submit">
                                Submit
                            </Button>
                        </Form.Item>
                    </Col>
                </Row>
            </Form>
        </>
    )
}

export default DatasourceInput