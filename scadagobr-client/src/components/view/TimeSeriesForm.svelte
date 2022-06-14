<script lang="ts">
    import {createEventDispatcher, onMount} from 'svelte';
    import {axiosJwt} from "../../shared/net/axios";
    import {PathsV1} from "../../shared/net/request";
    import {NotificationType, sendNotification} from "../../shared/stores/notifications";

    const dispatch = createEventDispatcher();

    let dataSetOptions: any[] = []

    let dataPointsOptions: any[] = []

    let currentDataSet: any = null

    let selectedDataPoints: any[] = []
    export let label = "Time Series"
    export let defaultPeriod = 60
    export let defaultWidth = 300
    export let defaultHeight = 200
    export let selectedDataPointsIds: string[] = []

    let selectedDataPoint = null

    onMount(async () => {
        if (selectedDataPointsIds) {
            const response = await Promise.all(selectedDataPointsIds.map(e => axiosJwt.get(PathsV1.DataPointGetById(e))));

            selectedDataPoints = [].concat(...response.map(e => e.data));
        }
    })

    $: {
        if (currentDataSet) {
            dataPointsOptions = currentDataSet.dataPoints
        } else {
            dataPointsOptions = null
        }
    }

    onMount(() => {
        axiosJwt.get(PathsV1.DataSourceGet)
            .then(res => {
                dataSetOptions = res.data
            })
            .catch(err => {
                sendNotification(err.message, "", NotificationType.Danger)
            })
    });

    function onChange() {
        const dataPointsIds = selectedDataPoints.map(dataPoint => dataPoint.id)
        dispatch('change', {dataPointsIds, period: defaultPeriod, width: defaultWidth})
    }

    function addDataPoint(e) {
        if (selectedDataPoint == null) {
            return
        }

        if (selectedDataPoints.findIndex(x => x.id === selectedDataPoint.id) == -1) {
            selectedDataPoints = [...selectedDataPoints, selectedDataPoint]
            onChange()
        } else {
            sendNotification("Data point already selected", "", NotificationType.Warning)
        }
    }

    let removeDataPoint = (e) => {
        selectedDataPoints = selectedDataPoints.filter(x => x.id !== e.id)
        onChange()
    };
</script>

<label class="label mt-3">Label</label>
<input class="input"
       bind:value={label}
       on:change={onChange}
       placeholder="Label"/>

<label class="label mt-3">Default period (minutes)</label>
<input class="input"
       bind:value={defaultPeriod}
       on:change={onChange}
       type="number"
       placeholder="Text"/>

<label class="label mt-3">Default width</label>
<input class="input"
       bind:value={defaultWidth}
       on:change={onChange}
       type="number"
       placeholder="Text"/>

<label class="label mt-3">Default height</label>
<input class="input"
       bind:value={defaultHeight}
       on:change={onChange}
       type="number"
       placeholder="Text"/>

<label class="label mt-3">Dataset</label>
<div class="select" style="width: 100%">
    <select style="width: 100%" bind:value={currentDataSet}>
        {#each dataSetOptions as dataSet (dataSet.id)}
            <option value={dataSet}>{dataSet.name}</option>
        {/each}
    </select>
</div>

{#if currentDataSet}
    <label class="label mt-3">Data Point</label>
    <div class="select" style="width: 100%">
        <select style="width: 100%" bind:value={selectedDataPoint}>
            {#each dataPointsOptions as dataPoint (dataPoint.id)}
                <option value={dataPoint}>{dataPoint.name}</option>
            {/each}
        </select>
    </div>
    <div class="has-text-right mt-2">
        <button class="button is-small"
                on:click={addDataPoint}>
            Add datapoint
        </button>
    </div>
{/if}
<hr>

<h5>Datapoints</h5>

<ul>
    {#each selectedDataPoints as dataPoint (dataPoint.id)}
        <li>
            <div class="columns">
                <div class="column is-two-thirds">
                    <p>{dataPoint.name}</p>
                </div>
                <div class="column has-text-right">
                    <button class="button is-small" on:click={() => removeDataPoint(dataPoint)}>Remove</button>
                </div>
            </div>
        </li>
    {/each}
</ul>