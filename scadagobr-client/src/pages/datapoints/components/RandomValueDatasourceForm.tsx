import {Form, FormInstance, Input, InputNumber} from "antd";

interface RandomValueDatasourceFormProps {
}


const validateRange = ({getFieldValue}: FormInstance): any => ({
    // @ts-ignore
    validator(_, __) {
        if (+getFieldValue("initialValue") > +getFieldValue("endValue")) {
            return Promise.reject(new Error('Initial value must be less than end value'));
        }

        return Promise.resolve();
    },
});

export default function RandomValueDatasourceForm(props: RandomValueDatasourceFormProps) {

    return (
        <>
            <Form.Item name="initialValue" initialValue={0} label="Initial Value"
                // @ts-ignore
                       rules={[{type: 'number', required: true}, validateRange]}>
                <InputNumber/>
            </Form.Item>
            <Form.Item name="endValue" initialValue={100} label="End Value"
                // @ts-ignore
                       rules={[{type: 'number', required: true}, validateRange]}>
                <InputNumber/>
            </Form.Item>
        </>
    )
}