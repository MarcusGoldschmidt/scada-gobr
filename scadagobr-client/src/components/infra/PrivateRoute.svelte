<script>
    import {onMount} from 'svelte'
    import {navigate} from "svelte-routing";
    import {authStore} from "../../shared/stores/user";

    onMount(() => {
        if (!$authStore) {

            navigate("/login", {
                state: {
                    from: window.location.pathname
                }
            });
        }
    })
</script>

{#if $authStore}
    <slot/>
{/if}