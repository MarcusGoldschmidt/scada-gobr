import {Col, Form, Input, Radio, Row, Select} from "antd";

interface SqlDatasourceFormProps {
}

export default function SqlDatasourceForm(props: SqlDatasourceFormProps) {
    return (
        <>
            <Form.Item name="name" label="Name" rules={[{required: true}]}>
                <Input placeholder="Name"/>
            </Form.Item>
        </>

    )
}