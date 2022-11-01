import {Form, FormInstance, Input, InputNumber} from "antd";
import InputMask from 'react-input-mask';
import {validateDuration} from "../../../../core/types/timespan";


interface RandomValueDatasourceFormProps {
    period?: number;
}

const validateRange = (_: FormInstance): any => ({
    // @ts-ignore
    validator(_, value) {
        if (!validateDuration(value)) {
            return Promise.reject(new Error('Invalid period'));
        }

        return Promise.resolve();
    },
});

export default function RandomValueDatasourceForm(props: RandomValueDatasourceFormProps) {
    return (
        <>
            <Form.Item name={["data", "period"]} initialValue={"01h00m00s"} label="Period"
                // @ts-ignore
                       rules={[{required: true}, validateRange]}>
                <InputMask mask="99h99m99s" autoComplete="off">
                    <Input/>
                </InputMask>
            </Form.Item>
        </>
    )
}