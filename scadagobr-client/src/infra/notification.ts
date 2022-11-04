import {notification} from "antd";
import {ArgsProps} from "antd/lib/notification";

type NotificationType = 'success' | 'info' | 'warning' | 'error';

export const openNotificationWithIcon = (props: ArgsProps, type: NotificationType = 'info') => {
    notification[type]({
        ...props,
    });
};

export const openNotificationCreated = (message: string = "Resource created") => {
    openNotificationWithIcon({
        message,
    }, 'success');
}

export const openNotificationUpdated = (message: string = "Resource updated") => {
    openNotificationWithIcon({
        message,
    }, 'success');
}

export const openNotificationDeleted = (message: string = "Resource deleted") => {
    openNotificationWithIcon({
        message,
    }, 'error');
}