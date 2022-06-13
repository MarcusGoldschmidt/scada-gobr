<script lang="ts">
    import * as yup from 'yup';
    import {Field, Input} from 'svelma'
    import {getErrorsFromSchema, inputError} from "../../shared/utils";
    import {axiosJwt} from "../../shared/net/axios";
    import {PathsV1} from "../../shared/net/request";
    import {onMount} from "svelte";
    import {fly} from 'svelte/transition';
    import FieldHorizontal from "../../components/FieldHorizontal.svelte";
    import {ViewComponentType} from "../../shared/types";
    import ViewComponentForm from "../../components/view/ViewComponentForm.svelte";
    import ViewsShow from "./ViewsShow.svelte";
    import {v4 as uuidv4} from 'uuid';
    import {NotificationType, sendNotification} from "../../shared/stores/notifications";


    export let viewId = null

    onMount(async () => {
        if (viewId) {
            const response = await axiosJwt.get(PathsV1.ViewGetById + viewId)
            const view = response.data

            name = view.name
            viewComponents = view.viewComponents.map(e => ({
                ...e,
                x: e.position.x,
                y: e.position.y,
                type: e.viewType
            }))
        }
    });

    // Common view fields
    let name = ""
    let viewComponents: any[] = []
    let errors = {}
    let schema = yup.object().shape({
        name: yup.string().required(),
    });

    // Handler Ui forms to the right
    let showViewForm = false
    let showViewComponentsForm = false
    let showHelper = true

    function resetHelper() {
        showHelper = true
        showViewForm = false
        showViewComponentsForm = false
    }

    function closeHelper() {
        showHelper = false
        selectedComponentId = null
        componentsFormData = {}
    }

    // form data for view components
    let selectedType: ViewComponentType = ViewComponentType.TimeSeries
    let componentsFormData = {}
    let selectedComponentId = null

    $: addComponent = () => {
        const x = window.innerWidth / 2;
        const y = window.innerHeight / 2;

        if (selectedComponentId) {
            viewComponents = viewComponents.map(component => {
                if (component.id === selectedComponentId) {
                    return {
                        ...component,
                        ...componentsFormData
                    }
                }
                return component
            })

            return
        }

        viewComponents = [...viewComponents, {type: selectedType, x, y, ...componentsFormData, id: uuidv4()}]
    }

    $: onSave = async () => {
        const err = await getErrorsFromSchema(schema, {name});

        if (err) {
            errors = err
            resetHelper()
            showViewForm = true
            return
        }
        errors = {}

        if (viewId) {
            axiosJwt.put(PathsV1.ViewUpdate + viewId, {
                name,
                viewComponents
            }).then(() => {
                resetHelper()
                showViewForm = false
                sendNotification("View edited successfully", "", NotificationType.Info)
            }).catch(() => {
                sendNotification("Error editing the view", "", NotificationType.Danger)
            })
        } else {
            axiosJwt.post(PathsV1.ViewCreate, {
                name,
                viewComponents
            }).then(() => {
                resetHelper()
                showViewForm = false
                sendNotification("View created successfully", "", NotificationType.Info)
            }).catch(() => {
                sendNotification("Error creating the view", "", NotificationType.Danger)
            })
        }


    }
    let removeComponent = async () => {
        if (selectedComponentId == null) {
            return
        }

        if (viewId) {
            try {
                await axiosJwt.delete(PathsV1.ViewComponentDelete(viewId, selectedComponentId))
                sendNotification("Data point deleted successfully", "", NotificationType.Info)
            } catch (e) {
                sendNotification("Error deleting the data point", "", NotificationType.Danger)
            }
        }

        viewComponents = viewComponents.filter(e => e.id !== selectedComponentId)
        resetHelper()
    };
</script>

<svelte:window on:keydown={(event) => {
    if (event.key === 'Escape') {
        resetHelper()
    }
}}/>


{#if showHelper}
    <div class="has-text-right"
         style="position: absolute; right: 0; z-index: 999; font-size: 2rem"
         in:fly="{{ x: 200, duration: 400 }}"
         out:fly="{{ x: 200, duration: 500 }}"
    >
        <i class="fa-solid fa-pencil mt-5"
           style="cursor: pointer;"
           on:click={() => {
                closeHelper()
                showViewForm = true
           }}
        ></i>
        <br>
        <i class="fa-solid fa-circle-plus mt-5"
           style="cursor: pointer;"
           on:click={() => {
                closeHelper()
                showViewComponentsForm = true
           }}
        ></i>
        <br>
        <i class="fa-solid fa-floppy-disk mt-5"
           style="cursor: pointer;"
           on:click={onSave}
        ></i>
    </div>
{/if}

{#if showViewForm}
    <div
            style="position: absolute; right: 0; z-index: 999"
            in:fly="{{ x: 200, duration: 400 }}"
            out:fly="{{ x: 200, duration: 500 }}">
        <div class="box">
            <i class="fa-solid fa-xmark"
               on:click={resetHelper}
               style="font-size: 2rem; cursor: pointer"></i>

            <div class="field has-text-left">
                <h2>View</h2>
            </div>

            <div class="columns">
                <div class="column is-full">
                    <FieldHorizontal label="Name">
                        <Field {...inputError(errors.name)}>
                            <Input class="input {inputError(errors.name).type}" bind:value={name} type=""
                                   placeholder="View Name"/>
                        </Field>
                    </FieldHorizontal>
                </div>
            </div>
        </div>
    </div>
{/if}

{#if showViewComponentsForm}
    <div
            style="position: absolute; right: 0; z-index: 999; min-width: 20%"
            in:fly="{{ x: 200, duration: 400 }}"
            out:fly="{{ x: 200, duration: 500 }}">
        <div class="box">
            <i class="fa-solid fa-xmark"
               on:click={resetHelper}
               style="font-size: 2rem"></i>

            <div class="columns">
                <div class="column is-two-thirds has-text-left">
                    <h2>Component</h2>
                </div>
                {#if selectedComponentId}
                    <div class="column has-text-right">
                        <button class="button is-danger" on:click={removeComponent}>Remove</button>
                    </div>
                {/if}
            </div>


            <div class="field has-text-left">
            </div>
            {#if !selectedComponentId}
                <div class="select" style="width: 100%">
                    <select bind:value={selectedType} style="width: 100%">
                        <option>Select dropdown</option>
                        <option value={ViewComponentType.Text}>Teste</option>
                        <option value={ViewComponentType.TimeSeries}>Time Series</option>
                        <option value={ViewComponentType.Graphical}>Graphical</option>
                    </select>
                </div>
            {/if}

            <br>
            <ViewComponentForm
                    on:change={e => componentsFormData.data = e.detail}
                    bind:data={componentsFormData.data}
                    bind:type={selectedType}></ViewComponentForm>
            <br>

            <div class="has-text-right">
                <button class="button is-primary mt-4"
                        on:click={addComponent}>
                    {selectedComponentId ? "Edit" : "Add"}
                </button>
            </div>
        </div>
    </div>
{/if}

<ViewsShow
        on:edit={(e) => {
            closeHelper()
            selectedComponentId = e.detail.view.id
            componentsFormData = e.detail.view
            selectedType = e.detail.view.type
            showViewComponentsForm = true
        }}
        bind:viewsComponents={viewComponents}
></ViewsShow>