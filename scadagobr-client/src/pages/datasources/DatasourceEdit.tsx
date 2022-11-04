import {Button, Col, Form, Input, Modal, Row, Select, Space, Switch, Typography} from "antd";
import {Link, useLocation} from "react-location";
import React from "react";
import {DataPointLoggingType, DatasourceType} from "../../core/enums/datasource_type";
import SqlDatasourceForm from "./components/forms/SqlDatasourceForm";
import RandomValueDatasourceForm from "./components/forms/RandomValueDatasourceForm";
import HttpRequestDatasourceForm from "./components/forms/HttpRequestDatasourceForm";
import HttpServerDatasourceForm from "./components/forms/HttpServerDatasourceForm";
import {openNotificationDeleted, openNotificationUpdated} from "../../infra/notification";
import {useDataSource, useUpdateDatasource} from "../../core/hooks/datasource";
import AppTable from "../../infra/components/AppTable";
import {randomValueDataPointColumn} from "./components/datapointColunms";
import Column from "antd/es/table/Column";
import RandomValueDataPointForm from "../datapoints/components/RandomValueDataPointForm";
import {useCreateDataPoint, useDataPoints, useDeleteDataPoint} from "../../core/hooks/datapoint";
import {ExclamationCircleOutlined} from "@ant-design/icons";

const mapComponent = {
    [DatasourceType.Sql]: <SqlDatasourceForm/>,
    [DatasourceType.RandomValue]: <RandomValueDatasourceForm/>,
    [DatasourceType.HttpRequest]: <HttpRequestDatasourceForm/>,
    [DatasourceType.HttpServer]: <HttpServerDatasourceForm/>,
}

const mapTableComponent: Record<DatasourceType, JSX.Element[]> = {
    [DatasourceType.RandomValue]: randomValueDataPointColumn,
    [DatasourceType.Sql]: [],
    [DatasourceType.HttpRequest]: [],
    [DatasourceType.HttpServer]: [],
}

const mapDatapointInputComponent = {
    [DatasourceType.RandomValue]: (e: any) => <RandomValueDataPointForm {...e}/>,
    [DatasourceType.Sql]: [],
    [DatasourceType.HttpRequest]: [],
    [DatasourceType.HttpServer]: [],
}


interface datasourceEditProps {
    datasourceId: string
}

function DatasourceEdit(props: datasourceEditProps) {
    const [form] = Form.useForm();
    const [formDataPoint] = Form.useForm();
    const {history} = useLocation()

    const [formOpen, setFormOpen] = React.useState(false)

    const updateDatasource = useUpdateDatasource({
        onSuccess: () => {
            openNotificationUpdated()
        }
    });

    const createDataPoint = useCreateDataPoint({
        onSuccess: () => {
            setFormOpen(false)
        }
    });
    const deleteDataPoint = useDeleteDataPoint({
        onSuccess: () => {
            openNotificationDeleted()
        }
    });

    const datapoints = useDataPoints(props.datasourceId);
    const dataSource = useDataSource(props.datasourceId)


    const onFinishDataPointForm = (e: any) => {

        createDataPoint.mutate({
            ...e,
            datasourceId: props.datasourceId,
            isEnable: true,
            type: dataSource.data?.type,
        })
    }

    if (dataSource.isLoading) {
        return <div>Loading...</div>
    }

    const {data} = dataSource

    const onChangeEnableDataPoint = () => {

    }

    const onDeleteDataPoint = (datasourceId: string, datapointId: string) => {
        Modal.confirm({
            icon: <ExclamationCircleOutlined/>,
            content: 'Are you sure you want to delete this datapoint?',
            async onOk() {
                deleteDataPoint.mutate({
                    datasourceId,
                    datapointId
                })
            }
        });
    }

    return (
        <>
            <Row>
                <Col span={24}>
                    <Typography.Title> Edit Datasource</Typography.Title>
                </Col>
            </Row>
            <Form form={form} layout="vertical" onFinish={(e) => {
            }}>
                <Row>
                    <Col xs={24} md={12}>
                        <Form.Item name="name" initialValue={data!.name} label="Name" rules={[{required: true}]}>
                            <Input placeholder="Name"/>
                        </Form.Item>
                        {mapComponent[data!.type] ?? <Typography.Title>Datasource type not found</Typography.Title>}
                    </Col>
                </Row>
                <Row>
                    <Col xs={12} md={6}>
                        {!formOpen &&
                            <Button type="primary" htmlType="submit" onClick={e => setFormOpen(true)}>
                                Create Datapoint
                            </Button>
                        }
                    </Col>
                    <Col xs={12} md={6} style={{textAlign: "right"}}>
                        <Button type="primary" htmlType="submit">
                            Save
                        </Button>
                    </Col>
                </Row>
            </Form>
            <br/>
            <hr/>

            {formOpen &&
                <Form form={formDataPoint} layout="vertical" onFinish={onFinishDataPointForm}>

                    <Row>
                        <Col xs={24} md={12}>
                            <Form.Item name="name" label="Name" rules={[{required: true}]}>
                                <Input placeholder="Name"/>
                            </Form.Item>
                            {/*// @ts-ignore*/}
                            {mapDatapointInputComponent[data!.type](data)}

                            <Form.Item initialValue={DataPointLoggingType.AllData}
                                       name="loggingType" label="Logging type"
                                       rules={[{required: true}]}>
                                <Select>
                                    {Object.keys(DataPointLoggingType).map((type) =>
                                        <Select.Option key={type} value={type}>{type}</Select.Option>)
                                    }
                                </Select>
                            </Form.Item>

                            <Form.Item name="unit" label="Unit of measure">
                                <Input/>
                            </Form.Item>

                            <Row justify={"end"}>
                                <Col xs={12} md={4} style={{textAlign: "right"}}>
                                    <Button danger onClick={e => {
                                        setFormOpen(false)
                                        formDataPoint.resetFields()
                                    }}>
                                        Cancel
                                    </Button>
                                </Col>
                                <Col xs={12} md={4} style={{textAlign: "right"}}>
                                    <Button type="primary" htmlType="submit">
                                        Create Datapoint
                                    </Button>
                                </Col>
                            </Row>
                        </Col>
                    </Row>
                </Form>
            }

            {!formOpen &&
                // @ts-ignore
                <AppTable
                    {...datapoints}
                >
                    {mapTableComponent[data!.type].map(e => e)}
                    <Column
                        title="Enable"
                        key="enable"
                        render={(_: any, record: any) => (
                            <Switch defaultChecked={record.isEnable} onChange={onChangeEnableDataPoint}/>
                        )}
                    />
                    <Column
                        title="Action"
                        key="action"
                        render={(_: any, record: any) => (
                            <Space size="middle">
                                <Link to={`datapoint/${record.id}`}>Edit</Link>
                                <a onClick={e => onDeleteDataPoint(props.datasourceId, record.id)}>Delete</a>
                            </Space>
                        )}
                    />
                </AppTable>}

        </>
    )
}

export default DatasourceEdit