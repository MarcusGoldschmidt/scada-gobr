<script lang="ts">

    import {ViewComponentType} from "../../shared/types";

    export let type: ViewComponentType = null;
    export let data = {};
</script>

{#if type === ViewComponentType.TimeSeries}
{/if}

{#if type === ViewComponentType.Graphical}
{/if}

{#if type === ViewComponentType.Text}
    <h1>{data.text || "-"}</h1>
{/if}