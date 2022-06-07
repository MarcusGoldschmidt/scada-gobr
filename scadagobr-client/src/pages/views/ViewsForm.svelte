<script lang="ts">

    import MultiTimeSeries from "../../components/charts/MultiTimeSeries.svelte";

    let positionY = 0;
    let positionX = 0;

    let dragging = false

    let reference

    let dragObj

    let m = {x: 0, y: 0};

    $: startDragging = () => {
        if (dragging) {
            return
        }
        dragging = true
        document.body.style.cursor = "grabbing"
    }

    $: handleMousemove = (event) => {
        if (!dragging) {
            return
        }

        m.x = event.clientX;
        m.y = event.clientY;

        const bodyRect = document.body.getBoundingClientRect()
        const elemRect = reference.getBoundingClientRect()
        positionY = m.y - (elemRect.top - bodyRect.top)
        positionX = m.x - (elemRect.left - bodyRect.left)
    }

    $: handleCursorHover = (value) => {
        if (dragging) {
            return
        }
        document.body.style.cursor = value
    }

</script>

<div bind:this={reference}
     class="drag-zone"
     on:mousemove={handleMousemove}>
    <div>
        <span class="fa-solid fa-arrows-up-down-left-right"
              style="padding: 3px;position: relative; top: {positionY}px; left: {positionX}px;"
              on:mouseenter={() => handleCursorHover('grab')}
              on:mouseout={() => handleCursorHover('')}
              on:mouseup={() => {dragging = false; document.body.style.cursor = ''}}
              on:mousedown={startDragging}>
        </span>
        <div bind:this={dragObj}
             style="position: relative; top: {positionY}px; left: {positionX}px;">
            <h1>Teste</h1>
        </div>
    </div>
</div>

<style lang="scss">
  .drag-zone {
    min-height: 94vh;
  }
</style>