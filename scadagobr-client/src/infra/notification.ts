import {notification} from "antd";
import {ArgsProps} from "antd/lib/notification";

type NotificationType = 'success' | 'info' | 'warning' | 'error';

export const openNotificationWithIcon = (props: ArgsProps, type: NotificationType = 'info') => {
    notification[type]({
        ...props,
    });
};