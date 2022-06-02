<script lang="ts">
    export let label
</script>

<div class="field is-horizontal">
    <div class="field-label is-normal">
        <label class="label">{label || ""}</label>
    </div>
    <div class="field-body">
        <slot></slot>
    </div>
</div>