import {Col, Form, Input, Radio, Row, Select} from "antd";

interface HttpServerDatasourceFormProps {
}

export default function HttpServerDatasourceForm(props: HttpServerDatasourceFormProps) {

    return (
        <>
            <Form.Item name="name" label="Name" rules={[{required: true}]}>
                <Input placeholder="Name"/>
            </Form.Item>
        </>

    )
}