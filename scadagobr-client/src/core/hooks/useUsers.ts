import {useQuery} from 'react-query'
import {axios} from "../../infra/axios";
import {PathsV1} from "../../infra/endpoints";

export default function useUsers() {
    return useQuery([PathsV1.UserGet], async () => {
        const {data} = await axios.get(
            PathsV1.UserGet
        )
        return data
    })
}