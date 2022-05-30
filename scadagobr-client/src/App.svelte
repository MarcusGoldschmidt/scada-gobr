<script lang="ts">
    import Login from "./components/auth/Login.svelte";

</script>

<main>

    <Login
    ></Login>
</main>

<style lang="scss">
  :global {
    @import "assets/global";

    @import '@fortawesome/fontawesome-free/css/all.css';
    @import 'bulma/bulma';

  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>