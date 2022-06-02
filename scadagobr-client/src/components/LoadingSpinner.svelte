<script lang="ts">
    export let size = 5;
    export let sizeUnit = "vw"
    export let density = "0.4rem"

    $: cssSize = size.toString() + sizeUnit

</script>

<div class="loader-wrapper">
    <div class="loader" style="height: {cssSize}; width: {cssSize}; border-width: {density}"></div>
</div>

<style lang="scss">
  @import "bulma/sass/utilities/derived-variables";

  .loader-wrapper {
    display: flex;
    justify-content: center;
    height: 100%;
    width: 100%;
    transition: opacity .1s;
    align-items: center;
    border-radius: 6px;

    .loader {
      $color: $primary;

      border-left-color: $primary;
      border-bottom-color: $color;
    }

  }
</style>