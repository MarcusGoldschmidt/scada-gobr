<script lang="ts">

    import {Field, Input} from 'svelma'
    import {authStore} from "../../shared/stores/user";
    import {later} from "../../shared/utils";

    let email = "";
    let password = ""

    $: onLogin = async () => {
        await later(1000)

        authStore.set({
            jts: {
                refreshToken: "",
                token: ""
            },
            user: {
                email: email,
                id: "id"
            }
        })
    }
</script>

<main>
    <i class="fa fa-circle-o-notch fa-spin" style="font-size:24px"></i>
    <section class="hero has-background-light is-fullheight">
        <div class="hero-body">
            <div class="container">
                <div class="column is-4 is-offset-4">
                    <h3 class="title has-text-black has-text-centered">Scada-gobr</h3>
                    <hr class="login-hr">
                    <p class="subtitle has-text-black has-text-centered">Please login to proceed.</p>
                    <div class="box">
                        <form on:submit|preventDefault={onLogin}>
                            <Field type="is-danger" message="Email is invalid">
                                <Input class="input" type="email" placeholder="Your Email" autofocus=""/>
                            </Field>

                            <Field>
                                <Input type="password" passwordReveal={true} placeholder="Your Password"/>
                            </Field>
                            <button class="button is-block is-info is-fullwidth">
                                Login <i class="fa fa-sign-in" aria-hidden="true"></i>
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </section>
</main>