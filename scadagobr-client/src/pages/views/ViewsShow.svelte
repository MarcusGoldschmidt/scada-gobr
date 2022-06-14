<script lang="ts">
    // Index for the current drag
    import {createEventDispatcher, onMount} from "svelte";
    import ViewComponentView from "../../components/view/ViewComponentView.svelte";

    const dispatch = createEventDispatcher();

    let currentDrag = null

    export let viewsComponents = []

    let reference

    let Xreference = 0
    let Yreference = 0

    onMount(() => {
        const bodyRect = document.body.getBoundingClientRect()
        const elemRect = reference.getBoundingClientRect()

        Xreference = (elemRect.top - bodyRect.top)
        Yreference = (elemRect.left - bodyRect.left)
    })

    $: dragging = currentDrag != null

    $: startDragging = (i) => {
        if (!!currentDrag) {
            return
        }
        currentDrag = i
        document.body.style.cursor = "grabbing"
    }

    $: stopDragging = (i) => {
        currentDrag = null
        document.body.style.cursor = ""
    }

    $: handleMousemove = (event) => {
        if (!dragging) {
            return
        }

        const bodyRect = document.body.getBoundingClientRect()
        const elemRect = reference.getBoundingClientRect()

        const positionY = event.clientY
        const positionX = event.clientX

        viewsComponents[currentDrag] = {
            ...viewsComponents[currentDrag],
            x: positionX,
            y: positionY
        }


        viewsComponents = viewsComponents
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
        {#each viewsComponents as view, i (view.id)}
            <div style="position: absolute; top: {view.y || 0}px; left: {view.x || 0}px; z-index: 999">
                <span class="relative fa-solid fa-arrows-up-down-left-right"
                      on:mouseenter={() => handleCursorHover('grab')}
                      on:mouseout={() => handleCursorHover('')}
                      on:mouseup={() => stopDragging(i)}
                      on:mousedown={() => startDragging(i)}>
                </span>

                <span class="fa-solid fa-pencil"
                      style="cursor: pointer;"
                      on:click={() => dispatch('edit', {view, index: i})}
                ></span>

                <ViewComponentView
                        type={view.type}
                        data={view.data}
                ></ViewComponentView>
            </div>
        {/each}
    </div>
</div>

<style lang="scss">
  .drag-zone {
    min-height: 94vh;
  }

  .relative {
    position: relative;
  }
</style>