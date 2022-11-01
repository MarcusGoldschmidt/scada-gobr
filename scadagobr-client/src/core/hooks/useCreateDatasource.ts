import {useMutation, useQuery} from "react-query";
import {PathsV1} from "../../infra/endpoints";
import {axios} from "../../infra/axios";

export default function useCreateDatasource() {
    return useMutation(async (event) => {
        const {data} = await axios.post(PathsV1.DataSourceCreate, event)
        return data
    });
}

