import {useMutation, useQuery, useQueryClient} from "react-query";
import {axios} from "../../infra/axios";
import {PathsV1} from "../../infra/endpoints";
import {Datapoint} from "./datapoint";
import {DatasourceType} from "../enums/datasource_type";

interface Datasource {
    createdAt: string
    data: any
    dataPoints: Datapoint[]
    id: string
    name: string
    type: DatasourceType
    updatedAt: string
}

export function useDeleteDatasource() {
    const queryClient = useQueryClient()

    return useMutation(async (id: string) => {
        const {data} = await axios.delete(PathsV1.DataSourceDelete(id))

        await queryClient.invalidateQueries([PathsV1.DataSourceGet])
        return data
    });
}

export function useCreateDatasource({onSuccess}: { onSuccess?: () => void }) {
    const queryClient = useQueryClient()

    return useMutation(async (event) => {
        const {data} = await axios.post(PathsV1.DataSourceCreate, event)
        await queryClient.invalidateQueries([PathsV1.DataSourceGet])
        return data
    }, {
        onSuccess
    });
}

export function useUpdateDatasource({onSuccess}: { onSuccess?: () => void }) {
    const queryClient = useQueryClient()

    return useMutation(async (event) => {
        const {data} = await axios.post(PathsV1.DataSourceUpdate, event)
        await queryClient.invalidateQueries([PathsV1.DataSourceGet])
        return data
    }, {
        onSuccess
    });
}

export function useDataSource(datasourceId: string) {
    return useQuery([PathsV1.DataSourceGetById(datasourceId)], async () => {
        const {data} = await axios.get<Datasource>(
            PathsV1.DataSourceGetById(datasourceId)
        )
        return data
    })
}

export function useDataSources() {
    return useQuery([PathsV1.DataSourceGet], async () => {
        const {data} = await axios.get(
            PathsV1.DataSourceGet
        )
        return data
    })
}