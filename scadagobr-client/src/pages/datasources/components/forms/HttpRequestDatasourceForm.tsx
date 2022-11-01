import {Col, Form, Input, Radio, Row, Select} from "antd";

interface HttpRequestDatasourceFormProps {
}

export default function HttpRequestDatasourceForm(props: HttpRequestDatasourceFormProps) {

    return (
        <>
            <Form.Item name="http" label="http" rules={[{required: true}]}>
                <Input placeholder="http"/>
            </Form.Item>
        </>

    )
}