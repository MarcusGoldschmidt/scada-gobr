<script lang="ts">
    import {createEventDispatcher, onMount} from 'svelte';
    import {axiosJwt} from "../../shared/net/axios";
    import {PathsV1} from "../../shared/net/request";
    import {NotificationType, sendNotification} from "../../shared/stores/notifications";

    const dispatch = createEventDispatcher();

    let dataSetOptions: any[] = []

    let dataPointsOptions: any[] = []

    let currentDataSet: any = null

    export let selectedDataPoints: any[] = []
    export let defaultPeriod = 60
    export let defaultWidth = 300

    let selectedDataPoint = null

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

        if (selectedDataPoints.findIndex(x => x === selectedDataPoint) == -1) {
            selectedDataPoints = [...selectedDataPoints, selectedDataPoint]
            onChange()
        } else {
            sendNotification("Data point already selected", "", NotificationType.Warning)
        }
    }

    let removeDataPoint = (e) => {
        if (selectedDataPoint == null) {
            return
        }

        selectedDataPoints = selectedDataPoints.filter(x => x !== selectedDataPoint)
        onChange()
    };
</script>

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
                    <button class="button is-small" on:click={removeDataPoint}>Remove</button>
                </div>
            </div>
        </li>
    {/each}
</ul>