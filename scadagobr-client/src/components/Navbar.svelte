<script lang="ts">
    import {authStore} from "../shared/stores/user";
    import {link, navigate} from "svelte-routing";

    let isLogin = false

    authStore.subscribe(value =>
        isLogin = !!value
    )
</script>

{#if (isLogin)}
    <nav class="navbar is-white">
        <div class="navbar-brand">
            <a class="navbar-item brand-text" href="/" use:link>
                GOBR
            </a>
        </div>
        <div id="navMenu" class="navbar-menu">
            <div class="navbar-start">
                <a class="navbar-item" href="/watch-list" use:link>
                    Watch list
                </a>
                <a class="navbar-item" href="/views" use:link>
                    Graphical views
                </a>
                <a class="navbar-item" href="/datasources" use:link>
                    Data sources
                </a>
                <a class="navbar-item" href="/users" use:link>
                    Users
                </a>
                <a class="navbar-item" href="/reports" use:link>
                    Reports
                </a>
            </div>

        </div>

        <div class="navbar-end">
            <div class="navbar-item">
                <div class="buttons">
                    <a class="button is-primary" on:click={() => {
                            authStore.set(null)
                            navigate('/login')
                        }}>
                        <strong>Logout</strong>
                    </a>
                </div>
            </div>
        </div>
    </nav>
{/if}