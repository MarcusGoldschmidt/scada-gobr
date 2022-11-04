import {Skeleton, Table} from "antd";
import {UseQueryResult} from "react-query";
import * as React from "react";

type AppTableProps<T> = UseQueryResult<T[]> & {
    children?: React.ReactNode,
    skeletonRows?: number,
    rowKey?: string
}

function AppTable<RecordType>({
                                  isLoading,
                                  data,
                                  children,
                                  skeletonRows = 10,
                                  rowKey = "id"
                              }: AppTableProps<RecordType>) {
    if (isLoading) {
        return <>
            <Skeleton active paragraph={{rows: skeletonRows}}></Skeleton>
        </>
    }

    return (
        <Table rowKey={rowKey} dataSource={data as any} pagination={{position: ['bottomRight']}}>
            {children}
        </Table>
    )
}

export default AppTable
