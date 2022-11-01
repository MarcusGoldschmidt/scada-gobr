import {Col, Form, Input, Radio, Row, Select} from "antd";
import {DatasourceType} from "../../../core/enums/datasource_type";

interface DatasourceCommonFormProps {
    onChangeType?: (type: DatasourceType) => void
}

export default function DatasourceCommonForm(
    {
        onChangeType = (_) => {
        }
    }: DatasourceCommonFormProps
) {

    return (
        <>
            <Form.Item name="name" label="Name" rules={[{required: true}]}>
                <Input placeholder="Name"/>
            </Form.Item>
            <Form.Item initialValue={DatasourceType.RandomValue} name="type" label="Type" rules={[{required: true}]}>
                <Select onChange={onChangeType}>
                    {Object.keys(DatasourceType).map((type) =>
                        <Select.Option key={type} value={type}>{type}</Select.Option>)
                    }
                </Select>
            </Form.Item>
        </>

    )
}