import {Form, FormInstance, InputNumber, Row} from "antd";
import {validateDuration} from "../../../core/types/timespan";
import TimeSpanInput from "../../../components/Input/TimeSpanInput";

interface RandomValueDataPointForm {
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

export default function RandomValueDataPointForm(props: RandomValueDataPointForm) {

    return (
        <>
            <Row>
                <Form.Item name={["data", "initialValue"]} initialValue={0} label="Initial Value"
                    // @ts-ignore
                           rules={[{type: 'number', required: true}, validateRange]}>
                    <InputNumber/>
                </Form.Item>
                <Form.Item name={["data", "endValue"]} initialValue={100} label="End Value"
                    // @ts-ignore
                           rules={[{type: 'number', required: true}, validateRange]}>
                    <InputNumber/>
                </Form.Item>
            </Row>
        </>
    )
}