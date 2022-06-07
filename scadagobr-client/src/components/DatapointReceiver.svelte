<script lang="ts">

    import {createWs} from "../shared/net/ws";
    import {onMount} from "svelte";
    import {PathsWsV1} from "../shared/net/request";

    export let dataPointId: string

    let socket

    $: connect = () => {
        socket = createWs(PathsWsV1.DataPoint + dataPointId);

        socket.onmessage = (event) => {
            console.log(JSON.parse(event.data));
        };

        socket.onclose = function (e) {
            console.error('Socket is closed. Reconnect will be attempted in 5 second.', e.reason);
            setTimeout(function () {
                connect();
            }, 5000);
        };
    }

    onMount(() => {
        connect()

        return () => {
            socket.close();
        }
    })
</script>

<slot></slot>