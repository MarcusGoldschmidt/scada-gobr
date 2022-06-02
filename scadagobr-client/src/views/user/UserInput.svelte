<script lang="ts">
    import {Field, Input, Switch} from 'svelma'
    import * as yup from 'yup';
    import {getErrorsFromSchema, inputError} from "../../shared/utils";
    import {navigate} from "svelte-routing";
    import FieldHorizontal from "../../components/FieldHorizontal.svelte";
    import {axiosJwt} from "../../shared/http/axios";
    import {PathsV1} from "../../shared/http/request";
    import {onMount} from "svelte";
    import {sendNotification} from "../../shared/stores/notifications";

    export let userId = null
    let username = ""
    let password = ""
    let email = ""
    let administrator = false
    let homeUrl = ""
    let errors = {}

    let schema = yup.object().shape({
        username: yup.string().required(),
        email: yup.string().email(),
        homeUrl: yup.string(),
        administrator: yup.boolean().required(),
        password: userId ? yup.string() : yup.string().required(),
    });

    // fetch user data from server
    onMount(async () => {
        if (userId) {
            const response = await axiosJwt.get(PathsV1.UserGetById + userId)
            const user = response.data

            // set user properties with the response of user
            username = user.name
            email = user.email
            homeUrl = user.homeUrl
            administrator = user.administrator
            password = user.password
        }
    });

    $: onSubmit = async () => {

        const err = await getErrorsFromSchema(schema, {username, password, email, administrator, homeUrl});

        if (err) {
            errors = err
            return
        }
        errors = {}

        const body = {
            name: username,
            password,
            email,
            administrator,
            homeUrl
        };

        let promise

        if (userId) {
            promise = axiosJwt.put(PathsV1.UserUpdate + userId, body)
        } else {
            promise = axiosJwt.post(PathsV1.UserCreate, body)
        }

        promise.then(() => {
            navigate("/users")
        }).catch(err => {
            sendNotification("error", err.response.data.message)
        })
    }
</script>

<section class="hero">
    <div class="hero-body">
        <div class="container">
            <div class="columns is-5-tablet is-4-desktop is-3-widescreen">
                <div class="column">
                    <form class="box" on:submit|preventDefault={onSubmit}>
                        <div class="field has-text-centered">
                            <h2>{userId ? "Edit" : "Create"} user</h2>
                        </div>

                        <FieldHorizontal label="Username">
                            <Field {...inputError(errors.username)}>
                                <Input class="input {inputError(errors.username).type}" bind:value={username} type=""
                                       placeholder="Username"/>
                            </Field>
                        </FieldHorizontal>

                        <FieldHorizontal label="Email">
                            <Field {...inputError(errors.email)}>
                                <Input class="input {inputError(errors.email).type}" bind:value={email} type=""
                                       placeholder="Email"/>
                            </Field>
                        </FieldHorizontal>

                        <FieldHorizontal label="Home Url">
                            <Field {...inputError(errors.homeUrl)}>
                                <Input class="input {inputError(errors.homeUrl).type}" bind:value={homeUrl} type=""
                                       placeholder=""/>
                            </Field>
                        </FieldHorizontal>

                        <FieldHorizontal label="Administrator">
                            <Switch bind:checked={administrator}></Switch>
                        </FieldHorizontal>

                        {#if (!userId)}
                            <FieldHorizontal label="Password">
                                <Field {...inputError(errors.password)}>
                                    <Input type="password" bind:value={password}
                                           class={inputError(errors.password).type}
                                           passwordReveal={true}
                                           placeholder="Password"/>
                                </Field>
                            </FieldHorizontal>
                        {/if}

                        <div class="columns">
                            <div class="column is-four-fifths">

                            </div>
                            <div class="column is-right">
                                <button class="button is-warning" on:click={() => navigate('/users')}>
                                    Cancel <i class="fa fa-sign-in" aria-hidden="true"></i>
                                </button>
                            </div>
                            <div class="column">
                                <button class="button is-success" type="submit">
                                    {userId ? "Edit" : "Create"} <i class="fa fa-sign-in" aria-hidden="true"></i>
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</section>