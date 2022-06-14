<script lang="ts">
    import {ViewComponentType} from "../../shared/types";
    import TextForm from "./TextForm.svelte";
    import {createEventDispatcher} from 'svelte';
    import TimeSeriesForm from "./TimeSeriesForm.svelte";

    const dispatch = createEventDispatcher();

    export let type: ViewComponentType = null;
    export let data: any = {};

    function onChange(e) {
        dispatch('change', e.detail);
    }

</script>

{#if type === ViewComponentType.TimeSeries}
    <TimeSeriesForm
            defaultPeriod={data.period}
            defaultWidth={data.width}
            selectedDataPointsIds={data.dataPointsIds}
            on:change={onChange}
    ></TimeSeriesForm>
{/if}

{#if type === ViewComponentType.Graphical}
{/if}

{#if type === ViewComponentType.Text}
    <TextForm
            text={data.text || ""}
            on:change={onChange}
    ></TextForm>
{/if}