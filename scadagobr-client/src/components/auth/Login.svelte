<script lang="ts">
    import {Field, Input} from 'svelma'
    import * as yup from 'yup';
    import {globalHistory} from "svelte-routing/src/history";
    import {loginUser} from "../../shared/http/axios";
    import {getErrorsFromSchema, inputError} from "../../shared/utils";
    import {navigate} from "svelte-routing";
    import {authStore} from "../../shared/stores/user";
    import {onMount} from "svelte";

    let schema = yup.object().shape({
        username: yup.string().required(),
        password: yup.string().required(),
    });

    let username = "";
    let password = ""

    let errors = {}

    onMount(() => {
        if ($authStore) {
            navigate("/", {replace: true})
        }
    })

    $: onLogin = async () => {

        const err = await getErrorsFromSchema(schema, {username, password})

        if (err) {
            errors = err
            return
        }
        errors = {}

        const success = await loginUser(username, password)

        if (success) {
            navigate(globalHistory.location.state.from)
        }

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
                            <Field {...inputError(errors.username)}>
                                <Input class="input {inputError(errors.username).type}" bind:value={username} type=""
                                       placeholder="Username"/>
                            </Field>

                            <Field {...inputError(errors.password)}>
                                <Input type="password" bind:value={password}
                                       class={inputError(errors.password).type}
                                       passwordReveal={true}
                                       placeholder="Your Password"/>
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