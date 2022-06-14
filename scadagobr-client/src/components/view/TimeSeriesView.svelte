<script lang="ts">
    import {LineChart} from "@carbon/charts-svelte";
    import {onMount} from "svelte";
    import {axiosJwt} from "../../shared/net/axios";
    import {PathsV1} from "../../shared/net/request";

    interface Series {
        group: string,
        date: string,
        value: number
    }

    export let name = "LineChart"
    export let height = 300
    export let width = 600
    export let dataPointsIds: string[] = []
    export let data: Series[] = []

    onMount(async () => {
        if (dataPointsIds.length > 0) {
            const response = await axiosJwt.get(PathsV1.DataSeriesGetByGroup, {
                params: {
                    dataPointsIds: dataPointsIds
                }
            })
            data = response.data
        }
    })
</script>

<LineChart
        data={data}
        options={{
	"title": name,
	"axes": {
		"bottom": {
			"mapsTo": "timestamp",
			"scaleType": "time",
			"title": "Date",
		},
		"left": {
			"mapsTo": "value",
			"scaleType": "linear",
			"title": "Value"
		}
	},
	"timeScale": {
        "addSpaceOnEdges": "0"
	},
	"curve": "curveMonotoneX",
	"height": `${height}px`,
	"width": `${width}px`,
	"toolbar": false,
	"tooltip": {
        "valueFormatter": (e) => {
            if (e instanceof Date){
                return e.toLocaleTimeString()
            }
            return e
        }
    }
}}
/>


<style lang="scss">

</style>