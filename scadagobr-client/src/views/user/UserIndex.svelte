<script lang="ts">
    import Table from "../../components/Table.svelte";
    import {ModalCard} from 'svelma'
    import {axiosJwt} from "../../shared/http/axios";
    import {PathsV1} from "../../shared/http/request";
    import {link} from "svelte-routing";
    import {NotificationType, sendNotification} from "../../shared/stores/notifications";

    $: fun = async (x) => {
        const {data: users} = await axiosJwt.get(PathsV1.UserGet, {
            params: {page: x, size: 20}
        })
        return users
    };

    function deleteUser(id) {
        axiosJwt.delete(PathsV1.UserDelete + id)
            .then(e => {
                fun = fun
                sendNotification("successful deleted user", "", NotificationType.Success)
            })
            .catch(e => {
                sendNotification(e.response.data.error, "", NotificationType.Danger)
            })
    }

    let userDeleteId: string | null = null;

    $: modelActive = !!userDeleteId
</script>

<ModalCard
        bind:active={modelActive} title="My Modal Title"
        on:success={() => userDeleteId = null}

>

</ModalCard>

<div class="content-container">
    <div class="columns">
        <div class="column">
            <a href="/users/create" use:link class="button is-primary">Create user</a>
        </div>
    </div>

    <Table paginateFunction={fun} useQueryString={true}>
        <tr slot="head">
            <th>Nome</th>
            <th>Email</th>
            <th>Home Url</th>
            <th></th>
            <th></th>
        </tr>

        <tr slot="items" let:item>
            <td>{item.name}</td>
            <td>{item.email || "-"}</td>
            <td>{item.homeUrl || "-"}</td>
            <td>
                <a class="button is-primary" href="/users/edit/{item.id}">Edit</a>
            </td>
            <td>
                <button class="button is-danger" on:click={() => deleteUser(item.id)}>Remove</button>
            </td>
        </tr>
    </Table>
</div>