import {useQuery} from 'react-query'
import {axios} from "../../infra/axios";
import {PathsV1} from "../../infra/endpoints";

export default function useDataSources() {
    return useQuery([PathsV1.DataSourceGet], async () => {
        const {data} = await axios.get(
            PathsV1.DataSourceGet
        )
        return data
    })
}