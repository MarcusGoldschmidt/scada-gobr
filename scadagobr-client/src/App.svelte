<script lang="ts">
    import {Route, Router} from "svelte-routing";

    import Login from "./pages/auth/Login.svelte";
    import PrivateRoute from "./components/infra/PrivateRoute.svelte";
    import Notifications from "./components/infra/Notifications.svelte";
    import Navbar from "./components/Navbar.svelte";
    import UserIndex from "./pages/user/UserIndex.svelte";
    import UserInput from "./pages/user/UserInput.svelte";

    import '@fortawesome/fontawesome-free/css/all.css';
    import "@carbon/charts/styles.min.css";
    import "carbon-components/css/carbon-components.min.css";
    import ViewsForm from "./pages/views/ViewsForm.svelte";
</script>

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

    <Route path="/views">
        <PrivateRoute>
            <ViewsForm/>
        </PrivateRoute>
    </Route>

    <Route path="/views/create">
        <PrivateRoute>
            <UserInput/>
        </PrivateRoute>
    </Route>

    <Route path="/views/edit/:id" let:params>
        <PrivateRoute>
            <UserInput userId={params.id}/>
        </PrivateRoute>
    </Route>
</Router>

<style lang="scss">
  :global {
    @import "assets/global";

    @import 'bulma/bulma';

    .content-container {
      padding-left: 2vw;
      padding-right: 2vw;
    }
  }

</style>