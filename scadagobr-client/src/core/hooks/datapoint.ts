import {DataPointLoggingType, DatasourceType} from "../enums/datasource_type";
import {QueryClient, useMutation, useQuery, useQueryClient} from "react-query";
import {axios} from "../../infra/axios";
import {PathsV1} from "../../infra/endpoints";

export interface Datapoint {
    name: string
    datasourceId: string
    isEnable: boolean
    unit: string
    purgeAfter: number
    type: DatasourceType
    typeData: any,
    createdAt: string
    loggingType: DataPointLoggingType
}

const invalidateDataSource = async (queryClient: QueryClient, datasourceId: string) => {
    await queryClient.invalidateQueries([PathsV1.DataSourceGetById(datasourceId)])
    await queryClient.invalidateQueries([PathsV1.DataPointsGet(datasourceId)])

}

export function useCreateDataPoint({onSuccess}: { onSuccess?: () => void }) {
    const queryClient = useQueryClient()

    return useMutation(async (datapoint: Datapoint) => {
        const {data} = await axios.post(PathsV1.DataPointCreate(datapoint.datasourceId), datapoint)

        await invalidateDataSource(queryClient, datapoint.datasourceId)
        return data
    }, {onSuccess});
}

export function useDeleteDataPoint({onSuccess}: { onSuccess?: () => void }) {
    const queryClient = useQueryClient()

    return useMutation(async ({datasourceId, datapointId}: { datasourceId: string, datapointId: string }) => {
        const {data} = await axios.delete(PathsV1.DataPointDelete(datasourceId, datapointId))

        await invalidateDataSource(queryClient, datasourceId)
        return data
    }, {onSuccess});
}

export function useDataPoints(datasourceId: string) {
    return useQuery([PathsV1.DataPointsGet(datasourceId)], async () => {
        const {data} = await axios.get<Datapoint[]>(
            PathsV1.DataPointsGet(datasourceId)
        )
        return data
    })
}