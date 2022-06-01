<script lang="ts">
    import {Route, Router} from "svelte-routing";

    import Login from "./components/auth/Login.svelte";
    import PrivateRoute from "./components/infra/PrivateRoute.svelte";
    import Notifications from "./components/infra/Notifications.svelte";
    import Navbar from "./components/Navbar.svelte";
</script>
<main>
    <Notifications/>
    <Navbar></Navbar>

    <Router>

        <Route path="/">
            <PrivateRoute>
                <h3>Home</h3>
                <p>Home sweet home...</p>
            </PrivateRoute>
        </Route>

        <Route path="/login">
            <Login/>
        </Route>
    </Router>

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