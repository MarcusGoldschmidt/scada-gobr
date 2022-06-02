<script lang="ts">
    import {Route, Router} from "svelte-routing";

    import Login from "./views/auth/Login.svelte";
    import PrivateRoute from "./components/infra/PrivateRoute.svelte";
    import Notifications from "./components/infra/Notifications.svelte";
    import Navbar from "./components/Navbar.svelte";
    import UserIndex from "./views/user/UserIndex.svelte";
    import UserInput from "./views/user/UserInput.svelte";
</script>
<main>
    <Notifications/>
    <Navbar/>

    <Router>
        <Route path="/login">
            <Login/>
        </Route>

        <Route path="/">
            <PrivateRoute>
                <h3>Home</h3>
                <p>Home sweet home...</p>
            </PrivateRoute>
        </Route>

        <Route path="/users">
            <PrivateRoute>
                <UserIndex></UserIndex>
            </PrivateRoute>
        </Route>

        <Route path="/users/create">
            <PrivateRoute>
                <UserInput/>
            </PrivateRoute>
        </Route>

        <Route path="/users/edit/:id" let:params>
            <PrivateRoute>
                <UserInput userId={params.id}/>
            </PrivateRoute>
        </Route>
    </Router>

</main>

<style lang="scss">
  :global {
    @import "assets/global";

    @import '@fortawesome/fontawesome-free/css/all.css';
    @import 'bulma/bulma';

    .content-container {
      padding-left: 2vw;
      padding-right: 2vw;
    }
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>