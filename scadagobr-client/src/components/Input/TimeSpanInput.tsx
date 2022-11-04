import {Form, FormInstance, Input} from "antd";
import InputMask from 'react-input-mask';
import {useMemo} from "react";
import {nanosecondToTimespanFormatted, validateDuration} from "../../core/types/timespan";


interface TimeSpanInputProps {
    period?: number;
    name: string | string[];
    label: string;
    required?: boolean;
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

export default function TimeSpanInput({
                                          period = undefined,
                                          name,
                                          required = false,
                                          label
                                      }: TimeSpanInputProps) {
    const timeSpan = useMemo(
        () => nanosecondToTimespanFormatted(period, "01h00m00s"),
        [period]
    )

    return (
        <>
            <Form.Item name={name} initialValue={timeSpan} label={label}
                // @ts-ignore
                       rules={[{required}, validateRange]}>
                <InputMask mask="99h99m99s" autoComplete="off">
                    <Input/>
                </InputMask>
            </Form.Item>
        </>
    )
}