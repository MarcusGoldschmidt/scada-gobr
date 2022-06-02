<script lang="ts" generic="T">

    import {LoadingStatus} from "../shared/types";
    import LoadingSpinner from "./LoadingSpinner.svelte";

    let items: T[] = []
    export let paginateFunction: (page: number) => Promise<T[]>
    export let useQueryString: boolean = false

    let currentPage = 1
    let loadingStatus: LoadingStatus = LoadingStatus.Fetching

    $: {
        if (useQueryString) {
            const searchParams = new URLSearchParams(window.location.search);
            currentPage = +searchParams.get("page")
            if (currentPage < 1) {
                currentPage = 1
            }
        }

        loadingStatus = LoadingStatus.Fetching
        paginateFunction(currentPage)
            .then(e => {
                loadingStatus = LoadingStatus.Success
                items = e
            })
            .catch(e => {
                loadingStatus = LoadingStatus.Error
            })
    }

    $: goToPage = async (page: number) => {
        if (useQueryString) {
            const searchParams = new URLSearchParams(window.location.search);
            searchParams.set("page", page.toString());

            const newRelativePathQuery = window.location.pathname + '?' + searchParams.toString();
            history.pushState(null, '', newRelativePathQuery);
        }

        currentPage = page
    }

    $: next = async () => {
        if (items.length === 0) {
            return
        }

        await goToPage(currentPage + 1)
    }

    $: previous = async () => {
        if (currentPage === 1) {
            return
        }

        await goToPage(currentPage - 1)
    }
</script>

{#if (loadingStatus === LoadingStatus.Fetching)}
    <LoadingSpinner
            size="{20}"
    ></LoadingSpinner>
{:else if (loadingStatus === LoadingStatus.Error)}
    <h1>Error</h1>
{:else}
    <table class="table">
        <thead>
        <slot name="head"/>
        </thead>
        <tbody>
        {#each items as item (item.id)}
            <slot name="items" {item}/>
        {/each}
        </tbody>
    </table>
    <div style="display: flex; justify-content: center">
        <a class="pagination-previous {currentPage === 1 ? 'is-disabled' : ''}"
           title="This is the first page"
           on:click={previous}>
            Previous
        </a>
        <a class="pagination-next {items.length === 0 ? 'is-disabled' : ''}" on:click={next}>Next page</a>
    </div>
{/if}

