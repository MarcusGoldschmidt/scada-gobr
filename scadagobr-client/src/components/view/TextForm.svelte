<script lang="ts">
    import {Field, Input} from 'svelma'
    import {createEventDispatcher} from 'svelte';

    const dispatch = createEventDispatcher();

    export let text = ""


</script>

<Field class="mt-5" label="Text">
    <Input class="input"
           bind:value={text}
           on:change={(e) => dispatch('change', {text: e.target.value})}
           type="text"
           placeholder="Text"/>
</Field>