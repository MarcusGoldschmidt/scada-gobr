<script lang="ts">
    import Table from "../../components/Table.svelte";
    import {ModalCard} from 'svelma'
    import {axiosJwt} from "../../shared/net/axios";
    import {PathsV1} from "../../shared/net/request";
    import {link} from "svelte-routing";
    import {NotificationType, sendNotification} from "../../shared/stores/notifications";

    $: fun = async (x) => {
        const {data: users} = await axiosJwt.get(PathsV1.ViewGet, {
            params: {page: x, size: 20}
        })
        return users
    };

    function deleteView(id) {
        axiosJwt.delete(PathsV1.ViewDelete + id)
            .then(e => {
                fun = fun
                sendNotification("successful deleted view", "", NotificationType.Success)
            })
            .catch(e => {
                sendNotification(e.response.data.error, "", NotificationType.Danger)
            })
    }

    let viewDeleteId: string | null = null;

    $: modelActive = !!viewDeleteId
</script>

<ModalCard
        bind:active={modelActive} title="My Modal Title"
        on:success={() => viewDeleteId = null}
>

</ModalCard>

<div class="content-container">
    <div class="columns">
        <div class="column">
            <a href="/views/create" use:link class="button is-primary">Create view</a>
        </div>
    </div>

    <Table paginateFunction={fun} useQueryString={true}>
        <tr slot="head">
            <th>Nome</th>
            <th></th>
            <th></th>
        </tr>

        <tr slot="items" let:item>
            <td>{item.name}</td>
            <td>
                <a class="button is-primary" href="/views/{item.id}">Edit</a>
            </td>
            <td>
                <button class="button is-danger" on:click={() => deleteView(item.id)}>Remove</button>
            </td>
        </tr>
    </Table>
</div>