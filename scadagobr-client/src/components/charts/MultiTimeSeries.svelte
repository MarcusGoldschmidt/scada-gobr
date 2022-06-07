<script>
    import {LineChart} from "@carbon/charts-svelte";

    let data = [
        {
            "group": "Dataset 1",
            "date": "2019-01-01T04:00:00.000Z",
            "value": -10000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-01T09:00:00.000Z",
            "value": -12000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-01T14:00:00.000Z",
            "value": -14000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-02T04:00:00.000Z",
            "value": -25000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-02T06:00:00.000Z",
            "value": -26000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-03T04:00:00.000Z",
            "value": -10000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-03T09:00:00.000Z",
            "value": 10000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-03T14:00:00.000Z",
            "value": 12000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-05T04:00:00.000Z",
            "value": 45000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-07T04:00:00.000Z",
            "value": 49000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-07T19:00:00.000Z",
            "value": 45000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-09T04:00:00.000Z",
            "value": 50000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-09T09:00:00.000Z",
            "value": 52000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-09T19:00:00.000Z",
            "value": 55000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-10T04:00:00.000Z",
            "value": 50000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-12T04:00:00.000Z",
            "value": 65000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-13T04:00:00.000Z",
            "value": 80000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-14T14:00:00.000Z",
            "value": 85000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-15T11:00:00.000Z",
            "value": 90000
        },
        {
            "group": "Dataset 1",
            "date": "2019-01-15T22:00:00.000Z",
            "value": 70000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-01T04:00:00.000Z",
            "value": 20000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-01T07:00:00.000Z",
            "value": 22000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-01T20:00:00.000Z",
            "value": 24000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-02T04:00:00.000Z",
            "value": 35000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-02T11:00:00.000Z",
            "value": 36000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-03T04:00:00.000Z",
            "value": 20000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-03T10:00:00.000Z",
            "value": 20000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-03T22:00:00.000Z",
            "value": 22000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-05T04:00:00.000Z",
            "value": 62000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-06T04:00:00.000Z",
            "value": 52000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-07T04:00:00.000Z",
            "value": 52000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-07T19:00:00.000Z",
            "value": 52000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-09T04:00:00.000Z",
            "value": 60000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-09T09:00:00.000Z",
            "value": 62000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-09T14:00:00.000Z",
            "value": 62000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-12T04:00:00.000Z",
            "value": 65000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-14T04:00:00.000Z",
            "value": 40000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-15T09:00:00.000Z",
            "value": 45000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-15T14:00:00.000Z",
            "value": 35000
        },
        {
            "group": "Dataset 2",
            "date": "2019-01-15T22:00:00.000Z",
            "value": 30000
        }
    ];

</script>

<LineChart
        data={data}
        options={{
	"title": "Line (dense time series)",
	"axes": {
		"bottom": {
			"mapsTo": "date",
			"scaleType": "time"
		},
		"left": {
			"mapsTo": "value",
			"scaleType": "linear"
		}
	},
	"curve": "curveMonotoneX",
	"height": "400px"
}}
/>


<style lang="scss">

</style>