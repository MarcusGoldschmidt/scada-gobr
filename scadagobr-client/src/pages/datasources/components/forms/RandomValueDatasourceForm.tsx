import {Form, FormInstance, Input} from "antd";
import InputMask from 'react-input-mask';
import {nanosecondToTimespanFormatted, validateDuration} from "../../../../core/types/timespan";
import {useMemo} from "react";


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

export default function RandomValueDatasourceForm({period}: RandomValueDatasourceFormProps) {
    const timeSpan = useMemo(
        () => nanosecondToTimespanFormatted(period, "01h00m00s"),
        [period]
    )

    return (
        <>
            <Form.Item name={["data", "period"]} initialValue={timeSpan} label="Period"
                // @ts-ignore
                       rules={[{required: true}, validateRange]}>
                <InputMask mask="99h99m99s" autoComplete="off">
                    <Input/>
                </InputMask>
            </Form.Item>
        </>
    )
}