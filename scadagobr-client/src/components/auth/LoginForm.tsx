import {Button, Checkbox, Form, Input, notification} from 'antd';
import React from 'react';
import {openNotificationWithIcon} from "../../infra/notification";

const openNotification = () => {
    notification.open({
        message: 'Notification Title',
        description:
            'This is the content of the notification. This is the content of the notification. This is the content of the notification.',
        onClick: () => {
            console.log('Notification Clicked!');
        },
    });
};

const App: React.FC = () => {
    const onFinish = async (values: any) => {

    };

    const onFinishFailed = (errorInfo: any) => {
        openNotificationWithIcon({
            message: 'Error',
        });
    };

    return (
        <Form
            name="basic"
            labelCol={{ span: 6 }}
            wrapperCol={{ span: 16 }}
            initialValues={{ remember: true }}
            onFinish={onFinish}
            onFinishFailed={onFinishFailed}
            autoComplete="off"
        >
            <Form.Item
                label="Username"
                name="username"
                rules={[{ required: true, message: 'Please input your username!' }]}
            >
                <Input />
            </Form.Item>

            <Form.Item
                label="Password"
                name="password"
                rules={[{ required: true, message: 'Please input your password!' }]}
            >
                <Input.Password />
            </Form.Item>

            <Form.Item wrapperCol={{ sm: {offset: 12}, md: {offset: 12}, lg: {offset: 6, span: 16} }}>
                <Button type="primary" htmlType="submit">
                    Submit
                </Button>
            </Form.Item>
        </Form>
    );
};

export default App;