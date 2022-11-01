import {useMutation, useQuery, useQueryClient} from "react-query";
import {PathsV1} from "../../infra/endpoints";
import {axios} from "../../infra/axios";

export default function useDeleteDatasource() {

    const queryClient = useQueryClient()

    return useMutation(async (id: string) => {
        const {data} = await axios.delete(PathsV1.DataSourceDelete(id))

        await queryClient.invalidateQueries([PathsV1.DataSourceGet])
        return data
    });
}

